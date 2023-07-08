package utils

import (
	"time"

	"gorm.io/gorm"
)

type Rental struct {
	gorm.Model
	Post          string
	Price         *int
	Latitude      *int
	Longitude     *int
	Garage        *bool
	AvailableFrom *time.Time
	AvailableTo   *time.Time
}

// Apartment, house, granny flat
type RentalType struct {
	gorm.Model
	Type string
}

// Male, female, all
type Gender struct {
	gorm.Model
	Preference string
}

// Young, middle aged, old
type Age struct {
	gorm.Model
	Preference string
}

// Short term, long term, all
type Duration struct {
	gorm.Model
	Preference string
}

// Singles, couples, all
type Tenant struct {
	gorm.Model
	Preference string
}

// Features
type Features struct {
}
