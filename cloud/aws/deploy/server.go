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

package deploy

import (
	_ "embed"

	"github.com/nitrictech/nitric/cloud/common/deploy/pulumi"
	deploy "github.com/nitrictech/nitric/core/pkg/api/nitric/deploy/v1"
)

type DeployServer struct {
	deploy.UnimplementedDeployServiceServer
}

// Embeds the runtime directly into the deploytime binary
// This way the versions will always match as they're always built and versioned together (as a single artifact)
// This should also help with docker build speeds as the runtime has already been "downloaded"
//
//go:embed runtime-aws
var runtime []byte

func NewServer() (*DeployServer, error) {
	err := pulumi.InstallResources()
	if err != nil {
		return nil, err
	}

	return &DeployServer{}, nil
}
