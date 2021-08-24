package azqueue_service_iface

import (
	"context"
	"time"

	"github.com/Azure/azure-storage-queue-go/azqueue"
)

func AdaptServiceUrl(c azqueue.ServiceURL) AzqueueServiceUrlIface {
	return serviceUrl{c}
}

func AdaptQueueUrl(c azqueue.QueueURL) AzqueueQueueUrlIface {
	return queueUrl{c}
}

func AdaptMessageUrl(c azqueue.MessagesURL) AzqueueMessageUrlIface {
	return messageUrl{c}
}

func AdaptMessageIdUrl(c azqueue.MessageIDURL) AzqueueMessageIdUrlIface {
	return messageIdUrl{c}
}

type (
	serviceUrl   struct{ c azqueue.ServiceURL }
	queueUrl     struct{ c azqueue.QueueURL }
	messageUrl   struct{ c azqueue.MessagesURL }
	messageIdUrl struct{ c azqueue.MessageIDURL }
)

func (c serviceUrl) NewQueueURL(queueName string) AzqueueQueueUrlIface {
	return AdaptQueueUrl(c.c.NewQueueURL(queueName))
}

func (c queueUrl) NewMessageURL() AzqueueMessageUrlIface {
	return AdaptMessageUrl(c.c.NewMessagesURL())
}

func (c messageUrl) Enqueue(ctx context.Context, messageText string, visibilityTimeout time.Duration, timeToLive time.Duration) (*azqueue.EnqueueMessageResponse, error) {
	return c.c.Enqueue(ctx, messageText, visibilityTimeout, timeToLive)
}

func (c messageUrl) Dequeue(ctx context.Context, maxMessages int32, visibilityTimeout time.Duration) (*azqueue.DequeuedMessagesResponse, error) {
	return c.c.Dequeue(ctx, maxMessages, visibilityTimeout)
}

func (c messageUrl) NewMessageIDURL(messageId azqueue.MessageID) AzqueueMessageIdUrlIface {
	return AdaptMessageIdUrl(c.c.NewMessageIDURL(messageId))
}

func (c messageIdUrl) Delete(ctx context.Context, popReceipt azqueue.PopReceipt) (*azqueue.MessageIDDeleteResponse, error) {
	return c.c.Delete(ctx, popReceipt)
}
