package lambdas

import (
	"caption_service/go-aws/pkg/db"
	"caption_service/go-aws/pkg/models"
	"context"
	"encoding/json"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

const (
	maxAttempts = 3               // maximum number of attempts
	delay       = 1 * time.Second // delay between attempts
)

// Polling endpoint for a job's info, stops polling dynamodb once the captionjob has been completed or failed.
func GetStatusHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	jobID := request.PathParameters["jobID"]

	var job *models.CaptionJob
	var err error
	for attempts := 0; attempts < maxAttempts; attempts++ {
		job, err = db.GetJob(jobID)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       err.Error(),
			}, nil
		}

		// Caption job has finished processing
		if job.Status == "Complete" || job.Status == "Failed" {
			break
		}

		time.Sleep(delay)
	}

	if job.Status != "Complete" && job.Status != "Failed" {
		return events.APIGatewayProxyResponse{
			StatusCode: 202,
			Body:       "Job is still in progress. Please try again later.",
		}, nil
	}

	body, err := json.Marshal(job)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(body),
	}, nil
}
