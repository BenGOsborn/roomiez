<script lang="ts">
	import { PUBLIC_API_ENDPOINT } from "$env/static/public";
	import { getRentals, subscribeEmail, type SearchParams } from "$lib/api";
	import { derived, type Stores } from "svelte/store";
	import Pagination from "../components/Pagination.svelte";
	import Query from "../components/Query/index.svelte";
	import Rental from "../components/Rental.svelte";
	import { page, age, duration, features, gender, location, radius, rentalType, tenant, price, bond, email } from "../stores";
	import Map from "../components/Map.svelte";

	let timeoutId: number | null = null;

	const searchParams = derived<Stores, SearchParams | undefined>(
		[page, age, duration, features, gender, location, radius, rentalType, tenant, price, bond],
		([$page, $age, $duration, $features, $gender, $location, $radius, $rentalType, $tenant, $price, $bond], set) => {
			if (timeoutId) clearTimeout(timeoutId);

			timeoutId = setTimeout(() => {
				const _age = !!$age ? $age : null;
				const _duration = !!$duration ? $duration : null;
				const _gender = !!$gender ? $gender : null;
				const _rentalType = !!$rentalType ? $rentalType : null;
				const _tenant = !!$tenant ? $tenant : null;

				const _location = !!$location ? $location : null;
				const _radius = !!$radius ? $radius * 1000 : null;

				const _price = $price <= 1500 ? $price : null;
				const _bond = $bond <= 3000 ? $bond : null;

				const _features = Object.entries($features).reduce((prev, value) => {
					if (value[1]) prev.push(value[0]);
					return prev;
				}, [] as string[]);

				const searchParams: SearchParams = {
					page: $page,
					age: _age,
					bond: _bond,
					duration: _duration,
					features: _features,
					gender: _gender,
					location: _location,
					radius: _radius,
					price: _price,
					rentalType: _rentalType,
					tenant: _tenant
				};

				set(searchParams);
			}, 500);
		}
	);

	const emailSubscription = derived<Stores, { params: SearchParams; email: string } | undefined>([email, searchParams], ([$email, $searchParams]) => {
		if (!$searchParams || !$email) return;

		return {
			params: $searchParams,
			email: $email
		};
	});

	emailSubscription.subscribe(async (value) => {
		if (!!value) await subscribeEmail(PUBLIC_API_ENDPOINT, value.email, value.params);
	});
</script>

<div class="mx-auto w-4/5 mt-8">
	<div class="flex justify-between md:space-x-8 md:space-y-0 space-y-8 flex-col md:flex-row">
		<div class="md:w-1/4">
			<Query />
		</div>
		<div class="md:w-3/4">
			{#if $searchParams}
				{#await getRentals(PUBLIC_API_ENDPOINT, $searchParams)}
					<p class="text-center font-medium text-gray-800">Loading rentals...</p>
				{:then rentals}
					<div class="flex flex-col space-y-8">
						<Map {rentals} />
						<div class="grid grid-cols-1 gap-8 xl:grid-cols-2">
							{#each rentals as rental (rental.id)}
								<Rental {rental} />
							{/each}
						</div>
						<Pagination />
					</div>
				{/await}
			{/if}
		</div>
	</div>
</div>
