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

package pulumix

import (
	"fmt"
	"strings"

	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	resourcespb "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/events"
	"github.com/pulumi/pulumi/sdk/v3/go/common/apitype"
)

type pulumiEventHandler struct {
	tree DataTree
}

func (p *pulumiEventHandler) handleResourcePreEvent(resourcePreEvent *apitype.ResourcePreEvent) (*deploymentspb.ResourceUpdate, error) {
	// Add as a root key to the map

	// Happens before resource updates are applied
	// attempt to locate the parent node
	meta := resourcePreEvent.Metadata.New
	if meta == nil {
		meta = resourcePreEvent.Metadata.Old
	}

	parentNode := p.tree.FindNode(meta.Parent)
	if parentNode == nil {
		parentNode = p.tree.Root
	}

	nitricResource := NitricResourceIdFromPulumiUrn(resourcePreEvent.Metadata.URN)

	// nitricAction := deploymentspb.ResourceDeploymentAction_SAME

	/*
		TODO: review unhandled cases (currently defaulting to update)
		OpRead
		OpReadReplacement
		OpRefresh
		OpReadDiscard
		OpDiscardReplaced
		OpRemovePendingReplace
		OpImport
		OpImportReplacement
	*/

	var nitricAction deploymentspb.ResourceDeploymentAction
	switch resourcePreEvent.Metadata.Op {
	case apitype.OpCreate:
		nitricAction = deploymentspb.ResourceDeploymentAction_CREATE
	case apitype.OpDelete, apitype.OpDeleteReplaced:
		nitricAction = deploymentspb.ResourceDeploymentAction_DELETE
	case apitype.OpReplace, apitype.OpCreateReplacement:
		nitricAction = deploymentspb.ResourceDeploymentAction_REPLACE
	case apitype.OpSame:
		nitricAction = deploymentspb.ResourceDeploymentAction_SAME
	case apitype.OpUpdate:
		nitricAction = deploymentspb.ResourceDeploymentAction_UPDATE
	default:
		nitricAction = deploymentspb.ResourceDeploymentAction_UPDATE
	}

	node := &DataNode{
		Id: resourcePreEvent.Metadata.URN,
		// TODO: Populate the nitric resource
		Data: &ResourceData{
			nitricResource: nitricResource,
			action:         nitricAction,
		},
		Children: make([]*DataNode, 0),
	}

	parentNode.AddChild(node)

	if nitricResource != nil {
		return &deploymentspb.ResourceUpdate{
			Id:          nitricResource,
			Action:      nitricAction,
			Status:      deploymentspb.ResourceDeploymentStatus_PENDING,
			SubResource: "",
		}, nil
	} else {
		nitricResourceNode := node.FindParent(func(n *DataNode) bool {
			return n.Data != nil && n.Data.nitricResource != nil
		})
		if nitricResourceNode != nil {
			nitricResource = nitricResourceNode.Data.nitricResource
		}

		urnParts := strings.Split(resourcePreEvent.Metadata.URN, "$")
		return &deploymentspb.ResourceUpdate{
			Id:          nitricResource,
			Action:      node.Data.action,
			Status:      deploymentspb.ResourceDeploymentStatus_IN_PROGRESS,
			SubResource: urnParts[len(urnParts)-1],
		}, nil
	}
}

func (p *pulumiEventHandler) handleResourceOutputsEvent(resOutputsEvent *apitype.ResOutputsEvent) (*deploymentspb.ResourceUpdate, error) {
	// Happens after resource updates are applied
	// Find the URN in the tree
	urn := resOutputsEvent.Metadata.URN
	// Happens after resource updates fail to apply
	resourceNode := p.tree.FindNode(urn)
	var nitricResource *resourcespb.ResourceIdentifier
	var subResource string

	if resourceNode.Data != nil && resourceNode.Data.nitricResource != nil {
		// we have a nitric resource
		nitricResource = resourceNode.Data.nitricResource
	} else {
		// just a regular pleb resource
		// find its nitric parent

		nitricParentNode := resourceNode.FindParent(func(n *DataNode) bool {
			return n.Data != nil && n.Data.nitricResource != nil
		})

		if nitricParentNode != nil {
			nitricResource = nitricParentNode.Data.nitricResource
		}

		urnParts := strings.Split(urn, "$")
		subResource = urnParts[len(urnParts)-1]
	}

	return &deploymentspb.ResourceUpdate{
		Id:          nitricResource,
		Action:      resourceNode.Data.action,
		Status:      deploymentspb.ResourceDeploymentStatus_SUCCESS,
		SubResource: subResource,
	}, nil
}

func (p *pulumiEventHandler) handleResourceFailedEvent(resOpFailedEvent *apitype.ResOpFailedEvent) (*deploymentspb.ResourceUpdate, error) {
	urn := resOpFailedEvent.Metadata.URN
	// Happens after resource updates fail to apply
	resourceNode := p.tree.FindNode(urn)
	var nitricResource *resourcespb.ResourceIdentifier
	var subResource string

	if resourceNode.Data != nil && resourceNode.Data.nitricResource != nil {
		// we have a nitric resource
		nitricResource = resourceNode.Data.nitricResource
	} else {
		// just a regular pleb resource
		// find its nitric parent

		nitricParentNode := resourceNode.FindParent(func(n *DataNode) bool {
			return n.Data != nil && n.Data.nitricResource != nil
		})

		if nitricParentNode != nil {
			nitricResource = nitricParentNode.Data.nitricResource
		}

		urnParts := strings.Split(urn, "$")
		subResource = urnParts[len(urnParts)-1]
	}

	return &deploymentspb.ResourceUpdate{
		Id:          nitricResource,
		Action:      resourceNode.Data.action,
		Status:      deploymentspb.ResourceDeploymentStatus_FAILED,
		SubResource: subResource,
	}, nil
}

func (p *pulumiEventHandler) engineEventToResourceUpdate(evt events.EngineEvent) (*deploymentspb.ResourceUpdate, error) {
	if evt.ResourcePreEvent != nil {
		p.handleResourcePreEvent(evt.ResourcePreEvent)
	} else if evt.ResOutputsEvent != nil {
		p.handleResourceOutputsEvent(evt.ResOutputsEvent)
	} else if evt.ResOpFailedEvent != nil {
		p.handleResourceFailedEvent(evt.ResOpFailedEvent)
	}

	return nil, fmt.Errorf("unknown event type")
}

type ResourceData struct {
	nitricResource *resourcespb.ResourceIdentifier
	action         deploymentspb.ResourceDeploymentAction
}

type (
	DataTree = Tree[ResourceData]
	DataNode = Node[ResourceData]
)

func StreamPulumiUpEngineEvents(stream deploymentspb.Deployment_UpServer, pulumiEventsChan <-chan events.EngineEvent) error {
	evtHandler := pulumiEventHandler{
		tree: DataTree{
			Root: &DataNode{
				Id:       "stack",
				Data:     nil,
				Children: make([]*DataNode, 0),
			},
		},
	}

	for evt := range pulumiEventsChan {
		// translate the engine event to a server message and send back to the CLIent
		updateDetails, err := evtHandler.engineEventToResourceUpdate(evt)
		if err != nil {
			fmt.Println("encountered an error: ", err.Error())
			continue
		}

		err = stream.Send(&deploymentspb.DeploymentUpEvent{
			Content: &deploymentspb.DeploymentUpEvent_Update{
				Update: updateDetails,
			},
		})

		if err != nil {
			return err
		}
	}
	return nil
}

func StreamPulumiDownEngineEvents(stream deploymentspb.Deployment_DownServer, pulumiEventsChan <-chan events.EngineEvent) error {
	evtHandler := pulumiEventHandler{
		tree: DataTree{
			Root: &DataNode{
				Id:       "stack",
				Data:     nil,
				Children: make([]*DataNode, 0),
			},
		},
	}

	for evt := range pulumiEventsChan {
		// translate the engine event to a server message and send back to the CLIent
		nitricEvent, err := evtHandler.engineEventToResourceUpdate(evt)
		if err != nil {
			fmt.Println("encountered an error: ", err.Error())
			continue
		}

		err = stream.Send(&deploymentspb.DeploymentDownEvent{
			Content: &deploymentspb.DeploymentDownEvent_Update{
				Update: nitricEvent,
			},
		})

		if err != nil {
			return err
		}
	}

	return nil
}
