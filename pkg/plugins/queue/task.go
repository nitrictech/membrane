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

// FailedTask - A task that has failed to be queued
type FailedTask struct {
	Task    *NitricTask
	Message string
}

// NitricTask - A task for asynchronous processing
type NitricTask struct {
	ID          string                 `json:"id,omitempty" log:"ID"`
	LeaseID     string                 `json:"leaseId,omitempty" log:"LeaseID"`
	PayloadType string                 `json:"payloadType,omitempty" log:"PayLoadType"`
	Payload     map[string]interface{} `json:"payload,omitempty"`
}
