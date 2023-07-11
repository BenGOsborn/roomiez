package utils

import "gorm.io/gorm"

type Fields struct {
	RentalType []string `json:"rentalType"`
	Gender     []string `json:"gender"`
	Age        []string `json:"age"`
	Duration   []string `json:"duration"`
	Tenant     []string `json:"tenant"`
	Feature    []string `json:"feature"`
}

func GetFields(db *gorm.DB) (*Fields, error) {
	fields := &Fields{}

	if err := db.Model(&RentalType{}).Pluck("type", &fields.RentalType).Error; err != nil {
		return nil, err
	}

	if err := db.Model(&Gender{}).Pluck("preference", &fields.Gender).Error; err != nil {
		return nil, err
	}

	if err := db.Model(&Age{}).Pluck("preference", &fields.Age).Error; err != nil {
		return nil, err
	}

	if err := db.Model(&Duration{}).Pluck("preference", &fields.Duration).Error; err != nil {
		return nil, err
	}

	if err := db.Model(&Tenant{}).Pluck("preference", &fields.Tenant).Error; err != nil {
		return nil, err
	}

	if err := db.Model(&Feature{}).Pluck("name", &fields.Feature).Error; err != nil {
		return nil, err
	}

	return fields, nil
}
