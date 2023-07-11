import axios from "axios";

interface Fields {
	rentalType: string[];
	gender: string[];
	age: string[];
	duration: string[];
	tenant: string[];
	feature: string[];
}

// Get the search fields
export async function getFields(apiBase: string): Promise<Fields> {
	const { data } = await axios.get<Fields>(`${apiBase}/rentals/fields`);

	return data;
}

interface SearchFields {
	page: number;
	location?: {
		latitude: number;
		longitude: number;
		radius: number;
	};
	price?: number;
	bond?: number;
	rentalType?: string;
	gender?: string;
	age?: string;
	duration?: string;
	tenant?: string;
	features?: string[];
}

// type RentalResult struct {
// 	ID                 uint    `json:"id"`
// 	URL                string  `json:"url"`
// 	Location           *string `json:"location"`
// 	Price              *int    `json:"price"`
// 	Bond               *int    `json:"bond"`
// 	RentalType         *string `json:"rentalType"`
// 	GenderPreference   *string `json:"gender"`
// 	AgePreference      *string `json:"age"`
// 	DurationPreference *string `json:"duration"`
// 	TenantPreference   *string `json:"tenant"`
// }

// type FeatureResult struct {
// 	RentalID    uint
// 	FeatureName string
// }

// type SearchResult struct {
// 	RentalResult
// 	Features *[]string `json:"features"`
// }

interface SearchResult {
	id: number;
	url: string;
	location?: string;
	price?: number;
	bond?: number;
	rentalType?: string;
	gender?: string;
	age?: string;
	duration?: string;
	tenant?: string;
	features?: string[];
}

// Get rentals matching the search criteria
export async function getRentals(apiBase: string) {
	// First convert the fields to a query string
}
