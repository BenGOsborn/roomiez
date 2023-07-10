package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bengosborn/roomiez/aws/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// Load requirements
	env, err := utils.LoadEnv(ctx)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(mysql.Open(env.DSN))
	if err != nil {
		return nil, err
	}

	// Return all fields
	fields, err := utils.GetFields(db)
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(fields)
	if err != nil {
		return nil, err
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
		Body: string(body),
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
