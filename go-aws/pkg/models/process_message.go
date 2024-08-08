package models

// Model for messages to the process image SQS queue
type ProcessMessage struct {
	Bucket string
	Key    string
}
