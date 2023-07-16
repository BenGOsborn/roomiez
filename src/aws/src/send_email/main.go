package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bengosborn/roomiez/aws/utils"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func HandleRequest(ctx context.Context, sqsEvent events.SQSEvent) error {
	// Load requirements
	logger := log.New(os.Stdout, "[SendEmail] ", log.Ldate|log.Ltime)

	env, err := utils.LoadEnv(ctx)
	if err != nil {
		logger.Println(err)

		return err
	}

	client := sendgrid.NewSendClient(env.SendGridAPIKey)

	for _, message := range sqsEvent.Records {
		record := &utils.SubscriptionRecord{}

		if err := json.Unmarshal([]byte(message.Body), record); err != nil {
			logger.Println(err)

			continue
		}

		from := mail.NewEmail("Roomiez", "rentals@roomiez.co")
		subject := ""
		to := mail.NewEmail(record.ID, record.Email)

		email := mail.NewV3MailInit(from, subject, to)
		email.SetTemplateID(env.SendGridTemplateId)

		// **** We need to search first and make sure we have sufficient users

		// substitutions := map[string]string{
		// 	"placeholder1": "value1",
		// 	"placeholder2": "value2",
		// }

		if _, err := client.Send(email); err != nil {
			logger.Println(err)

			continue
		}
	}

	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
