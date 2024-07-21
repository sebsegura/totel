package main

import (
	"context"
	"github.com/Bancar/lambda-go"
	"net/http"
	"totel/fn1/pkg/service"
	"totel/utel"
)

func main() {
	ctx := context.Background()
	cfg := &utel.Config{
		ServiceName: "fn1",
		Owner:       "goala",
		Flow:        "myflow",
	}
	utel.SetUtelConfig(cfg)
	tp := utel.InitTracer(ctx, cfg)
	defer func() { _ = tp.Shutdown(ctx) }()

	svc := service.New(&http.Client{})
	lambda.EnableLocalHTTP("9000")
	lambda.SyncStart(svc.Handle)
}
