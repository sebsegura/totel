package main

import (
	"context"
	"github.com/Bancar/lambda-go"
	"totel/fn2/pkg/service"
	"totel/utel"
)

func main() {
	ctx := context.Background()
	cfg := &utel.Config{
		ServiceName: "fn2",
		Owner:       "goala",
		Flow:        "myflow",
	}
	utel.SetUtelConfig(cfg)
	tp := utel.InitTracer(ctx, cfg)
	defer func() { _ = tp.Shutdown(ctx) }()

	client := service.NewDDBClient(ctx)
	svc := service.New(client)
	lambda.EnableLocalHTTP("9001")
	lambda.SyncStart(svc.Handle)
}
