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

type Params struct {
	Page       uint      `json:"page"`
	Location   *string   `json:"location"`
	Radius     *uint     `json:"radius"`
	Price      *int      `json:"price"`
	Bond       *int      `json:"bond"`
	RentalType *string   `json:"rentalType"`
	Gender     *string   `json:"gender"`
	Age        *string   `json:"age"`
	Duration   *string   `json:"duration"`
	Tenant     *string   `json:"tenant"`
	Features   *[]string `json:"features"`
}

type Body struct {
	Email  string `json:"email"`
	Params Params `json:"params"`
}

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// Load requirements
	logger := log.New(os.Stdout, "[Subscribe] ", log.Ldate|log.Ltime)

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

	// Extract features
	body := &Body{}
	if err := json.Unmarshal([]byte(request.Body), body); err != nil {
		logger.Println(err)

		return nil, err
	}

	latitude, longitude, err := utils.CoordsFromAddress(ctx, *body.Params.Location, env.AWSLocationPlaceIndex)
	if err != nil {
		logger.Println(err)

		return nil, err
	}

	searchParams := &utils.SearchParams{
		Page:       body.Params.Page,
		Latitude:   &latitude,
		Longitude:  &longitude,
		Radius:     body.Params.Radius,
		Price:      body.Params.Price,
		Bond:       body.Params.Bond,
		RentalType: body.Params.RentalType,
		Gender:     body.Params.Gender,
		Age:        body.Params.Age,
		Duration:   body.Params.Duration,
		Tenant:     body.Params.Tenant,
		Features:   body.Params.Features,
	}

	// Store data in dynamodb
	key, err := utils.GenerateKey(body.Email, env.SecretKey)
	if err != nil {
		logger.Println(err)

		return nil, err
	}

	record := &utils.SubscriptionRecord{
		ID:           key,
		SearchParams: searchParams,
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
