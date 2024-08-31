package main

import (
	"github.com/Bancar/goala/utel"
	"github.com/Bancar/lambda-go"
	"totel/fn2/pkg/service"
)

func main() {
	cfg := &utel.Config{
		ServiceName: "fn2",
		Owner:       "goala",
		Flow:        "myflow",
	}
	utel.SetUtelConfig(cfg)

	svc := service.New(service.NewInMemClient())
	//lambda.EnableLocalHTTP("9001")
	lambda.HTTPStart(svc.Handle, service.Instrument)
}
