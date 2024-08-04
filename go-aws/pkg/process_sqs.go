package pkg

import (
	"caption_service/go-aws/pkg/models"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

// When a message is sent to SQS, this lambda will be triggered. We will
// receive the message with the image key details and then retrieve the image
// from the s3 bucket.
// Next, we pass the image through our docker image with the model, and then
// upload the caption result to the caption S3 bucket.
func ProcessSQSHandler(ctx context.Context, sqsEvent events.SQSEvent) error {
	// Load AWS SDK configuration
	sdkConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Printf("failed to load default config: %s", err)
		return err
	}

	for _, message := range sqsEvent.Records {
		// unmarshall our json string
		var sqsMessage models.SQSMessage
		err := json.Unmarshal([]byte(message.Body), &sqsMessage)
		if err != nil {
			return fmt.Errorf("error unmarshalling SQS message: %v", err)
		}

		fmt.Printf("Processing message for image: %s\n", sqsMessage.Key)
		fmt.Printf("Image in bucket: %s\n", sqsMessage.Bucket)

		// process the message
		s3Client := s3.NewFromConfig(sdkConfig)

		// fetch image from S3
		s3Resp, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(sqsMessage.Bucket),
			Key:    aws.String(sqsMessage.Key),
		})
		if err != nil {
			return fmt.Errorf("error getting object from S3: %v", err)
		}

		defer s3Resp.Body.Close()

		// Process the data from S3
		body, err := io.ReadAll(s3Resp.Body)
		if err != nil {
			return fmt.Errorf("error reading S3 object body: %v", err)
		}

		fmt.Printf("Object data: %s\n", string(body))

		// TODO: pass image into docker with our model
		// put the caption into a caption S3 bucket with the same key

	}

	return nil
}
