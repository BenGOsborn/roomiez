package utils_test

import (
	"context"
	"testing"

	"github.com/bengosborn/roomiez/aws/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestModels(t *testing.T) {
	ctx := context.Background()

	env, err := utils.LoadEnv(ctx)
	if err != nil {
		t.Error(err)
	}

	db, err := gorm.Open(mysql.Open(env.DSN), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
	if err != nil {
		t.Error(err)
	}

	t.Run("Seed data", func(t *testing.T) {
		t.Helper()

		// Initialize tables
		if err := db.AutoMigrate(&utils.Rental{}); err != nil {
			t.Error(err)
		}

		if err := db.AutoMigrate(&utils.RentalType{}); err != nil {
			t.Error(err)
		}

		if err := db.AutoMigrate(&utils.Gender{}); err != nil {
			t.Error(err)
		}

		if err := db.AutoMigrate(&utils.Age{}); err != nil {
			t.Error(err)
		}

		if err := db.AutoMigrate(&utils.Duration{}); err != nil {
			t.Error(err)
		}

		if err := db.AutoMigrate(&utils.Tenant{}); err != nil {
			t.Error(err)
		}

		// Create seed data
		if res := db.Create(&utils.RentalType{Type: "Apartment"}); res.Error != nil {
			t.Log(res.Error)
		}
		if res := db.Create(&utils.RentalType{Type: "House"}); res.Error != nil {
			t.Log(res.Error)
		}

		if res := db.Create(&utils.Gender{Preference: "Male"}); res.Error != nil {
			t.Log(res.Error)
		}
		if res := db.Create(&utils.Gender{Preference: "Female"}); res.Error != nil {
			t.Log(res.Error)
		}

		if res := db.Create(&utils.Age{Preference: "Young"}); res.Error != nil {
			t.Log(res.Error)
		}
		if res := db.Create(&utils.Age{Preference: "Middle Aged"}); res.Error != nil {
			t.Log(res.Error)
		}
		if res := db.Create(&utils.Age{Preference: "Old"}); res.Error != nil {
			t.Log(res.Error)
		}

		if res := db.Create(&utils.Duration{Preference: "Short Term"}); res.Error != nil {
			t.Log(res.Error)
		}
		if res := db.Create(&utils.Duration{Preference: "Long Term"}); res.Error != nil {
			t.Log(res.Error)
		}

		if res := db.Create(&utils.Tenant{Preference: "Singles"}); res.Error != nil {
			t.Log(res.Error)
		}
		if res := db.Create(&utils.Tenant{Preference: "Couples"}); res.Error != nil {
			t.Log(res.Error)
		}

		if res := db.Create(&utils.Feature{Name: "Garage"}); res.Error != nil {
			t.Log(res.Error)
		}
		if res := db.Create(&utils.Feature{Name: "WiFi"}); res.Error != nil {
			t.Log(res.Error)
		}
		if res := db.Create(&utils.Feature{Name: "Bills Included"}); res.Error != nil {
			t.Log(res.Error)
		}
		if res := db.Create(&utils.Feature{Name: "Furnished"}); res.Error != nil {
			t.Log(res.Error)
		}
		if res := db.Create(&utils.Feature{Name: "Pets Allowed"}); res.Error != nil {
			t.Log(res.Error)
		}
		if res := db.Create(&utils.Feature{Name: "Garage"}); res.Error != nil {
			t.Log(res.Error)
		}
		if res := db.Create(&utils.Feature{Name: "Mattress"}); res.Error != nil {
			t.Log(res.Error)
		}
		if res := db.Create(&utils.Feature{Name: "Pool"}); res.Error != nil {
			t.Log(res.Error)
		}
		if res := db.Create(&utils.Feature{Name: "Gym"}); res.Error != nil {
			t.Log(res.Error)
		}
	})
}
