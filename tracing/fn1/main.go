package main

import (
	"github.com/Bancar/goala/utel"
	"github.com/Bancar/lambda-go"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net/http"
	"totel/fn1/pkg/service"
)

func main() {
	cfg := &utel.Config{
		ServiceName: "fn1",
		Owner:       "goala",
		Flow:        "myflow",
	}
	utel.SetUtelConfig(cfg)

	svc := service.New(&http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	})
	//lambda.EnableLocalHTTP("9000")
	lambda.SyncStart(svc.Handle, service.Instrument)
}
