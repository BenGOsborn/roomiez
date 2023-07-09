package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bengosborn/roomiez/aws/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	PageSize = 25
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

	// Seach for rentals
	searchParams, err := ParseQueryString(&request.MultiValueQueryStringParameters)
	if err != nil {
		return nil, err
	}

	rentals, err := utils.SearchRentals(db, searchParams, PageSize)
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(rentals)
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

// Parse the query string
func ParseQueryString(queryString *map[string][]string) (*utils.SearchParams, error) {
	searchParams := &utils.SearchParams{Page: 1}

	if latitude, ok := (*queryString)["latitude"]; ok {
		temp, err := strconv.ParseFloat(latitude[0], 64)
		if err != nil {
			return nil, err
		}

		searchParams.Latitude = &temp
	}

	if longitude, ok := (*queryString)["longitude"]; ok {
		temp, err := strconv.ParseFloat(longitude[0], 64)
		if err != nil {
			return nil, err
		}

		searchParams.Longitude = &temp
	}

	if radius, ok := (*queryString)["radius"]; ok {
		tmp, err := strconv.Atoi(radius[0])
		if err != nil {
			return nil, err
		}

		temp := uint(tmp)
		searchParams.Radius = &temp
	}

	if price, ok := (*queryString)["price"]; ok {
		temp, err := strconv.Atoi(price[0])
		if err != nil {
			return nil, err
		}

		searchParams.Price = &temp
	}

	if bond, ok := (*queryString)["bond"]; ok {
		temp, err := strconv.Atoi(bond[0])
		if err != nil {
			return nil, err
		}

		searchParams.Bond = &temp
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
		temp, err := strconv.Atoi(page[0])
		if err != nil {
			return nil, err
		}

		searchParams.Page = uint(temp)
	}

	return searchParams, nil
}

func main() {
	lambda.Start(HandleRequest)
}
