package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
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

	// Process the post
	post := event.Post

	hash := md5.New()
	hash.Write([]byte(post))
	hashString := hex.EncodeToString(hash.Sum(nil))

	if err := db.Where("post_hash = ?", hashString).First(&utils.Rental{}).Error; err != nil {
		return nil, err
	}

	// Parse the post
	rental, err := utils.ProcessPost(ctx, llm, post)

	headers := map[string]string{
		"Content-Type":                "text/plain",
		"Access-Control-Allow-Origin": "*",
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    headers,
		Body:       "OK",
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
