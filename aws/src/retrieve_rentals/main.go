package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bengosborn/roomiez/aws/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	PageSize = 25
)

type SearchParams struct {
	Latitude   *string
	Longitude  *string
	Radius     *string
	Price      *string
	Bond       *string
	RentalType *string
	Gender     *string
	Age        *string
	Duration   *string
	Tenant     *string
	Features   *[]string
	Page       string
}

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

	// Make search query
	searchParams := ParseQueryString(&request.MultiValueQueryStringParameters)

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
		Body: "Hello world",
	}, nil
}

// Parse the query string
func ParseQueryString(queryString *map[string][]string) *SearchParams {
	searchParams := &SearchParams{Page: "1"}

	if latitude, ok := (*queryString)["latitude"]; ok {
		searchParams.Latitude = &latitude[0]
	}

	if longitude, ok := (*queryString)["longitude"]; ok {
		searchParams.Longitude = &longitude[0]
	}

	if radius, ok := (*queryString)["radius"]; ok {
		searchParams.Radius = &radius[0]
	}

	if price, ok := (*queryString)["price"]; ok {
		searchParams.Price = &price[0]
	}

	if bond, ok := (*queryString)["bond"]; ok {
		searchParams.Bond = &bond[0]
	}

	if rentalType, ok := (*queryString)["rentalType"]; ok {
		searchParams.RentalType = &rentalType[0]
	}

	if gender, ok := (*queryString)["gender"]; ok {
		searchParams.Gender = &gender[0]
	}

	if age, ok := (*queryString)["age"]; ok {
		searchParams.Age = &age[0]
	}

	if duration, ok := (*queryString)["duration"]; ok {
		searchParams.Duration = &duration[0]
	}

	if tenant, ok := (*queryString)["tenant"]; ok {
		searchParams.Tenant = &tenant[0]
	}

	if features, ok := (*queryString)["features"]; ok {
		searchParams.Features = &features
	}

	if page, ok := (*queryString)["page"]; ok {
		searchParams.Page = page[0]
	}

	return searchParams
}

func main() {
	lambda.Start(HandleRequest)
}
