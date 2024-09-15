package main

import (
	"github.com/Bancar/goala/utel"
	"github.com/Bancar/lambda-go"
	"totel/fn2/pkg/service"
)

func main() {
	utel.SetUtelConfig("goala", "myflow")

	svc := service.New(service.NewInMemClient())
	//lambda.EnableLocalHTTP("9001")
	lambda.HTTPStart(svc.Handle, service.Instrument)
}
