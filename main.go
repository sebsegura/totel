package main

import (
	"context"
	"github.com/Bancar/goala/ulog"
	"github.com/Bancar/lambda-go"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	met "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"log"
)

type Event struct{}

const _ip = "18.208.231.129:4317"

// const _ip = "localhost:4317"

// const _ip = "otel-collector.otel-collector.local:4317"

func deltaSelector(kind metric.InstrumentKind) metricdata.Temporality {
	switch kind {
	case metric.InstrumentKindCounter,
		metric.InstrumentKindHistogram,
		metric.InstrumentKindObservableGauge,
		metric.InstrumentKindObservableCounter:
		return metricdata.DeltaTemporality
	case metric.InstrumentKindUpDownCounter,
		metric.InstrumentKindObservableUpDownCounter:
		return metricdata.CumulativeTemporality
	}
	panic("unknown instrument kind")
}

func InitMeter() *metric.MeterProvider {
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("totel")))
	if err != nil {
		log.Fatalf("cannot set resource: %v", err)
	}

	exporter, err := otlpmetricgrpc.New(
		context.Background(),
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(_ip),
		otlpmetricgrpc.WithTemporalitySelector(deltaSelector),
	)
	if err != nil {
		log.Fatalf("cannot create exporter: %v", err)
	}

	provider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(exporter)))

	return provider
}

func IncrementCounter(ctx context.Context, mp *metric.MeterProvider) {
	meter := mp.Meter("delta")
	counter, err := meter.Int64Counter(
		"delta.counter",
		met.WithDescription("OTEL test"),
		met.WithUnit("1"),
	)
	if err != nil {
		log.Fatalf("cannot create metric: %v", err)
	}

	counter.Add(ctx, 1)
}

func Do(ctx context.Context, _ *Event) error {
	ulog.Info("starting meter provider")
	mp := InitMeter()
	defer func() {
		if err := mp.Shutdown(ctx); err != nil {
			log.Fatalf("cannot shutdown meter provider: %v", err)
		}
	}()
	ulog.Info("started meter provider")

	for i := 0; i < 2; i++ {
		IncrementCounter(ctx, mp)
	}

	ulog.Info("success!")

	return nil
}

func main() {
	lambda.EnableLocalHTTP("9000")
	lambda.AsyncStart(Do, false)
	//lambda.Start(otellambda.InstrumentHandler(Do))
}
