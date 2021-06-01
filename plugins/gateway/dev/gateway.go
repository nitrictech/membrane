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

// The AWS HTTP gateway plugin
package gateway_plugin

import (
	"strings"

	"github.com/nitric-dev/membrane/handler"
	"github.com/nitric-dev/membrane/sdk"
	"github.com/nitric-dev/membrane/triggers"
	"github.com/nitric-dev/membrane/utils"
	"github.com/valyala/fasthttp"
)

type HttpGateway struct {
	address string
	server  *fasthttp.Server
	sdk.UnimplementedGatewayPlugin
}

// TODO: Lets bind this to a struct...
func httpHandler(handler handler.TriggerHandler) func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		var triggerTypeString = string(ctx.Request.Header.Peek("x-nitric-source-type"))

		// Handle Event/Subscription Request Types
		if strings.ToUpper(triggerTypeString) == triggers.TriggerType_Subscription.String() {
			trigger := string(ctx.Request.Header.Peek("x-nitric-source"))
			requestId := string(ctx.Request.Header.Peek("x-nitric-request-id"))
			payload := ctx.Request.Body()

			err := handler.HandleEvent(&triggers.Event{
				ID:      requestId,
				Topic:   trigger,
				Payload: payload,
			})

			if err != nil {
				// TODO: Make this more informative
				ctx.Error("There was an error processing the event", 500)
			} else {
				ctx.SuccessString("text/plain", "Successfully Handled the Event")
			}

			// return here...
			return
		}

		httpReq := triggers.FromHttpRequest(ctx)
		// Handle HTTP Request Types
		response, err := handler.HandleHttpRequest(httpReq)

		if err != nil {
			// TODO: Redact message in production
			ctx.Error(err.Error(), 500)
			return
		}

		if response.Header != nil {
			response.Header.CopyTo(&ctx.Response.Header)
		}

		ctx.Response.Header.Del("Content-Length")
		ctx.Response.Header.Del("Connection")

		ctx.Response.SetBody(response.Body)
		ctx.Response.SetStatusCode(response.StatusCode)
	}
}

func (s *HttpGateway) Start(handler handler.TriggerHandler) error {
	// Start the fasthttp server
	s.server = &fasthttp.Server{
		Handler: httpHandler(handler),
	}

	return s.server.ListenAndServe(s.address)
}

func (s *HttpGateway) Stop() error {
	if s.server != nil {
		return s.server.Shutdown()
	}
	return nil
}

// Create new HTTP gateway
// XXX: No External Args for function atm (currently the plugin loader does not pass any argument information)
func New() (sdk.GatewayService, error) {
	address := utils.GetEnv("GATEWAY_ADDRESS", ":9001")

	return &HttpGateway{
		address: address,
	}, nil
}
