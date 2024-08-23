package utel

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
)

func InitMeter(ctx context.Context, res *resource.Resource) (*metric.MeterProvider, error) {
	exporter, err := newMeterExporter(ctx)
	if err != nil {
		return nil, err
	}

	mp := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(exporter)),
	)

	return mp, nil
}

// deltaSelector function to obtain the corresponding metric temporality
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

func newMeterExporter(ctx context.Context) (*otlpmetricgrpc.Exporter, error) {
	return otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(url()),
		otlpmetricgrpc.WithTemporalitySelector(deltaSelector),
	)
}

func IncrementCounter(ctx context.Context, name string, attributes ...attribute.KeyValue) error {
	cfg := GetUtelConfig()
	mp := otel.GetMeterProvider()
	meter := mp.Meter(cfg.Owner)
	counter, err := meter.Int64Counter(name)
	if err != nil {
		return err
	}
	attributes = append(attributes, attribute.String("Flow", cfg.Flow))
	counter.Add(ctx, 1, otelmetric.WithAttributes(attributes...))
	return nil
}
