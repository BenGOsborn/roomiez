package utils

import (
	"gorm.io/gorm"
)

type Rental struct {
	gorm.Model
	PostHash     string `gorm:"uniqueIndex,varchar(255)"`
	Coordinates  *string
	Price        *int
	Bond         *int
	RentalTypeID *uint
	GenderID     *uint
	AgeID        *uint
	DurationID   *uint
	TenantID     *uint
	Features     []Feature `gorm:"many2many:user_languages;"`
}

// Apartment, house, granny flat
type RentalType struct {
	gorm.Model
	Type    string `gorm:"uniqueIndex,varchar(255)"`
	Rentals []Rental
}

// Male, female, all
type Gender struct {
	gorm.Model
	Preference string `gorm:"uniqueIndex,varchar(255)"`
	Rentals    []Rental
}

// Young, middle aged, old
type Age struct {
	gorm.Model
	Preference string `gorm:"uniqueIndex,varchar(255)"`
	Rentals    []Rental
}

// Short term, long term, all
type Duration struct {
	gorm.Model
	Preference string `gorm:"uniqueIndex,varchar(255)"`
	Rentals    []Rental
}

// Singles, couples, all
type Tenant struct {
	gorm.Model
	Preference string `gorm:"uniqueIndex,varchar(255)"`
	Rentals    []Rental
}

// e.g. garage, bills included, furnished, wifi
type Feature struct {
	gorm.Model
	Name string `gorm:"uniqueIndex,varchar(255)"`
}
