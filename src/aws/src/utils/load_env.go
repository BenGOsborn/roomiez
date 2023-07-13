package utils

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/joho/godotenv"
)

type Env struct {
	DSN                   string `json:"DSN"`
	OpenAIAPIKey          string `json:"OPENAI_API_KEY"`
	AWSLocationPlaceIndex string `json:"AWS_LOCATION_PLACE_INDEX"`
	SendGridAPIKey        string `json:"SENDGRID_API_KEY"`
}

// Load environment variables from secrets
func LoadEnv(ctx context.Context) (*Env, error) {
	if os.Getenv("ENV") != "production" {
		if err := godotenv.Load("../../.env"); err != nil {
			return nil, err
		}
	}

	secretArn := os.Getenv("SECRETS_ARN")

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	client := secretsmanager.NewFromConfig(cfg)

	resp, err := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{SecretId: &secretArn})
	if err != nil {
		return nil, err
	}

	if resp.SecretString != nil {
		env := &Env{}
		if err := json.Unmarshal([]byte(*resp.SecretString), env); err != nil {
			return nil, err
		}

		return env, nil
	}

	return nil, errors.New("missing secret string")
}
