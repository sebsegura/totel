package utel

import (
	"context"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"log"
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

func InitTracer(ctx context.Context, cfg *Config) *sdktrace.TracerProvider {
	res, err := NewResource(ctx, cfg)
	if err != nil {
		log.Fatalf("cannot start trace provider: %v", err)
	}

	exporter, err := newTraceExporter(ctx)
	if err != nil {
		log.Fatalf("cannot set trace exporter: %v", err)
	}

	idg := xray.NewIDGenerator()

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithIDGenerator(idg),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(xray.Propagator{})

	return tp
}
