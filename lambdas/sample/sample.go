package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/glassechidna/go-emf/emf"
)

func main() {
	lambda.Start(handle)
}

func handle(ctx context.Context, input json.RawMessage) error {
	fmt.Println("hello world")
	fmt.Println(`{"some": "json"}`)
	emf.Emit(emf.MSI{
		"Input": input,
		"Field": "Value",
	})

	return nil
}
