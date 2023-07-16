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

	// **** We need a different way of getting this field and parsing it into our search params - we need to form a custom struct

	env, err := utils.LoadEnv(ctx)
	if err != nil {
		logger.Println(err)

		return nil, err
	}

	table := os.Getenv("TABLE")

	sess, err := session.NewSession()
	if err != nil {
		logger.Println(err)

		return nil, err
	}

	ddb := dynamodb.New(sess)

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
		Email:        body.Email,
	}

	av, err := dynamodbattribute.MarshalMap(record)
	if err != nil {
		logger.Println(err)

		return nil, err
	}

	if _, err := ddb.PutItem(&dynamodb.PutItemInput{Item: av, TableName: &table}); err != nil {
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
