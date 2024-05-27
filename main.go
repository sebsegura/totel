package main

import (
	"context"
	"github.com/Bancar/goala/ulog"
	"github.com/Bancar/lambda-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	met "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"log"
	"time"
)

type Event struct{}

// const _ip = "34.207.131.19:4317"
const _ip = "localhost:4317"

// const _ip = "otel-collector.otel-collector.local:4317"

func Do(ctx context.Context, _ *Event) error {
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("falopa")))
	if err != nil {
		return err
	}

	metricExporter, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(_ip),
	)
	if err != nil {
		ulog.With(ulog.Str("error", err.Error())).Error("connection failed")
		return err
	}

	mp := metric.NewMeterProvider(
		metric.WithResource(r),
		metric.WithReader(metric.NewPeriodicReader(metricExporter, metric.WithInterval(1*time.Second))))
	defer func() {
		if err := mp.Shutdown(ctx); err != nil {
			log.Fatalf("error shutting down meter provider: %v", err)
		}
	}()

	otel.SetMeterProvider(mp)

	meter := otel.Meter("falopa")
	counter, err := meter.Int64Counter(
		"falopa.counter",
		met.WithDescription("OTEL test"),
		met.WithUnit("{count}"),
	)
	if err != nil {
		ulog.With(ulog.Str("error", err.Error())).Error("metric creation failed")
		return err
	}

	counter.Add(ctx, 1)
	ulog.Info("success!")

	return nil
}

func main() {
	lambda.EnableLocalHTTP("9000")
	lambda.AsyncStart(Do, false)
}
