package main

import (
	"caption_service/pkg"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(pkg.SendSQSHandler)
}
