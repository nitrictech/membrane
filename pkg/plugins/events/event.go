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
package events

import "fmt"

// NitricEvent - An event for asynchronous processing and reactive programming
type NitricEvent struct {
	ID          string                 `json:"id,omitempty"`
	PayloadType string                 `json:"payloadType,omitempty"`
	Payload     map[string]interface{} `json:"payload,omitempty"`
}

func (e *NitricEvent) String() string {
	return fmt.Sprintf("{ID: %s, PayloadType: %s}", e.ID, e.PayloadType)
}
