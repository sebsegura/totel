package service

import (
	"github.com/Bancar/goala/utel"
	"go.opentelemetry.io/otel"
	"net/http"
	ut "totel/utel"
)

func Instrument(next func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg := &ut.Config{
			ServiceName: "fn2",
			Owner:       "goala",
			Flow:        "myflow",
		}
		ut.SetUtelConfig(cfg)

		ctx := r.Context()
		mp := utel.EnableOTELMetric(ctx, cfg.Owner)
		otel.SetMeterProvider(mp)
		defer func() { _ = mp.Shutdown(ctx) }()

		next(w, r)
	}
}
