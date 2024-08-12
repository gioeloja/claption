package db

import (
	"caption_service/go-aws/pkg/models"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Gets a CaptionJob from DynamoDB given a job id
func GetJob(jobID string) (*models.CaptionJob, error) {
	// AWS client
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("REGION")),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %v", err)
	}

	// dynamoDB client
	svc := dynamodb.New(sess)

	input := &dynamodb.GetItemInput{
		TableName: aws.String("ImageCaptionJobs"),
		Key: map[string]*dynamodb.AttributeValue{
			"JobID": {
				S: aws.String(jobID),
			},
		},
	}

	// get item from dynamodb
	result, err := svc.GetItem(input)
	if err != nil {
		return nil, fmt.Errorf("failed to get item from dynamoDB: %v", err)
	}

	if result.Item == nil {
		return nil, fmt.Errorf("job not found")
	}

	// unmarshall
	var captionJob models.CaptionJob
	err = dynamodbattribute.UnmarshalMap(result.Item, &captionJob)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall item: %v", err)
	}

	return &captionJob, nil
}
