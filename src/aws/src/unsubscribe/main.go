package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bengosborn/roomiez/aws/utils"
)

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	logger := log.New(os.Stdout, "[ProcessRental] ", log.Ldate|log.Ltime)

	env, err := utils.LoadEnv(ctx)
	if err != nil {
		logger.Println(err)

		return nil, err
	}

	return nil, nil
}

func main() {
	lambda.Start(HandleRequest)
}
