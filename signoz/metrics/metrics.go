package metrics

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	met "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
	"log"
)

const (
	_collectorURL = "localhost:4317"
)

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
	exporter, err := otlpmetricgrpc.New(
		context.Background(),
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(_collectorURL),
		otlpmetricgrpc.WithTemporalitySelector(deltaSelector),
	)
	if err != nil {
		log.Fatalf("failed to create exporter: %v", err)
	}

	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", "mentat"),
			attribute.String("service.version", "1.0.0")))
	if err != nil {
		log.Fatalf("failed to set resource: %v", err)
	}

	provider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(exporter)))

	return provider
}

func IncrementCounter(ctx context.Context, mp *metric.MeterProvider, name string) error {
	meter := mp.Meter("mentat")

	counter, err := meter.Int64Counter(
		name,
		met.WithDescription("number of runs"),
		met.WithUnit("1"),
	)
	if err != nil {
		return err
	}

	counter.Add(ctx, 1)

	return nil
}
