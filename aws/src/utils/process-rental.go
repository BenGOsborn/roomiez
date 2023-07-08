package utils

type RentalTypeSchema string

const (
	RentalTypeApartment RentalTypeSchema = "Apartment"
	RentalTypeHouse     RentalTypeSchema = "House"
)

type GenderSchema string

const (
	GenderMale   GenderSchema = "Male"
	GenderFemale GenderSchema = "Female"
)

type AgeSchema string

const (
	AgeYoung       AgeSchema = "Young"
	AggeMiddleAged AgeSchema = "Middle Aged"
	AgeOld         AgeSchema = "Old"
)

type DurationSchema string

const (
	DurationShortTerm DurationSchema = "Short Term"
	DurationLongTerm  DurationSchema = "Long Term"
)

type TenantSchema string

const (
	TenantSingles TenantSchema = "Singles"
	TenantCouples TenantSchema = "Couples"
)

type FeatureSchema string

const (
	FeatureGarage        FeatureSchema = "Garage"
	FeatureWiFi          FeatureSchema = "WiFi"
	FeatureBillsIncluded FeatureSchema = "Bills Included"
	FeatureFurnished     FeatureSchema = "Furnished"
)

type RentalSchema struct {
	Price      *int              `json:"price"`
	Location   *string           `json:"location"`
	RentalType *RentalTypeSchema `json:"rentalType"`
	Gender     *GenderSchema     `json:"gender"`
	Age        *AgeSchema        `json:"age"`
	Duration   *DurationSchema   `json:"duration"`
	Tenant     *TenantSchema     `json:"tenant"`
	Features   []FeatureSchema   `json:"features"`
}

// Process a post
func ProcessRental(rentalPost string, url string) {

}
