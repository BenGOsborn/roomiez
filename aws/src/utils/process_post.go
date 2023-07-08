package utils

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
)

type RentalTypeSchema string

const (
	RentalTypeApartment RentalTypeSchema = "Apartment"
	RentalTypeHouse     RentalTypeSchema = "House"
)

type GenderSchema string

const (
	GenderMale   GenderSchema = "Male"
	GenderFemale GenderSchema = "Female"
)

type AgeSchema string

const (
	AgeYoung       AgeSchema = "Young"
	AggeMiddleAged AgeSchema = "Middle Aged"
	AgeOld         AgeSchema = "Old"
)

type DurationSchema string

const (
	DurationShortTerm DurationSchema = "Short Term"
	DurationLongTerm  DurationSchema = "Long Term"
)

type TenantSchema string

const (
	TenantSingles TenantSchema = "Singles"
	TenantCouples TenantSchema = "Couples"
)

type FeatureSchema string

const (
	FeatureGarage        FeatureSchema = "Garage"
	FeatureWiFi          FeatureSchema = "WiFi"
	FeatureBillsIncluded FeatureSchema = "Bills Included"
	FeatureFurnished     FeatureSchema = "Furnished"
)

type RentalSchema struct {
	Price      *int              `json:"price"`
	Bond       *int              `json:"bond"`
	Location   *string           `json:"location"`
	RentalType *RentalTypeSchema `json:"rentalType"`
	Gender     *GenderSchema     `json:"gender"`
	Age        *AgeSchema        `json:"age"`
	Duration   *DurationSchema   `json:"duration"`
	Tenant     *TenantSchema     `json:"tenant"`
	Features   []FeatureSchema   `json:"features"`
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

	rental := &RentalSchema{}
	if err := json.Unmarshal([]byte(rawData), rental); err != nil {
		return nil, err
	}

	return rental, nil
}
