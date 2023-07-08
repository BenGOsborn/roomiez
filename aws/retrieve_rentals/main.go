package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context) (events.APIGatewayProxyResponse, error) {
	responseBody := "Hello, world!"
	headers := map[string]string{
		"Content-Type": "text/plain",
	}

	response := events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    headers,
		Body:       responseBody,
	}

	return response, nil
}

func main() {
	lambda.Start(HandleRequest)
}
