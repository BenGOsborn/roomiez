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

	t.Run("Pagination search", func(t *testing.T) {
		t.Helper()

		var perPage uint = 1

		rentalsFirst, err := utils.SearchRentals(db, &utils.SearchParams{Page: 1}, perPage)
		if err != nil {
			t.Error(err)
		} else if len(*rentalsFirst) > int(perPage) {
			t.Error("page size exceeded")
		}

		rentalsSecond, err := utils.SearchRentals(db, &utils.SearchParams{Page: 2}, perPage)
		if err != nil {
			t.Error(err)
		} else if len(*rentalsSecond) > int(perPage) {
			t.Error("page size exceeded")
		} else if (*rentalsFirst)[0].URL == (*rentalsSecond)[0].URL {
			t.Error("pages contain overlap")
		}
	})

	t.Run("Geo search", func(t *testing.T) {
		t.Helper()

		latitude := 0.0
		longitude := 0.0
		var radius uint = 0

		rentals, err := utils.SearchRentals(db, &utils.SearchParams{Page: 1, Latitude: &latitude, Longitude: &longitude, Radius: &radius}, 10)
		if err != nil {
			t.Error(err)
		} else if len(*rentals) != 0 {
			t.Error("out of bounds records included")
		}

	})

	t.Run("Features", func(t *testing.T) {
		t.Helper()

		rentals, err := utils.SearchRentals(db, &utils.SearchParams{Page: 1, Features: &[]string{"Mattress", "WiFi"}}, 10)
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

		if _, err := utils.SearchRentals(db, &utils.SearchParams{Page: 1, Age: &age, Tenant: &tenantType, Duration: &duration, RentalType: &rentalType, Gender: &gender}, 10); err != nil {
			t.Error(err)
		}
	})
}
