package utils

import (
	"gorm.io/gorm"
)

type SearchParams struct {
	Page       uint
	Latitude   *float64
	Longitude  *float64
	Radius     *uint // metres
	Price      *int
	Bond       *int
	RentalType *string
	Gender     *string
	Age        *string
	Duration   *string
	Tenant     *string
	Features   *[]string
}

type RentalResult struct {
	ID                 uint    `json:"id"`
	URL                string  `json:"url"`
	Coordinates        *string `json:"coordinates"`
	Price              *int    `json:"price"`
	Bond               *int    `json:"bond"`
	RentalType         *string `json:"rentalType"`
	GenderPreference   *string `json:"gender"`
	AgePreference      *string `json:"age"`
	DurationPreference *string `json:"duration"`
	TenantPreference   *string `json:"tenant"`
}

type FeatureResult struct {
	RentalID    uint
	FeatureName string
}

type SearchResult struct {
	RentalResult
	Features *[]string `json:"features"`
}

// Find a list of rentals that match the search params
func SearchRentals(db *gorm.DB, searchParams *SearchParams, perPage uint) (*[]SearchResult, error) {
	// Retrieve all matching rentals
	query := db.Table("rentals")

	query = query.Select("rentals.id", "rentals.url", "rentals.coordinates", "rentals.price", "rentals.bond",
		"rental_types.type AS rental_type",
		"genders.preference AS gender_preference",
		"ages.preference AS age_preference",
		"durations.preference AS duration_preference",
		"tenants.preference AS tenant_preference",
	)

	query = query.Joins("LEFT JOIN rental_types ON rentals.rental_type_id = rental_types.id").
		Joins("LEFT JOIN genders ON rentals.gender_id = genders.id").
		Joins("LEFT JOIN ages ON rentals.age_id = ages.id").
		Joins("LEFT JOIN durations ON rentals.duration_id = durations.id").
		Joins("LEFT JOIN tenants ON rentals.tenant_id = tenants.id")

	if searchParams.Price != nil {
		query = query.Where("price = ?", searchParams.Price)
	}

	if searchParams.Bond != nil {
		query = query.Where("bond = ?", searchParams.Bond)
	}

	if searchParams.RentalType != nil {
		query = query.Where("rental_type = ?", searchParams.RentalType)
	}

	if searchParams.Gender != nil {
		query = query.Where("gender_preference = ?", searchParams.Gender)
	}

	if searchParams.Age != nil {
		query = query.Where("age_preference = ?", searchParams.Age)
	}

	if searchParams.Duration != nil {
		query = query.Where("duration_preference = ?", searchParams.Duration)
	}

	if searchParams.Tenant != nil {
		query = query.Where("tenant_preference = ?", searchParams.Tenant)
	}

	if searchParams.Features != nil {
		for _, feature := range *searchParams.Features {
			query = query.Where("EXISTS (SELECT 1 FROM rental_features JOIN features ON features.id = rental_features.feature_id WHERE rental_features.rental_id = rentals.id AND features.name = ?)", feature)
		}
	}

	if searchParams.Latitude != nil && searchParams.Longitude != nil && searchParams.Radius != nil {
		query = query.Where("ST_Distance(ST_GeomFromText(coordinates), POINT(?, ?)) <= ?", searchParams.Latitude, searchParams.Longitude, searchParams.Radius)
	}

	query = query.Offset((int(searchParams.Page) - 1) * int(perPage)).Limit(int(perPage))

	rentalResults := &[]RentalResult{}
	if err := query.Find(rentalResults).Error; err != nil {
		return nil, err
	}

	// Retrieve associated features
	query = db.Table("rental_features")

	query = query.Select(
		"rental_features.rental_id AS rental_id",
		"features.name AS feature_name",
	)

	query = query.Joins("JOIN features ON rental_features.feature_id = features.id")

	rentalIds := &[]uint{}
	for _, result := range *rentalResults {
		*rentalIds = append(*rentalIds, result.ID)
	}
	query = query.Where("rental_features.rental_id IN ?", *rentalIds)

	featureResults := &[]FeatureResult{}
	if err := query.Find(featureResults).Error; err != nil {
		return nil, err
	}

	// Join rentals and features
	features := make(map[uint][]string)

	for _, result := range *featureResults {
		features[result.RentalID] = append(features[result.RentalID], result.FeatureName)
	}

	searchResults := &[]SearchResult{}
	for _, result := range *rentalResults {
		temp, ok := features[result.ID]
		if ok {
			*searchResults = append(*searchResults, SearchResult{RentalResult: result, Features: &temp})
		} else {
			*searchResults = append(*searchResults, SearchResult{RentalResult: result})
		}

	}

	return searchResults, nil
}
