package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context) (events.APIGatewayProxyResponse, error) {
	headers := map[string]string{
		"Content-Type":                "text/plain",
		"Access-Control-Allow-Origin": "*",
	}

	response := events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    headers,
		Body:       "Hello world",
	}

	return response, nil
}

func main() {
	lambda.Start(HandleRequest)
}
