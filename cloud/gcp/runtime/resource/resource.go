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

package resource

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/nitrictech/nitric/cloud/common/deploy/tags"
	"github.com/nitrictech/nitric/cloud/gcp/runtime/env"

	apigateway "cloud.google.com/go/apigateway/apiv1"
	"cloud.google.com/go/apigateway/apiv1/apigatewaypb"

	commonenv "github.com/nitrictech/nitric/cloud/common/runtime/env"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
)

type GcpResourceProvider interface {
	// GetServiceAccountEmail for google cloud projects
	GetServiceAccountEmail() (string, error)
	GetProjectID() (string, error)
}

type GcpResourceService struct {
	apiClient           *apigateway.Client
	stackName           string
	serviceAccountEmail string
	projectID           string
	region              string
}

var _ GcpResourceProvider = &GcpResourceService{}
var _ v1.ResourcesServer = &GcpResourceService{}

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

func (g *GcpResourceService) getApiGatewayDetails(ctx context.Context, name string) (*v1.ResourceDetailsResponse, error) {
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
		return &v1.ResourceDetailsResponse{
			Id:       gw.Name,
			Provider: "gcp",
			Service:  "ApiGateway",
			Details: &v1.ResourceDetailsResponse_Api{
				Api: &v1.ApiResourceDetails{
					Url: fmt.Sprintf("https://%s", gw.DefaultHostname),
				},
			},
		}, nil
	} else {
		return nil, err
	}
}

func (g *GcpResourceService) Declare(ctx context.Context, req *v1.ResourceDeclareRequest) (*v1.ResourceDeclareResponse, error) {
	return &v1.ResourceDeclareResponse{}, nil
}

func (g *GcpResourceService) Details(ctx context.Context, req *v1.ResourceDetailsRequest) (*v1.ResourceDetailsResponse, error) {
	resourceName := req.Id.Name
	resourceType := req.Id.Type

	switch resourceType {
	case v1.ResourceType_Api:
		return g.getApiGatewayDetails(ctx, resourceName)
	default:
		return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
	}
}

func (g *GcpResourceService) GetProjectID() (string, error) {
	if g.projectID == "" {
		if env := env.GOOGLE_PROJECT_ID.String(); env != "" {
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

func (g *GcpResourceService) GetServiceAccountEmail() (string, error) {
	if g.serviceAccountEmail == "" {
		if env := env.SERVICE_ACCOUNT_EMAIL.String(); env != "" {
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

func New() (*GcpResourceService, error) {
	stack := commonenv.NITRIC_STACK_ID.String()
	region := env.GCP_REGION.String()

	apiClient, err := apigateway.NewClient(context.TODO())
	if err != nil {
		return nil, err
	}

	return &GcpResourceService{
		stackName: stack,
		apiClient: apiClient,
		region:    region,
	}, nil
}
