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
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc/codes"

	pb "github.com/nitrictech/nitric/pkg/api/nitric/v1"
	"github.com/nitrictech/nitric/pkg/plugins/queue"
	"github.com/nitrictech/nitric/pkg/span"
	"github.com/nitrictech/protoutils"
)

// GRPC Interface for registered Nitric Storage Plugins
type QueueServiceServer struct {
	pb.UnimplementedQueueServiceServer
	plugin queue.QueueService
}

func (s *QueueServiceServer) checkPluginRegistered() error {
	if s.plugin == nil {
		return NewPluginNotRegisteredError("Queue")
	}

	return nil
}

func (s *QueueServiceServer) Send(ctx context.Context, req *pb.QueueSendRequest) (*pb.QueueSendResponse, error) {
	if err := s.checkPluginRegistered(); err != nil {
		return nil, err
	}

	if err := req.ValidateAll(); err != nil {
		return nil, newGrpcErrorWithCode(codes.InvalidArgument, "QueueService.Send", err)
	}

	task := req.GetTask()

	// auto generate an ID if we did not receive one
	ID := task.GetId()
	if ID == "" {
		ID = uuid.New().String()
	}

	nitricTask := queue.NitricTask{
		ID:          ID,
		PayloadType: task.GetPayloadType(),
		Payload:     task.GetPayload().AsMap(),
	}

	sp := span.FromHeaders(ctx, "queue-"+req.GetQueue(), map[string][]string{})
	sp.SetAttributes(
		semconv.CodeFunctionKey.String("Queue.Send"),
		semconv.MessagingDestinationKindQueue,
		semconv.MessagingDestinationKey.String(req.GetQueue()),
		semconv.MessagingMessageIDKey.String(nitricTask.ID),
		attribute.Key("messaging.message_type").String(nitricTask.PayloadType),
	)

	defer sp.End()

	if err := s.plugin.Send(req.GetQueue(), nitricTask); err != nil {
		sp.RecordError(err)

		return nil, err
	}

	// Success
	return &pb.QueueSendResponse{}, nil
}

func (s *QueueServiceServer) SendBatch(ctx context.Context, req *pb.QueueSendBatchRequest) (*pb.QueueSendBatchResponse, error) {
	if err := s.checkPluginRegistered(); err != nil {
		return nil, err
	}

	if err := req.ValidateAll(); err != nil {
		return nil, newGrpcErrorWithCode(codes.InvalidArgument, "QueueService.SendBatch", err)
	}

	// Translate tasks
	tasks := make([]queue.NitricTask, len(req.GetTasks()))
	for i, task := range req.GetTasks() {
		// auto generate an ID if we did not receive one
		ID := task.GetId()
		if ID == "" {
			ID = uuid.New().String()
		}

		tasks[i] = queue.NitricTask{
			ID:          ID,
			PayloadType: task.GetPayloadType(),
			Payload:     task.GetPayload().AsMap(),
		}
	}

	if resp, err := s.plugin.SendBatch(req.GetQueue(), tasks); err == nil {
		failedTasks := make([]*pb.FailedTask, len(resp.FailedTasks))
		for i, failedTask := range resp.FailedTasks {
			st, _ := protoutils.NewStruct(failedTask.Task.Payload)
			failedTasks[i] = &pb.FailedTask{
				Message: failedTask.Message,
				Task: &pb.NitricTask{
					Id:          failedTask.Task.ID,
					PayloadType: failedTask.Task.PayloadType,
					Payload:     st,
				},
			}
		}
		return &pb.QueueSendBatchResponse{
			FailedTasks: failedTasks,
		}, nil
	} else {
		return nil, NewGrpcError("QueueService.SendBatch", err)
	}
}

func (s *QueueServiceServer) Receive(ctx context.Context, req *pb.QueueReceiveRequest) (*pb.QueueReceiveResponse, error) {
	if err := s.checkPluginRegistered(); err != nil {
		return nil, err
	}

	if err := req.ValidateAll(); err != nil {
		return nil, newGrpcErrorWithCode(codes.InvalidArgument, "QueueService.Receive", err)
	}

	// Convert gRPC request to plugin params
	depth := uint32(req.GetDepth())
	popOptions := queue.ReceiveOptions{
		QueueName: req.GetQueue(),
		Depth:     &depth,
	}

	// Perform the Queue Receive operation
	tasks, err := s.plugin.Receive(popOptions)
	if err != nil {
		return nil, NewGrpcError("QueueService.Receive", err)
	}

	// Convert the NitricTasks to the gRPC type
	grpcTasks := make([]*pb.NitricTask, 0, len(tasks))
	for _, task := range tasks {
		st, _ := protoutils.NewStruct(task.Payload)
		grpcTasks = append(grpcTasks, &pb.NitricTask{
			Id:          task.ID,
			Payload:     st,
			LeaseId:     task.LeaseID,
			PayloadType: task.PayloadType,
		})
	}

	// Return the tasks
	res := pb.QueueReceiveResponse{
		Tasks: grpcTasks,
	}
	return &res, nil
}

func (s *QueueServiceServer) Complete(ctx context.Context, req *pb.QueueCompleteRequest) (*pb.QueueCompleteResponse, error) {
	if err := s.checkPluginRegistered(); err != nil {
		return nil, err
	}

	if err := req.ValidateAll(); err != nil {
		return nil, newGrpcErrorWithCode(codes.InvalidArgument, "QueueService.Complete", err)
	}
	// Convert gRPC request to plugin params
	queueName := req.GetQueue()
	leaseId := req.GetLeaseId()

	// Perform the Queue Complete operation
	err := s.plugin.Complete(queueName, leaseId)
	if err != nil {
		return nil, NewGrpcError("QueueService.Complete", err)
	}

	// Return a successful response
	return &pb.QueueCompleteResponse{}, nil
}

func NewQueueServiceServer(plugin queue.QueueService) pb.QueueServiceServer {
	return &QueueServiceServer{
		plugin: plugin,
	}
}
