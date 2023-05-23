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

package grpc

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
	"github.com/nitrictech/nitric/core/pkg/worker"
	"github.com/nitrictech/nitric/core/pkg/worker/adapter"
	"github.com/nitrictech/nitric/core/pkg/worker/pool"
)

type FaasServer struct {
	pb.UnimplementedFaasServiceServer
	pool pool.WorkerPool
}

// Starts a new stream
// A reference to this stream will be passed on to a new worker instance
// This represents a new server that is ready to begin processing
func (s *FaasServer) TriggerStream(stream pb.FaasService_TriggerStreamServer) error {
	cm, err := stream.Recv()
	if err != nil {
		return status.Errorf(codes.Internal, "error reading message from stream: %v", err)
	}

	ir := cm.GetInitRequest()

	if ir == nil {
		// SHUT IT DOWN!!!!
		// The first message must be an init request from the prospective FaaS worker
		return status.Error(codes.FailedPrecondition, "first message must be InitRequest")
	}

	var wrkr worker.Worker
	adapter := adapter.NewGrpcAdapter(stream)

	if api := ir.GetApi(); api != nil {
		// Create a new route worker
		wrkr = worker.NewRouteWorker(adapter, &worker.RouteWorkerOptions{
			Api:     api.Api,
			Path:    api.Path,
			Methods: api.Methods,
		})
	} else if subscription := ir.GetSubscription(); subscription != nil {
		wrkr = worker.NewSubscriptionWorker(adapter, &worker.SubscriptionWorkerOptions{
			Topic: subscription.Topic,
		})
	} else if schedule := ir.GetSchedule(); schedule != nil {
		wrkr = worker.NewScheduleWorker(adapter, &worker.ScheduleWorkerOptions{
			Key: schedule.Key,
		})
	} else if notification := ir.GetBucketNotification(); notification != nil {
		wrkr = worker.NewBucketNotificationWorker(adapter, &worker.BucketNotificationWorkerOptions{
			Notification: notification,
		})
	} else {
		return fmt.Errorf("unable to register incompatible worker. This can be caused by out of date Nitric CLI versions, an upgrade may resolve this issue")
	}

	// Add it to our new pool
	if err := s.pool.AddWorker(wrkr); err != nil {
		// Worker could not be added
		// Cancel the stream by returning an error
		// This should cause the spawned child process to exit
		return err
	}

	// We're good to go
	errchan := make(chan error)

	// Start the worker
	go adapter.Start(errchan)

	// block here on error returned from the worker
	err = <-errchan
	log.Default().Println("FaaS stream closed, removing worker")

	// Worker is done so we can remove it from the pool
	rwErr := s.pool.RemoveWorker(wrkr)
	if rwErr != nil {
		if err != nil {
			return errors.Wrap(err, rwErr.Error())
		}

		return rwErr
	}
	return err
}

func NewFaasServer(workerPool pool.WorkerPool) *FaasServer {
	return &FaasServer{
		pool: workerPool,
	}
}
