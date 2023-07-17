package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bengosborn/roomiez/aws/utils"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func HandleRequest(ctx context.Context, sqsEvent events.SQSEvent) error {
	// Load requirements
	logger := log.New(os.Stdout, "[SendEmail] ", log.Ldate|log.Ltime)

	env, err := utils.LoadEnv(ctx)
	if err != nil {
		logger.Println(err)

		return err
	}

	unsubscribeUrl := os.Getenv("UNSUBSCRIBE_URL")

	client := sendgrid.NewSendClient(env.SendGridAPIKey)

	db, err := gorm.Open(mysql.Open(env.DSN))
	if err != nil {
		logger.Println(err)

		return err
	}

	for _, message := range sqsEvent.Records {
		record := &utils.SubscriptionRecord{}

		if err := json.Unmarshal([]byte(message.Body), record); err != nil {
			logger.Println(err)

			return err
		}

		from := mail.NewEmail("Roomiez", "rentals@roomiez.co")
		subject := ""
		to := mail.NewEmail(record.ID, record.Email)

		email := mail.NewV3MailInit(from, subject, to)
		email.SetTemplateID(env.SendGridTemplateId)

		rentals, err := utils.SearchRentals(db, record.SearchParams)
		if err != nil {
			logger.Println(err)

			return err
		} else if len(*rentals) < 4 {
			logger.Println("not enough rentals found")

			return err
		}

		for i := 0; i < 4; i++ {
			body := fmt.Sprint("property_", i+1)
			url := fmt.Sprint("url_", i+1)

			email.Personalizations[0].SetDynamicTemplateData(body, (*rentals)[i].Description)
			email.Personalizations[0].SetDynamicTemplateData(url, (*rentals)[i].URL)
		}

		email.Personalizations[0].SetDynamicTemplateData("unsubscribe", fmt.Sprint(unsubscribeUrl, "?id=", record.ID))

		resp, err := client.Send(email)
		if err != nil {
			logger.Println(err)

			return err
		} else if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			err = errors.New(fmt.Sprint("bad response ", resp.StatusCode, " ", resp.Body))

			logger.Println(err)

			return err
		}

		logger.Println("sent email to ", record.ID)
	}

	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
