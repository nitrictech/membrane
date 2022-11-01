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

package worker

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"github.com/nitrictech/nitric/pkg/triggers"
)

// A Nitric HTTP worker
type HttpWorker struct {
	address string
}

func (s *HttpWorker) HandlesHttpRequest(trigger *triggers.HttpRequest) bool {
	return true
}

func (s *HttpWorker) HandlesEvent(trigger *triggers.Event) bool {
	return true
}

// HandleEvent - Handles an event from a subscription by converting it to an HTTP request.
func (h *HttpWorker) HandleEvent(ctx context.Context, trigger *triggers.Event) error {
	address := fmt.Sprintf("http://%s/subscriptions/%s", h.address, trigger.Topic)

	httpRequest := fasthttp.AcquireRequest()
	httpRequest.SetRequestURI(address)
	httpRequest.Header.Add("x-nitric-request-id", trigger.ID)
	httpRequest.Header.Add("x-nitric-source-type", triggers.TriggerType_Subscription.String())
	httpRequest.Header.Add("x-nitric-source", trigger.Topic)

	var resp fasthttp.Response

	httpRequest.SetBody(trigger.Payload)
	httpRequest.Header.SetContentLength(len(trigger.Payload))

	// TODO: Handle response or error and respond appropriately
	err := fasthttp.Do(httpRequest, &resp)
	if err == nil && resp.StatusCode() >= 200 && resp.StatusCode() <= 299 {
		return nil
	}
	if err != nil {
		return errors.Wrapf(err, "Error processing event (%d): %s", resp.StatusCode(), string(resp.Body()))
	}
	return errors.Errorf("Error processing event (%d): %s", resp.StatusCode(), string(resp.Body()))
}

// HandleHttpRequest - Handles an HTTP request by forwarding it as an HTTP request.
func (h *HttpWorker) HandleHttpRequest(ctx context.Context, trigger *triggers.HttpRequest) (*triggers.HttpResponse, error) {
	address := fmt.Sprintf("http://%s%s", h.address, trigger.Path)

	httpRequest := fasthttp.AcquireRequest()
	httpRequest.SetRequestURI(address)

	for key, val := range trigger.Query {
		for _, v := range val {
			httpRequest.URI().QueryArgs().Add(key, v)
		}
	}

	for key, val := range trigger.Header {
		for _, v := range val {
			httpRequest.Header.Add(key, v)
		}
	}

	httpRequest.Header.Del("Content-Length")
	httpRequest.SetBody(trigger.Body)
	httpRequest.Header.SetContentLength(len(trigger.Body))

	var resp fasthttp.Response
	err := fasthttp.Do(httpRequest, &resp)
	if err != nil {
		return nil, err
	}

	return triggers.FromHttpResponse(&resp), nil
}

// Creates a new HttpWorker
// Will wait to ensure that the provided address is dialable
// before proceeding
func NewHttpWorker(address string) (*HttpWorker, error) {
	// Dial the child port to see if it's open and ready...
	maxWaitTime := time.Duration(5) * time.Second
	// Longer poll times, e.g. 200 milliseconds results in slow lambda cold starts (15s+)
	pollInterval := time.Duration(15) * time.Millisecond

	waitedTime := time.Duration(0)
	for {
		conn, _ := net.Dial("tcp", address)
		if conn != nil {
			conn.Close()
			break
		} else {
			if waitedTime < maxWaitTime {
				time.Sleep(pollInterval)
				waitedTime += pollInterval
			} else {
				return nil, fmt.Errorf("Unable to dial http worker, does it expose a http server at: %s?", address)
			}
		}
	}

	// Dial the provided address to ensure its availability
	return &HttpWorker{
		address: address,
	}, nil
}
