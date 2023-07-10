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
	Features   []string
}

type SearchResult struct {
	ID                 uint
	URL                string
	Coordinates        *string
	Price              *int
	Bond               *int
	RentalType         *string
	GenderPreference   *string
	AgePreference      *string
	DurationPreference *string
	TenantPreference   *string
	// Features           []string
}

// **** I need 3 structs - the final combines the two for the array and the searchresult here - we need to fetch and group by separately for our records then assign the correct features from the rental ids (group by rental id)

// Find a list of rentals that match the search params
func SearchRentals(db *gorm.DB, searchParams *SearchParams, perPage uint) (*[]SearchResult, error) {
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
		for _, feature := range searchParams.Features {
			query = query.Where("EXISTS (SELECT 1 FROM user_languages JOIN features ON features.id = user_languages.feature_id WHERE user_languages.rental_id = rentals.id AND features.name = ?)", feature)
		}
	}

	if searchParams.Latitude != nil && searchParams.Longitude != nil && searchParams.Radius != nil {
		query = query.Where("ST_Distance(ST_GeomFromText(coordinates), POINT(?, ?)) <= ?", searchParams.Latitude, searchParams.Longitude, searchParams.Radius)
	}

	query = query.Offset((int(searchParams.Page) - 1) * int(perPage)).Limit(int(perPage))

	results := &[]SearchResult{}
	if err := query.Find(results).Error; err != nil {
		return nil, err
	}

	return results, nil
}
