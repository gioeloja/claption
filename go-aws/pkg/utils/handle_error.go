package utils

import (
	"caption_service/go-aws/pkg/models"
	"context"
	"log"
	"time"
)

// parameters:
//   - jobID: caption job id
//   - err: error to process
//   - stage: high lvl stage at which the error happened
//   - errContext: more specific details about what we were attempting when this error happened
func HandleError(jobID string, err error, stage string, errContext string) error {
	// Send update status message to status queue so we store the error in dynamoDB.
	if err != nil {
		captionError := &models.CaptionError{
			Stage:        stage,
			Context:      errContext,
			ErrorMessage: err.Error(),
		}

		captionJob := models.CaptionJob{
			JobID:     jobID,
			Status:    "Failed",
			Error:     captionError,
			Timestamp: int(time.Now().Unix()),
		}

		messageSendErr := SendUpdateStatusMessage(context.Background(), captionJob)
		if messageSendErr != nil {
			log.Fatalf("Failed to send status SQS message: %v", messageSendErr)
		}
	}

	// return the original err passed in
	return err
}
