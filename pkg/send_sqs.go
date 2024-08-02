package pkg

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go/aws"
)

func SendSQSHandler(ctx context.Context, s3Event events.S3Event) error {
	// Load AWS SDK configuration
	sdkConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Printf("failed to load default config: %s", err)
		return err
	}

	// Create S3 and SQS clients
	s3Client := s3.NewFromConfig(sdkConfig)
	sqsClient := sqs.NewFromConfig(sdkConfig)

	// get env var defined in AWS
	sqsQueueURL := os.Getenv("IMAGE_PROCESSING_SQS_URL")

	for _, record := range s3Event.Records {
		bucket := record.S3.Bucket.Name
		key := record.S3.Object.URLDecodedKey

		// get object data type
		_, err := s3Client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: &bucket,
			Key:    &key,
		})
		if err != nil {
			log.Printf("error getting head of object %s/%s: %s", bucket, key, err)
			return err
		}

		// message model
		message := map[string]string{
			"bucket": bucket,
			"key":    key,
		}

		messageJSON, err := json.Marshal(message)
		if err != nil {
			log.Printf("error marshalling message body: %s", err)
			return err
		}

		// sending SQS msg
		_, err = sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
			QueueUrl:       &sqsQueueURL,
			MessageBody:    aws.String((string(messageJSON))),
			MessageGroupId: aws.String("image-processing-group"),
		})
		if err != nil {
			log.Printf("error sending message to SQS: %s", err)
			return err
		}

		log.Printf("successfully sent message to SQS for object %s/%s", bucket, key)
	}

	return nil
}
