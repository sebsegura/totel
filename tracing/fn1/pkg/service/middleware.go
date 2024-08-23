package service

import (
	"context"
	"log"
	"totel/utel"
)

func Instrument(next func(ctx context.Context, in *Request) (*Response, error)) func(context.Context, *Request) (*Response, error) {
	return func(ctx context.Context, in *Request) (*Response, error) {
		shutdown, err := utel.EnableOTELInstrumentation(ctx)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			_ = shutdown(ctx)
		}()

		return next(ctx, in)
	}
}
