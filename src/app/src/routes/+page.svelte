<script lang="ts">
	import { PUBLIC_API_ENDPOINT } from "$env/static/public";
	import { getRentals, type SearchFields } from "$lib/api";
	import { derived, type Stores } from "svelte/store";
	import Pagination from "../components/Pagination.svelte";
	import Query from "../components/Query/index.svelte";
	import Rental from "../components/Rental.svelte";
	import { page, age, duration, features, gender, location, radius, rentalType, tenant, price, bond } from "../stores";

	let timeoutId: number | null = null;

	const searchFields = derived<Stores, SearchFields | undefined>(
		[page, age, duration, features, gender, location, radius, rentalType, tenant, price, bond],
		([$page, $age, $duration, $features, $gender, $location, $radius, $rentalType, $tenant, $price, $bond], set) => {
			if (timeoutId) clearTimeout(timeoutId);

			timeoutId = setTimeout(() => {
				const _age = !!$age ? $age : null;
				const _duration = !!$duration ? $duration : null;
				const _gender = !!$gender ? $gender : null;
				const _rentalType = !!$rentalType ? $rentalType : null;
				const _tenant = !!$tenant ? $tenant : null;

				const _location = !!$location ? { location: $location, radius: $radius } : null;

				const _price = $price <= 1500 ? $price : null;
				const _bond = $bond <= 3000 ? $bond : null;

				const _features = Object.entries($features).reduce((prev, value) => {
					if (value[1]) prev.push(value[0]);
					return prev;
				}, [] as string[]);

				const searchFields: SearchFields = {
					page: $page,
					age: _age,
					bond: _bond,
					duration: _duration,
					features: _features,
					gender: _gender,
					location: _location,
					price: _price,
					rentalType: _rentalType,
					tenant: _tenant
				};

				set(searchFields);
			}, 500);
		}
	);
</script>

<div class="mx-auto w-4/5 mt-8">
	<div class="flex justify-between space-x-8">
		<div class="w-1/4">
			<Query />
		</div>
		<div class="w-3/4">
			{#if $searchFields}
				{#await getRentals(PUBLIC_API_ENDPOINT, $searchFields)}
					<p class="text-center font-medium text-gray-800">Loading rentals...</p>
				{:then rentals}
					<div class="flex flex-col space-y-8">
						<div class="grid grid-cols-2 gap-8">
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
