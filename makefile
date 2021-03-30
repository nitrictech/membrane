membrane: install
	@echo Building Go Project...
	@CGO_ENABLED=1 GOOS=linux go build -o bin/membrane pluggable_membrane.go

install:
	@echo installing go dependencies
	@go mod download

install-tools: install
	@echo Installing tools from tools.go
	@cat ./tools/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

clean:
	@rm -rf ./bin/
	@rm -rf ./lib/
	@rm -rf ./interfaces/

# Run all tests
test: test-membrane test-aws-plugins test-gcp-plugins test-dev-plugins
	@echo Done.

# Generate interfaces
generate-proto:
	@echo Generating Proto Sources
	@mkdir -p ./interfaces/
	@protoc --go_out=./interfaces/ --go-grpc_out=./interfaces/ -I ./contracts/proto ./contracts/proto/*/**/*.proto

# Build all service factory plugins
plugins: aws-plugin gcp-plugin dev-plugin
	@echo Done.

# Test the membrane
test-membrane: install-tools generate-proto 
	@echo Testing Membrane
	@go run github.com/onsi/ginkgo/ginkgo -cover ./membrane/...

# BEGIN AWS Plugins
test-aws-plugins:
	@echo Testing AWS Plugins
	@go run github.com/onsi/ginkgo/ginkgo -cover ./plugins/aws/...

aws-static: generate-proto
	@echo Building static AWS membrane
	@CGO_ENABLED=0 GOOS=linux go build -o bin/membrane -ldflags="-extldflags=-static" ./plugins/aws/static_membrane.go

# Cross-platform Build
aws-static-xp: generate-proto
	@echo Building static AWS membrane
	@CGO_ENABLED=0 go build -o bin/membrane -ldflags="-extldflags=-static" ./plugins/aws/static_membrane.go

# Service Factory Plugin for Pluggable Membrane
aws-plugin:
	@echo Building AWS Service Factory Plugin
	@go build -buildmode=plugin -o lib/plugins/aws.so ./plugins/aws/plugin.go

aws-docker-static:
	@docker build . -f ./plugins/aws/aws.dockerfile -t nitricimages/membrane-aws

aws-docker: aws-docker-static
	@echo Built AWS Docker Images
# END AWS Plugins

# BEGIN Azure Plugins
test-azure-plugins:
	@echo Testing Azure Plugins
	@go run github.com/onsi/ginkgo/ginkgo -cover ./plugins/azure/...

azure-static: generate-proto
	@echo Building static Azure membrane
	@CGO_ENABLED=0 GOOS=linux go build -o bin/membrane -ldflags="-extldflags=-static" ./plugins/azure/static_membrane.go

# Cross-platform Build
azure-static-xp: generate-proto
	@echo Building static Azure membrane
	@CGO_ENABLED=0 go build -o bin/membrane -ldflags="-extldflags=-static" ./plugins/azure/static_membrane.go

# Service Factory Plugin for Pluggable Membrane
azure-plugin:
	@echo Building Azure Service Factory Plugin
	@go build -buildmode=plugin -o lib/plugins/azure.so ./plugins/azure/plugin.go

azure-docker-static:
	@docker build . -f ./plugins/azure/azure.dockerfile -t nitricimages/membrane-azure

azure-docker: azure-docker-static # azure-docker-alpine azure-docker-debian
	@echo Built Azure Docker Images
# END Azure Plugins

# BEGIN GCP Plugins
test-gcp-plugins:
	@echo Testing GCP Plugins
	@go run github.com/onsi/ginkgo/ginkgo -cover ./plugins/gcp/...

gcp-static: generate-proto
	@echo Building static GCP membrane
	@CGO_ENABLED=0 GOOS=linux go build -o bin/membrane -ldflags="-extldflags=-static" ./plugins/gcp/static_membrane.go

# Cross-platform Build
gcp-static-xp: generate-proto
	@echo Building static GCP membrane
	@CGO_ENABLED=0 go build -o bin/membrane -ldflags="-extldflags=-static" ./plugins/gcp/static_membrane.go

# Service Factory Plugin for Pluggable Membrane
gcp-plugin:
	@echo Building GCP Service Factory Plugin
	@go build -buildmode=plugin -o lib/plugins/gcp.so ./plugins/gcp/plugin.go

gcp-docker-static:
	@docker build . -f ./plugins/gcp/gcp.dockerfile -t nitricimages/membrane-gcp

gcp-docker: gcp-docker-static # gcp-docker-alpine gcp-docker-debian
	@echo Built GCP Docker Images
# END GCP Plugins

# BEGIN Local Plugins
# Cross-platform build only, this membrane is not for production use.
dev-static: generate-proto
	@echo Building static Local membrane
	@CGO_ENABLED=0 go build -o bin/membrane -ldflags="-extldflags=-static" ./plugins/dev/static_membrane.go

# Service Factory Plugin for Pluggable Membrane
dev-plugin:
	@echo Building Dev Service Factory Plugin
	@go build -buildmode=plugin -o lib/plugins/dev.so ./plugins/dev/plugin.go

dev-docker-static:
	@docker build . -f ./plugins/dev/dev.dockerfile -t nitricimages/membrane-local

dev-docker: dev-docker-static
	@echo Built Local Docker Images

test-dev-plugins:
	@echo Testing Dev Plugins
	@go run github.com/onsi/ginkgo/ginkgo -cover ./plugins/dev/...
# END Local Plugins

membrane-docker-alpine: generate-proto
	@docker build . -f alpine.dockerfile -t nitric:membrane-alpine
membrane-docker-debian: generate-proto
	@docker build . -f debian.dockerfile -t nitric:membrane-debian

# Generate proto files locally before building docker images
# TODO: Get alpine image generating its own sources
membrane-docker: generate-proto membrane-docker-alpine membrane-docker-debian
	@echo Built Docker Images