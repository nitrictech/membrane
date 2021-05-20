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
	"fmt"
	"log"
	"os"
	"plugin"
	"strconv"
	"strings"

	"github.com/nitric-dev/membrane/membrane"
	"github.com/nitric-dev/membrane/plugins/sdk"
	"github.com/nitric-dev/membrane/utils"
)

// Pluggable version of the Nitric membrane
func main() {
	serviceAddress := utils.GetEnv("SERVICE_ADDRESS", "127.0.0.1:50051")
	childAddress := utils.GetEnv("CHILD_ADDRESS", "127.0.0.1:8080")
	pluginDir := utils.GetEnv("PLUGIN_DIR", "./plugins")
	serviceFactoryPluginFile := utils.GetEnv("SERVICE_FACTORY_PLUGIN", "default.so")

	var childCommand []string
	// Get the command line arguments, minus the program name in index 0.
	if len(os.Args) > 1 && len(os.Args[1:]) > 0 {
		childCommand = os.Args[1:]
	} else {
		childCommand = strings.Fields(utils.GetEnv("INVOKE", ""))
		if len(childCommand) > 0 {
			fmt.Println("Warning: use of INVOKE environment variable is deprecated and may be removed in a future version")
		}
	}

	tolerateMissingServices := utils.GetEnv("TOLERATE_MISSING_SERVICES", "false")

	tolerateMissing, err := strconv.ParseBool(tolerateMissingServices)
	// Set tolerate missing to false by default so missing plugins will cause a fatal error for safety.
	if err != nil {
		log.Println(fmt.Sprintf("failed to parse TOLERATE_MISSING_SERVICES environment variable with value [%s], defaulting to false", tolerateMissingServices))
		tolerateMissing = false
	}
	var serviceFactory sdk.ServiceFactory = nil

	// Load the Plugin Factory
	if plug, err := plugin.Open(fmt.Sprintf("%s/%s", pluginDir, serviceFactoryPluginFile)); err == nil {
		if symbol, err := plug.Lookup("New"); err == nil {
			if newFunc, ok := symbol.(func() (sdk.ServiceFactory, error)); ok {
				if serviceFactoryPlugin, err := newFunc(); err == nil {
					serviceFactory = serviceFactoryPlugin
				}
			}
		}
	}
	if serviceFactory == nil {
		log.Fatalf("failed to load Provider Factory Plugin: %s", serviceFactoryPluginFile)
	}

	// Load the concrete service implementations
	var eventingService sdk.EventService = nil
	var gatewayService sdk.GatewayService = nil
	var keyValueService sdk.KeyValueService = nil
	var queueService sdk.QueueService = nil
	var storageService sdk.StorageService = nil

	// Load the eventing service
	if eventingService, err = serviceFactory.NewEventService(); err != nil {
		log.Fatal(err)
	}
	// Load the gateway service
	if gatewayService, err = serviceFactory.NewGatewayService(); err != nil {
		log.Fatal(err)
	}
	// Load the key value service
	if keyValueService, err = serviceFactory.NewKeyValueService(); err != nil {
		log.Fatal(err)
	}
	// Load the queue service
	if queueService, err = serviceFactory.NewQueueService(); err != nil {
		log.Fatal(err)
	}
	// Load the storage service
	if storageService, err = serviceFactory.NewStorageService(); err != nil {
		log.Fatal(err)
	}

	// Construct and validate the membrane server
	membraneServer, err := membrane.New(&membrane.MembraneOptions{
		ServiceAddress:          serviceAddress,
		ChildAddress:            childAddress,
		ChildCommand:            childCommand,
		EventingPlugin:          eventingService,
		KvPlugin:                keyValueService,
		StoragePlugin:           storageService,
		GatewayPlugin:           gatewayService,
		QueuePlugin:             queueService,
		TolerateMissingServices: tolerateMissing,
	})

	if err != nil {
		log.Fatalf("There was an error initialising the membraneServer server: %v", err)
	}

	// Start the Membrane server
	membraneServer.Start()
}
