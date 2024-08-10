package lambdas

import (
	"bytes"
	"caption_service/go-aws/pkg/models"
	"caption_service/go-aws/pkg/utils"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

// When a message is sent to the process image SQS queue, this lambda will be triggered. We will
// receive the message with the image key details and then retrieve the image from the s3 bucket.
// Next, we pass the image through our docker image with the model, and then
// upload the caption result to the caption S3 bucket.
func ProcessSQSHandler(ctx context.Context, sqsEvent events.SQSEvent) error {
	// Load AWS SDK configuration
	sdkConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Printf("failed to load default config: %s", err)
		return err
	}

	for _, message := range sqsEvent.Records {
		// unmarshall our json string
		var sqsMessage models.ProcessMessage
		err := json.Unmarshal([]byte(message.Body), &sqsMessage)
		if err != nil {
			return fmt.Errorf("error unmarshalling SQS message: %v", err)
		}

		fmt.Printf("Processing message for image: %s\n", sqsMessage.Key)
		fmt.Printf("Image in bucket: %s\n", sqsMessage.Bucket)

		// process the message
		s3Client := s3.NewFromConfig(sdkConfig)

		// fetch image from S3
		s3Resp, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(sqsMessage.Bucket),
			Key:    aws.String(sqsMessage.Key),
		})
		if err != nil {
			return fmt.Errorf("error getting object from S3: %v", err)
		}

		defer s3Resp.Body.Close()

		// process the data from S3
		imageData, err := io.ReadAll(s3Resp.Body)
		if err != nil {
			return fmt.Errorf("error reading S3 object body: %v", err)
		}

		// send image to our flask endpoint with the model
		caption, err := SendImageToDocker(imageData)
		if err != nil {
			return fmt.Errorf("error when sending image to model endpoint: %v", err)
		}
		fmt.Printf("Caption is: %s", caption)

		// send job status update to status SQS queue
		err = utils.SendUpdateStatusMessage(context.Background(), models.CaptionJob{
			JobID:     sqsMessage.Key,
			Status:    "Complete",
			Caption:   &caption,
			Timestamp: int(time.Now().Unix()),
		})
		if err != nil {
			log.Fatalf("Failed to send status SQS message: %v", err)
		}
	}

	return nil
}

// Send image data to our dockerized python model
func SendImageToDocker(imageData []byte) (string, error) {
	// buffer to store our multipart form data
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// create form file field for 'file'
	fileWriter, err := w.CreateFormFile("file", "image.png")
	if err != nil {
		return "", fmt.Errorf("error creating form file field: %v", err)
	}

	// write image data to form field
	_, err = fileWriter.Write(imageData)
	if err != nil {
		return "", fmt.Errorf("error writing image data to form field: %v", err)
	}

	w.Close()

	// POST req
	url := os.Getenv("CAPTION_IMAGE_EC2_ENDPOINT")
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return "", fmt.Errorf("error creating HTTP request: %v", err)
	}

	// specify multipart
	req.Header.Set("Content-Type", w.FormDataContentType())

	// send req
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// process response
	var responseJson map[string]string
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	// unmarshall response back into JSON
	err = json.Unmarshal(responseBody, &responseJson)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected response status: %v, response body: %s", resp.Status, responseBody)
	}

	return responseJson["message"], nil
}
