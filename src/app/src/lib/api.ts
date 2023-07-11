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
export async function getRentals(apiBase: string, searchFields: SearchFields): Promise<SearchResult[]> {
	// Convert the fields to a query string
	let queryParams = `?page=${searchFields.page}`;

	if (searchFields.location)
		queryParams += `&latitude=${searchFields.location.latitude}&longitude=${searchFields.location.longitude}&radius=${searchFields.location.radius}`;
	if (searchFields.price) queryParams += `&price=${searchFields.price}`;
	if (searchFields.bond) queryParams += `&bond=${searchFields.bond}`;
	if (searchFields.rentalType) queryParams += `&rentalType=${searchFields.rentalType}`;
	if (searchFields.gender) queryParams += `&gender=${searchFields.gender}`;
	if (searchFields.age) queryParams += `&age=${searchFields.age}`;
	if (searchFields.duration) queryParams += `&duration=${searchFields.duration}`;
	if (searchFields.tenant) queryParams += `&tenant=${searchFields.tenant}`;
	if (searchFields.features) searchFields.features.forEach((feature) => (queryParams += `&feature=${feature}`));

	// Make request
	const { data } = await axios.get<SearchResult[]>(`${apiBase}/rentals${queryParams}`);

	return data;
}
