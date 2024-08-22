package service

import (
	"context"
	"totel/utel"
)

func Instrument(next func(ctx context.Context, in *Request) (*Response, error)) func(context.Context, *Request) (*Response, error) {
	return func(ctx context.Context, in *Request) (*Response, error) {
		cfg := &utel.Config{
			ServiceName: "fn1",
			Owner:       "goala",
			Flow:        "myflow",
		}
		utel.SetUtelConfig(cfg)
		tp := utel.InitTracer(ctx, cfg)
		defer func() { _ = tp.Shutdown(ctx) }()

		return next(ctx, in)
	}
}
