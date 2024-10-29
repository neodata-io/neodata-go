package tracing

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type Tracer struct {
	TracerProvider trace.TracerProvider
}

func NewTracer(serviceName string) (*Tracer, error) {
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}

	provider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			"service.name", serviceName,
		)),
	)
	otel.SetTracerProvider(provider)

	return &Tracer{TracerProvider: provider}, nil
}

func (n *NeoCtx) GetTracer(name string) trace.Tracer {
	return n.TracerProvider.Tracer(name)
}
