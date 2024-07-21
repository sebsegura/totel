package utel

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
)

func NewResource(ctx context.Context, cfg *Config) (*resource.Resource, error) {
	opts := resource.WithAttributes(
		attribute.String("service.owner", cfg.Owner),
		attribute.String("service.flow", cfg.Flow),
	)
	res, err := resource.New(ctx, opts)
	if err != nil {
		return nil, err
	}

	return resource.Merge(
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(cfg.ServiceName)),
		res,
	)
}
