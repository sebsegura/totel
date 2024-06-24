package main

import (
	"context"
	"fmt"
	"github.com/Bancar/lambda-go"
)

type Request struct {
	Hello string `json:"hello"`
}

func Do(ctx context.Context, in *Request) error {
	fmt.Println(in)
	return nil
}

func main() {
	//lambda.EnableLocalHTTP("9000")
	lambda.AsyncStart(Do, false)
}
