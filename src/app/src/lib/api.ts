import axios from "axios";
import { z } from "zod";

const fieldsSchema = z.object({
	rentalType: z.array(z.string()),
	gender: z.array(z.string()),
	age: z.array(z.string()),
	duration: z.array(z.string()),
	tenant: z.array(z.string()),
	feature: z.array(z.string())
});

type Fields = z.infer<typeof fieldsSchema>;

// Get the search fields
export async function getFields(apiBase: string): Promise<Fields> {
	const { data } = await axios.get<Fields>(`${apiBase}/rentals/fields`);

	return fieldsSchema.parse(data);
}

export interface SearchFields {
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

const searchResultSchema = z.object({
	id: z.number(),
	url: z.string(),
	location: z.string().nullable(),
	price: z.number().nullable(),
	bond: z.number().nullable(),
	rentalType: z.string().nullable(),
	gender: z.string().nullable(),
	age: z.string().nullable(),
	duration: z.string().nullable(),
	tenant: z.string().nullable(),
	features: z.array(z.string()).nullable()
});

export type SearchResult = z.infer<typeof searchResultSchema>;

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

	return z.array(searchResultSchema).parse(data);
}
