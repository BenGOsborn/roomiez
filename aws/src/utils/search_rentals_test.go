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
		} else if (*rentalsFirst)[0].ID == (*rentalsSecond)[0].ID {
			t.Error("pages contain overlap")
		}
	})

	t.Run("Geo search", func(t *testing.T) {
		var perPage uint = 1

		latitude := 0.0
		longitude := 0.0
		var radius uint = 0

		rentals, err := utils.SearchRentals(db, &utils.SearchParams{Page: 1, Latitude: &latitude, Longitude: &longitude, Radius: &radius}, perPage)
		if err != nil {
			t.Error(err)
		} else if len(*rentals) != 0 {
			t.Error("out of bounds records included")
		}

		longitude = -40.0
		latitude = 160.0
		radius = 100

		rentals, err = utils.SearchRentals(db, &utils.SearchParams{Page: 1, Latitude: &latitude, Longitude: &longitude, Radius: &radius}, perPage)
		if err != nil {
			t.Error(err)
		} else if len(*rentals) == 0 {
			t.Error("no records found")
		}
	})
}