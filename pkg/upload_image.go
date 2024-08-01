package pkg

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

const (
	bucketName = "cs-image-data-bucket"
	region     = "us-east-1"
)

// POST upload image handler. Request body should be an image's base64 encoding. { "body": "xyz" }
func UploadHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	godotenv.Load()

	// Make sure base64 was provided
	base64Data := request.Body
	if base64Data == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "No image data provided",
		}, nil
	}

	// decode base64 data
	fileData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf("Unable to decode file data: %s", err.Error()),
		}, nil
	}

	// create S3 session
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))

	s3Client := s3.New(sess)

	// generate unique ID for the object
	id := uuid.New().String() + ".png"

	// put image in s3
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(id),
		Body:        bytes.NewReader(fileData),
		ContentType: aws.String("image/png"),
	})
	if err != nil {
		log.Printf("Failed to upload to S3: %s", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Failed to upload data to S3: " + err.Error()),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Image successfully uploaded",
	}, nil
}
