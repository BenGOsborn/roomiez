package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"gorm.io/gorm"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/location"
)

const (
	CentreLongitude = 151.2093
	CentreLatitude  = -33.8688
)

type PostSchema struct {
	Price      *int      `json:"price"`
	Bond       *int      `json:"bond"`
	Location   *string   `json:"location"`
	RentalType *string   `json:"rentalType"`
	Gender     *string   `json:"gender"`
	Age        *string   `json:"age"`
	Duration   *string   `json:"duration"`
	Tenant     *string   `json:"tenant"`
	Features   *[]string `json:"features"`
}

type RentalSchema struct {
	PostSchema
	Description string `json:"description"`
}

func (r *RentalSchema) String() string {
	data, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}

	return string(data)
}

// Extract all data from a post
func ProcessPost(ctx context.Context, llm *openai.Chat, post string) (*RentalSchema, error) {
	validationChain := NewPostValidation(llm)
	validation, err := chains.Run(ctx, validationChain, post, chains.WithTemperature(0.4))
	if err != nil {
		return nil, err
	} else if validation != "yes" {
		return nil, errors.New("invalid post")
	}

	extractionChain := NewPostExtraction(llm)
	rawData, err := chains.Run(ctx, extractionChain, post, chains.WithTemperature(0.4))
	if err != nil {
		return nil, err
	}

	postData := PostSchema{}
	if err := json.Unmarshal([]byte(rawData), &postData); err != nil {
		return nil, err
	}

	descriptionChain := NewPostDescription(llm)
	description, err := chains.Run(ctx, descriptionChain, post, chains.WithTemperature(0.4))
	if err != nil {
		return nil, err
	}

	rental := &RentalSchema{PostSchema: postData, Description: description}

	return rental, nil
}

// Get coords from an address
func CoordsFromAddress(ctx context.Context, address string, placeIndexName string) (float64, float64, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return -1, -1, err
	}

	svc := location.NewFromConfig(cfg)

	response, err := svc.SearchPlaceIndexForText(ctx, &location.SearchPlaceIndexForTextInput{
		IndexName:    &placeIndexName,
		Text:         &address,
		BiasPosition: []float64{CentreLongitude, CentreLatitude},
		MaxResults:   1,
	})
	if err != nil {
		return -1, -1, err
	}

	if len(response.Results) == 0 {
		return -1, -1, errors.New("invalid address")
	}

	place := response.Results[0].Place
	latitude := place.Geometry.Point[1]
	longitude := place.Geometry.Point[0]

	return latitude, longitude, nil
}

// Save a rental to database
func SaveRental(ctx context.Context, db *gorm.DB, rental *RentalSchema, url string, placeIndexName string) error {
	newRental := &Rental{}

	newRental.URL = url
	newRental.Price = rental.Price
	newRental.Bond = rental.Bond
	newRental.Description = rental.Description

	if rental.Location != nil {
		latitude, longitude, err := CoordsFromAddress(ctx, *rental.Location, placeIndexName)
		if err != nil {
			return err
		}

		coordinates := fmt.Sprintf("POINT(%f %f)", longitude, latitude)

		newRental.Coordinates = &coordinates
		newRental.Location = rental.Location
	}

	if rental.RentalType != nil {
		rentalType := &RentalType{}

		if err := db.First(rentalType, "type = ?", *rental.RentalType).Error; err != nil {
			return err
		}

		newRental.RentalTypeID = &rentalType.ID
	}

	if rental.Gender != nil {
		gender := &Gender{}
		if err := db.First(gender, "preference = ?", *rental.Gender).Error; err != nil {
			return err
		}

		newRental.GenderID = &gender.ID
	}

	if rental.Age != nil {
		age := &Age{}
		if err := db.First(age, "preference = ?", *rental.Age).Error; err != nil {
			return err
		}

		newRental.AgeID = &age.ID
	}

	if rental.Duration != nil {
		duration := &Duration{}
		if err := db.First(duration, "preference = ?", *rental.Duration).Error; err != nil {
			return err
		}

		newRental.DurationID = &duration.ID
	}

	if rental.Tenant != nil {
		tenant := &Tenant{}
		if err := db.First(tenant, "preference = ?", *rental.Tenant).Error; err != nil {
			return err
		}

		newRental.TenantID = &tenant.ID
	}

	if len(*rental.Features) > 0 {
		existingFeatures := &[]Feature{}
		if err := db.Find(existingFeatures, "name IN ?", *rental.Features).Error; err != nil {
			return err
		}

		newRental.Features = append(newRental.Features, *existingFeatures...)
	}

	return db.Create(newRental).Error
}
