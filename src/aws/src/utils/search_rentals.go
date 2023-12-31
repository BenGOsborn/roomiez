package utils

import (
	"gorm.io/gorm"
)

const (
	PerPage = 10
)

type SearchParams struct {
	Page       uint      `json:"page"`
	Latitude   *float64  `json:"latitude"`
	Longitude  *float64  `json:"longitude"`
	Radius     *uint     `json:"radius"` // metres
	Price      *int      `json:"price"`
	Bond       *int      `json:"bond"`
	RentalType *string   `json:"rentalType"`
	Gender     *string   `json:"gender"`
	Age        *string   `json:"age"`
	Duration   *string   `json:"duration"`
	Tenant     *string   `json:"tenant"`
	Features   *[]string `json:"features"`
}

type RentalResult struct {
	ID                 uint    `json:"id"`
	URL                string  `json:"url"`
	Description        string  `json:"description"`
	Location           *string `json:"location"`
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
func SearchRentals(db *gorm.DB, searchParams *SearchParams) (*[]SearchResult, error) {
	// Retrieve all matching rentals
	query := db.Table("rentals")

	query = query.Select("rentals.id", "rentals.url", "rentals.description", "rentals.location", "rentals.coordinates", "rentals.price", "rentals.bond", "rentals.coordinates",
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
		query = query.Where("price <= ?", searchParams.Price)
	}

	if searchParams.Bond != nil {
		query = query.Where("bond <= ?", searchParams.Bond)
	}

	if searchParams.RentalType != nil {
		query = query.Where("rental_types.type = ?", searchParams.RentalType)
	}

	if searchParams.Gender != nil {
		query = query.Where("genders.preference = ?", searchParams.Gender)
	}

	if searchParams.Age != nil {
		query = query.Where("ages.preference = ?", searchParams.Age)
	}

	if searchParams.Duration != nil {
		query = query.Where("durations.preference = ?", searchParams.Duration)
	}

	if searchParams.Tenant != nil {
		query = query.Where("tenants.preference = ?", searchParams.Tenant)
	}

	if searchParams.Features != nil {
		for _, feature := range *searchParams.Features {
			query = query.Where("EXISTS (SELECT 1 FROM rental_features JOIN features ON features.id = rental_features.feature_id WHERE rental_features.rental_id = rentals.id AND features.name = ?)", feature)
		}
	}

	if searchParams.Latitude != nil && searchParams.Longitude != nil && searchParams.Radius != nil {
		query = query.Where("ST_Distance_Sphere(ST_GeomFromText(coordinates), POINT(?, ?)) <= ?", searchParams.Longitude, searchParams.Latitude, searchParams.Radius)
	}

	query = query.Offset((int(searchParams.Page) - 1) * int(PerPage)).Limit(int(PerPage)).Order("rentals.created_at DESC")

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

	featureResults := []FeatureResult{}
	if err := query.Find(&featureResults).Error; err != nil {
		return nil, err
	}

	// Join rentals and features
	features := make(map[uint]*[]string)

	for _, result := range featureResults {
		arr, ok := features[result.RentalID]
		if ok {
			*arr = append(*arr, result.FeatureName)
		} else {
			features[result.RentalID] = &[]string{result.FeatureName}
		}
	}

	searchResults := []SearchResult{}
	for _, result := range *rentalResults {
		temp, ok := features[result.ID]
		if ok {
			searchResults = append(searchResults, SearchResult{RentalResult: result, Features: temp})
		} else {
			searchResults = append(searchResults, SearchResult{RentalResult: result})
		}

	}

	return &searchResults, nil
}
