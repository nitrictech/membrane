package lambda_service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	events "github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nitric-dev/membrane/handler"
	"github.com/nitric-dev/membrane/sdk"
	"github.com/nitric-dev/membrane/triggers"
)

type eventType int

const (
	unknown eventType = iota
	sns
	httpEvent
)

type LambdaRuntimeHandler func(handler interface{})

//Event incoming event
type Event struct {
	Requests []triggers.Trigger
}

func (event *Event) getEventType(data []byte) eventType {
	tmp := make(map[string]interface{})
	// Unmarshal so we can get just enough info about the type of event to fully deserialize it
	json.Unmarshal(data, &tmp)

	// If our event is a HTTP request
	if _, ok := tmp["rawPath"]; ok {
		return httpEvent
	} else if records, ok := tmp["Records"]; ok {
		recordsList, _ := records.([]interface{})
		record, _ := recordsList[0].(map[string]interface{})
		// We have some kind of event here...
		// we'll assume its an SNS
		var eventSource string
		if es, ok := record["EventSource"]; ok {
			eventSource = es.(string)
		} else if es, ok := record["eventSource"]; ok {
			eventSource = es.(string)
		}

		switch eventSource {
		case "aws:sns":
			return sns
		}
	}

	return unknown
}

// implement the unmarshal interface in order to handle multiple event types
func (event *Event) UnmarshalJSON(data []byte) error {
	var err error

	event.Requests = make([]triggers.Trigger, 0)

	switch event.getEventType(data) {
	case sns:
		snsEvent := &events.SNSEvent{}
		err = json.Unmarshal(data, snsEvent)

		if err == nil {
			// Map over the records and return
			for _, snsRecord := range snsEvent.Records {
				messageString := snsRecord.SNS.Message
				// FIXME: What about non-nitric SNS events???
				messageJson := &sdk.NitricEvent{}

				// Populate the JSON
				err = json.Unmarshal([]byte(messageString), messageJson)

				topicArn := snsRecord.SNS.TopicArn
				topicParts := strings.Split(topicArn, ":")
				trigger := topicParts[len(topicParts)-1]
				// get the topic name from the full ARN.
				// Get the topic name from the arn

				if err == nil {
					// Decode the message to see if it's a Nitric message
					payloadMap := messageJson.Payload
					payloadBytes, err := json.Marshal(&payloadMap)

					if err == nil {
						event.Requests = append(event.Requests, &triggers.Event{
							ID:      messageJson.ID,
							Topic:   trigger,
							Payload: payloadBytes,
						})
					}
				}
			}
		}
		break
	case httpEvent:
		evt := &events.APIGatewayV2HTTPRequest{}

		err = json.Unmarshal(data, evt)

		if err == nil {
			// Copy the headers and re-write for the proxy
			headerCopy := make(map[string]string)

			for key, val := range evt.Headers {
				if strings.ToLower(key) == "host" {
					headerCopy["x-forwarded-for"] = val
				} else {
					headerCopy[key] = val
				}
			}

			event.Requests = append(event.Requests, &triggers.HttpRequest{
				// FIXME: Translate to http.Header
				Header: headerCopy,
				Body:   []byte(evt.Body),
				Method: evt.RequestContext.HTTP.Method,
				Path:   evt.RawPath,
				Query:  evt.QueryStringParameters,
			})
		}

		break
	default:
		jsonEvent := make(map[string]interface{})

		err = json.Unmarshal(data, &jsonEvent)

		if err != nil {
			return err
		}

		err = fmt.Errorf("Unhandled Event Type: %v", data)
	}

	return err
}

type LambdaGateway struct {
	handler handler.TriggerHandler
	runtime LambdaRuntimeHandler
	sdk.UnimplementedGatewayPlugin
}

func (s *LambdaGateway) handle(ctx context.Context, event Event) (interface{}, error) {
	for _, request := range event.Requests {
		// TODO: Build up an array of responses?
		//in some cases we won't need to send a response as well...
		// resp := s.handler(&request)

		switch request.GetTriggerType() {
		case triggers.TriggerType_Request:
			if httpEvent, ok := request.(*triggers.HttpRequest); ok {
				response, err := s.handler.HandleHttpRequest(httpEvent)

				if err != nil {
					return events.APIGatewayProxyResponse{
						StatusCode: 500,
						Body:       "Error processing lambda request",
						// TODO: Need to determine best case when to use this...
						IsBase64Encoded: true,
					}, nil
				}

				lambdaHTTPHeaders := make(map[string]string)

				if response.Header != nil {
					response.Header.VisitAll(func(key []byte, val []byte) {
						lambdaHTTPHeaders[string(key)] = string(val)
					})
				}

				responseString := base64.StdEncoding.EncodeToString(response.Body)

				// We want to sniff the content type of the body that we have here as lambda cannot gzip it...
				return events.APIGatewayProxyResponse{
					StatusCode: response.StatusCode,
					Headers:    lambdaHTTPHeaders,
					Body:       responseString,
					// TODO: Need to determine best case when to use this...
					IsBase64Encoded: true,
				}, nil
			} else {
				return nil, fmt.Errorf("Error!: Found non HttpRequest in event with trigger type: %s", triggers.TriggerType_Request.String())
			}
			break
		case triggers.TriggerType_Subscription:
			if event, ok := request.(*triggers.Event); ok {
				if err := s.handler.HandleEvent(event); err != nil {
					return nil, err
				}
			} else {
				return nil, fmt.Errorf("Error!: Found non Event in event with trigger type: %s", triggers.TriggerType_Subscription.String())
			}
			break
		}
	}
	return nil, nil
}

// Start the lambda gateway handler
func (s *LambdaGateway) Start(handler handler.TriggerHandler) error {
	s.handler = handler
	// Here we want to begin polling lambda for incoming requests...
	// Assuming that this is blocking
	s.runtime(s.handle)

	return fmt.Errorf("Something went wrong causing the lambda runtime to stop")
}

func New() (sdk.GatewayService, error) {
	return &LambdaGateway{
		runtime: lambda.Start,
	}, nil
}

func NewWithRuntime(runtime LambdaRuntimeHandler) (sdk.GatewayService, error) {
	return &LambdaGateway{
		runtime: runtime,
	}, nil
}
