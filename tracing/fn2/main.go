package main

import (
	"github.com/Bancar/lambda-go"
	"totel/fn2/pkg/service"
)

func main() {
	svc := service.New(service.NewInMemClient())
	//lambda.EnableLocalHTTP("9001")
	lambda.HTTPStart(svc.Handle, service.Instrument)
}
