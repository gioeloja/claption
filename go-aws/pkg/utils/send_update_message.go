package utils

import (
	"caption_service/go-aws/pkg/models"
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go/aws"
)

// Sends an update status message to the status SQS queue. This will trigger an update in dynamoDB
func SendUpdateStatusMessage(ctx context.Context, job models.CaptionJob) error {
	// Load AWS SDK configuration
	sdkConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Printf("failed to load default config: %s", err)
		return err
	}

	// Create SQS client
	sqsClient := sqs.NewFromConfig(sdkConfig)

	// get env var defined in AWS
	sqsQueueURL := os.Getenv("STATUS_UPDATES_SQS_URL")

	messageJSON, err := json.Marshal(job)
	if err != nil {
		log.Printf("error marshalling message body: %s", err)
		return err
	}

	// sending SQS msg
	_, err = sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:       &sqsQueueURL,
		MessageBody:    aws.String((string(messageJSON))),
		MessageGroupId: aws.String("job-status-group"),
	})
	if err != nil {
		log.Printf("error sending message to SQS: %s", err)
		return err
	}

	log.Printf("successfully sent status update message for job %s", job.JobID)

	return nil
}
