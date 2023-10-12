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

package exec

import (
	"fmt"

	"github.com/nitrictech/nitric/cloud/common/deploy/image"
	"github.com/nitrictech/nitric/cloud/gcp/deploy/config"
	"github.com/nitrictech/nitric/cloud/gcp/deploy/iam"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/deploy/v1"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/cloudrun"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/cloudtasks"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/projects"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/serviceaccount"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type CloudRunner struct {
	pulumi.ResourceState

	Name       string
	Service    *cloudrun.Service
	Url        pulumi.StringInput
	Invoker    *serviceaccount.Account
	EventToken pulumi.StringOutput
}

type CloudRunnerArgs struct {
	Location        pulumi.StringInput
	ProjectId       string
	Compute         *v1.ExecutionUnit
	Image           *image.Image
	EnvMap          map[string]string
	DelayQueue      *cloudtasks.Queue
	BaseComputeRole *projects.IAMCustomRole
	ServiceAccount  *serviceaccount.Account
	Config          config.GcpCloudRunConfig

	StackID string
}

func GetPerms(telemetry int) []string {
	perms := []string{
		"storage.buckets.list",
		"storage.buckets.get",
		"cloudtasks.queues.get",
		"cloudtasks.tasks.create",
		"cloudtrace.traces.patch",
		"monitoring.timeSeries.create",
		// permission for blob signing
		// this is safe as only permissions this account has are delegated
		"iam.serviceAccounts.signBlob",
		// Basic list permissions
		"pubsub.topics.list",
		"pubsub.topics.get",
		"pubsub.snapshots.list",
		"pubsub.subscriptions.get",
		"resourcemanager.projects.get",
		"secretmanager.secrets.list",
		"apigateway.gateways.list",
	}

	// Add perms for telemetry
	if telemetry > 0 {
		perms = append(perms, []string{
			"monitoring.metricDescriptors.create",
			"monitoring.metricDescriptors.get",
			"monitoring.metricDescriptors.list",
			"monitoring.monitoredResourceDescriptors.get",
			"monitoring.monitoredResourceDescriptors.list",
			"monitoring.timeSeries.create",
		}...)
	}

	return perms
}

func NewCloudRunner(ctx *pulumi.Context, name string, args *CloudRunnerArgs, opts ...pulumi.ResourceOption) (*CloudRunner, error) {
	res := &CloudRunner{
		Name: name,
	}

	err := ctx.RegisterComponentResource("nitric:exec:GCPCloudRunner", name, res, opts...)
	if err != nil {
		return nil, err
	}

	// generate a token for internal application events to authenticate themeselves
	// https://cloud.google.com/appengine/docs/flexible/writing-and-responding-to-pub-sub-messages?tab=go#top
	token, err := random.NewRandomPassword(ctx, res.Name+"-event-token", &random.RandomPasswordArgs{
		Special: pulumi.Bool(false),
		Length:  pulumi.Int(32),
		Keepers: pulumi.ToMap(map[string]interface{}{
			"name": name,
		}),
	})
	if err != nil {
		return nil, errors.WithMessage(err, "service event token")
	}

	res.EventToken = token.Result

	// apply basic project level permissions for nitric resource discovery
	_, err = projects.NewIAMMember(ctx, res.Name+"-project-member", &projects.IAMMemberArgs{
		Project: pulumi.String(args.ProjectId),
		Member:  pulumi.Sprintf("serviceAccount:%s", args.ServiceAccount.Email),
		Role:    args.BaseComputeRole.Name,
	})
	if err != nil {
		return nil, errors.WithMessage(err, "function project membership "+res.Name)
	}

	// give the service account permission to use itself
	_, err = serviceaccount.NewIAMMember(ctx, res.Name+"-acct-member", &serviceaccount.IAMMemberArgs{
		ServiceAccountId: args.ServiceAccount.Name,
		Member:           pulumi.Sprintf("serviceAccount:%s", args.ServiceAccount.Email),
		Role:             pulumi.String("roles/iam.serviceAccountUser"),
	})
	if err != nil {
		return nil, errors.WithMessage(err, "service account self membership "+res.Name)
	}

	env := getCloudRunnerEnvs(args)

	env = append(env, cloudrun.ServiceTemplateSpecContainerEnvArgs{
		Name:  pulumi.String("EVENT_TOKEN"),
		Value: res.EventToken,
	})

	if args.DelayQueue != nil {
		env = append(env, cloudrun.ServiceTemplateSpecContainerEnvArgs{
			Name:  pulumi.String("DELAY_QUEUE_NAME"),
			Value: pulumi.Sprintf("projects/%s/locations/%s/queues/%s", args.DelayQueue.Project, args.DelayQueue.Location, args.DelayQueue.Name),
		})
	}

	for k, v := range args.EnvMap {
		env = append(env, cloudrun.ServiceTemplateSpecContainerEnvArgs{
			Name:  pulumi.String(k),
			Value: pulumi.String(v),
		})
	}

	res.Service, err = cloudrun.NewService(ctx, name, &cloudrun.ServiceArgs{
		AutogenerateRevisionName: pulumi.BoolPtr(true),
		Location:                 args.Location,
		Project:                  pulumi.String(args.ProjectId),
		Template: cloudrun.ServiceTemplateArgs{
			Metadata: cloudrun.ServiceTemplateMetadataArgs{
				Annotations: pulumi.StringMap{
					"autoscaling.knative.dev/minScale": pulumi.Sprintf("%d", args.Config.MinInstances),
					"autoscaling.knative.dev/maxScale": pulumi.Sprintf("%d", args.Config.MaxInstances),
				},
			},
			Spec: cloudrun.ServiceTemplateSpecArgs{
				ServiceAccountName:   args.ServiceAccount.Email,
				ContainerConcurrency: pulumi.Int(args.Config.Concurrency),
				TimeoutSeconds:       pulumi.Int(args.Config.Timeout),
				Containers: cloudrun.ServiceTemplateSpecContainerArray{
					cloudrun.ServiceTemplateSpecContainerArgs{
						Envs:  env,
						Image: args.Image.URI(),
						Ports: cloudrun.ServiceTemplateSpecContainerPortArray{
							cloudrun.ServiceTemplateSpecContainerPortArgs{
								ContainerPort: pulumi.Int(9001),
							},
						},
						Resources: cloudrun.ServiceTemplateSpecContainerResourcesArgs{
							Limits: pulumi.StringMap{
								"memory": pulumi.Sprintf("%dMi", args.Config.Memory),
							},
						},
					},
				},
			},
		},
	}, append(opts, pulumi.Parent(res))...)
	if err != nil {
		return nil, errors.WithMessage(err, "cloud run "+name)
	}

	svcAcct, err := iam.NewServiceAccount(ctx, name+"-cloudrun-invoker", &iam.GcpIamServiceAccountArgs{
		AccountId: name,
	})
	if err != nil {
		return nil, errors.WithMessage(err, "invokerAccount "+name)
	}

	res.Invoker = svcAcct.ServiceAccount

	_, err = cloudrun.NewIamMember(ctx, name+"-invoker", &cloudrun.IamMemberArgs{
		Member:   pulumi.Sprintf("serviceAccount:%s", res.Invoker.Email),
		Role:     pulumi.String("roles/run.invoker"),
		Service:  res.Service.Name,
		Location: res.Service.Location,
	}, append(opts, pulumi.Parent(res))...)
	if err != nil {
		return nil, errors.WithMessage(err, "iam member "+name)
	}

	res.Url = res.Service.Statuses.ApplyT(func(ss []cloudrun.ServiceStatus) (string, error) {
		if len(ss) == 0 {
			return "", errors.New("serviceStatus is empty")
		}

		return *ss[0].Url, nil
	}).(pulumi.StringInput)

	return res, ctx.RegisterResourceOutputs(res, pulumi.Map{
		"name":    pulumi.String(res.Name),
		"service": res.Service,
		"url":     res.Url,
	})
}

func getCloudRunnerEnvs(args *CloudRunnerArgs) cloudrun.ServiceTemplateSpecContainerEnvArray {
	return cloudrun.ServiceTemplateSpecContainerEnvArray{
		cloudrun.ServiceTemplateSpecContainerEnvArgs{
			Name:  pulumi.String("NITRIC_ENVIRONMENT"),
			Value: pulumi.String("cloud"),
		},
		cloudrun.ServiceTemplateSpecContainerEnvArgs{
			Name:  pulumi.String("MIN_WORKERS"),
			Value: pulumi.String(fmt.Sprintf("%d", args.Compute.Workers)),
		},
		cloudrun.ServiceTemplateSpecContainerEnvArgs{
			Name:  pulumi.String("NITRIC_STACK"),
			Value: pulumi.String(args.StackID),
		},
		cloudrun.ServiceTemplateSpecContainerEnvArgs{
			Name:  pulumi.String("SERVICE_ACCOUNT_EMAIL"),
			Value: args.ServiceAccount.Email,
		},
		cloudrun.ServiceTemplateSpecContainerEnvArgs{
			Name:  pulumi.String("GCP_REGION"),
			Value: args.Location,
		},
		cloudrun.ServiceTemplateSpecContainerEnvArgs{
			Name:  pulumi.String("NITRIC_HTTP_PROXY_PORT"),
			Value: pulumi.String(fmt.Sprint(3000)),
		},
	}
}
