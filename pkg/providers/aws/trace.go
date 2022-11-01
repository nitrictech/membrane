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

package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/pkg/errors"
	lambdadetector "go.opentelemetry.io/contrib/detectors/aws/lambda"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/nitrictech/nitric/pkg/span"
	"github.com/nitrictech/nitric/pkg/utils"
)

func newTracerProvider(ctx context.Context) (*sdktrace.TracerProvider, error) {
	span.FunctionName = lambdacontext.FunctionName
	span.UseFuncNameAsSpanName = true

	exp, err := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx,
		resource.WithDetectors(lambdadetector.NewResourceDetector()),
		resource.WithTelemetrySDK(),
		resource.WithAttributes(
			semconv.CloudProviderAWS,
			semconv.CloudPlatformAWSLambda,
			semconv.ServiceNameKey.String(span.FunctionName),
			semconv.ServiceNamespaceKey.String(utils.GetEnv("NITRIC_STACK", "")),
		),
	)
	if err != nil {
		return nil, err
	}

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			xray.Propagator{},
			propagation.TraceContext{},
			propagation.Baggage{},
		))

	rate, err := utils.PercentFromIntString(utils.GetEnv("NITRIC_TRACE_SAMPLE_PERCENT", "10"))
	if err != nil {
		return nil, errors.WithMessagef(err, "NITRIC_TRACE_SAMPLE_PERCENT should be an int not %s", rate)
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(rate))),
		sdktrace.WithBatcher(exp),
		sdktrace.WithIDGenerator(xray.NewIDGenerator()),
		sdktrace.WithResource(res),
	), nil
}
