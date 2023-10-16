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

package gateway

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/nitrictech/nitric/cloud/aws/runtime/core"
	"github.com/nitrictech/nitric/cloud/common/deploy/tags"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
	"github.com/nitrictech/nitric/core/pkg/plugins/gateway"
	"github.com/nitrictech/nitric/core/pkg/utils"
	"github.com/nitrictech/nitric/core/pkg/worker"
	"github.com/nitrictech/nitric/core/pkg/worker/pool"
)

type LambdaRuntimeHandler func(interface{})

func (s *LambdaGateway) getTopicNameForArn(ctx context.Context, topicArn string) (string, error) {
	topics, err := s.provider.GetResources(ctx, core.AwsResource_Topic)
	if err != nil {
		return "", fmt.Errorf("error retrieving topics: %w", err)
	}

	for name, arn := range topics {
		if arn == topicArn {
			return name, nil
		}
	}

	return "", fmt.Errorf("could not find topic for arn %s", topicArn)
}

func (s *LambdaGateway) getBucketNameForArn(ctx context.Context, bucketArn string) (string, error) {
	buckets, err := s.provider.GetResources(ctx, core.AwsResource_Bucket)
	if err != nil {
		return "", fmt.Errorf("error retrieving topics: %w", err)
	}

	for name, arn := range buckets {
		if arn == bucketArn {
			return name, nil
		}
	}

	return "", fmt.Errorf("could not find topic for arn %s", bucketArn)
}

type LambdaGateway struct {
	pool     pool.WorkerPool
	provider core.AwsProvider
	runtime  LambdaRuntimeHandler
	gateway.UnimplementedGatewayPlugin
	finished chan int
}

// Handle websocket events
func (s *LambdaGateway) handleWebsocketEvent(ctx context.Context, evt events.APIGatewayWebsocketProxyRequest) (interface{}, error) {
	// Use the routekey to get the event type

	wsEvent := v1.WebsocketEvent_Message
	switch evt.RequestContext.RouteKey {
	case "$connect":
		wsEvent = v1.WebsocketEvent_Connect
	case "$disconnect":
		wsEvent = v1.WebsocketEvent_Disconnect
	}

	api, err := s.provider.GetApiGatewayById(ctx, evt.RequestContext.APIID)
	if err != nil {
		return nil, err
	}

	stackID := utils.GetEnv("NITRIC_STACK_ID", "")
	nitricName, ok := api.Tags[tags.GetResourceNameKey(stackID)]
	if !ok {
		return nil, fmt.Errorf("received websocket trigger from non-nitric API gateway")
	}

	queryParams := map[string]*v1.QueryValue{}
	for k, v := range evt.QueryStringParameters {
		queryParams[k] = &v1.QueryValue{
			Value: []string{v},
		}
	}

	req := &v1.TriggerRequest{
		Data: []byte(evt.Body),
		Context: &v1.TriggerRequest_Websocket{
			Websocket: &v1.WebsocketTriggerContext{
				ConnectionId: evt.RequestContext.ConnectionID,
				Event:        wsEvent,
				// Get the API gateways nitric name
				Socket:      nitricName,
				QueryParams: queryParams,
			},
		},
	}

	wrk, err := s.pool.GetWorker(&pool.GetWorkerOptions{
		Trigger: req,
	})
	if err != nil {
		return nil, err
	}

	_, err = wrk.HandleTrigger(ctx, req)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error processing lambda request",
			// TODO: Need to determine best case when to use this...
			IsBase64Encoded: false,
		}, nil
	}

	// if response.GetWebsocket() == nil || !response.GetWebsocket().Success {
	// 	return events.APIGatewayProxyResponse{
	// 		StatusCode: 500,
	// 	}, nil
	// }

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}

// Handle API events
func (s *LambdaGateway) handleApiEvent(ctx context.Context, evt events.APIGatewayV2HTTPRequest) (interface{}, error) {
	// Copy the headers and re-write for the proxy
	headerCopy := map[string]*v1.HeaderValue{}

	for key, val := range evt.Headers {
		if strings.ToLower(key) == "host" {
			headerCopy[xforwardHeader] = &v1.HeaderValue{
				Value: []string{val},
			}
		} else {
			if headerCopy[key] == nil {
				headerCopy[key] = &v1.HeaderValue{}
			}
			headerCopy[key].Value = append(headerCopy[key].Value, val)
		}
	}

	// Copy the cookies over
	headerCopy["Cookie"] = &v1.HeaderValue{
		Value: evt.Cookies,
	}

	// Parse the raw query string
	qVals, err := url.ParseQuery(evt.RawQueryString)
	if err != nil {
		return nil, fmt.Errorf("error parsing query for httpEvent: %w", err)
	}
	query := map[string]*v1.QueryValue{}
	for k, v := range qVals {
		query[k] = &v1.QueryValue{
			Value: v,
		}
	}

	data := []byte(evt.Body)
	if evt.IsBase64Encoded {
		data, err = base64.StdEncoding.DecodeString(evt.Body)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode:      400,
				Body:            "Error processing lambda request",
				IsBase64Encoded: false,
			}, nil
		}
	}

	req := &v1.TriggerRequest{
		Data: data,
		Context: &v1.TriggerRequest_Http{
			Http: &v1.HttpTriggerContext{
				Method:      evt.RequestContext.HTTP.Method,
				Path:        evt.RawPath,
				QueryParams: query,
				Headers:     headerCopy,
			},
		},
	}

	wrk, err := s.pool.GetWorker(&pool.GetWorkerOptions{
		Trigger: req,
	})
	if err != nil {
		return nil, err
	}

	response, err := wrk.HandleTrigger(ctx, req)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error processing lambda request",
			// TODO: Need to determine best case when to use this...
			IsBase64Encoded: false,
		}, nil
	}

	lambdaHTTPHeaders := make(map[string]string)
	if response.GetHttp().Headers != nil {
		for k, v := range response.GetHttp().Headers {
			lambdaHTTPHeaders[k] = v.Value[0]
		}
	}

	responseString := base64.StdEncoding.EncodeToString(response.Data)

	return events.APIGatewayProxyResponse{
		StatusCode:      int(response.GetHttp().Status),
		Headers:         lambdaHTTPHeaders,
		Body:            responseString,
		IsBase64Encoded: true,
	}, nil
}

type ScheduleMessage struct {
	Schedule string
}

func (s *LambdaGateway) handleScheduleEvent(ctx context.Context, evt nitricScheduleEvent) (interface{}, error) {
	if evt.Schedule == "" {
		return nil, fmt.Errorf("unable to identify source nitric schedule")
	}

	request := &v1.TriggerRequest{
		// Send empty data for now (no reason to send data for schedules at the moment)
		Data: nil,
		Context: &v1.TriggerRequest_Topic{
			Topic: &v1.TopicTriggerContext{
				Topic: worker.ScheduleKeyToTopicName(evt.Schedule),
			},
		},
	}

	wrkr, err := s.pool.GetWorker(&pool.GetWorkerOptions{
		Trigger: request,
		// Only send Cloudwatch events to schedule workers
		Filter: func(w worker.Worker) bool {
			_, ok := w.(*worker.ScheduleWorker)
			return ok
		},
	})
	if err != nil {
		return nil, fmt.Errorf("no worker available to handle schedule %s", evt.Schedule)
	}

	resp, err := wrkr.HandleTrigger(context.TODO(), request)
	if err != nil {
		return nil, err
	}

	if !resp.GetTopic().Success {
		return nil, fmt.Errorf("schedule execution failed")
	}

	return nil, nil
}

func (s *LambdaGateway) handleSnsEvents(ctx context.Context, records []Record) (interface{}, error) {
	for _, snsRecord := range records {
		messageString := snsRecord.SNS.Message
		// var id string
		attrs := map[string]string{}

		for k, v := range snsRecord.SNS.MessageAttributes {
			sv, ok := v.(string)
			if ok {
				attrs[k] = sv
			}
		}

		tName, err := s.getTopicNameForArn(ctx, snsRecord.SNS.TopicArn)
		if err != nil {
			log.Default().Printf("unable to find nitric topic: %v", err)
			continue
		}

		request := &v1.TriggerRequest{
			Data: []byte(messageString),
			Context: &v1.TriggerRequest_Topic{
				Topic: &v1.TopicTriggerContext{
					Topic: tName,
				},
			},
		}

		wrkr, err := s.pool.GetWorker(&pool.GetWorkerOptions{
			Trigger: request,
			// Only send SNS events to subscription workers
			Filter: func(w worker.Worker) bool {
				_, ok := w.(*worker.SubscriptionWorker)
				return ok
			},
		})
		if err != nil {
			return nil, fmt.Errorf("unable to get worker to event trigger")
		}

		var mc propagation.MapCarrier = attrs

		trigResp, err := wrkr.HandleTrigger(xray.Propagator{}.Extract(ctx, mc), request)
		if err != nil {
			return nil, err
		}

		if trigResp.GetTopic() == nil || !trigResp.GetTopic().Success {
			return nil, fmt.Errorf("event handler return non success")
		}
	}

	return nil, nil
}

func (s *LambdaGateway) handleHealthCheck(ctx context.Context, evt healthCheckEvent) (interface{}, error) {
	return map[string]interface{}{
		"healthy": true,
	}, nil
}

// Converts the GCP event type to our abstract event type
func notificationEventToEventType(eventType string) (v1.BucketNotificationType, error) {
	if ok := strings.Contains(eventType, "ObjectCreated:"); ok {
		return v1.BucketNotificationType_Created, nil
	} else if ok := strings.Contains(eventType, "ObjectRemoved:"); ok {
		return v1.BucketNotificationType_Deleted, nil
	}
	return v1.BucketNotificationType_All, fmt.Errorf("unsupported bucket notification event type %s", eventType)
}

func (s *LambdaGateway) handleS3Event(ctx context.Context, records []Record) (interface{}, error) {
	for _, s3Record := range records {
		bucketName, err := s.getBucketNameForArn(ctx, s3Record.EventSourceArn)
		if err != nil {
			log.Default().Println("unable to locate nitric bucket")
			return nil, fmt.Errorf("unable to find nitric bucket: %w", err)
		}

		eventType, err := notificationEventToEventType(s3Record.EventName)
		if err != nil {
			return nil, err
		}

		request := &v1.TriggerRequest{
			Context: &v1.TriggerRequest_Notification{
				Notification: &v1.NotificationTriggerContext{
					Source: bucketName,
					Notification: &v1.NotificationTriggerContext_Bucket{
						Bucket: &v1.BucketNotification{
							Key:  s3Record.S3.Object.Key,
							Type: eventType,
						},
					},
				},
			},
		}

		wrkr, err := s.pool.GetWorker(&pool.GetWorkerOptions{
			Trigger: request,
			// Only send S3 events to bucket notification workers
			Filter: func(w worker.Worker) bool {
				_, ok := w.(*worker.BucketNotificationWorker)
				return ok
			},
		})
		if err != nil {
			return nil, fmt.Errorf("unable to get worker to event trigger")
		}

		var mc propagation.MapCarrier = s3Record.ResponseElements

		_, err = wrkr.HandleTrigger(xray.Propagator{}.Extract(ctx, mc), request)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (s *LambdaGateway) routeEvent(ctx context.Context, evt Event) (interface{}, error) {
	switch evt.Type() {
	case websocketEvent:
		return s.handleWebsocketEvent(ctx, evt.APIGatewayWebsocketProxyRequest)
	case httpEvent:
		return s.handleApiEvent(ctx, evt.APIGatewayV2HTTPRequest)
	case healthcheck:
		return s.handleHealthCheck(ctx, evt.healthCheckEvent)
	case sns:
		return s.handleSnsEvents(ctx, evt.Records)
	case s3:
		return s.handleS3Event(ctx, evt.Records)
	case schedule:
		return s.handleScheduleEvent(ctx, evt.nitricScheduleEvent)
	default:
		return nil, fmt.Errorf("unhandled lambda event type: %+v", evt)
	}
}

// Start the lambda gateway handler
func (s *LambdaGateway) Start(pool pool.WorkerPool) error {
	s.pool = pool
	// Here we want to begin polling lambda for incoming requests...
	s.runtime(func(ctx context.Context, evt Event) (interface{}, error) {
		a, err := s.routeEvent(ctx, evt)

		tp, ok := otel.GetTracerProvider().(*sdktrace.TracerProvider)
		if ok {
			_ = tp.ForceFlush(ctx)
		}

		return a, err
	})
	// Unblock the 'Stop' function if it's waiting.
	go func() { s.finished <- 1 }()
	return nil
}

func (s *LambdaGateway) Stop() error {
	// XXX: This is a NO_OP Process, as this is a pull based system
	// We don't need to stop listening to anything
	log.Default().Println("gateway 'Stop' called, waiting for lambda runtime to finish")
	// Lambda can't be stopped, need to wait for it to finish
	<-s.finished
	return nil
}

func New(provider core.AwsProvider) (gateway.GatewayService, error) {
	return &LambdaGateway{
		provider: provider,
		runtime:  lambda.Start,
		finished: make(chan int),
	}, nil
}

func NewWithRuntime(provider core.AwsProvider, runtime LambdaRuntimeHandler) (gateway.GatewayService, error) {
	return &LambdaGateway{
		provider: provider,
		runtime:  runtime,
		finished: make(chan int),
	}, nil
}
