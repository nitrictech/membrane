// Copyright 2021 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package exec

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/eventgrid/eventgrid"
	"github.com/pulumi/pulumi-azure-native-sdk/app"
	"github.com/pulumi/pulumi-azure-native-sdk/authorization"
	"github.com/pulumi/pulumi-azure-native-sdk/containerregistry"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/nitrictech/nitric/cloud/azure/deploy/policy"
	"github.com/nitrictech/nitric/cloud/azure/deploy/utils"
	common "github.com/nitrictech/nitric/cloud/common/deploy/tags"
	deploy "github.com/nitrictech/nitric/core/pkg/api/nitric/deploy/v1"
)

type ContainerAppArgs struct {
	ResourceGroupName             pulumi.StringInput
	Location                      pulumi.StringInput
	StackID                       pulumi.StringInput
	SubscriptionID                pulumi.StringInput
	Registry                      *containerregistry.Registry
	RegistryUser                  pulumi.StringPtrInput
	RegistryPass                  pulumi.StringPtrInput
	ManagedEnv                    *app.ManagedEnvironment
	ImageUri                      pulumi.StringInput
	Env                           app.EnvironmentVarArray
	ExecutionUnit                 *deploy.ExecutionUnit
	ManagedIdentityID             pulumi.StringOutput
	MongoDatabaseName             pulumi.StringInput
	MongoDatabaseConnectionString pulumi.StringInput
}

type ContainerApp struct {
	pulumi.ResourceState

	Name    string
	hostUrl *pulumi.StringOutput
	Sp      *policy.ServicePrincipal
	App     *app.ContainerApp
}

// HostUrl - Returns the HostURL of the application
// this will also ensure that the application has been successfully deployed
func (c *ContainerApp) HostUrl() (pulumi.StringOutput, error) {
	if c.hostUrl == nil {
		// Set the hostUrl from the App FQDN
		hostUrl := c.App.LatestRevisionFqdn.ApplyT(func(fqdn string) (string, error) {
			// Get the full URL of the deployed container
			hostUrl := "https://" + fqdn

			hCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()

			// Poll the URL until the host has started.
			for {
				// Provide data in the expected shape. The content is current not important.
				empty := ""
				dummyEvgt := eventgrid.Event{
					ID:          &empty,
					Data:        &empty,
					EventType:   &empty,
					Subject:     &empty,
					DataVersion: &empty,
				}

				jsonStr, err := dummyEvgt.MarshalJSON()
				if err != nil {
					return "", err
				}

				body := bytes.NewBuffer(jsonStr)

				req, err := http.NewRequestWithContext(hCtx, "POST", hostUrl, body)
				if err != nil {
					return "", err
				}

				// TODO: Implement a membrane health check handler in the Membrane and trigger that instead.
				// Set event type header to simulate a subscription validation event.
				// These events are automatically resolved by the Membrane and won't be processed by handlers.
				req.Header.Set("aeg-event-type", "SubscriptionValidation")
				req.Header.Set("Content-Type", "application/json")
				client := &http.Client{
					Timeout: 10 * time.Second,
				}

				resp, err := client.Do(req)
				if err == nil {
					resp.Body.Close()
					break
				}
			}

			return hostUrl, nil
		}).(pulumi.StringOutput)

		c.hostUrl = &hostUrl
	}

	return *c.hostUrl, nil
}

// Built in role definitions for Azure
// See below URL for mapping
// https://docs.microsoft.com/en-us/azure/role-based-access-control/built-in-roles
var RoleDefinitions = map[string]string{
	"KVSecretsOfficer":    "b86a8fe4-44ce-4948-aee5-eccb2c155cd7",
	"BlobDataContrib":     "ba92f5b4-2d11-453d-a403-e96b0029c9fe",
	"QueueDataContrib":    "974c5e8b-45b9-4653-ba55-5f855dd0fb88",
	"EventGridDataSender": "d5a91429-5739-47e2-a06b-3470a27159e7",
	// Access for locating resources
	"TagContributor": "4a9ae827-6dc8-4573-8ac7-8239d42aa03f",
}

func NewContainerApp(ctx *pulumi.Context, name string, args *ContainerAppArgs, opts ...pulumi.ResourceOption) (*ContainerApp, error) {
	res := &ContainerApp{
		Name: name,
	}

	err := ctx.RegisterComponentResource("nitric:func:ContainerApp", name, res, opts...)
	if err != nil {
		return nil, err
	}

	res.Sp, err = policy.NewServicePrincipal(ctx, name, &policy.ServicePrincipalArgs{}, pulumi.Parent(res))
	if err != nil {
		return nil, err
	}

	scope := pulumi.Sprintf("subscriptions/%s/resourceGroups/%s", args.SubscriptionID, args.ResourceGroupName)

	// Assign roles to the new SP
	for defName, id := range RoleDefinitions {
		_ = ctx.Log.Info("Assignment "+utils.ResourceName(ctx, name+defName, utils.AssignmentRT)+" roleDef "+id, &pulumi.LogArgs{Ephemeral: true})

		_, err = authorization.NewRoleAssignment(ctx, utils.ResourceName(ctx, name+defName, utils.AssignmentRT), &authorization.RoleAssignmentArgs{
			PrincipalId:      res.Sp.ServicePrincipalId,
			PrincipalType:    pulumi.StringPtr("ServicePrincipal"),
			RoleDefinitionId: pulumi.Sprintf("/subscriptions/%s/providers/Microsoft.Authorization/roleDefinitions/%s", args.SubscriptionID, id),
			Scope:            scope,
		}, pulumi.Parent(res))
		if err != nil {
			return nil, err
		}
	}

	env := app.EnvironmentVarArray{
		app.EnvironmentVarArgs{
			Name:  pulumi.String("NITRIC_ENVIRONMENT"),
			Value: pulumi.String("cloud"),
		},
		app.EnvironmentVarArgs{
			Name:  pulumi.String("MIN_WORKERS"),
			Value: pulumi.String(fmt.Sprint(args.ExecutionUnit.Workers)),
		},
		app.EnvironmentVarArgs{
			Name:  pulumi.String("AZURE_SUBSCRIPTION_ID"),
			Value: args.SubscriptionID,
		},
		app.EnvironmentVarArgs{
			Name:  pulumi.String("AZURE_RESOURCE_GROUP"),
			Value: args.ResourceGroupName,
		},
		app.EnvironmentVarArgs{
			Name:      pulumi.String("AZURE_CLIENT_ID"),
			SecretRef: pulumi.String("client-id"),
		},
		app.EnvironmentVarArgs{
			Name:      pulumi.String("AZURE_TENANT_ID"),
			SecretRef: pulumi.String("tenant-id"),
		},
		app.EnvironmentVarArgs{
			Name:      pulumi.String("AZURE_CLIENT_SECRET"),
			SecretRef: pulumi.String("client-secret"),
		},
		app.EnvironmentVarArgs{
			Name:  pulumi.String("TOLERATE_MISSING_SERVICES"),
			Value: pulumi.String("true"),
		},
		app.EnvironmentVarArgs{
			Name:  pulumi.String("MONGODB_CONNECTION_STRING"),
			Value: args.MongoDatabaseConnectionString,
		},
		app.EnvironmentVarArgs{
			Name:  pulumi.String("MONGODB_DATABASE"),
			Value: args.MongoDatabaseName,
		},
	}

	// memory := common.IntValueOrDefault(args.Compute.Unit().Memory, 128)
	// we can't define memory without defining the cpu..
	appName := utils.ResourceName(ctx, name, utils.ContainerAppRT)

	res.App, err = app.NewContainerApp(ctx, appName, &app.ContainerAppArgs{
		ResourceGroupName:    args.ResourceGroupName,
		Location:             args.Location,
		ManagedEnvironmentId: args.ManagedEnv.ID(),
		Configuration: app.ConfigurationArgs{
			Ingress: app.IngressArgs{
				External:   pulumi.BoolPtr(true),
				TargetPort: pulumi.Int(9001),
			},
			Registries: app.RegistryCredentialsArray{
				app.RegistryCredentialsArgs{
					Server:            args.Registry.LoginServer,
					Username:          args.RegistryUser,
					PasswordSecretRef: pulumi.String("pwd"),
				},
			},
			Secrets: app.SecretArray{
				app.SecretArgs{
					Name:  pulumi.String("pwd"),
					Value: args.RegistryPass,
				},
				app.SecretArgs{
					Name:  pulumi.String("client-id"),
					Value: res.Sp.ClientID,
				},
				app.SecretArgs{
					Name:  pulumi.String("tenant-id"),
					Value: res.Sp.TenantID,
				},
				app.SecretArgs{
					Name:  pulumi.String("client-secret"),
					Value: res.Sp.ClientSecret,
				},
			},
		},
		Tags: common.Tags(ctx, args.StackID, name),
		Template: app.TemplateArgs{
			Containers: app.ContainerArray{
				app.ContainerArgs{
					Name:  pulumi.String("myapp"),
					Image: args.ImageUri,
					Env:   append(env, args.Env...),
				},
			},
		},
	}, pulumi.Parent(res))
	if err != nil {
		return nil, err
	}

	authName := fmt.Sprintf("%s-auth", appName)

	_, err = app.NewContainerAppsAuthConfig(ctx, authName, &app.ContainerAppsAuthConfigArgs{
		AuthConfigName:   pulumi.String("current"),
		ContainerAppName: res.App.Name,
		GlobalValidation: &app.GlobalValidationArgs{
			UnauthenticatedClientAction: app.UnauthenticatedClientActionV2Return401,
		},
		IdentityProviders: &app.IdentityProvidersArgs{
			AzureActiveDirectory: &app.AzureActiveDirectoryArgs{
				Enabled: pulumi.Bool(true),
				Registration: &app.AzureActiveDirectoryRegistrationArgs{
					ClientId:                res.Sp.ClientID,
					ClientSecretSettingName: pulumi.String("client-secret"),
					OpenIdIssuer:            pulumi.Sprintf("https://sts.windows.net/%s/v2.0", res.Sp.TenantID),
				},
				Validation: &app.AzureActiveDirectoryValidationArgs{
					AllowedAudiences: pulumi.StringArray{args.ManagedIdentityID},
				},
			},
		},
		Platform: &app.AuthPlatformArgs{
			Enabled: pulumi.Bool(true),
		},
		ResourceGroupName: args.ResourceGroupName,
	}, pulumi.Parent(res.App))
	if err != nil {
		return nil, err
	}

	return res, ctx.RegisterResourceOutputs(res, pulumi.Map{
		"name":         pulumi.StringPtr(res.Name),
		"containerApp": res.App,
	})
}
