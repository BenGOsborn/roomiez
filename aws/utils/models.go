package utils

import (
	"time"

	"gorm.io/gorm"
)

type Rental struct {
	gorm.Model
	Post          string `gorm:"unique"`
	Price         *int
	Latitude      *int
	Longitude     *int
	AvailableFrom *time.Time
	AvailableTo   *time.Time
	Features      []Feature `gorm:"many2many:user_languages;"`
}

// Apartment, house, granny flat
type RentalType struct {
	gorm.Model
	Type    string
	Rentals []Rental
}

// Male, female, all
type Gender struct {
	gorm.Model
	Preference string
	Rentals    []Rental
}

// Young, middle aged, old
type Age struct {
	gorm.Model
	Preference string
	Rentals    []Rental
}

// Short term, long term, all
type Duration struct {
	gorm.Model
	Preference string
	Rentals    []Rental
}

// Singles, couples, all
type Tenant struct {
	gorm.Model
	Preference string
	Rentals    []Rental
}

// e.g. garage, bills included, furnished, wifi
type Feature struct {
	gorm.Model
	Name string
}
