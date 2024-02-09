// Copyright 2021 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
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

	resourcespb "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
	"github.com/pulumi/pulumi-azure-native-sdk/authorization"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type StackRolesArgs struct{}

type StackRoles struct {
	pulumi.ResourceState

	Name               string
	ClientID           pulumi.StringOutput
	TenantID           pulumi.StringOutput
	ServicePrincipalId pulumi.StringOutput
	ClientSecret       pulumi.StringOutput
}

type RoleDefinition struct {
	Description      pulumi.StringInput
	Permissions      authorization.PermissionArray
	AssignableScopes pulumi.StringArray
}

var roleDefinitions = map[resourcespb.Action]RoleDefinition{
	resourcespb.Action_KeyValueStoreRead: {
		Description: pulumi.String("keyvalue read access"),
		Permissions: authorization.PermissionArray{
			authorization.PermissionArgs{
				Actions: pulumi.StringArray{},
				DataActions: pulumi.StringArray{
					pulumi.String("Microsoft.Storage/storageAccounts/tableServices/tables/entities/read"),
				},
				NotActions: pulumi.StringArray{},
			},
		},
		AssignableScopes: pulumi.ToStringArray([]string{
			"/",
		}),
	},
	resourcespb.Action_QueueList: {
		Description: pulumi.String("queue list access"),
		Permissions: authorization.PermissionArray{
			authorization.PermissionArgs{
				Actions: pulumi.StringArray{
					pulumi.String("Microsoft.Storage/storageAccounts/queueServices/queues/read"),
				},
				DataActions: pulumi.StringArray{},
				NotActions:  pulumi.StringArray{},
			},
		},
		AssignableScopes: pulumi.ToStringArray([]string{
			"/",
		}),
	},
	resourcespb.Action_QueueDetail: {
		Description: pulumi.String("queue detail access"),
		Permissions: authorization.PermissionArray{
			authorization.PermissionArgs{
				Actions: pulumi.StringArray{
					pulumi.String("Microsoft.Storage/storageAccounts/queueServices/queues/read"),
				},
				DataActions: pulumi.StringArray{},
				NotActions:  pulumi.StringArray{},
			},
		},
		AssignableScopes: pulumi.ToStringArray([]string{
			"/",
		}),
	},
	resourcespb.Action_QueueSend: {
		Description: pulumi.String("queue send access"),
		Permissions: authorization.PermissionArray{
			authorization.PermissionArgs{
				Actions: pulumi.StringArray{},
				DataActions: pulumi.StringArray{
					pulumi.String("Microsoft.Storage/storageAccounts/queueServices/queues/messages/write"),
				},
				NotActions: pulumi.StringArray{},
			},
		},
		AssignableScopes: pulumi.ToStringArray([]string{
			"/",
		}),
	},
	resourcespb.Action_QueueReceive: {
		Description: pulumi.String("queue receive access"),
		Permissions: authorization.PermissionArray{
			authorization.PermissionArgs{
				Actions: pulumi.StringArray{
					pulumi.String("Microsoft.Storage/storageAccounts/queueServices/queues/read"),
				},
				DataActions: pulumi.StringArray{
					pulumi.String("Microsoft.Storage/storageAccounts/queueServices/queues/messages/read"),
					pulumi.String("Microsoft.Storage/storageAccounts/queueServices/queues/messages/delete"),
					// pulumi.String("Microsoft.Storage/storageAccounts/queueServices/queues/messages/update"),
				},
				NotActions: pulumi.StringArray{},
			},
		},
		AssignableScopes: pulumi.ToStringArray([]string{
			"/",
		}),
	},
	resourcespb.Action_KeyValueStoreWrite: {
		Description: pulumi.String("keyvalue write access"),
		Permissions: authorization.PermissionArray{
			authorization.PermissionArgs{
				Actions: pulumi.StringArray{},
				DataActions: pulumi.StringArray{
					pulumi.String("Microsoft.Storage/storageAccounts/tableServices/tables/entities/write"),
				},
				NotActions: pulumi.StringArray{},
			},
		},
		AssignableScopes: pulumi.ToStringArray([]string{
			"/",
		}),
	},
	resourcespb.Action_KeyValueStoreDelete: {
		Description: pulumi.String("keyvalue delete access"),
		Permissions: authorization.PermissionArray{
			authorization.PermissionArgs{
				Actions: pulumi.StringArray{},
				DataActions: pulumi.StringArray{
					pulumi.String("Microsoft.Storage/storageAccounts/tableServices/tables/entities/delete"),
				},
				NotActions: pulumi.StringArray{},
			},
		},
		AssignableScopes: pulumi.ToStringArray([]string{
			"/",
		}),
	},
	resourcespb.Action_BucketFileGet: {
		Description: pulumi.String("bucket read access"),
		Permissions: authorization.PermissionArray{
			authorization.PermissionArgs{
				Actions: pulumi.StringArray{
					pulumi.String("Microsoft.Storage/storageAccounts/blobServices/containers/read"),
				},
				DataActions: pulumi.StringArray{
					pulumi.String("Microsoft.Storage/storageAccounts/blobServices/containers/blobs/read"),
				},
				NotActions: pulumi.StringArray{},
			},
		},
		AssignableScopes: pulumi.ToStringArray([]string{
			"/",
		}),
	},
	resourcespb.Action_BucketFilePut: {
		Description: pulumi.String("bucket file write access"),
		Permissions: authorization.PermissionArray{
			authorization.PermissionArgs{
				Actions: pulumi.StringArray{},
				DataActions: pulumi.StringArray{
					pulumi.String("Microsoft.Storage/storageAccounts/blobServices/containers/blobs/write"),
				},
				NotActions: pulumi.StringArray{},
			},
		},
		AssignableScopes: pulumi.ToStringArray([]string{
			"/",
		}),
	},
	resourcespb.Action_BucketFileList: {
		Description: pulumi.String("bucket file list access"),
		Permissions: authorization.PermissionArray{
			authorization.PermissionArgs{
				DataActions: pulumi.StringArray{
					pulumi.String("Microsoft.Storage/storageAccounts/blobServices/containers/blobs/read"),
				},
				NotActions: pulumi.StringArray{},
			},
		},
		AssignableScopes: pulumi.ToStringArray([]string{
			"/",
		}),
	},
	resourcespb.Action_BucketFileDelete: {
		Description: pulumi.String("bucket file delete access"),
		Permissions: authorization.PermissionArray{
			authorization.PermissionArgs{
				Actions: pulumi.StringArray{},
				DataActions: pulumi.StringArray{
					pulumi.String("Microsoft.Storage/storageAccounts/blobServices/containers/blobs/delete"),
				},
				NotActions: pulumi.StringArray{},
			},
		},
		AssignableScopes: pulumi.ToStringArray([]string{
			"/",
		}),
	},
	resourcespb.Action_TopicDetail: {
		Description: pulumi.String("topic detail access"),
		Permissions: authorization.PermissionArray{
			authorization.PermissionArgs{
				Actions: pulumi.StringArray{
					pulumi.String("Microsoft.EventGrid/topics/read"),
				},
				DataActions: pulumi.StringArray{},
				NotActions:  pulumi.StringArray{},
			},
		},
		AssignableScopes: pulumi.ToStringArray([]string{
			"/",
		}),
	},
	resourcespb.Action_TopicList: {
		Description: pulumi.String("topic list access"),
		Permissions: authorization.PermissionArray{
			authorization.PermissionArgs{
				Actions: pulumi.StringArray{
					pulumi.String("Microsoft.EventGrid/topics/read"),
				},
				DataActions: pulumi.StringArray{},
				NotActions:  pulumi.StringArray{},
			},
		},
		AssignableScopes: pulumi.ToStringArray([]string{
			"/",
		}),
	},
	resourcespb.Action_TopicEventPublish: {
		Description: pulumi.String("topic event publish access"),
		Permissions: authorization.PermissionArray{
			authorization.PermissionArgs{
				Actions: pulumi.StringArray{
					pulumi.String("Microsoft.EventGrid/topics/*/write"),
				},
				DataActions: pulumi.StringArray{
					pulumi.String("Microsoft.EventGrid/events/send/action"),
				},
				NotActions: pulumi.StringArray{},
			},
		},
		AssignableScopes: pulumi.ToStringArray([]string{
			"/",
		}),
	},
	resourcespb.Action_SecretAccess: {
		Description: pulumi.String("keyvault secret read access"),
		Permissions: authorization.PermissionArray{
			authorization.PermissionArgs{
				Actions: pulumi.StringArray{},
				DataActions: pulumi.StringArray{
					pulumi.String("Microsoft.KeyVault/vaults/secrets/getSecret/action"),
				},
				NotActions: pulumi.StringArray{},
			},
		},
		AssignableScopes: pulumi.ToStringArray([]string{
			"/",
		}),
	},
	resourcespb.Action_SecretPut: {
		Description: pulumi.String("keyvault secret write access"),
		Permissions: authorization.PermissionArray{
			authorization.PermissionArgs{
				Actions: pulumi.StringArray{
					pulumi.String("Microsoft.KeyVault/vaults/secrets/write"),
				},
				DataActions: pulumi.StringArray{
					pulumi.String("Microsoft.KeyVault/vaults/secrets/setSecret/action"),
				},
				NotActions: pulumi.StringArray{},
			},
		},
		AssignableScopes: pulumi.ToStringArray([]string{
			"/",
		}),
	},
}

type Roles struct {
	pulumi.ResourceState

	Name            string
	RoleDefinitions map[resourcespb.Action]*authorization.RoleDefinition
}

var actionNames = map[resourcespb.Action]string{
	resourcespb.Action_BucketFileGet:       "BucketFileGet",
	resourcespb.Action_BucketFilePut:       "BucketFilePut",
	resourcespb.Action_BucketFileDelete:    "BucketFileDelete",
	resourcespb.Action_BucketFileList:      "BucketFileList",
	resourcespb.Action_TopicDetail:         "TopicDetail",
	resourcespb.Action_TopicEventPublish:   "TopicPublish",
	resourcespb.Action_TopicList:           "TopicList",
	resourcespb.Action_SecretAccess:        "SecretAccess",
	resourcespb.Action_SecretPut:           "SecretPut",
	resourcespb.Action_KeyValueStoreDelete: "KeyValueStoreDelete",
	resourcespb.Action_KeyValueStoreRead:   "KeyValueStoreRead",
	resourcespb.Action_KeyValueStoreWrite:  "KeyValueStoreWrite",
	resourcespb.Action_QueueSend:           "QueueSend",
	resourcespb.Action_QueueReceive:        "QueueReceive",
	resourcespb.Action_QueueDetail:         "QueueDetail",
	resourcespb.Action_QueueList:           "QueueList",
}

func CreateRoles(ctx *pulumi.Context, stackId string, subscriptionId string, rgName pulumi.StringInput) (*Roles, error) {
	res := &Roles{Name: "nitric-roles", RoleDefinitions: map[resourcespb.Action]*authorization.RoleDefinition{}}

	err := ctx.RegisterComponentResource("nitricazure:AzureADRoles", "nitric-roles", res)
	if err != nil {
		return nil, err
	}

	for id, roleDef := range roleDefinitions {
		name := actionNames[id]

		roleGuid, err := random.NewRandomUuid(ctx, name, &random.RandomUuidArgs{
			Keepers: pulumi.ToMap(map[string]interface{}{
				"subscriptionId": subscriptionId,
			}),
		}, pulumi.Parent(res))
		if err != nil {
			return nil, err
		}

		roleName := fmt.Sprintf("%s-%s", stackId, name)

		createdRole, err := authorization.NewRoleDefinition(ctx, name, &authorization.RoleDefinitionArgs{
			RoleDefinitionId: roleGuid.Result,
			RoleName:         pulumi.String(roleName),
			Description:      roleDef.Description,
			Permissions:      roleDef.Permissions,
			Scope:            pulumi.Sprintf("/subscriptions/%s/resourceGroups/%s", subscriptionId, rgName),
			AssignableScopes: pulumi.StringArray{
				pulumi.Sprintf("/subscriptions/%s/resourceGroups/%s", subscriptionId, rgName),
			},
		}, pulumi.Parent(res))
		if err != nil {
			return nil, err
		}

		res.RoleDefinitions[id] = createdRole
	}

	return res, nil
}
