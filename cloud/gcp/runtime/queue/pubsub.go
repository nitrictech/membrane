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

package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nitrictech/nitric/cloud/common/deploy/resources"
	"github.com/nitrictech/nitric/cloud/common/deploy/tags"
	"github.com/nitrictech/nitric/cloud/gcp/runtime/env"
	grpccodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cloud.google.com/go/pubsub"
	pubsubbase "cloud.google.com/go/pubsub/apiv1"
	pubsubpb "cloud.google.com/go/pubsub/apiv1/pubsubpb"
	"github.com/GoogleCloudPlatform/opentelemetry-operations-go/propagator"
	"go.opentelemetry.io/otel/propagation"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	ifaces_pubsub "github.com/nitrictech/nitric/cloud/gcp/ifaces/pubsub"
	"github.com/nitrictech/nitric/core/pkg/plugins/errors"
	"github.com/nitrictech/nitric/core/pkg/plugins/errors/codes"
	"github.com/nitrictech/nitric/core/pkg/plugins/queue"
)

type PubsubQueueService struct {
	queue.UnimplementedQueuePlugin
	client              ifaces_pubsub.PubsubClient
	newSubscriberClient func(ctx context.Context, opts ...option.ClientOption) (ifaces_pubsub.SubscriberClient, error)
	projectId           string
	cache               map[string]ifaces_pubsub.Topic
}

// Retrieves the Nitric "Queue Topic" for the specified queue (PubSub Topic).
//
// This retrieves the default Nitric Queue for the Topic based on tagging conventions.
func (s *PubsubQueueService) getPubsubTopicFromName(queue string) (ifaces_pubsub.Topic, error) {
	if s.cache == nil {
		topics := s.client.Topics(context.Background())
		s.cache = make(map[string]ifaces_pubsub.Topic)
		stackID := env.GetNitricStackID()
		for {
			t, err := topics.Next()
			if errors.Is(err, iterator.Done) {
				break
			}
			if err != nil {
				return nil, fmt.Errorf("an error occurred finding queue: %s; %w", queue, err)
			}

			labels, err := t.Labels(context.TODO())
			if err != nil {
				return nil, fmt.Errorf("an error occurred finding queue labels: %s; %w", queue, err)
			}

			resType, hasType := labels[tags.GetResourceTypeKey(stackID)]

			if name, ok := labels[tags.GetResourceNameKey(stackID)]; ok && name == queue && hasType && resType == "queue" {
				s.cache[name] = t
			}
		}
	}

	if t, ok := s.cache[queue]; ok {
		return t, nil
	}

	return nil, fmt.Errorf("queue not found")
}

// Retrieves the Nitric "Queue Subscription" for the specified queue (PubSub Topic).
//
// GCP PubSub requires a Subscription in order to Pull messages from a Topic.
// we use this behavior to emulate a queue.
//
// This retrieves the default Nitric Pull subscription for the Topic base on convention.
func (s *PubsubQueueService) getQueueSubscription(ctx context.Context, queueName string) (ifaces_pubsub.Subscription, error) {
	// We'll be using pubsub with pull subscribers to facilitate queue functionality
	topic, err := s.getPubsubTopicFromName(queueName)
	if err != nil {
		return nil, err
	}

	subsIt := topic.Subscriptions(ctx)

	for {
		sub, err := subsIt.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve pull subscription for topic: %s\n%w", topic.ID(), err)
		}

		labels, err := sub.Labels(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve pull subscription labels for topic: %s\n%w", topic.ID(), err)
		}

		resourceType, hasType := labels[tags.GetResourceTypeKey(env.GetNitricStackID())]
		if name, ok := labels[tags.GetResourceNameKey(env.GetNitricStackID())]; hasType && ok && resourceType == string(resources.Queue) {
			if name == queueName {
				return sub, nil
			}
		}
	}

	return nil, fmt.Errorf("pull subscription not found, pull subscribers may not be configured for this topic")
}

func (s *PubsubQueueService) Send(ctx context.Context, queue string, task queue.NitricTask) error {
	newErr := errors.ErrorsWithScope(
		"PubsubQueueService.Send",
		map[string]interface{}{
			"queue": queue,
			"task":  task,
		},
	)

	// We'll be using pubsub with pull subscribers to facilitate queue functionality
	topic, err := s.getPubsubTopicFromName(queue)
	if err != nil {
		return newErr(
			codes.NotFound,
			"queue not found",
			err,
		)
	}

	if taskBytes, err := json.Marshal(task); err == nil {
		attributes := propagation.MapCarrier{}

		propagator.CloudTraceFormatPropagator{}.Inject(ctx, attributes)

		msg := ifaces_pubsub.AdaptPubsubMessage(&pubsub.Message{
			Attributes: attributes,
			Data:       taskBytes,
		})

		result := topic.Publish(ctx, msg)

		if _, err := result.Get(ctx); err != nil {
			return newErr(
				codes.Internal,
				"error retrieving publish result",
				err,
			)
		}
	} else {
		return newErr(
			codes.Internal,
			"error marshalling the task",
			err,
		)
	}

	return nil
}

func (s *PubsubQueueService) SendBatch(ctx context.Context, q string, tasks []queue.NitricTask) (*queue.SendBatchResponse, error) {
	newErr := errors.ErrorsWithScope(
		"PubsubQueueService.SendBatch",
		map[string]interface{}{
			"queue":     q,
			"tasks.len": len(tasks),
		},
	)

	// We'll be using pubsub with pull subscribers to facilitate queue functionality
	topic, err := s.getPubsubTopicFromName(q)
	if err != nil {
		return nil, newErr(
			codes.NotFound,
			"queue not found",
			err,
		)
	}

	// SendBatch once we've published all tasks to the client
	// TODO: We may want to revisit this, and chunk up our publishing in a way that makes more sense...
	results := make([]ifaces_pubsub.PublishResult, 0)
	failedTasks := make([]*queue.FailedTask, 0)
	publishedTasks := make([]queue.NitricTask, 0)

	attributes := propagation.MapCarrier{}

	propagator.CloudTraceFormatPropagator{}.Inject(ctx, attributes)

	for _, task := range tasks {
		t := task
		if taskBytes, err := json.Marshal(t); err == nil {
			msg := ifaces_pubsub.AdaptPubsubMessage(&pubsub.Message{
				Data:       taskBytes,
				Attributes: attributes,
			})

			results = append(results, topic.Publish(ctx, msg))
			publishedTasks = append(publishedTasks, t)
		} else {
			failedTasks = append(failedTasks, &queue.FailedTask{
				Task:    &t,
				Message: "Error unmarshalling message for queue",
			})
		}
	}

	for idx, result := range results {
		// Iterate over the results to check for successful publishing...
		if _, err := result.Get(ctx); err != nil {
			// Add this to our failures list in our results...
			failedTasks = append(failedTasks, &queue.FailedTask{
				Task:    &publishedTasks[idx],
				Message: err.Error(),
			})
		}
	}

	return &queue.SendBatchResponse{
		FailedTasks: failedTasks,
	}, nil
}

// Receives a collection of tasks off a given queue.
func (s *PubsubQueueService) Receive(ctx context.Context, options queue.ReceiveOptions) ([]queue.NitricTask, error) {
	newErr := errors.ErrorsWithScope(
		"PubsubQueueService.Receive",
		map[string]interface{}{
			"options": options,
		},
	)

	if err := options.Validate(); err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"invalid receive options provided",
			err,
		)
	}

	// Find the generic pull subscription for the provided topic (queue)
	queueSubscription, err := s.getQueueSubscription(ctx, options.QueueName)
	if err != nil {
		return nil, newErr(
			codes.NotFound,
			"could not find queue subscription",
			err,
		)
	}

	// Using base client, so that asynchronous message acknowledgement can take place without needing to keep messages
	// in a stateful service. Standard PubSub go library doesn't provide access to the 'acknowledge' ID of the messages
	// or an independent acknowledge function. It's only provided as a method on message objects.
	client, err := s.newSubscriberClient(ctx)
	if err != nil {
		return nil, newErr(
			codes.Internal,
			"failed to create subscriber client",
			err,
		)
	}
	defer client.Close()

	// Retrieve the requested number of messages from the subscription (queue)
	req := pubsubpb.PullRequest{
		Subscription: queueSubscription.String(),
		MaxMessages:  int32(*options.Depth),
	}
	res, err := client.Pull(ctx, &req)
	if err != nil {
		errStatus, _ := status.FromError(err)
		if errStatus.Code() == grpccodes.PermissionDenied {
			return nil, newErr(
				codes.PermissionDenied,
				"permission denied, have you requested access to this queue?", err)
		}

		return nil, newErr(
			codes.Internal,
			"failed to pull messages",
			err,
		)
	}

	// An empty list is returned from PubSub if no messages are available
	// we return our own empty list in turn.
	if len(res.ReceivedMessages) == 0 {
		return []queue.NitricTask{}, nil
	}

	// Convert the PubSub messages into Nitric tasks
	var tasks []queue.NitricTask
	for _, m := range res.ReceivedMessages {
		var nitricTask queue.NitricTask
		err := json.Unmarshal(m.Message.Data, &nitricTask)
		if err != nil {
			// TODO: append error to error list and Nack the message.
			continue
		}

		tasks = append(tasks, queue.NitricTask{
			ID:          nitricTask.ID,
			Payload:     nitricTask.Payload,
			PayloadType: nitricTask.PayloadType,
			LeaseID:     m.AckId,
		})
	}

	return tasks, nil
}

// Completes a previously popped queue item
func (s *PubsubQueueService) Complete(ctx context.Context, queue string, leaseId string) error {
	newErr := errors.ErrorsWithScope(
		"PubsubQueueService.Complete",
		map[string]interface{}{
			"queue":   queue,
			"leaseId": leaseId,
		},
	)

	// Find the generic pull subscription for the provided topic (queue)
	queueSubscription, err := s.getQueueSubscription(ctx, queue)
	if err != nil {
		return newErr(
			codes.NotFound,
			"could not find queue subscription",
			err,
		)
	}

	// Using base client, so that asynchronous message acknowledgement can take place without needing to keep messages
	// in a stateful service. Standard PubSub go library is stateful and don't provide access to the acknowledge ID of
	// the messages or an independent acknowledge function. It's only provided as a method on message objects.
	client, err := s.newSubscriberClient(ctx)
	if err != nil {
		return newErr(
			codes.Internal,
			"failed to create subscriberclient",
			err,
		)
	}
	defer client.Close()

	// Acknowledge the queue item, so it's removed from the queue
	req := pubsubpb.AcknowledgeRequest{
		Subscription: queueSubscription.String(),
		AckIds:       []string{leaseId},
	}
	err = client.Acknowledge(ctx, &req)
	if err != nil {
		errStatus, _ := status.FromError(err)
		if errStatus.Code() == grpccodes.PermissionDenied {
			return newErr(
				codes.PermissionDenied,
				"permission denied, have you requested access to the queue?",
				err)
		}

		return newErr(
			codes.Internal,
			"failed to de-queue task",
			err,
		)
	}

	return nil
}

// adaptNewClient - Adapts the pubsubbase.NewSubscriberClient func to one that implements the SubscriberClient
// interface. This is used to enable substitution of the base pubsub client, primarily for mocking support.
func adaptNewClient(f func(context.Context, ...option.ClientOption) (*pubsubbase.SubscriberClient, error)) func(ctx context.Context, opts ...option.ClientOption) (ifaces_pubsub.SubscriberClient, error) {
	return func(c context.Context, opts ...option.ClientOption) (ifaces_pubsub.SubscriberClient, error) {
		return f(c, opts...)
	}
}

// New - Constructs a new GCP pubsub client with defaults
func New() (queue.QueueService, error) {
	ctx := context.Background()

	credentials, credentialsError := google.FindDefaultCredentials(ctx, pubsub.ScopeCloudPlatform)
	if credentialsError != nil {
		return nil, fmt.Errorf("GCP credentials error: %w", credentialsError)
	}
	client, clientError := pubsub.NewClient(ctx, credentials.ProjectID)
	if clientError != nil {
		return nil, fmt.Errorf("pubsub client error: %w", clientError)
	}

	return &PubsubQueueService{
		client: ifaces_pubsub.AdaptPubsubClient(client),
		// TODO: replace this with a better mechanism for mocking the client.
		newSubscriberClient: adaptNewClient(pubsubbase.NewSubscriberClient),
		projectId:           credentials.ProjectID,
	}, nil
}

func NewWithClient(client ifaces_pubsub.PubsubClient) queue.QueueService {
	return &PubsubQueueService{
		client:              client,
		newSubscriberClient: nil,
	}
}

// *pubsubbase.SubscriberClient
func NewWithClients(client ifaces_pubsub.PubsubClient, subscriberClientGenerator func(ctx context.Context, opts ...option.ClientOption) (ifaces_pubsub.SubscriberClient, error)) queue.QueueService {
	return &PubsubQueueService{
		client:              client,
		newSubscriberClient: subscriberClientGenerator,
	}
}
