package models

// This is the format for caption job data in DynamoDB
// This model can also be used for the update status SQS queue
type CaptionJob struct {
	JobID string
	// UserID    string
	// ImageBase64Encoding string
	Status       string
	Caption      *string
	ErrorMessage *string
	Timestamp    int
}
