package utils

import (
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
	query := db.Model(&Rental{})

	if searchParams.Price != nil {
		query = query.Where("price = ?", *searchParams.Price)
	}

	if searchParams.Bond != nil {
		query = query.Where("bond = ?", *searchParams.Bond)
	}

	if searchParams.RentalType != nil {
		query = query.Joins("JOIN rental_types ON rentals.rental_type_id = rental_types.id").Where("rental_types.type = ?", *searchParams.RentalType)
	}

	if searchParams.Gender != nil {
		query = query.Joins("JOIN genders ON rentals.gender_id = genders.id").Where("genders.preference = ?", *searchParams.Gender)
	}

	if searchParams.Age != nil {
		query = query.Joins("JOIN ages ON rentals.age_id = ages.id").Where("ages.preference = ?", *searchParams.Age)
	}

	if searchParams.Duration != nil {
		query = query.Joins("JOIN durations ON rentals.duration_id = durations.id").Where("durations.preference = ?", *searchParams.Duration)
	}

	if searchParams.Tenant != nil {
		query = query.Joins("JOIN tenants ON rentals.tenant_id = tenants.id").Where("tenants.preference = ?", *searchParams.Tenant)
	}

	if searchParams.Latitude != nil && searchParams.Longitude != nil && searchParams.Radius != nil {
		query = query.Where("ST_Distance(ST_GeomFromText(coordinates), POINT(?, ?)) <= ?", searchParams.Latitude, searchParams.Longitude, searchParams.Radius)
	}

	query = query.Offset((int(searchParams.Page) - 1) * int(perPage)).Limit(int(perPage))

	rentals := &[]Rental{}
	if err := query.Find(rentals).Error; err != nil {
		return nil, err
	}

	return rentals, nil
}
