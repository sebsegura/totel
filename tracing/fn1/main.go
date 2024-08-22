package main

import (
	"github.com/Bancar/lambda-go"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net/http"
	"totel/fn1/pkg/service"
)

func main() {
	svc := service.New(&http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	})
	//lambda.EnableLocalHTTP("9000")
	lambda.SyncStart(svc.Handle, service.Instrument)
}
