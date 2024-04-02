package deploy

import (
	"fmt"
	"strings"

	"github.com/nitrictech/nitric/cloud/common/deploy/image"
	"github.com/nitrictech/nitric/cloud/common/deploy/provider"
	"github.com/nitrictech/nitric/cloud/common/deploy/pulumix"
	"github.com/nitrictech/nitric/cloud/common/deploy/resources"
	"github.com/nitrictech/nitric/cloud/common/deploy/tags"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/pulumi/pulumi-oci/sdk/go/oci/artifacts"
	"github.com/pulumi/pulumi-oci/sdk/go/oci/functions"
	"github.com/pulumi/pulumi-oci/sdk/go/oci/identity"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// This represents a unique unit of execution at the moment this is a container but could also be many things e.g. WASM, Binary, source zip etc.
type Service struct {
	pulumi.ResourceState
	Name string
}

func (n *NitricOCIPulumiProvider) createContainerRepository(ctx *pulumi.Context, parent pulumi.Resource, name string) (*artifacts.ContainerRepository, error) {
	return artifacts.NewContainerRepository(ctx, name, &artifacts.ContainerRepositoryArgs{
		DisplayName:   pulumi.String(name),
		CompartmentId: n.compartment.CompartmentId,
		FreeformTags:  pulumi.ToMap(tags.TagsAsInterface(n.stackId, name, resources.Service)),
		IsPublic:      pulumi.Bool(false),
	}, pulumi.Parent(parent))
}

func (n *NitricOCIPulumiProvider) createImage(ctx *pulumi.Context, parent pulumi.Resource, name string, repo *artifacts.ContainerRepository, config *deploymentspb.Service, runtime provider.RuntimeProvider) (*image.Image, error) {
	if config.GetImage() == nil {
		return nil, fmt.Errorf("oci provider can only deploy service with an image source")
	}

	if config.GetImage().GetUri() == "" {
		return nil, fmt.Errorf("oci provider can only deploy service with an image source")
	}

	if config.Type == "" {
		config.Type = "default"
	}

	authToken, err := identity.NewAuthToken(ctx, name, &identity.AuthTokenArgs{
		UserId:      n.serviceAccount.ID(),
		Description: pulumi.String("AuthToken for OCI CLI authentication"),
	})
	if err != nil {
		return nil, err
	}

	return image.NewImage(ctx, name, &image.ImageArgs{
		SourceImage:   config.GetImage().GetUri(),
		RepositoryUrl: pulumi.Sprintf("%s.ocir.io/%s/%s:latest", n.config.Region, repo.Namespace, repo.DisplayName),
		Username:      pulumi.Sprintf("%s/%s", repo.Namespace, n.serviceAccount.Email),
		Password:      authToken.Token,
		Runtime:       runtime(),
	}, pulumi.Parent(parent), pulumi.DependsOn([]pulumi.Resource{repo}))
}

func (a *NitricOCIPulumiProvider) Service(ctx *pulumi.Context, parent pulumi.Resource, name string, config *pulumix.NitricPulumiServiceConfig, runtime provider.RuntimeProvider) error {
	// Create the ECR repository to push the image to
	repo, err := a.createContainerRepository(ctx, parent, name)
	if err != nil {
		return err
	}

	image, err := a.createImage(ctx, parent, name, repo, config.Service, runtime)
	if err != nil {
		return err
	}

	appName := fmt.Sprintf("application-%s", name)
	app, err := functions.NewApplication(ctx, appName, &functions.ApplicationArgs{
		DisplayName:   pulumi.String(appName),
		CompartmentId: a.compartment.CompartmentId,
		FreeformTags:  pulumi.ToMap(tags.TagsAsInterface(a.stackId, name, resources.Service)),
		SubnetIds:     pulumi.StringArray{a.subnet.ID()},
	})
	if err != nil {
		return err
	}

	// Get the image uri in the correct form
	imageUri := image.URI().ApplyT(func(uri string) (string, error) {
		parts := strings.Split(uri, "@")
		if len(parts) < 2 {
			return "", fmt.Errorf("could not extract image uri")
		}

		return fmt.Sprintf("%s:latest", parts[0]), nil
	}).(pulumi.StringOutput)

	// Get the image digest by itself
	imageDigest := image.URI().ApplyT(func(uri string) (string, error) {
		parts := strings.Split(uri, "@")
		if len(parts) < 2 {
			return "", fmt.Errorf("could not extract image digest from uri")
		}

		return parts[1], nil
	}).(pulumi.StringOutput)

	function, err := functions.NewFunction(ctx, name, &functions.FunctionArgs{
		DisplayName:   pulumi.String(name),
		ApplicationId: app.ID(),
		MemoryInMbs:   pulumi.String("512"),
		FreeformTags:  pulumi.ToMap(tags.TagsAsInterface(a.stackId, name, resources.Service)),
		Image:         imageUri,
		ImageDigest:   imageDigest,
	})
	if err != nil {
		return err
	}

	a.functions[name] = function

	return nil
}
