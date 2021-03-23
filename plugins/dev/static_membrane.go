package main

import (
	"log"
	"strconv"

	"github.com/nitric-dev/membrane/membrane"
	auth "github.com/nitric-dev/membrane/plugins/dev/auth"
	eventing "github.com/nitric-dev/membrane/plugins/dev/eventing"
	gateway "github.com/nitric-dev/membrane/plugins/dev/gateway"
	kv "github.com/nitric-dev/membrane/plugins/dev/kv"
	queue "github.com/nitric-dev/membrane/plugins/dev/queue"
	storage "github.com/nitric-dev/membrane/plugins/dev/storage"
	"github.com/nitric-dev/membrane/utils"
)

func main() {
	serviceAddress := utils.GetEnv("SERVICE_ADDRESS", "127.0.0.1:50051")
	childAddress := utils.GetEnv("CHILD_ADDRESS", "127.0.0.1:8080")
	childCommand := utils.GetEnv("INVOKE", "")
	tolerateMissingServices := utils.GetEnv("TOLERATE_MISSING_SERVICES", "false")

	tolerateMissing, err := strconv.ParseBool(tolerateMissingServices)

	if err != nil {
		log.Fatalf("There was an error initialising the membrane server: %v", err)
	}

	eventingPlugin, _ := eventing.New()
	kvPlugin, _ := kv.New()
	storagePlugin, _ := storage.New()
	gatewayPlugin, _ := gateway.New()
	queuePlugin, _ := queue.New()
	authPlugin, _ := auth.New()

	membrane, err := membrane.New(&membrane.MembraneOptions{
		ServiceAddress:          serviceAddress,
		ChildAddress:            childAddress,
		ChildCommand:            childCommand,
		EventingPlugin:          eventingPlugin,
		KvPlugin:                kvPlugin,
		StoragePlugin:           storagePlugin,
		QueuePlugin:             queuePlugin,
		AuthPlugin:              authPlugin,
		GatewayPlugin:           gatewayPlugin,
		TolerateMissingServices: tolerateMissing,
	})

	if err != nil {
		log.Fatalf("There was an error initialising the membrane server: %v", err)
	}

	// Start the Membrane server
	membrane.Start()
}
