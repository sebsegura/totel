package utel

import (
	"context"
	"go.opentelemetry.io/otel"
)

func EnableOTELInstrumentation(ctx context.Context) (func(context.Context) error, error) {
	res, err := NewResource(ctx, GetUtelConfig())
	if err != nil {
		return nil, err
	}

	mp, err := InitMeter(ctx, res)
	if err != nil {
		return nil, err
	}

	otel.SetMeterProvider(mp)

	tp, err := InitTracer(ctx, res)
	if err != nil {
		return nil, err
	}

	otel.SetTracerProvider(tp)

	return func(ctx context.Context) error {
		if err = mp.Shutdown(ctx); err != nil {
			return err
		}

		if err = tp.Shutdown(ctx); err != nil {
			return err
		}

		return nil
	}, nil
}
