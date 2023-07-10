package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bengosborn/roomiez/aws/utils"
	"github.com/tmc/langchaingo/llms/openai"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Body struct {
	Post string `json:"post"`
	URL  string `json:"url"`
}

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// Load requirements
	logger := log.New(os.Stdout, "[ProcessRental] ", log.Ldate|log.Ltime)

	env, err := utils.LoadEnv(ctx)
	if err != nil {
		logger.Println(err)

		return nil, err
	}

	db, err := gorm.Open(mysql.Open(env.DSN))
	if err != nil {
		logger.Println(err)

		return nil, err
	}

	llm, err := openai.NewChat(openai.WithModel("gpt-4"), openai.WithToken(env.OpenAIAPIKey))
	if err != nil {
		logger.Println(err)

		return nil, err
	}

	// Process the post and ensure no duplicates
	body := &Body{}
	if err := json.Unmarshal([]byte(request.Body), body); err != nil {
		logger.Println(err)

		return nil, err
	}

	if err := db.Where("url = ?", body.URL).First(&utils.Rental{}).Error; err == nil {
		err = errors.New("post already exists")
		logger.Println(err)

		return nil, err
	}

	rental, err := utils.ProcessPost(ctx, llm, body.Post)
	if err != nil {
		logger.Println(err)

		return nil, err
	}

	// Save the post
	if err := utils.SaveRental(ctx, db, rental, body.URL, env.AWSLocationPlaceIndex); err != nil {
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
