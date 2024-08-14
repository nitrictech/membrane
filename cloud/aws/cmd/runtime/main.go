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

package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/nitrictech/nitric/cloud/aws/runtime/api"
	"github.com/nitrictech/nitric/cloud/aws/runtime/batch"
	"github.com/nitrictech/nitric/cloud/aws/runtime/env"
	lambda_service "github.com/nitrictech/nitric/cloud/aws/runtime/gateway"
	dynamodb_service "github.com/nitrictech/nitric/cloud/aws/runtime/keyvalue"
	sqs_service "github.com/nitrictech/nitric/cloud/aws/runtime/queue"
	"github.com/nitrictech/nitric/cloud/aws/runtime/resource"
	secrets_manager_secret_service "github.com/nitrictech/nitric/cloud/aws/runtime/secret"
	sql_service "github.com/nitrictech/nitric/cloud/aws/runtime/sql"
	s3_service "github.com/nitrictech/nitric/cloud/aws/runtime/storage"
	sns_service "github.com/nitrictech/nitric/cloud/aws/runtime/topic"
	"github.com/nitrictech/nitric/cloud/aws/runtime/websocket"
	commonenv "github.com/nitrictech/nitric/cloud/common/runtime/env"
	base_http "github.com/nitrictech/nitric/cloud/common/runtime/gateway"
	"github.com/nitrictech/nitric/core/pkg/logger"
	"github.com/nitrictech/nitric/core/pkg/membrane"
	"github.com/nitrictech/nitric/core/pkg/membrane/job"
)

func launchService() {
	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	signal.Notify(term, os.Interrupt, syscall.SIGINT)

	logger.SetLogLevel(logger.INFO)

	gatewayEnv := env.GATEWAY_ENVIRONMENT.String()

	membraneOpts := membrane.DefaultMembraneOptions()

	provider, err := resource.New()
	if err != nil {
		logger.Fatalf("could not create aws provider: %v", err)
		return
	}

	// Load the appropriate gateway based on the environment.
	switch gatewayEnv {
	case "lambda":
		membraneOpts.GatewayPlugin, _ = lambda_service.New(provider)
	default:
		membraneOpts.GatewayPlugin, _ = base_http.NewHttpGateway(nil)
	}

	membraneOpts.ApiPlugin = api.NewAwsApiGatewayProvider(provider)
	membraneOpts.SecretManagerPlugin, _ = secrets_manager_secret_service.New(provider)
	membraneOpts.KeyValuePlugin, _ = dynamodb_service.New(provider)
	membraneOpts.TopicsPlugin, _ = sns_service.New(provider)
	membraneOpts.StoragePlugin, _ = s3_service.New(provider)
	membraneOpts.ResourcesPlugin = provider
	membraneOpts.WebsocketPlugin, _ = websocket.NewAwsApiGatewayWebsocket(provider)
	membraneOpts.QueuesPlugin, _ = sqs_service.New(provider)
	membraneOpts.SqlPlugin = sql_service.NewRdsSqlService()
	membraneOpts.BatchPlugin, _ = batch.New()

	m, err := membrane.New(membraneOpts)
	if err != nil {
		logger.Fatalf("There was an error initializing the membrane server: %v", err)
	}

	errChan := make(chan error)
	// Start the Membrane server
	go func(chan error) {
		errChan <- m.Start()
	}(errChan)

	select {
	case membraneError := <-errChan:
		logger.Errorf("Membrane Error: %v, exiting\n", membraneError)
	case sigTerm := <-term:
		logger.Infof("Received %v, exiting\n", sigTerm)
	}

	m.Stop()
}

func launchBatch(jobName string) {
	provider, err := resource.New()
	if err != nil {
		logger.Fatalf("could not create aws provider: %v", err)
		return
	}

	queueService, _ := sqs_service.New(provider)
	topicService, _ := sns_service.New(provider)
	storageService, _ := s3_service.New(provider)

	// launch a job membrane
	membrane := job.NewJobMembrane(
		strings.Join(os.Args[1:], " "),
		job.WithQueueServer(queueService),
		job.WithTopicServer(topicService),
		job.WithStorageServer(storageService),
		job.WithResourcesServer(provider),
	)

	err = membrane.Run()
	if err != nil {
		log.Fatalf("error executing nitric job %s: %s", jobName, err)
	}
}

func main() {
	jobName := commonenv.NITRIC_JOB_NAME.String()
	if jobName != "" {
		launchBatch(jobName)
	} else {
		launchService()
	}
}
