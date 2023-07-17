package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
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

func HandleRequest(ctx context.Context, sqsEvent events.SQSEvent) error {
	// Load requirements
	logger := log.New(os.Stdout, "[ProcessRental] ", log.Ldate|log.Ltime)

	env, err := utils.LoadEnv(ctx)
	if err != nil {
		logger.Println(err)

		return err
	}

	db, err := gorm.Open(mysql.Open(env.DSN))
	if err != nil {
		logger.Println(err)

		return err
	}

	llm, err := openai.NewChat(openai.WithModel("gpt-4"), openai.WithToken(env.OpenAIAPIKey))
	if err != nil {
		logger.Println(err)

		return err
	}

	// Process the post and ensure no duplicates
	for _, message := range sqsEvent.Records {
		body := &Body{}
		if err := json.Unmarshal([]byte(message.Body), body); err != nil {
			logger.Println(err)

			return err
		}

		if err := db.Where("url = ?", body.URL).First(&utils.Rental{}).Error; err == nil {
			err = errors.New("post already exists")
			logger.Println(err)

			return err
		}

		rental, err := utils.ProcessPost(ctx, llm, body.Post)
		if err != nil {
			logger.Println(err)

			return err
		}

		// Save the post
		if err := utils.SaveRental(ctx, db, rental, body.URL, env.AWSLocationPlaceIndex); err != nil {
			logger.Println(err)

			return err
		}

		logger.Println("processed ", body.URL)
	}

	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
