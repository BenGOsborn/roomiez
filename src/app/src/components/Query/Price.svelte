<script lang="ts">
	import { z } from "zod";
	import { bond, price } from "../../stores";

	const STORE_KEY = "STORE:PRICE";

	const stateSchema = z.object({
		price: z.number(),
		bond: z.number()
	});
	type State = z.infer<typeof stateSchema>;

	const saved = localStorage.getItem(STORE_KEY);
	if (saved) {
		const state = stateSchema.parse(JSON.parse(saved));

		price.set(state.price);
		bond.set(state.bond);
	}

	price.subscribe((value) => localStorage.setItem(STORE_KEY, JSON.stringify({ price: value, bond: $bond } as State)));
	bond.subscribe((value) => localStorage.setItem(STORE_KEY, JSON.stringify({ price: $price, bond: value } as State)));
</script>

<div class="bg-white p-4 rounded-md flex flex-col space-y-4 drop-shadow">
	<p class="text-gray-800 font-bold">Price</p>
	<div class="flex flex-col space-y-2">
		<label for="rent" class="text-gray-600 font-medium text-sm">Maximum weekly rent (AUD)</label>
		<div class="flex items-center justify-between space-x-2">
			<input id="rent" class="w-4/5" placeholder="Radius" type="range" min="0" max="1501" bind:value={$price} />
			{#if $price <= 1500}
				<span class="font-medium text-gray-600">${$price.toLocaleString()}</span>
			{:else}
				<span class="font-medium text-gray-600">Any</span>
			{/if}
		</div>
	</div>
	<div class="flex flex-col space-y-2">
		<label for="rent" class="text-gray-600 font-medium text-sm">Maximum bond (AUD)</label>
		<div class="flex items-center justify-between space-x-2">
			<input class="w-4/5" placeholder="Radius" type="range" min="0" max="3001" bind:value={$bond} />
			{#if $bond <= 3000}
				<span class="font-medium text-gray-600">${$bond.toLocaleString()}</span>
			{:else}
				<span class="font-medium text-gray-600">Any</span>
			{/if}
		</div>
	</div>
</div>
