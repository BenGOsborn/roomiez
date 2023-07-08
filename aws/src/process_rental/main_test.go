package main_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/bengosborn/roomiez/aws/utils"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestMain(t *testing.T) {

	t.Run("Seed data", func(t *testing.T) {
		if os.Getenv("ENV") != "production" {
			if err := godotenv.Load("../../.env"); err != nil {
				fmt.Println(err)
			}
		}

		dsn := os.Getenv("DSN")

		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
		if err != nil {
			fmt.Println(err)
			t.FailNow()
		}

		// Initialize tables
		if err := db.AutoMigrate(&utils.Rental{}); err != nil {
			fmt.Println(err)
			t.FailNow()
		}

		if err := db.AutoMigrate(&utils.RentalType{}); err != nil {
			fmt.Println(err)
			t.FailNow()
		}

		if err := db.AutoMigrate(&utils.Gender{}); err != nil {
			fmt.Println(err)
			t.FailNow()
		}

		if err := db.AutoMigrate(&utils.Age{}); err != nil {
			fmt.Println(err)
			t.FailNow()
		}

		if err := db.AutoMigrate(&utils.Duration{}); err != nil {
			fmt.Println(err)
			t.FailNow()
		}

		if err := db.AutoMigrate(&utils.Tenant{}); err != nil {
			fmt.Println(err)
			t.FailNow()
		}

		// Create seed data
		if res := db.Create(&utils.RentalType{Type: "Apartment"}); res.Error != nil {
			fmt.Println(res.Error)
		}
		if res := db.Create(&utils.RentalType{Type: "House"}); res.Error != nil {
			fmt.Println(res.Error)
		}

		if res := db.Create(&utils.Gender{Preference: "Male"}); res.Error != nil {
			fmt.Println(res.Error)
		}
		if res := db.Create(&utils.Gender{Preference: "Female"}); res.Error != nil {
			fmt.Println(res.Error)
		}

		if res := db.Create(&utils.Age{Preference: "Young"}); res.Error != nil {
			fmt.Println(res.Error)
		}
		if res := db.Create(&utils.Age{Preference: "Middle Aged"}); res.Error != nil {
			fmt.Println(res.Error)
		}
		if res := db.Create(&utils.Age{Preference: "Old"}); res.Error != nil {
			fmt.Println(res.Error)
		}

		if res := db.Create(&utils.Duration{Preference: "Short Term"}); res.Error != nil {
			fmt.Println(res.Error)
		}
		if res := db.Create(&utils.Duration{Preference: "Long Term"}); res.Error != nil {
			fmt.Println(res.Error)
		}

		if res := db.Create(&utils.Tenant{Preference: "Singles"}); res.Error != nil {
			fmt.Println(res.Error)
		}
		if res := db.Create(&utils.Tenant{Preference: "Couples"}); res.Error != nil {
			fmt.Println(res.Error)
		}

		if res := db.Create(&utils.Feature{Name: "Garage"}); res.Error != nil {
			fmt.Println(res.Error)
		}
		if res := db.Create(&utils.Feature{Name: "WiFi"}); res.Error != nil {
			fmt.Println(res.Error)
		}
		if res := db.Create(&utils.Feature{Name: "Bills Included"}); res.Error != nil {
			fmt.Println(res.Error)
		}
		if res := db.Create(&utils.Feature{Name: "Furnished"}); res.Error != nil {
			fmt.Println(res.Error)
		}
	})
}
