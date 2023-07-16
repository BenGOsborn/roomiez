package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/bengosborn/roomiez/aws/utils"
	"github.com/google/uuid"
)

func HandleRequest(ctx context.Context, request events.CloudWatchEvent) (*events.APIGatewayProxyResponse, error) {
	// Load requirements
	logger := log.New(os.Stdout, "[ScheduleEmail] ", log.Ldate|log.Ltime)

	utils.LoadEnv(ctx)

	table := os.Getenv("TABLE")
	sqsUrl := os.Getenv("SQS_URL")

	sess, err := session.NewSession()
	if err != nil {
		logger.Println(err)

		return nil, err
	}

	ddb := dynamodb.New(sess)
	q := sqs.New(sess)

	// Schedule emails to be sent
	entries := []*sqs.SendMessageBatchRequestEntry{}

	if err := ddb.ScanPages(&dynamodb.ScanInput{TableName: &table}, func(page *dynamodb.ScanOutput, lastPage bool) bool {
		for _, item := range page.Items {
			record := &utils.SubscriptionRecord{}

			if err := dynamodbattribute.UnmarshalMap(item, record); err != nil {
				logger.Println(err)

				continue
			}

			data, err := json.Marshal(record)
			if err != nil {
				logger.Println(err)

				continue
			}

			id := uuid.NewString()
			body := string(data)

			entries = append(entries, &sqs.SendMessageBatchRequestEntry{Id: &id, MessageBody: &body})
		}

		return !lastPage
	}); err != nil {
		logger.Println(err)

		return nil, err
	}

	if _, err := q.SendMessageBatch(&sqs.SendMessageBatchInput{QueueUrl: &sqsUrl, Entries: entries}); err != nil {
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
