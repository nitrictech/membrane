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
	"context"
	"fmt"
	"runtime/debug"
	"strings"

	apiv1 "cloud.google.com/go/firestore/apiv1/admin"
	"cloud.google.com/go/firestore/apiv1/admin/adminpb"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/nitrictech/nitric/cloud/common/deploy/image"
	"github.com/nitrictech/nitric/cloud/common/deploy/telemetry"
	"github.com/nitrictech/nitric/cloud/gcp/deploy/collection"
	"github.com/nitrictech/nitric/cloud/gcp/deploy/config"
	"github.com/nitrictech/nitric/cloud/gcp/deploy/exec"
	"github.com/nitrictech/nitric/cloud/gcp/deploy/gateway"
	"github.com/nitrictech/nitric/cloud/gcp/deploy/iam"
	"github.com/nitrictech/nitric/cloud/gcp/deploy/policy"
	"github.com/nitrictech/nitric/cloud/gcp/deploy/schedule"
	"github.com/nitrictech/nitric/cloud/gcp/deploy/secret"
	"github.com/nitrictech/nitric/cloud/gcp/deploy/storage"
	events "github.com/nitrictech/nitric/cloud/gcp/deploy/topic"
	deploy "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/cloudscheduler"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/cloudtasks"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/organizations"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/projects"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/serviceaccount"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewUpProgram(ctx context.Context, details *StackDetails, config *config.GcpConfig, spec *deploy.Spec) (auto.Stack, error) {
	return auto.UpsertStackInlineSource(context.TODO(), details.FullStackName, details.Project, func(ctx *pulumi.Context) (err error) {
		defer func() {
			if r := recover(); r != nil {
				stack := string(debug.Stack())
				err = fmt.Errorf("recovered panic: %+v\n Stack: %s", r, stack)
			}
		}()

		// Get Websockets
		websockets := lo.Filter[*deploy.Resource](spec.Resources, func(item *deploy.Resource, index int) bool {
			return item.GetWebsocket() != nil
		})
		if len(websockets) > 0 {
			return fmt.Errorf("websockets currently in preview not supported in the GCP provider.")
		}

		project, err := organizations.LookupProject(ctx, &organizations.LookupProjectArgs{
			ProjectId: &details.ProjectId,
		}, nil)
		if err != nil {
			return err
		}

		nitricProj, err := NewProject(ctx, "project", &ProjectArgs{
			ProjectId:     details.ProjectId,
			ProjectNumber: project.Number,
		})
		if err != nil {
			return err
		}

		defaultResourceOptions := pulumi.DependsOn([]pulumi.Resource{nitricProj})

		// Calculate unique stackID
		stackRandId, err := random.NewRandomString(ctx, fmt.Sprintf("%s-stack-name", ctx.Stack()), &random.RandomStringArgs{
			Special: pulumi.Bool(false),
			Length:  pulumi.Int(8),
			Upper:   pulumi.Bool(false),
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

		stackID := <-stackIdChan

		collections := lo.Filter[*deploy.Resource](spec.Resources, func(res *deploy.Resource, _ int) bool {
			return res.Id.Type == v1.ResourceType_Collection
		})

		// If collections are required we need to instansiate a default database if one doesn't exist
		if len(collections) > 0 {
			fsAdminClient, err := apiv1.NewFirestoreAdminClient(context.TODO())
			if err != nil {
				return err
			}

			defaultDb, _ := fsAdminClient.GetDatabase(context.TODO(), &adminpb.GetDatabaseRequest{
				Name: fmt.Sprintf("projects/%s/databases/(default)", *project.ProjectId),
			})

			defaultDbExists := defaultDb != nil

			// create a firestore database for the stack or adopt the default database
			// TODO: Determine if we can create multiple databases
			_, err = collection.NewFirestoreCollectionDatabase(ctx, fmt.Sprintf("nitric-%s", ctx.Stack()), &collection.FirestoreCollectionDatabaseArgs{
				Location:      details.Region,
				DefaultExists: defaultDbExists,
			})
			if err != nil {
				fmt.Println("got a database err")
				return err
			}
		}

		// create a shared queue for enabling delayed messaging
		// cloud run functions will create OIDC tokens for their own service accounts
		// to apply to push actions to pubsub, so their scope should still be limited to that
		topicDelayQueue, err := cloudtasks.NewQueue(ctx, "delay-queue", &cloudtasks.QueueArgs{
			Location: pulumi.String(details.Region),
		}, defaultResourceOptions)
		if err != nil {
			return err
		}

		// Deploy all execution units
		authToken, err := getGCPToken(ctx)
		if err != nil {
			return err
		}

		execs := map[string]*exec.CloudRunner{}

		baseCustomRoleId, err := random.NewRandomString(ctx, fmt.Sprintf("%s-base-role", details.FullStackName), &random.RandomStringArgs{
			Special: pulumi.Bool(false),
			Length:  pulumi.Int(8),
			Keepers: pulumi.ToMap(map[string]interface{}{
				"stack-name": details.FullStackName,
			}),
		})
		if err != nil {
			return errors.WithMessage(err, "base customRole id")
		}

		principalMap := make(policy.PrincipalMap)
		principalMap[v1.ResourceType_Function] = make(map[string]*serviceaccount.Account)

		// setup a basic IAM role for general access and resource discovery
		var baseComputeRole *projects.IAMCustomRole

		for _, res := range spec.Resources {
			switch eu := res.Config.(type) {
			case *deploy.Resource_ExecutionUnit:
				if eu.ExecutionUnit.GetImage() == nil {
					return fmt.Errorf("gcp provider can only deploy execution with an image source")
				}

				if eu.ExecutionUnit.GetImage().GetUri() == "" {
					return fmt.Errorf("gcp provider can only deploy execution with an image source")
				}

				if eu.ExecutionUnit.Type == "" {
					eu.ExecutionUnit.Type = "default"
				}

				// get config for execution unit
				unitConfig, hasConfig := config.Config[eu.ExecutionUnit.Type]
				if !hasConfig {
					return status.Errorf(codes.InvalidArgument, "unable to find config %s in stack config %+v", eu.ExecutionUnit.Type, config.Config)
				}

				// Set here because we need access to the config
				if baseComputeRole == nil {
					baseComputeRole, err = projects.NewIAMCustomRole(ctx, "base-role", &projects.IAMCustomRoleArgs{
						Title:       pulumi.String(details.FullStackName + "-functions-base-role"),
						Permissions: pulumi.ToStringArray(exec.GetPerms(unitConfig.Telemetry)),
						RoleId:      baseCustomRoleId.ID(),
					}, defaultResourceOptions)
					if err != nil {
						return errors.WithMessage(err, "base customRole")
					}
				}

				// Get the image name:tag from the uri
				imageUriSplit := strings.Split(eu.ExecutionUnit.GetImage().GetUri(), "/")
				imageName := imageUriSplit[len(imageUriSplit)-1]

				image, err := image.NewImage(ctx, res.Id.Name, &image.ImageArgs{
					SourceImage:   eu.ExecutionUnit.GetImage().GetUri(),
					RepositoryUrl: pulumi.Sprintf("gcr.io/%s/%s", details.ProjectId, imageName),
					Username:      pulumi.String("oauth2accesstoken"),
					Password:      pulumi.String(authToken.AccessToken),
					Server:        pulumi.String("https://gcr.io"),
					Runtime:       runtime,
					Telemetry: &telemetry.TelemetryConfigArgs{
						TraceSampling:       unitConfig.Telemetry,
						TraceName:           "googlecloud",
						MetricName:          "googlecloud",
						TraceExporterConfig: `{"retry_on_failure": {"enabled": false}}`,
						Extensions:          []string{},
					},
				}, defaultResourceOptions)
				if err != nil {
					return err
				}

				sa, err := iam.NewServiceAccount(ctx, res.Id.Name+"-cloudrun-exec-acct", &iam.GcpIamServiceAccountArgs{
					AccountId: res.Id.Name + "-exec",
				}, defaultResourceOptions)
				if err != nil {
					return err
				}

				if eu.ExecutionUnit.Type == "" {
					eu.ExecutionUnit.Type = "default"
				}

				if unitConfig.CloudRun != nil {
					execs[res.Id.Name], err = exec.NewCloudRunner(ctx, res.Id.Name, &exec.CloudRunnerArgs{
						StackID:         stackID,
						Location:        pulumi.String(details.Region),
						ProjectID:       details.ProjectId,
						Compute:         res.GetExecutionUnit(),
						Image:           image,
						EnvMap:          eu.ExecutionUnit.Env,
						DelayQueue:      topicDelayQueue,
						ServiceAccount:  sa.ServiceAccount,
						BaseComputeRole: baseComputeRole,
						Config:          *unitConfig.CloudRun,
					}, defaultResourceOptions)
					if err != nil {
						return err
					}
				} else {
					return status.Errorf(codes.InvalidArgument, "unsupported target for function config %+v", unitConfig)
				}

				principalMap[v1.ResourceType_Function][res.Id.Name] = sa.ServiceAccount
			}
		}

		// Deploy all buckets
		buckets := map[string]*storage.CloudStorageBucket{}
		for _, res := range spec.Resources {
			switch b := res.Config.(type) {
			case *deploy.Resource_Bucket:
				buckets[res.Id.Name], err = storage.NewCloudStorageBucket(ctx, res.Id.Name, &storage.CloudStorageBucketArgs{
					StackID:  stackID,
					Bucket:   b.Bucket,
					Location: details.Region,
				}, defaultResourceOptions)
				if err != nil {
					return err
				}

				for _, notification := range b.Bucket.Notifications {
					// Get the deployed execution unit
					unit, ok := execs[notification.GetExecutionUnit()]
					if !ok {
						return fmt.Errorf("invalid execution unit %s given for topic subscription", notification.GetExecutionUnit())
					}

					notificationName := fmt.Sprintf("%s-%s-%s-notify", res.Id.Name, strings.ToLower(notification.Config.BlobEventType.String()), notification.GetExecutionUnit())
					_, err = storage.NewCloudStorageNotification(ctx, notificationName, &storage.CloudStorageNotificationArgs{
						StackID:  stackID,
						Location: details.Region,
						Bucket:   buckets[res.Id.Name],
						Config:   notification.Config,
						Function: unit,
					}, defaultResourceOptions)
					if err != nil {
						return err
					}
				}
			}
		}

		apis := map[string]*gateway.ApiGateway{}
		for _, res := range spec.Resources {
			switch t := res.Config.(type) {
			case *deploy.Resource_Api:
				if t.Api.GetOpenapi() == "" {
					return fmt.Errorf("gcp provider can only deploy OpenAPI specs")
				}

				doc := &openapi3.T{}
				err := doc.UnmarshalJSON([]byte(t.Api.GetOpenapi()))
				if err != nil {
					return fmt.Errorf("invalid document suppled for api: %s", res.Id.Name)
				}

				apis[res.Id.Name], err = gateway.NewApiGateway(ctx, res.Id.Name, &gateway.ApiGatewayArgs{
					StackID:     stackID,
					ProjectId:   details.ProjectId,
					Functions:   execs,
					OpenAPISpec: doc,
				}, defaultResourceOptions)
				if err != nil {
					return err
				}
			}
		}

		httpProxies := map[string]*gateway.HttpProxy{}
		for _, res := range spec.Resources {
			switch t := res.Config.(type) {
			case *deploy.Resource_Http:
				fun := execs[t.Http.Target.GetExecutionUnit()]

				httpProxies[res.Id.Name], err = gateway.NewHttpProxy(ctx, res.Id.Name, &gateway.HttpProxyArgs{
					StackID:   stackID,
					ProjectId: details.ProjectId,
					Function:  fun,
				})
				if err != nil {
					return err
				}
			}
		}

		secrets := map[string]*secret.SecretManagerSecret{}
		for _, res := range spec.Resources {
			switch t := res.Config.(type) {
			case *deploy.Resource_Secret:
				secrets[res.Id.Name], err = secret.NewSecretManagerSecret(ctx, res.Id.Name, &secret.SecretManagerSecretArgs{
					StackID:   stackID,
					StackName: details.Stack,
					Secret:    t.Secret,
				}, defaultResourceOptions)
				if err != nil {
					return err
				}
			}
		}

		topics := map[string]*events.PubSubTopic{}
		for _, res := range spec.Resources {
			switch t := res.Config.(type) {
			case *deploy.Resource_Topic:
				topics[res.Id.Name], err = events.NewPubSubTopic(ctx, res.Id.Name, &events.PubSubTopicArgs{
					Topic:     t.Topic,
					Location:  details.Region,
					ProjectId: details.ProjectId,
					StackID:   stackID,
				}, defaultResourceOptions)
				if err != nil {
					return err
				}

				for _, sub := range t.Topic.Subscriptions {
					subName := events.GetSubName(sub.GetExecutionUnit(), res.Id.Name)

					// Get the deployed execution unit
					unit, ok := execs[sub.GetExecutionUnit()]
					if !ok {
						return fmt.Errorf("invalid execution unit %s given for topic subscription", sub.GetExecutionUnit())
					}

					_, err = events.NewPubSubPushSubscription(ctx, subName, &events.PubSubSubscriptionArgs{
						Topic:    topics[res.Id.Name],
						Function: unit,
					}, defaultResourceOptions)
					if err != nil {
						return err
					}
				}
			}
		}

		schedules := map[string]*cloudscheduler.Job{}
		for _, res := range spec.Resources {
			switch t := res.Config.(type) {
			case *deploy.Resource_Schedule:
				// get the target of the schedule

				execUnitName := t.Schedule.Target.GetExecutionUnit()
				execUnit, ok := execs[execUnitName]
				if !ok {
					return fmt.Errorf("no execution unit with name %s", execUnitName)
				}

				// Create schedule targeting a given lambda
				job, err := schedule.NewCloudSchedulerJob(ctx, res.Id.Name, &schedule.CloudSchedulerArgs{
					Exec:     execUnit,
					Schedule: t.Schedule,
					Tz:       config.ScheduleTimezone,
				})
				if err != nil {
					return err
				}

				schedules[res.Id.Name] = job.Job
			}
		}

		// Create policies
		for _, res := range spec.Resources {
			switch t := res.Config.(type) {
			case *deploy.Resource_Policy:
				_, err := policy.NewIAMPolicy(ctx, res.Id.Name, &policy.PolicyArgs{
					Policy: t.Policy,
					Resources: &policy.StackResources{
						Buckets: buckets,
						Topics:  topics,
						Secrets: secrets,
					},
					Principals: principalMap,
					ProjectID:  pulumi.String(details.ProjectId),
				})
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
}
