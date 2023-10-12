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

package schedule

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/nitrictech/nitric/cloud/gcp/deploy/exec"
	deploy "github.com/nitrictech/nitric/core/pkg/api/nitric/deploy/v1"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/cloudscheduler"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type CloudScheduler struct {
	pulumi.ResourceState

	Name string
	Job  *cloudscheduler.Job
}

type CloudSchedulerArgs struct {
	Location string

	Exec     *exec.CloudRunner
	Schedule *deploy.Schedule
	Tz       string
}

type ScheduleEvent struct {
	PayloadType string                 `yaml:"payloadType"`
	Payload     map[string]interface{} `yaml:"payload,omitempty"`
}

func NewCloudSchedulerJob(ctx *pulumi.Context, name string, args *CloudSchedulerArgs, opts ...pulumi.ResourceOption) (*CloudScheduler, error) {
	res := &CloudScheduler{
		Name: name,
	}

	err := ctx.RegisterComponentResource("nitric:schedule:GCPCloudSchedulerSchedule", name, res, opts...)
	if err != nil {
		return nil, err
	}

	eventJSON, err := json.Marshal(map[string]interface{}{
		"schedule": name,
	})
	if err != nil {
		return nil, err
	}

	payload := base64.StdEncoding.EncodeToString(eventJSON)

	res.Job, err = cloudscheduler.NewJob(ctx, name, &cloudscheduler.JobArgs{
		TimeZone: pulumi.String(args.Tz),
		HttpTarget: &cloudscheduler.JobHttpTargetArgs{
			Uri: pulumi.Sprintf("%s/x-nitric-schedule/%s?token=%s", args.Exec.Url, name, args.Exec.EventToken),
			OidcToken: &cloudscheduler.JobHttpTargetOidcTokenArgs{
				ServiceAccountEmail: args.Exec.Invoker.Email,
			},
			Body: pulumi.String(payload),
		},
		Schedule: pulumi.String(strings.ReplaceAll(args.Schedule.Cron, "'", "")),
	}, append(opts, pulumi.Parent(res))...)
	if err != nil {
		return nil, err
	}

	return res, err
}
