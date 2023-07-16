package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/bengosborn/roomiez/aws/utils"
)

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// Load requirements
	logger := log.New(os.Stdout, "[Unsubscribe] ", log.Ldate|log.Ltime)

	utils.LoadEnv(ctx)

	table := os.Getenv("TABLE")

	sess, err := session.NewSession()
	if err != nil {
		logger.Println(err)

		return nil, err
	}

	ddb := dynamodb.New(sess)

	// Extract email
	id, ok := request.QueryStringParameters["id"]
	if !ok {
		err := errors.New("id not included in query string")
		logger.Panicln(err)

		return nil, err
	}

	// Delete record from dynamodb
	if _, err := ddb.DeleteItem(&dynamodb.DeleteItemInput{Key: map[string]*dynamodb.AttributeValue{"id": {S: &id}}, TableName: &table}); err != nil {
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
