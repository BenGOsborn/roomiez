package main_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	main "github.com/bengosborn/roomiez/aws/process_rental"
	"github.com/bengosborn/roomiez/aws/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestHandleRequest(t *testing.T) {
	ctx := context.Background()

	env, err := utils.LoadEnv(ctx)
	if err != nil {
		t.Error(err)
	}

	db, err := gorm.Open(mysql.Open(env.DSN), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
	if err != nil {
		t.Error(err)
	}

	post := `Hello!
		We have found a new place to call 'home' in beautiful Drummoyne and are now looking for a third female to join us!
		THE PROPERTY:
		A tastefully remodeled 3-bedroom apartment located in a well-maintained building, featuring a unique townhouse-style layout spread across two floors.
		🚏3 mins to bus stop 
		⛴️ 5 mins to ferry wharf (Circular Quay/Olympic Park 
		less than 10 mins to IGA & Harris Farm 
		🚌 5 min bus ride to Birkenhead point 
		🏃‍♀️ quick stroll to the Bay Run 
		THE ROOM:
		Cosy and can comfortably fit a bed of any size
		- Unfurnished
		- Room 2 on floor plan 
		- Built-in wardrobe 
		- Balcony access
		- Shared bathroom 
		Rent: $270 per week 
		Bond: $1080 
		Bills to be spilt evenly (internet & electricity)
		Approx: $70 per month
		Lease start date: 18 July 2023 
		YOU: 
		-A friendly person in their mid to late 20/early 30s.  
		- Someone who values cleanliness and organization and is happy to stick to a cleaning schedule.
		- A respectful individual who understands the importance of personal space.
		- Ideally, someone who works full-time and limited wfh 
		- non smoker and not a party animal 
		-minimal stay 6 months and you will need to be approved by the agent. 
		-someone who is looking for a place to call home 
		ABOUT US:
		Tabitha (26): Sri Lankan law grad, trying hard to keep my plants alive. I am a big foodie and enjoy exploring cultures through food. I love farmers' markets, art galleries, museums and theatre.
		Laura (29): I am a social worker who loves a margarita or a glass of wine on a Friday night. I love exploring new restaurants with friends. I am a skincare, haircare and makeup enthusiast. I enjoy reading and love to travel with my next trip already planned.
		Open for inspections this Saturday (8th July)☺️
		If you feel like this could be your new home, please message me with a little bit about yourself!
		P.S: The property was offered unfurnished. We have some furniture, including a sofa, kitchen appliances, and a TV. You are welcome to bring your own furniture or contribute towards getting any additional items needed.`
	url := "https://www.facebook.com/groups/2280085492006745/permalink/7148701135145132/"

	t.Run("Save to database", func(t *testing.T) {
		price := 270
		bond := 1080
		location := "Drummoyne"
		var rentalType utils.RentalTypeSchema = "Apartment"
		var gender utils.GenderSchema = "Female"
		var age utils.AgeSchema = "Young"
		var duration utils.DurationSchema = "Long Term"
		var tenant utils.TenantSchema = "Singles"
		features := []utils.FeatureSchema{"WiFi"}

		rental := &utils.RentalSchema{
			Price:      &price,
			Bond:       &bond,
			Location:   &location,
			RentalType: &rentalType,
			Gender:     &gender,
			Age:        &age,
			Duration:   &duration,
			Tenant:     &tenant,
			Features:   features,
		}

		if err := utils.SaveRental(ctx, db, rental, url, env.AWSLocationPlaceIndex); err != nil {
			t.Error(err)
		}
	})

	t.Run("Create post", func(t *testing.T) {
		body, err := json.Marshal(&main.Body{Post: post, URL: url})
		if err != nil {
			t.Error(err)
		}

		if _, err := main.HandleRequest(ctx, events.APIGatewayProxyRequest{Body: string(body)}); err != nil {
			t.Error(err)
		}
	})
}
