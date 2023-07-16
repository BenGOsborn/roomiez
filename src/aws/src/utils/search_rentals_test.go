package utils_test

import (
	"context"
	"testing"

	"github.com/bengosborn/roomiez/aws/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestSearchRentals(t *testing.T) {
	ctx := context.Background()

	env, err := utils.LoadEnv(ctx)
	if err != nil {
		t.Error(err)
	}

	db, err := gorm.Open(mysql.Open(env.DSN))
	if err != nil {
		t.Error(err)
	}

	t.Run("Geo search", func(t *testing.T) {
		t.Helper()

		var radius uint = 1
		latitude, longitude, err := utils.CoordsFromAddress(ctx, "Surry Hills, new south wales", env.AWSLocationPlaceIndex)
		if err != nil {
			t.Error(err)
		}

		t.Log("Latitude", latitude)
		t.Log("Longitude", longitude)

		rentals, err := utils.SearchRentals(db, &utils.SearchParams{Page: 1, Latitude: &latitude, Longitude: &longitude, Radius: &radius})
		if err != nil {
			t.Error(err)
		} else if len(*rentals) == 0 {
			t.Error("no records found")
		}

	})

	t.Run("Features", func(t *testing.T) {
		t.Helper()

		rentals, err := utils.SearchRentals(db, &utils.SearchParams{Page: 1, Features: &[]string{"Mattress", "WiFi"}})
		if err != nil {
			t.Error(err)
		} else if len(*rentals) == 0 {
			t.Error("no records found")
		}
	})

	t.Run("Filters", func(t *testing.T) {
		t.Helper()

		age := "Young"
		tenantType := "Singles"
		duration := "Long Term"
		rentalType := "Apartment"
		gender := "Female"

		if _, err := utils.SearchRentals(db, &utils.SearchParams{Page: 1, Age: &age, Tenant: &tenantType, Duration: &duration, RentalType: &rentalType, Gender: &gender}); err != nil {
			t.Error(err)
		}
	})
}
