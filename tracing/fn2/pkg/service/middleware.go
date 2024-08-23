package service

import (
	"go.opentelemetry.io/otel/attribute"
	"log"
	"net/http"
	"totel/utel"
)

func Instrument(next func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		shutdown, err := utel.EnableOTELInstrumentation(ctx)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			_ = shutdown(ctx)
		}()

		wr := &statusRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next(wr, r)

		labels := []attribute.KeyValue{
			attribute.String("method", r.Method),
			attribute.Int("status_code", wr.statusCode),
		}

		_ = utel.IncrementCounter(ctx, "fn2-requests", labels...)
	}
}

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *statusRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *statusRecorder) Write(body []byte) (int, error) {
	return r.ResponseWriter.Write(body)
}
