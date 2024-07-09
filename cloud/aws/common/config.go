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

package common

import (
	"github.com/imdario/mergo"
	"github.com/mitchellh/mapstructure"
	"github.com/nitrictech/nitric/cloud/common/deploy/config"
)

type ApiConfig struct {
	Domains []string
}

type AwsImports struct {
	// A map of nitric names to ARNs
	Secrets map[string]string
}

type BatchComputeEnvConfig struct {
	MinCpus       int      `mapstructure:"min-cpus"`
	MaxCpus       int      `mapstructure:"max-cpus"`
	InstanceTypes []string `mapstructure:"instance-types"`
}

type AwsJobsConfig struct {
	Cpus   int
	Memory int
	Gpus   int
}

type AwsConfig struct {
	ScheduleTimezone                      string `mapstructure:"schedule-timezone,omitempty"`
	Import                                AwsImports
	Refresh                               bool
	Apis                                  map[string]*ApiConfig
	BatchComputeEnvConfig                 *BatchComputeEnvConfig `mapstructure:"batch-compute-env,omitempty"`
	Jobs                                  map[string]*AwsJobsConfig
	config.AbstractConfig[*AwsConfigItem] `mapstructure:"config,squash"`
}

type AwsConfigItem struct {
	Lambda    *AwsLambdaConfig `mapstructure:",omitempty"`
	Telemetry int
}

type AwsLambdaVpcConfig struct {
	SubnetIds        []string `mapstructure:"subnet-ids"`
	SecurityGroupIds []string `mapstructure:"security-group-ids"`
}

type AwsLambdaConfig struct {
	Memory                int
	Timeout               int
	ProvisionedConcurreny int                 `mapstructure:"provisioned-concurrency"`
	Vpc                   *AwsLambdaVpcConfig `mapstructure:"vpc,omitempty"`
}

var defaultLambdaConfig = &AwsLambdaConfig{
	Memory:                128,
	Timeout:               15,
	ProvisionedConcurreny: 0,
}

var defaultJobConfig = &AwsJobsConfig{
	Memory: 512,
	Cpus:   1,
	Gpus:   0,
}

var defaultBatchComputeEnvConfig = &BatchComputeEnvConfig{
	MinCpus:       0,
	MaxCpus:       32,
	InstanceTypes: []string{"optimal"},
}

var defaultAwsConfigItem = AwsConfigItem{
	Telemetry: 0,
}

// Return AwsConfig from stack attributes
func ConfigFromAttributes(attributes map[string]interface{}) (*AwsConfig, error) {
	// get config attributes
	err := config.ValidateRawConfigKeys(attributes, []string{"lambda"})
	if err != nil {
		return nil, err
	}

	awsConfig := &AwsConfig{}
	err = mapstructure.Decode(attributes, awsConfig)
	if err != nil {
		return nil, err
	}

	// Default timezone if not specified
	if awsConfig.ScheduleTimezone == "" {
		// default to UTC
		awsConfig.ScheduleTimezone = "UTC"
	}

	if awsConfig.Apis == nil {
		awsConfig.Apis = map[string]*ApiConfig{}
	}

	if awsConfig.Config == nil {
		awsConfig.Config = map[string]*AwsConfigItem{}
	}

	// if no default then set provider level defaults
	if _, hasDefault := awsConfig.Config["default"]; !hasDefault {
		awsConfig.Config["default"] = &defaultAwsConfigItem
		awsConfig.Config["default"].Lambda = defaultLambdaConfig
	}

	for configName, configVal := range awsConfig.Config {
		// Add omitted values from default configs where needed.
		err := mergo.Merge(configVal, defaultAwsConfigItem)
		if err != nil {
			return nil, err
		}

		if configVal.Lambda == nil { // check if no runtime config provided, default to Lambda.
			configVal.Lambda = defaultLambdaConfig
		} else {
			err := mergo.Merge(configVal.Lambda, defaultLambdaConfig)
			if err != nil {
				return nil, err
			}
		}

		awsConfig.Config[configName] = configVal
	}

	if awsConfig.BatchComputeEnvConfig == nil {
		awsConfig.BatchComputeEnvConfig = defaultBatchComputeEnvConfig
	}

	// merge in default values
	err = mergo.Merge(awsConfig.BatchComputeEnvConfig, defaultBatchComputeEnvConfig)
	if err != nil {
		return nil, err
	}

	// Default job config
	if awsConfig.Jobs == nil {
		awsConfig.Jobs = map[string]*AwsJobsConfig{}
	}

	if _, hasDefault := awsConfig.Jobs["default"]; !hasDefault {
		awsConfig.Jobs["default"] = defaultJobConfig
	}

	for configName, configVal := range awsConfig.Jobs {
		// Add omitted values from default configs where needed.
		err := mergo.Merge(configVal, defaultJobConfig)
		if err != nil {
			return nil, err
		}

		awsConfig.Jobs[configName] = configVal
	}

	return awsConfig, nil
}
