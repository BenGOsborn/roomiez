package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bengosborn/roomiez/aws/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// Load requirements
	logger := log.New(os.Stdout, "[RetrieveRentals] ", log.Ldate|log.Ltime)

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

	// Seach for rentals
	searchParams, err := ParseQueryString(ctx, env.AWSLocationPlaceIndex, &request.MultiValueQueryStringParameters)
	if err != nil {
		logger.Println(err)

		return nil, err
	}

	rentals, err := utils.SearchRentals(db, searchParams)
	if err != nil {
		logger.Println(err)

		return nil, err
	}

	body, err := json.Marshal(rentals)
	if err != nil {
		return nil, err
	}

	logger.Println(string(body))

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
func ParseQueryString(ctx context.Context, placeIndexName string, queryString *map[string][]string) (*utils.SearchParams, error) {
	searchParams := &utils.SearchParams{Page: 1}

	if location, ok := (*queryString)["location"]; ok {
		latitude, longitude, err := utils.CoordsFromAddress(ctx, location[0], placeIndexName)
		if err != nil {
			return nil, err
		}

		searchParams.Latitude = &latitude
		searchParams.Longitude = &longitude
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

	if features, ok := (*queryString)["feature"]; ok {
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
