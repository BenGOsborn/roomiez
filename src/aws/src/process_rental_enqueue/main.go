package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/bengosborn/roomiez/aws/utils"
)

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// Load requirements
	logger := log.New(os.Stdout, "[ProcessRentalEnqueue] ", log.Ldate|log.Ltime)

	utils.LoadEnv(ctx)

	sqsUrl := os.Getenv("SQS_URL")

	sess, err := session.NewSession()
	if err != nil {
		logger.Println(err)

		return nil, err
	}

	q := sqs.New(sess)

	// Add the post
	if _, err := q.SendMessage(&sqs.SendMessageInput{QueueUrl: &sqsUrl, MessageBody: &request.Body}); err != nil {
		logger.Println(err)

		return nil, err
	}

	logger.Println("OK")

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type":                "text/plain",
			"Access-Control-Allow-Origin": "*",
		},
		Body: "OK",
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
