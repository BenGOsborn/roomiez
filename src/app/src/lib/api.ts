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

export type Fields = z.infer<typeof fieldsSchema>;

// Get the search fields
export async function getFields(apiBase: string): Promise<Fields> {
	const { data } = await axios.get<Fields>(`${apiBase}/rentals/fields`);

	return fieldsSchema.parse(data);
}

export interface SearchParams {
	page: number;
	location: {
		location: string;
		radius: number;
	} | null;
	price: number | null;
	bond: number | null;
	rentalType: string | null;
	gender: string | null;
	age: string | null;
	duration: string | null;
	tenant: string | null;
	features: string[] | null;
}

const searchResultSchema = z.object({
	id: z.number(),
	url: z.string(),
	description: z.string(),
	location: z.string().nullable(),
	coordinates: z.string().nullable(),
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
export async function getRentals(apiBase: string, searchParams: SearchParams): Promise<SearchResult[]> {
	// Convert the fields to a query string
	let queryParams = `?page=${searchParams.page}`;

	if (searchParams.location) queryParams += `&location=${searchParams.location.location}&radius=${searchParams.location.radius}`;
	if (searchParams.price) queryParams += `&price=${searchParams.price}`;
	if (searchParams.bond) queryParams += `&bond=${searchParams.bond}`;
	if (searchParams.rentalType) queryParams += `&rentalType=${searchParams.rentalType}`;
	if (searchParams.gender) queryParams += `&gender=${searchParams.gender}`;
	if (searchParams.age) queryParams += `&age=${searchParams.age}`;
	if (searchParams.duration) queryParams += `&duration=${searchParams.duration}`;
	if (searchParams.tenant) queryParams += `&tenant=${searchParams.tenant}`;
	if (searchParams.features) searchParams.features.forEach((feature) => (queryParams += `&feature=${feature}`));

	// Make request
	const { data } = await axios.get<SearchResult[]>(`${apiBase}/rentals${queryParams}`);

	return z.array(searchResultSchema).parse(data);
}

// Subscribe a user
export async function subscribeEmail(apiBase: string, email: string, searchParams: SearchParams): Promise<void> {
	const data = {
		email,
		searchParams
	};

	// Make request
	await axios.put(`${apiBase}/subscribe`, data);
}
