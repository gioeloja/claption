package pkg

type SQSMessage struct {
	Bucket string
	Key    string
}

func NewSQSMessage(bucket string, key string) *SQSMessage {
	return &SQSMessage{
		Bucket: bucket,
		Key:    key,
	}
}
