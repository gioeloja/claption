package main

import (
	"caption_service/go-aws/pkg"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(pkg.UploadHandler)
}
