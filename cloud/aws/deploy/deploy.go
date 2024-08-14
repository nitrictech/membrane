// Copyright Nitric Pty Ltd.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package deploy

import (
	"fmt"
	"strings"

	_ "embed"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	lambdaClient "github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
	"github.com/nitrictech/nitric/cloud/aws/common"
	"github.com/nitrictech/nitric/cloud/common/deploy"
	"github.com/nitrictech/nitric/cloud/common/deploy/provider"
	"github.com/nitrictech/nitric/cloud/common/deploy/pulumix"
	"github.com/nitrictech/nitric/cloud/common/deploy/tags"
	resourcespb "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/apigatewayv2"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/batch"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/dynamodb"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ecr"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/rds"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/resourcegroups"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/s3"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/secretsmanager"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/sqs"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/codebuild"
	awsec2 "github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi-awsx/sdk/go/awsx/ec2"
	"github.com/pulumi/pulumi-docker/sdk/v4/go/docker"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NitricAwsPulumiProvider struct {
	*deploy.CommonStackDetails

	StackId   string
	AwsConfig *AwsConfig

	SqlDatabases map[string]*RdsDatabase

	DockerProvider   *docker.Provider
	Vpc              *ec2.Vpc
	VpcAzs           []string
	VpcSecurityGroup *awsec2.SecurityGroup

	ComputeEnvironment *batch.ComputeEnvironment
	JobQueue           *batch.JobQueue
	BatchExecutionRole *iam.Role

	// A codebuild job for creating the requested databases for a single database cluster
	DbMasterPassword      *random.RandomPassword
	CreateDatabaseProject *codebuild.Project
	CodeBuildRole         *iam.Role
	// A map of unique image keys to database migration codebuild projects
	DatabaseMigrationJobs map[string]*codebuild.Project
	DatabaseCluster       *rds.Cluster
	RdsPrxoy              *rds.Proxy
	EcrAuthToken          *ecr.GetAuthorizationTokenResult
	Lambdas               map[string]*lambda.Function
	LambdaRoles           map[string]*iam.Role
	BatchRoles            map[string]*iam.Role
	HttpProxies           map[string]*apigatewayv2.Api
	Apis                  map[string]*apigatewayv2.Api
	Secrets               map[string]*secretsmanager.Secret
	Buckets               map[string]*s3.Bucket
	BucketNotifications   map[string]*s3.BucketNotification
	Topics                map[string]*topic
	Queues                map[string]*sqs.Queue
	Websockets            map[string]*apigatewayv2.Api
	KeyValueStores        map[string]*dynamodb.Table

	provider.NitricDefaultOrder

	ResourceTaggingClient *resourcegroupstaggingapi.ResourceGroupsTaggingAPI
	LambdaClient          lambdaiface.LambdaAPI
}

var _ provider.NitricPulumiProvider = (*NitricAwsPulumiProvider)(nil)

const pulumiAwsVersion = "6.6.0"

func (a *NitricAwsPulumiProvider) Config() (auto.ConfigMap, error) {
	return auto.ConfigMap{
		"aws:region":     auto.ConfigValue{Value: a.Region},
		"aws:version":    auto.ConfigValue{Value: pulumiAwsVersion},
		"docker:version": auto.ConfigValue{Value: deploy.PulumiDockerVersion},
	}, nil
}

func (a *NitricAwsPulumiProvider) Init(attributes map[string]interface{}) error {
	var err error

	a.CommonStackDetails, err = deploy.CommonStackDetailsFromAttributes(attributes)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, err.Error())
	}

	a.AwsConfig, err = ConfigFromAttributes(attributes)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "Bad stack configuration: %s", err)
	}

	return nil
}

func (a *NitricAwsPulumiProvider) Pre(ctx *pulumi.Context, resources []*pulumix.NitricPulumiResource[any]) error {
	// make our random stackId
	stackRandId, err := random.NewRandomString(ctx, fmt.Sprintf("%s-stack-name", ctx.Stack()), &random.RandomStringArgs{
		Special: pulumi.Bool(false),
		Length:  pulumi.Int(8),
		Keepers: pulumi.ToMap(map[string]interface{}{
			"stack-name": ctx.Stack(),
		}),
	})
	if err != nil {
		return err
	}

	stackIdChan := make(chan string)
	pulumi.Sprintf("%s-%s", ctx.Stack(), stackRandId.Result).ApplyT(func(id string) string {
		stackIdChan <- id
		return id
	})

	a.StackId = <-stackIdChan

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Region: aws.String(a.Region)},
		SharedConfigState: session.SharedConfigEnable,
	}))

	a.ResourceTaggingClient = resourcegroupstaggingapi.New(sess)

	a.LambdaClient = lambdaClient.New(sess, &aws.Config{Region: aws.String(a.Region)})

	a.EcrAuthToken, err = ecr.GetAuthorizationToken(ctx, &ecr.GetAuthorizationTokenArgs{})
	if err != nil {
		return err
	}

	a.DockerProvider, err = docker.NewProvider(ctx, "docker-auth-provider", &docker.ProviderArgs{
		RegistryAuth: &docker.ProviderRegistryAuthArray{
			docker.ProviderRegistryAuthArgs{
				Address:  pulumi.String(a.EcrAuthToken.ProxyEndpoint),
				Username: pulumi.String(a.EcrAuthToken.UserName),
				Password: pulumi.String(a.EcrAuthToken.Password),
			},
		},
	})
	if err != nil {
		return err
	}

	// Create AWS Resource groups with our tags
	_, err = resourcegroups.NewGroup(ctx, "stack-resource-group", &resourcegroups.GroupArgs{
		Name:        pulumi.String(a.FullStackName),
		Description: pulumi.Sprintf("All deployed resources for the %s nitric stack", a.FullStackName),
		ResourceQuery: &resourcegroups.GroupResourceQueryArgs{
			Query: pulumi.Sprintf(`{
				"ResourceTypeFilters":["AWS::AllSupported"],
				"TagFilters":[{"Key":"%s"}]
			}`, tags.GetResourceNameKey(a.StackId)),
		},
	})

	jobs := lo.Filter(resources, func(item *pulumix.NitricPulumiResource[any], idx int) bool {
		return item.Id.Type == resourcespb.ResourceType_Batch
	})

	databases := lo.Filter(resources, func(item *pulumix.NitricPulumiResource[any], idx int) bool {
		return item.Id.Type == resourcespb.ResourceType_SqlDatabase
	})
	// Create a shared database cluster if we have more than one database
	if len(databases) > 0 {
		// deploy a VPC and security groups for this database cluster
		err := a.vpc(ctx)
		if err != nil {
			return err
		}

		// deploy the RDS cluster
		err = a.rds(ctx)
		if err != nil {
			return err
		}
	}

	if len(jobs) > 0 {
		var subnets pulumi.StringArrayOutput
		var vpcId pulumi.StringOutput

		if a.Vpc != nil {
			subnets = a.Vpc.PrivateSubnetIds
			vpcId = a.Vpc.VpcId
		} else {
			vpc, err := ec2.NewDefaultVpc(ctx, "default-vpc", nil)
			if err != nil {
				return fmt.Errorf("could not resolve default VPC")
			}
			subnets = vpc.PublicSubnetIds
			vpcId = vpc.VpcId
		}

		batchSecurityGroup, err := awsec2.NewSecurityGroup(ctx, "VpcSecurityGroup", &awsec2.SecurityGroupArgs{
			VpcId: vpcId,
			Egress: awsec2.SecurityGroupEgressArray{
				awsec2.SecurityGroupEgressArgs{
					Protocol:   pulumi.String("-1"),
					FromPort:   pulumi.Int(0),
					ToPort:     pulumi.Int(0),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
			},
			Ingress: awsec2.SecurityGroupIngressArray{},
		})
		if err != nil {
			return err
		}

		a.BatchExecutionRole, err = iam.NewRole(ctx, "BatchExecutionRole", &iam.RoleArgs{
			AssumeRolePolicy: pulumi.String(`{
				"Version": "2012-10-17",
				"Statement": [
					{
						"Action": "sts:AssumeRole",
						"Principal": {
							"Service": "batch.amazonaws.com"
						},
						"Effect": "Allow",
						"Sid": ""
					}
				]
			}`),
		})
		if err != nil {
			return err
		}

		_, err = iam.NewRolePolicy(ctx, "BatchExecutionRolePolicy", &iam.RolePolicyArgs{
			Role: a.BatchExecutionRole.ID(),
			Policy: pulumi.String(`{
                "Version": "2012-10-17",
                "Statement": [
                    {
                        "Effect": "Allow",
                        "Action": [
                            "logs:CreateLogGroup",
                            "logs:CreateLogStream",
                            "logs:PutLogEvents",
							"ecr:GetAuthorizationToken",
                            "ecr:BatchCheckLayerAvailability",
                            "ecr:GetDownloadUrlForLayer",
                            "ecr:BatchGetImage",
                            "ecr:DescribeRepositories",
                            "ecr:ListImages",
							"batch:SubmitJob"
                        ],
                        "Resource": "*"
                    }
                ]
            }`),
		})
		if err != nil {
			return err
		}

		instanceProfile, err := iam.NewInstanceProfile(ctx, "BatchInstanceProfile", &iam.InstanceProfileArgs{
			Role: a.BatchExecutionRole.Name,
		})
		if err != nil {
			return err
		}

		computeResourceOptions := &batch.ComputeEnvironmentComputeResourcesArgs{
			// AllocationStrategy: pulumi.String("BEST_FIT"),
			// MinVcpus:           pulumi.Int(0),
			MaxVcpus: pulumi.Int(8),
			// TODO Determine EC2 configuration

			// TODO: Make launchable instance types configurable from stack configuration
			InstanceTypes: pulumi.StringArray{
				// Small long running instances
				pulumi.String("optimal"),
				// Allow GPU capable instances
				pulumi.String("g3s.xlarge"),
			},
			Type:             pulumi.String("EC2"),
			Subnets:          subnets,
			SecurityGroupIds: pulumi.StringArray{batchSecurityGroup.ID()},
			InstanceRole:     instanceProfile.Arn,
		}

		// TODO: Possibly use a shared stack compute environment?
		a.ComputeEnvironment, err = batch.NewComputeEnvironment(ctx, "compute-environment", &batch.ComputeEnvironmentArgs{
			ComputeEnvironmentName: pulumi.Sprintf("%s-compute-environment", a.StackId),
			ComputeResources:       computeResourceOptions,
			Type:                   pulumi.String("MANAGED"),
		})
		if err != nil {
			return fmt.Errorf("error creating compute environment: %w", err)
		}
		// TODO: Possibly use a shared stack job queue?
		a.JobQueue, err = batch.NewJobQueue(ctx, "job-queue", &batch.JobQueueArgs{
			ComputeEnvironments: pulumi.StringArray{
				a.ComputeEnvironment.Arn,
			},
			State:    pulumi.String("ENABLED"),
			Priority: pulumi.Int(1),
			// TODO: Set tags for job definition discovery
			Tags: pulumi.ToStringMap(tags.Tags(a.StackId, "job-queue", "job-queue")),
		})
		if err != nil {
			return err
		}
	}

	return err
}

func (a *NitricAwsPulumiProvider) Post(ctx *pulumi.Context) error {
	return nil
}

func (a *NitricAwsPulumiProvider) Result(ctx *pulumi.Context) (pulumi.StringOutput, error) {
	outputs := []interface{}{}

	// Add APIs outputs
	if len(a.Apis) > 0 {
		outputs = append(outputs, pulumi.Sprintf("API Endpoints:\n──────────────"))
		for apiName, api := range a.Apis {
			outputs = append(outputs, pulumi.Sprintf("%s: %s", apiName, api.ApiEndpoint))
		}
	}

	// Add HTTP Proxy outputs
	if len(a.HttpProxies) > 0 {
		if len(outputs) > 0 {
			outputs = append(outputs, "\n")
		}
		outputs = append(outputs, pulumi.Sprintf("HTTP Proxies:\n──────────────"))
		for proxyName, proxy := range a.HttpProxies {
			outputs = append(outputs, pulumi.Sprintf("%s: %s", proxyName, proxy.ApiEndpoint))
		}
	}

	// Add Websocket outputs
	if len(a.Websockets) > 0 {
		if len(outputs) > 0 {
			outputs = append(outputs, "\n")
		}
		outputs = append(outputs, pulumi.Sprintf("Websockets:\n──────────────"))
		for wsName, ws := range a.Websockets {
			outputs = append(outputs, pulumi.Sprintf("%s: %s/%s", wsName, ws.ApiEndpoint, common.DefaultWsStageName))
		}
	}

	output, ok := pulumi.All(outputs...).ApplyT(func(deets []interface{}) string {
		stringyOutputs := make([]string, len(deets))
		for i, d := range deets {
			stringyOutputs[i] = d.(string)
		}

		return strings.Join(stringyOutputs, "\n")
	}).(pulumi.StringOutput)

	if !ok {
		return pulumi.StringOutput{}, fmt.Errorf("failed to generate pulumi output")
	}

	return output, nil
}

func NewNitricAwsProvider() *NitricAwsPulumiProvider {
	return &NitricAwsPulumiProvider{
		BatchRoles:            make(map[string]*iam.Role),
		Lambdas:               make(map[string]*lambda.Function),
		LambdaRoles:           make(map[string]*iam.Role),
		Apis:                  make(map[string]*apigatewayv2.Api),
		HttpProxies:           make(map[string]*apigatewayv2.Api),
		Secrets:               make(map[string]*secretsmanager.Secret),
		Buckets:               make(map[string]*s3.Bucket),
		BucketNotifications:   make(map[string]*s3.BucketNotification),
		Websockets:            make(map[string]*apigatewayv2.Api),
		Topics:                make(map[string]*topic),
		Queues:                make(map[string]*sqs.Queue),
		KeyValueStores:        make(map[string]*dynamodb.Table),
		DatabaseMigrationJobs: make(map[string]*codebuild.Project),
		SqlDatabases:          make(map[string]*RdsDatabase),
	}
}
