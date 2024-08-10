package lambdas

import (
	"caption_service/go-aws/pkg/models"
	"caption_service/go-aws/pkg/utils"
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go/aws"
)

// When an image is uploaded to S3, this lambda is triggered. This lambda will
// get the bucket/key details for the image and send a message to SQS.
func SendProcessMessage(ctx context.Context, s3Event events.S3Event) error {
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
	processQueueURL := os.Getenv("IMAGE_PROCESSING_SQS_URL")

	for _, record := range s3Event.Records {
		bucket := record.S3.Bucket.Name
		key := record.S3.Object.URLDecodedKey

		// get object data type
		_, err := s3Client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: &bucket,
			Key:    &key,
		})

		// TODO: maybe verify that this is a png using data type?

		if err != nil {
			log.Printf("error getting head of object %s/%s: %s", bucket, key, err)
			return err
		}

		// initialize SQS message
		message := models.ProcessMessage{
			Bucket: bucket,
			Key:    key,
		}

		messageJSON, err := json.Marshal(message)
		if err != nil {
			log.Printf("error marshalling message body: %s", err)
			return err
		}

		// sending process SQS msg
		_, err = sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
			QueueUrl:       &processQueueURL,
			MessageBody:    aws.String((string(messageJSON))),
			MessageGroupId: aws.String("image-processing-group"),
		})
		if err != nil {
			log.Printf("error sending message to process image SQS: %s", err)
			return err
		}

		// send job update
		// send caption job SQS msg
		err = utils.SendUpdateStatusMessage(context.Background(), models.CaptionJob{
			JobID:     key,
			Status:    "Image Processing Requested",
			Timestamp: int(time.Now().Unix()),
		})
		if err != nil {
			log.Fatalf("Failed to send status SQS message: %v", err)
		}

		log.Printf("successfully sent message to SQS for object %s/%s", bucket, key)
	}

	return nil
}
