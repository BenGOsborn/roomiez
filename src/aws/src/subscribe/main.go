package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bengosborn/roomiez/aws/utils"
)

type Body struct {
	Email        string             `json:"email"`
	SearchParams utils.SearchParams `json:"searchParams"`
}

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// Load requirements
	logger := log.New(os.Stdout, "[Subscribe] ", log.Ldate|log.Ltime)

	table := os.Getenv("TABLE")

	// Extract email
	body := &Body{}
	if err := json.Unmarshal([]byte(request.Body), body); err != nil {
		logger.Println(err)

		return nil, err
	}

	// Store data in dynamodb

	// **** We need to hash the key for the unsubscribe option (needs a GSI), insert the params as json, and insert the email

	logger.Println(table)

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
