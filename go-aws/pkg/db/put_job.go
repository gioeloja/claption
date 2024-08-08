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

// Inserts or updates a CaptionJob into our DynamoDB table
func PutJob(job models.CaptionJob) error {
	// AWS client
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("REGION")),
	})
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}

	// dynamoDB client
	svc := dynamodb.New(sess)

	// marshall job into map of DynamoDB AttributeValues
	av, err := dynamodbattribute.MarshalMap(job)
	if err != nil {
		return fmt.Errorf("failed to marshal job struct: %v", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("ImageCaptionJobs"),
		Item:      av,
	}

	// put item in db
	_, err = svc.PutItem(input)
	if err != nil {
		return fmt.Errorf("failed to put item: %v", err)
	}

	return nil
}
