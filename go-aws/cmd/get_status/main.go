package main

import (
	"caption_service/go-aws/pkg/lambdas"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(lambdas.GetStatusHandler)
}
