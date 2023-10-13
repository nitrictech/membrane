// Copyright 2021 Nitric Pty Ltd.
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

package core

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/nitrictech/nitric/cloud/common/deploy/tags"

	apigateway "cloud.google.com/go/apigateway/apiv1"
	"cloud.google.com/go/apigateway/apiv1/apigatewaypb"

	"github.com/nitrictech/nitric/core/pkg/plugins/resource"
	"github.com/nitrictech/nitric/core/pkg/utils"
)

type GcpProvider interface {
	// GetServiceAccountEmail for google cloud projects
	GetServiceAccountEmail() (string, error)
	GetProjectID() (string, error)
	resource.ResourceService
}

type gcpProviderImpl struct {
	apiClient           *apigateway.Client
	stackName           string
	serviceAccountEmail string
	projectID           string
	region              string
}

var _ resource.ResourceService = &gcpProviderImpl{}

const (
	metadataFlavorKey      = "Metadata-Flavor"
	metadataFlavorValue    = "Google"
	serviceAccountEnv      = "SERVICE_ACCOUNT_EMAIL"
	serviceAccountEmailUri = "http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/default/email"
	projectIdEnv           = "GOOGLE_PROJECT_ID"
	projectIdUri           = "http://metadata.google.internal/computeMetadata/v1/project/project-id"
)

func createMetadataRequest(uri string) (*http.Request, error) {
	req, err := http.NewRequest("GET", projectIdUri, nil)
	if err != nil {
		return nil, err
	}

	// Add the correct header
	req.Header.Add(metadataFlavorKey, metadataFlavorValue)

	return req, nil
}

type GenericIterator[T any] interface {
	Next() (T, error)
}

func filter(stackName string, name string) string {
	return fmt.Sprintf("labels.%s:%s", tags.GetResourceNameKey(stackName), name)
}

func (g *gcpProviderImpl) getApiGatewayDetails(ctx context.Context, name string) (*resource.DetailsResponse[any], error) {
	projectName, err := g.GetProjectID()
	if err != nil {
		return nil, err
	}

	gws := g.apiClient.ListGateways(ctx, &apigatewaypb.ListGatewaysRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectName, g.region),
		Filter: filter(g.stackName, name),
	})

	// there should only be a single entry in this array, we'll grab the first and then break
	if gw, err := gws.Next(); gw != nil && err == nil {
		return &resource.DetailsResponse[any]{
			Id:       gw.Name,
			Provider: "gcp",
			Service:  "ApiGateway",
			Detail: resource.ApiDetails{
				URL: fmt.Sprintf("https://%s", gw.DefaultHostname),
			},
		}, nil
	} else {
		return nil, err
	}
}

func (g *gcpProviderImpl) Declare(ctx context.Context, req resource.ResourceDeclareRequest) error {
	return nil
}

func (g *gcpProviderImpl) Details(ctx context.Context, typ resource.ResourceType, name string) (*resource.DetailsResponse[any], error) {
	switch typ {
	case resource.ResourceType_Api:
		return g.getApiGatewayDetails(ctx, name)
	default:
		return nil, fmt.Errorf("unsupported resource type: %s", typ)
	}
}

func (g *gcpProviderImpl) GetProjectID() (string, error) {
	if g.projectID == "" {
		if env := utils.GetEnv(projectIdEnv, ""); env != "" {
			return env, nil
		}

		req, err := createMetadataRequest(projectIdUri)
		if err != nil {
			return "", err
		}

		// read the response as the service account email
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return "", err
		}

		projectIdBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return "", err
		}

		g.projectID = string(projectIdBytes)
	}

	return g.projectID, nil
}

func (g *gcpProviderImpl) GetServiceAccountEmail() (string, error) {
	if g.serviceAccountEmail == "" {
		if env := utils.GetEnv(serviceAccountEnv, ""); env != "" {
			return env, nil
		}

		req, err := createMetadataRequest(serviceAccountEmailUri)
		if err != nil {
			return "", err
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return "", err
		}

		// read the response as the service account email
		emailBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return "", err
		}

		g.serviceAccountEmail = string(emailBytes)
	}

	return g.serviceAccountEmail, nil
}

func New() (GcpProvider, error) {
	stack := utils.GetEnv("NITRIC_STACK_ID", "")
	region := utils.GetEnv("GCP_REGION", "")

	apiClient, err := apigateway.NewClient(context.TODO())
	if err != nil {
		return nil, err
	}

	return &gcpProviderImpl{
		stackName: stack,
		apiClient: apiClient,
		region:    region,
	}, nil
}
