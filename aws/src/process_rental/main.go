package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bengosborn/roomiez/aws/utils"
	"github.com/tmc/langchaingo/llms/openai"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Event struct {
	Post string `json:"post"`
	URL  string `json:"url"`
}

func HandleRequest(ctx context.Context, event Event) (*events.APIGatewayProxyResponse, error) {
	// Load requirements
	env, err := utils.LoadEnv(ctx)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(mysql.Open(env.DSN), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
	if err != nil {
		return nil, err
	}

	llm, err := openai.NewChat(openai.WithModel("gpt-4"), openai.WithToken(env.OpenAIAPIKey))
	if err != nil {
		return nil, err
	}

	// Process the post and ensure no duplicates
	post := event.Post
	url := event.URL

	if err := db.Where("url = ?", url).First(&utils.Rental{}).Error; err == nil {
		return nil, errors.New("post already exists")
	}

	rental, err := utils.ProcessPost(ctx, llm, post)
	if err != nil {
		return nil, err
	}

	// Save the post
	if err := utils.SaveRental(ctx, db, rental, url, env.AWSLocationPlaceIndex); err != nil {
		return nil, err
	}

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
