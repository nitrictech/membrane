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

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/nitrictech/nitric/cloud/azure/deploy/api"
	"github.com/nitrictech/nitric/cloud/azure/deploy/bucket"
	"github.com/nitrictech/nitric/cloud/azure/deploy/collection"
	"github.com/nitrictech/nitric/cloud/azure/deploy/config"
	"github.com/nitrictech/nitric/cloud/azure/deploy/exec"
	"github.com/nitrictech/nitric/cloud/azure/deploy/schedule"
	"github.com/nitrictech/nitric/cloud/azure/deploy/topic"
	"github.com/nitrictech/nitric/cloud/azure/deploy/utils"
	"github.com/nitrictech/nitric/cloud/common/deploy/image"
	nitricresources "github.com/nitrictech/nitric/cloud/common/deploy/resources"
	common "github.com/nitrictech/nitric/cloud/common/deploy/tags"
	deploy "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-azure-native-sdk/authorization"
	"github.com/pulumi/pulumi-azure-native-sdk/keyvault"
	"github.com/pulumi/pulumi-azure-native-sdk/resources"
	"github.com/pulumi/pulumi-azure-native-sdk/storage"
	azureStorage "github.com/pulumi/pulumi-azure-native-sdk/storage"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewUpProgram(ctx context.Context, details *StackDetails, config *config.AzureConfig, spec *deploy.Spec) (auto.Stack, error) {
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
			return fmt.Errorf("websockets currently in preview not supported in the Azure provider.")
		}

		// Get Execution units
		executionUnits := lo.Filter[*deploy.Resource](spec.Resources, func(item *deploy.Resource, index int) bool {
			return item.GetExecutionUnit() != nil
		})

		// Get Collections
		collections := lo.Filter[*deploy.Resource](spec.Resources, func(item *deploy.Resource, index int) bool {
			return item.GetCollection() != nil
		})

		// Get Buckets
		buckets := lo.Filter[*deploy.Resource](spec.Resources, func(item *deploy.Resource, index int) bool {
			return item.GetBucket() != nil
		})

		// Get Topics
		topics := lo.Filter[*deploy.Resource](spec.Resources, func(item *deploy.Resource, index int) bool {
			return item.GetTopic() != nil
		})

		// Get Topics
		schedules := lo.Filter[*deploy.Resource](spec.Resources, func(item *deploy.Resource, index int) bool {
			return item.GetSchedule() != nil
		})
		apis := lo.Filter[*deploy.Resource](spec.Resources, func(item *deploy.Resource, index int) bool {
			return item.GetApi() != nil
		})

		// Calculate unique stackID
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

		stackID := <-stackIdChan

		clientConfig, err := authorization.GetClientConfig(ctx)
		if err != nil {
			return err
		}

		rg, err := resources.NewResourceGroup(ctx, utils.ResourceName(ctx, "", utils.ResourceGroupRT), &resources.ResourceGroupArgs{
			Location: pulumi.String(details.Region),
			Tags:     pulumi.ToStringMap(common.Tags(stackID, ctx.Stack(), nitricresources.Stack)),
		})
		if err != nil {
			return errors.WithMessage(err, "resource group create")
		}

		envMap := map[string]string{}
		contEnvArgs := &exec.ContainerEnvArgs{
			ResourceGroupName: rg.Name,
			Location:          rg.Location,
			EnvMap:            envMap,
			StackID:           stackID,
		}

		// Create a stack level keyvault if secrets are enabled
		// At the moment secrets have no config level setting
		kvName := utils.ResourceName(ctx, "", utils.KeyVaultRT)

		kv, err := keyvault.NewVault(ctx, kvName, &keyvault.VaultArgs{
			Location:          rg.Location,
			ResourceGroupName: rg.Name,
			Properties: &keyvault.VaultPropertiesArgs{
				EnableSoftDelete:        pulumi.Bool(false),
				EnableRbacAuthorization: pulumi.Bool(true),
				Sku: &keyvault.SkuArgs{
					Family: pulumi.String("A"),
					Name:   keyvault.SkuNameStandard,
				},
				TenantId: pulumi.String(clientConfig.TenantId),
			},
			Tags: pulumi.ToStringMap(common.Tags(stackID, ctx.Stack(), nitricresources.Stack)),
		})
		if err != nil {
			return err
		}

		contEnvArgs.KVaultName = kv.Name

		// Create a storage account if buckets or queues are required
		var storageAccount *azureStorage.StorageAccount
		if len(buckets) > 0 {
			accName := utils.ResourceName(ctx, details.FullStackName, utils.StorageAccountRT)
			storageAccount, err = azureStorage.NewStorageAccount(ctx, accName, &storage.StorageAccountArgs{
				AccessTier:        azureStorage.AccessTierHot,
				ResourceGroupName: rg.Name,
				Kind:              pulumi.String("StorageV2"),
				Sku: azureStorage.SkuArgs{
					Name: pulumi.String(storage.SkuName_Standard_LRS),
				},
				Tags: pulumi.ToStringMap(common.Tags(stackID, kvName, nitricresources.Stack)),
			})
			if err != nil {
				return err
			}

			contEnvArgs.StorageAccountBlobEndpoint = storageAccount.PrimaryEndpoints.Blob()
			contEnvArgs.StorageAccountQueueEndpoint = storageAccount.PrimaryEndpoints.Queue()
		}

		var mongoCollections *collection.MongoCollections
		if len(collections) > 0 {
			mongoCollections, err = collection.NewMongoCollections(ctx, "mongodb", &collection.MongoCollectionsArgs{
				ResourceGroup: rg,
				Collections:   collections,
			})
			if err != nil {
				return err
			}
		}

		deployedTopics := map[string]*topic.AzureEventGridTopic{}

		var contEnv *exec.ContainerEnv

		apps := map[string]*exec.ContainerApp{}

		if len(executionUnits) > 0 {
			contEnv, err = exec.NewContainerEnv(ctx, "containerEnv", contEnvArgs)
			if err != nil {
				return errors.WithMessage(err, "containerApps")
			}

			for _, eu := range executionUnits {
				repositoryUrl := pulumi.Sprintf("%s/%s-%s-%s", contEnv.Registry.LoginServer, details.Project, eu.Id.Name, "azure")

				image, err := image.NewImage(ctx, eu.Id.Name, &image.ImageArgs{
					SourceImage:   eu.GetExecutionUnit().GetImage().GetUri(),
					RepositoryUrl: repositoryUrl,
					Username:      contEnv.RegistryUser.Elem(),
					Password:      contEnv.RegistryPass.Elem(),
					Server:        contEnv.Registry.LoginServer,
					Runtime:       runtime,
				}, pulumi.Parent(contEnv))
				if err != nil {
					return err
				}

				mongodbName := pulumi.String("").ToStringOutput()
				mongoConnectionString := pulumi.String("").ToStringOutput()
				if mongoCollections != nil {
					mongodbName = mongoCollections.MongoDB.Name
					mongoConnectionString = mongoCollections.ConnectionString
				}

				if eu.GetExecutionUnit().Type == "" {
					eu.GetExecutionUnit().Type = "default"
				}

				euConfig, hasConfig := config.Config[eu.GetExecutionUnit().Type]
				if !hasConfig {
					return status.Errorf(codes.InvalidArgument, "Could not find type %s in config %+v", eu.GetExecutionUnit().Type, config)
				}

				if euConfig.ContainerApps != nil {
					schedules := lo.Filter(schedules, func(item *deploy.Resource, idx int) bool {
						return item.GetSchedule().Target.GetExecutionUnit() == eu.Id.Name
					})

					apps[eu.Id.Name], err = exec.NewContainerApp(ctx, eu.Id.Name, &exec.ContainerAppArgs{
						ResourceGroupName:             rg.Name,
						Location:                      pulumi.String(details.Region),
						SubscriptionID:                pulumi.String(clientConfig.SubscriptionId),
						Registry:                      contEnv.Registry,
						RegistryUser:                  contEnv.RegistryUser,
						RegistryPass:                  contEnv.RegistryPass,
						ManagedEnv:                    contEnv.ManagedEnv,
						ImageUri:                      image.URI(),
						Env:                           contEnv.Env,
						ExecutionUnit:                 eu.GetExecutionUnit(),
						ManagedIdentityID:             contEnv.ManagedUser.ClientId,
						MongoDatabaseName:             mongodbName,
						MongoDatabaseConnectionString: mongoConnectionString,
						Config:                        *euConfig.ContainerApps,
						Schedules:                     schedules,
						StackID:                       stackID,
					}, pulumi.Parent(contEnv))
					if err != nil {
						return status.Errorf(codes.Internal, "error occurred whilst creating container app %s", err.Error())
					}
				} else {
					return status.Errorf(codes.InvalidArgument, "unsupported target for function config %s", eu.Id.Name)
				}
			}
		}

		for _, s := range schedules {
			cAppTarget, ok := apps[s.GetSchedule().Target.GetExecutionUnit()]
			if !ok {
				return fmt.Errorf("could not find target %s for schedule %s", s.GetSchedule().Target, s.Id.Name)
			}

			_, err := schedule.NewDaprCronBindingSchedule(ctx, s.Id.Name, &schedule.ScheduleArgs{
				ResourceGroupName: rg.Name,
				Target:            cAppTarget,
				Environment:       contEnv.ManagedEnv,
				Schedule:          s.GetSchedule(),
			})
			if err != nil {
				return err
			}
		}

		// For each bucket create a new bucket
		for _, b := range buckets {
			azBucket, err := bucket.NewAzureStorageBucket(ctx, b.Id.Name, &bucket.AzureStorageBucketArgs{
				Account:       storageAccount,
				ResourceGroup: rg,
			})
			if err != nil {
				return err
			}

			for _, notification := range b.GetBucket().Notifications {
				unit, ok := apps[notification.GetExecutionUnit()]
				if !ok {
					return fmt.Errorf("invalid execution unit %s given for bucket subscription", notification.GetExecutionUnit())
				}

				notificationName := fmt.Sprintf("%s-%s-%s-notify", b.Id.Name, strings.ToLower(notification.Config.BlobEventType.String()), notification.GetExecutionUnit())
				_, err := bucket.NewAzureBucketNotification(ctx, notificationName, &bucket.AzureBucketNotificationArgs{
					Bucket:         azBucket,
					StorageAccount: storageAccount,
					Config:         notification.Config,
					Target:         unit,
				})
				if err != nil {
					return err
				}
			}
		}

		for _, t := range topics {
			deployedTopics[t.Id.Name], err = topic.NewAzureEventGridTopic(ctx, t.Id.Name, &topic.AzureEventGridTopicArgs{
				ResourceGroup: rg,
				StackID:       stackID,
			})
			if err != nil {
				return err
			}

			for _, s := range t.GetTopic().Subscriptions {
				err := deployedTopics[t.Id.Name].AddSubscription(ctx, utils.ResourceName(ctx, fmt.Sprintf("%s-%s", t.Id.Name, s.GetExecutionUnit()), utils.EventSubscriptionRT), &topic.AzureEventGridTopicSubscriptionArgs{
					Target: apps[s.GetExecutionUnit()],
				})
				if err != nil {
					return err
				}
			}
		}

		for _, a := range apis {
			if a.GetApi().GetOpenapi() == "" {
				return fmt.Errorf("azure provider can only deploy OpenAPI specs")
			}

			doc := &openapi3.T{}
			err := doc.UnmarshalJSON([]byte(a.GetApi().GetOpenapi()))
			if err != nil {
				return fmt.Errorf("invalid document suppled for api: %s", a.Id.Name)
			}

			_, err = api.NewAzureApiManagement(ctx, a.Id.Name, &api.AzureApiManagementArgs{
				ResourceGroupName: rg.Name,
				OrgName:           pulumi.String(details.Org),
				AdminEmail:        pulumi.String(details.AdminEmail),
				OpenAPISpec:       doc,
				Apps:              apps,
				ManagedIdentity:   contEnv.ManagedUser,
				StackID:           stackID,
			})
			if err != nil {
				return err
			}
		}

		// Add all HTTP proxies
		// httpProxies := map[string]*api.AzureHttpProxy{}
		for _, res := range spec.Resources {
			switch t := res.Config.(type) {
			case *deploy.Resource_Http:
				app := apps[t.Http.Target.GetExecutionUnit()]

				_, err = api.NewAzureHttpProxy(ctx, res.Id.Name, &api.AzureHttpProxyArgs{
					ResourceGroupName: rg.Name,
					OrgName:           pulumi.String(details.Org),
					AdminEmail:        pulumi.String(details.AdminEmail),
					App:               app,
					ManagedIdentity:   contEnv.ManagedUser,
					StackID:           stackID,
				})
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
}
