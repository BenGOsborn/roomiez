package utils

import (
	"fmt"

	"gorm.io/gorm"
)

type SearchParams struct {
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
	Page       uint
}

// Find a list of rentals that match the search params
func SearchRentals(db *gorm.DB, searchParams *SearchParams, perPage uint) (*[]Rental, error) {
	query := db.Table("rentals")

	query = query.Select("rentals.url", "rentals.coordinates", "rentals.price", "rentals.bond",
		"rental_types.type", "genders.preference AS gender_preference",
		"ages.preference AS age_preference", "durations.preference AS duration_preference", "tenants.preference AS tenant_preference",
	)

	query = query.Joins("JOIN rental_types ON rentals.rental_type_id = rental_types.id").
		Joins("JOIN genders ON rentals.gender_id = genders.id").
		Joins("JOIN ages ON rentals.age_id = ages.id").
		Joins("JOIN durations ON rentals.duration_id = durations.id").
		Joins("JOIN tenants ON rentals.tenant_id = tenants.id")

	if searchParams.Price != nil {
		query = query.Where("price = ?", *searchParams.Price)
	}

	if searchParams.Bond != nil {
		query = query.Where("bond = ?", *searchParams.Bond)
	}

	if searchParams.RentalType != nil {
		query = query.Where("rental_types.type = ?", *searchParams.RentalType)
	}

	if searchParams.Gender != nil {
		query = query.Where("genders.preference = ?", *searchParams.Gender)
	}

	if searchParams.Age != nil {
		query = query.Where("ages.preference = ?", *searchParams.Age)
	}

	if searchParams.Duration != nil {
		query = query.Where("durations.preference = ?", *searchParams.Duration)
	}

	if searchParams.Tenant != nil {
		query = query.Where("tenants.preference = ?", *searchParams.Tenant)
	}

	if searchParams.Latitude != nil && searchParams.Longitude != nil && searchParams.Radius != nil {
		query = query.Where("ST_Distance(ST_GeomFromText(coordinates), POINT(?, ?)) <= ?", searchParams.Latitude, searchParams.Longitude, searchParams.Radius)
	}

	query = query.Offset((int(searchParams.Page) - 1) * int(perPage)).Limit(int(perPage))

	rentals := &[]Rental{}

	result := make(map[string]interface{})

	if err := query.Find(result).Error; err != nil {
		return nil, err
	}

	fmt.Println(result)

	return rentals, nil
}
