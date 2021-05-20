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
	pubsub_service "github.com/nitric-dev/membrane/plugins/eventing/pubsub"
	http_service "github.com/nitric-dev/membrane/plugins/gateway/cloudrun"
	firestore_service "github.com/nitric-dev/membrane/plugins/kv/firestore"
	pubsub_queue_service "github.com/nitric-dev/membrane/plugins/queue/pubsub"
	storage_service "github.com/nitric-dev/membrane/plugins/storage/storage"
	"github.com/nitric-dev/membrane/sdk"
)

type GCPServiceFactory struct {
}

func New() sdk.ServiceFactory {
	return &GCPServiceFactory{}
}

// NewKeyValueService - Returns Google Cloud Firestore based kv service
func (p *GCPServiceFactory) NewKeyValueService() (sdk.KeyValueService, error) {
	return firestore_service.New()
}

// NewEventService - Returns Google Cloud Pubsub based eventing service
func (p *GCPServiceFactory) NewEventService() (sdk.EventService, error) {
	return pubsub_service.New()
}

// NewGatewayService - Google Cloud Http Gateway service
func (p *GCPServiceFactory) NewGatewayService() (sdk.GatewayService, error) {
	return http_service.New()
}

// NewQueueService - Returns Google Cloud Pubsub based queue service
func (p *GCPServiceFactory) NewQueueService() (sdk.QueueService, error) {
	return pubsub_queue_service.New()
}

// NewStorageService - Returns Google Cloud Storage based storage service
func (p *GCPServiceFactory) NewStorageService() (sdk.StorageService, error) {
	return storage_service.New()
}
