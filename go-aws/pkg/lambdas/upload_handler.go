package lambdas

import (
	"bytes"
	"caption_service/go-aws/pkg/models"
	"caption_service/go-aws/pkg/utils"
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

// POST upload image handler. Request body should be an image's base64 encoding. { "body": "xyz" }
func UploadHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// make sure base64 was provided
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

	// retrieve bucket name and region from environment variables
	bucketName := os.Getenv("IMAGE_DATA_BUCKET_NAME")
	region := os.Getenv("REGION")

	// AWS session
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))

	s3Client := s3.New(sess)

	// unique ID for the object
	// TODO: put file extension here
	id := uuid.New().String()

	// Put image in S3
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(id),
		Body:   bytes.NewReader(fileData),
		//TODO: put content type
		//ContentType: aws.String("image/png"),
	})
	if err != nil {
		log.Printf("Failed to upload to S3: %s", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Failed to upload data to S3: " + err.Error()),
		}, nil
	}

	// send caption job SQS msg
	err = utils.SendUpdateStatusMessage(context.Background(), models.CaptionJob{
		JobID:     id,
		Status:    "Uploaded",
		Timestamp: int(time.Now().Unix()),
	})
	if err != nil {
		log.Fatalf("Failed to send status SQS message: %v", err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       id,
	}, nil
}
