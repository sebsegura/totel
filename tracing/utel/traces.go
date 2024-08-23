package utel

import (
	"context"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"os"
)

func url() string {
	u := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if u == "" {
		u = "localhost:4317"
	}
	return u
}

func newTraceExporter(ctx context.Context) (*otlptrace.Exporter, error) {
	return otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(url()),
		otlptracegrpc.WithInsecure())
}

func InitTracer(ctx context.Context, res *resource.Resource) (*sdktrace.TracerProvider, error) {
	exporter, err := newTraceExporter(ctx)
	if err != nil {
		return nil, err
	}

	idg := xray.NewIDGenerator()

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithIDGenerator(idg),
	)
	otel.SetTextMapPropagator(xray.Propagator{})

	return tp, nil
}

func NewServerSpan(ctx context.Context, spanName string) (context.Context, trace.Span) {
	tracer := otel.Tracer(GetUtelConfig().ServiceName)
	return tracer.Start(ctx, spanName, trace.WithSpanKind(trace.SpanKindServer))
}

func NewClientSpan(ctx context.Context, spanName string) (context.Context, trace.Span) {
	tracer := otel.Tracer(GetUtelConfig().ServiceName)
	return tracer.Start(ctx, spanName, trace.WithSpanKind(trace.SpanKindClient))
}
