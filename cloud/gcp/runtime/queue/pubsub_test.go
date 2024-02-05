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

package queue_test

import (
	// "encoding/json"
	// "fmt"

	. "github.com/onsi/ginkgo"
	// . "github.com/onsi/gomega"
	// "golang.org/x/net/context"
	// "google.golang.org/api/option"
	// "github.com/nitrictech/nitric/core/pkg/plugins/queue"
	// ifaces_pubsub "github.com/nitrictech/nitric/cloud/gcp/ifaces/pubsub"
	// mock_pubsub "github.com/nitrictech/nitric/cloud/gcp/mocks/pubsub"
	// pubsub_queue_service "github.com/nitrictech/nitric/cloud/gcp/runtime/queue"
)

var _ = Describe("Pubsub", func() {
	Context("Send", func() {
		// When("Publishing to a queue that exists", func() {
		// 	mockPubsubClient := mock_pubsub.NewMockPubsubClient(
		// 		mock_pubsub.MockPubsubOptions{
		// 			Topics: []string{"test"},
		// 		})
		// 	queuePlugin := pubsub_queue_service.NewWithClient(mockPubsubClient)

		// 	It("Should queue the Nitric Task", func() {
		// 		err := queuePlugin.Send(context.TODO(), "test", queue.NitricTask{
		// 			ID:          "1234",
		// 			PayloadType: "test-payload",
		// 			Payload: map[string]interface{}{
		// 				"Test": "Test",
		// 			},
		// 		})

		// 		By("Not returning an error")
		// 		Expect(err).ShouldNot(HaveOccurred())

		// 		By("The queue containing the published message")
		// 		Expect(mockPubsubClient.PublishedMessages["test"]).To(HaveLen(1))
		// 	})
		// })

		// 	When("Publishing to a queue that does not exist", func() {
		// 		mockPubsubClient := mock_pubsub.NewMockPubsubClient(mock_pubsub.MockPubsubOptions{
		// 			Topics: []string{},
		// 		})
		// 		queuePlugin := pubsub_queue_service.NewWithClient(mockPubsubClient)
		// 		It("Should return the error that failed to publish", func() {
		// 			err := queuePlugin.Send(context.TODO(), "test", queue.NitricTask{
		// 				ID:          "mockrequestid",
		// 				PayloadType: "mockpayloadtype",
		// 				LeaseID:     "MockId",
		// 				Payload: map[string]interface{}{
		// 					"Test": "Test",
		// 				},
		// 			})

		// 			// It should still attempt
		// 			By("Returning an error")
		// 			Expect(err).Should(HaveOccurred())
		// 		})
		// 	})
		// })

		// Context("BatchPush", func() {
		// 	When("Publishing to a queue that exists", func() {
		// 		mockPubsubClient := mock_pubsub.NewMockPubsubClient(
		// 			mock_pubsub.MockPubsubOptions{
		// 				Topics: []string{"test"},
		// 			})
		// 		queuePlugin := pubsub_queue_service.NewWithClient(mockPubsubClient)

		// 		It("Should queue the Nitric Task", func() {
		// 			resp, err := queuePlugin.SendBatch(context.TODO(), "test", []queue.NitricTask{{
		// 				ID:          "1234",
		// 				PayloadType: "test-payload",
		// 				Payload: map[string]interface{}{
		// 					"Test": "Test",
		// 				},
		// 			}})

		// 			By("Not returning an error")
		// 			Expect(err).ShouldNot(HaveOccurred())

		// 			By("Returning no failed messages")
		// 			Expect(resp.FailedTasks).To(HaveLen(0))

		// 			By("The queue containing the published message")
		// 			Expect(mockPubsubClient.PublishedMessages["test"]).To(HaveLen(1))
		// 		})
		// 	})

		// 	When("Publishing to a queue that does not exist", func() {
		// 		mockPubsubClient := mock_pubsub.NewMockPubsubClient(mock_pubsub.MockPubsubOptions{
		// 			Topics: []string{},
		// 		})
		// 		queuePlugin := pubsub_queue_service.NewWithClient(mockPubsubClient)
		// 		It("Should return the messages that failed to publish", func() {
		// 			_, err := queuePlugin.SendBatch(context.TODO(), "test", []queue.NitricTask{{
		// 				ID:          "1234",
		// 				PayloadType: "test-payload",
		// 				Payload: map[string]interface{}{
		// 					"Test": "Test",
		// 				},
		// 			}})

		// 			// It should still attempt
		// 			By("Returning an error")
		// 			Expect(err).Should(HaveOccurred())
		// 		})
		// 	})
		// })

		// Context("Receive", func() {
		// 	When("Popping from a queue that exists", func() {
		// 		When("There is a message on the queue", func() {
		// 			mockId := "mockmessageid"
		// 			mockReceiptHandle := "mockreceipthandle"
		// 			jsonBytes, _ := json.Marshal(queue.NitricTask{
		// 				ID:          "mockrequestid",
		// 				PayloadType: "mockpayloadtype",
		// 				Payload:     map[string]interface{}{},
		// 			})

		// 			var mockMessage ifaces_pubsub.Message = mock_pubsub.MockPubsubMessage{
		// 				Id:        mockId,
		// 				AckId:     mockReceiptHandle,
		// 				DataBytes: jsonBytes,
		// 			}

		// 			mockPubsubClient := mock_pubsub.NewMockPubsubClient(mock_pubsub.MockPubsubOptions{
		// 				Topics: []string{"mock-queue"},
		// 				Messages: map[string][]ifaces_pubsub.Message{
		// 					"mock-queue": {
		// 						mockMessage,
		// 					},
		// 				},
		// 			})
		// 			queuePlugin := pubsub_queue_service.NewWithClients(mockPubsubClient, func(ctx context.Context, opts ...option.ClientOption) (ifaces_pubsub.SubscriberClient, error) {
		// 				return mock_pubsub.MockBaseClient{
		// 					Messages: map[string][]ifaces_pubsub.Message{
		// 						"mock-queue": {
		// 							mockMessage,
		// 						},
		// 					},
		// 				}, nil
		// 			})

		// 			It("Should receive the message", func() {
		// 				items, err := queuePlugin.Receive(context.TODO(), queue.ReceiveOptions{
		// 					QueueName: "mock-queue",
		// 					Depth:     nil,
		// 				})

		// 				By("Not returning an error")
		// 				Expect(err).ShouldNot(HaveOccurred())

		// 				By("Returning no failed messages")
		// 				Expect(len(items)).To(Equal(1))
		// 			})
		// 		})
		// 	})

		When("Publishing to a queue that does not exist", func() {
			//mockPubsubClient := mocks.NewMockPubsubClient([]string{})
			//queuePlugin := pubsub_queue_plugin.NewWithClient(mockPubsubClient)
			//It("Should return the messages that failed to publish", func() {
			//	_, err := queuePlugin.BatchPush("test", []queue.NitricTask{{
			//		ID:   "1234",
			//		PayloadType: "test-payload",
			//		Payload: map[string]interface{}{
			//			"Test": "Test",
			//		},
			//	}})
			//
			//	// It should still attempt
			//	By("Returning an error")
			//	Expect(err).Should(HaveOccurred())
			//})
		})
	})

	// Context("Complete", func() {
	// 	When("Pubsub acknowledge request succeeds", func() {
	// 		mockPubsubClient := mock_pubsub.NewMockPubsubClient(mock_pubsub.MockPubsubOptions{
	// 			Topics: []string{"mock-queue"},
	// 		})
	// 		queuePlugin := pubsub_queue_service.NewWithClients(mockPubsubClient, func(ctx context.Context, opts ...option.ClientOption) (ifaces_pubsub.SubscriberClient, error) {
	// 			return mock_pubsub.MockBaseClient{}, nil
	// 		})

	// 		It("Should not return an error", func() {
	// 			err := queuePlugin.Complete(context.TODO(), "mock-queue", "test-id")

	// 			By("Not returning an error")
	// 			Expect(err).ShouldNot(HaveOccurred())
	// 		})
	// 	})
	// 	When("Pubsub acknowledge request errors", func() {
	// 		mockPubsubClient := mock_pubsub.NewMockPubsubClient(mock_pubsub.MockPubsubOptions{
	// 			Topics: []string{"mock-queue"},
	// 		})
	// 		queuePlugin := pubsub_queue_service.NewWithClients(mockPubsubClient, func(ctx context.Context, opts ...option.ClientOption) (ifaces_pubsub.SubscriberClient, error) {
	// 			return mock_pubsub.MockBaseClient{
	// 				CompleteError: fmt.Errorf("mock complete error"),
	// 			}, nil
	// 		})

	// 		It("Should return an error", func() {
	// 			err := queuePlugin.Complete(context.TODO(), "mock-queue", "test-id")

	// 			By("Not returning an error")
	// 			Expect(err).Should(HaveOccurred())
	// 		})
	// 	})
	// })
})
