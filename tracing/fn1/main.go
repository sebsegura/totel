package main

import (
	"github.com/Bancar/goala/utel"
	"github.com/Bancar/lambda-go"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net/http"
	"totel/fn1/pkg/service"
)

func main() {
	utel.SetUtelConfig("goala", "myflow")

	svc := service.New(&http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	})
	//lambda.EnableLocalHTTP("9000")
	lambda.SyncStart(svc.Handle, service.Instrument)
}
