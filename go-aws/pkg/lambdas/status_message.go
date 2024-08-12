package lambdas

import (
	"caption_service/go-aws/pkg/db"
	"caption_service/go-aws/pkg/models"
	"caption_service/go-aws/pkg/utils"
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

// When a message is sent to job status SQS queue, this lambda will be triggered. We will
// receive the message with the CaptionJob and then put it into dynamoDB.
func StatusSQSHandler(ctx context.Context, sqsEvent events.SQSEvent) error {
	currStage := "Job Status SQS Trigger"
	for _, message := range sqsEvent.Records {
		// unmarshall our json string
		var sqsMessage models.CaptionJob
		err := json.Unmarshal([]byte(message.Body), &sqsMessage)
		if err != nil {
			return fmt.Errorf("error unmarshalling SQS message: %v", err)
		}

		fmt.Printf("Processing message for image: %s\n", sqsMessage.JobID)
		fmt.Printf("Job status: %s\n", sqsMessage.Status)

		// Put data into dynamoDB
		err = db.PutJob(sqsMessage)
		if err != nil {
			err = utils.HandleError(sqsMessage.JobID, err, currStage, "putting job into dynamodb")
			return fmt.Errorf("error putting job into dynamodb %v", err)
		}
	}

	return nil
}
