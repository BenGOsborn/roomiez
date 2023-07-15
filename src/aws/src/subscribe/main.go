package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
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

	env, err := utils.LoadEnv(ctx)
	if err != nil {
		logger.Println(err)

		return nil, err
	}

	sess, err := session.NewSession()
	if err != nil {
		logger.Println(err)

		return nil, err
	}

	svc := dynamodb.New(sess)

	// Extract email
	body := &Body{}
	if err := json.Unmarshal([]byte(request.Body), body); err != nil {
		logger.Println(err)

		return nil, err
	}

	// Store data in dynamodb
	key, err := utils.GenerateKey(body.Email, env.SecretKey)
	if err != nil {
		logger.Println(err)

		return nil, err
	}

	record := &utils.SubscriptionRecord{
		ID:           key,
		SearchParams: &body.SearchParams,
		Timestamp:    time.Now(),
	}

	av, err := dynamodbattribute.MarshalMap(record)
	if err != nil {
		logger.Println(err)

		return nil, err
	}

	if _, err := svc.PutItem(&dynamodb.PutItemInput{Item: av, TableName: &table}); err != nil {
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
