package trace

import (
	"go-5m3Micro/pkg/log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

const (
	jaegerName = "jaeger"
	zipkinName = "zipkin"
)

var agents map[string]struct{}

func InitAgent(o *Options) error {
	// 检查键是否存在
	_, exists := agents[o.Endpoint]
	if exists {
		return nil
	} else {
		return startAgent(o)
	}
}

func startAgent(o *Options) error {
	opts := []trace.TracerProviderOption{
		trace.WithSampler(trace.ParentBased(trace.TraceIDRatioBased(float64(o.Sampler)))),
		trace.WithResource(resource.NewSchemaless(attribute.String("service.name", o.Name))),
	}
	var exporter trace.SpanExporter
	var err error
	if len(o.Batcher) > 0 {
		switch o.Batcher {
		case jaegerName:
			exporter, err = jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(o.Endpoint)))
			if err != nil {
				return err
			}
		case zipkinName:
			exporter, err = zipkin.New(o.Endpoint)
			if err != nil {
				return err
			}
		}
		opts = append(opts, trace.WithBatcher(exporter))
	}

	// 创建 TracerProvider
	tp := trace.NewTracerProvider(opts...)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		log.Errorf("Error happened during trace handling: %v", err)
	}))
	agents[o.Endpoint] = struct{}{}
	return nil
}
