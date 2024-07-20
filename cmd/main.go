package main

import (
	"caption_service/pkg"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
)

func main() {
	// Get env vars
	godotenv.Load()

	// Initialize AWS session
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(os.Getenv("AWS_REGION")),
		Endpoint:         aws.String(os.Getenv("LOCALSTACK_ENDPOINT")),
		S3ForcePathStyle: aws.Bool(true), // needed for localstack to work
	})

	if err != nil {
		fmt.Printf("Failed to initialize new session: %v", err)
		return
	}

	// Initialize S3 service
	s3Client := s3.New(sess)

	bucketName := "test-bucket"
	err = pkg.CreateBucket(s3Client, bucketName)
	if err != nil {
		fmt.Printf("Couldn't create bucket: %v", err)
		return
	}

	fmt.Println("Successfully created bucket")
}
