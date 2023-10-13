// Copyright Nitric Pty Ltd.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package events

import (
	"fmt"

	"github.com/nitrictech/nitric/cloud/common/deploy/resources"

	common "github.com/nitrictech/nitric/cloud/common/deploy/tags"
	"github.com/nitrictech/nitric/cloud/gcp/deploy/exec"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/deploy/v1"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/pubsub"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type PubSubTopic struct {
	pulumi.ResourceState

	Name   string
	PubSub *pubsub.Topic
}

type PubSubTopicArgs struct {
	Location  string
	StackID   string
	ProjectId string

	Topic *v1.Topic
}

func NewPubSubTopic(ctx *pulumi.Context, name string, args *PubSubTopicArgs, opts ...pulumi.ResourceOption) (*PubSubTopic, error) {
	res := &PubSubTopic{
		Name: name,
	}

	err := ctx.RegisterComponentResource("nitric:topic:GCPPubSubTopic", name, res, opts...)
	if err != nil {
		return nil, err
	}

	res.PubSub, err = pubsub.NewTopic(ctx, name, &pubsub.TopicArgs{
		Labels: pulumi.ToStringMap(common.Tags(args.StackID, name, resources.Topic)),
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

type PubSubSubscription struct {
	pulumi.ResourceState

	Name string

	Subscription *pubsub.Subscription
}

type PubSubSubscriptionArgs struct {
	Function *exec.CloudRunner
	Topic    *PubSubTopic
}

func GetSubName(executionName string, topicName string) string {
	return fmt.Sprintf("%s-%s-sub", executionName, topicName)
}

func NewPubSubPushSubscription(ctx *pulumi.Context, name string, args *PubSubSubscriptionArgs, opts ...pulumi.ResourceOption) (*PubSubSubscription, error) {
	res := &PubSubSubscription{
		Name: name,
	}

	err := ctx.RegisterComponentResource("nitric:topic:GCPPubSubTopicSubscription", name, res, opts...)
	if err != nil {
		return nil, err
	}

	s, err := pubsub.NewSubscription(ctx, name, &pubsub.SubscriptionArgs{
		Topic:              args.Topic.PubSub.Name, // The GCP topic name
		AckDeadlineSeconds: pulumi.Int(300),
		RetryPolicy: pubsub.SubscriptionRetryPolicyArgs{
			MinimumBackoff: pulumi.String("15s"),
			MaximumBackoff: pulumi.String("600s"),
		},
		PushConfig: pubsub.SubscriptionPushConfigArgs{
			OidcToken: pubsub.SubscriptionPushConfigOidcTokenArgs{
				ServiceAccountEmail: args.Function.Invoker.Email,
			},
			// https://cloud.google.com/appengine/docs/flexible/writing-and-responding-to-pub-sub-messages?tab=go#top
			PushEndpoint: pulumi.Sprintf("%s/x-nitric-topic/%s?token=%s", args.Function.Url, args.Topic.Name, args.Function.EventToken),
		},
		ExpirationPolicy: &pubsub.SubscriptionExpirationPolicyArgs{
			Ttl: pulumi.String(""),
		},
	}, append(opts, pulumi.Parent(args.Function))...)
	if err != nil {
		return nil, errors.WithMessage(err, "subscription "+name+"-sub")
	}

	res.Subscription = s

	return res, nil
}
