<script lang="ts">
	import { PUBLIC_API_ENDPOINT } from "$env/static/public";
	import { getRentals } from "$lib/api";
	import Query from "../components/Query/index.svelte";
	import Rental from "../components/Rental.svelte";
</script>

<div class="mx-auto w-4/5 mt-8">
	<div class="flex justify-between space-x-8">
		<div class="w-1/4">
			<Query />
		</div>
		<div class="w-3/4">
			{#await getRentals(PUBLIC_API_ENDPOINT, { page: 1 })}
				<p class="text-center font-medium text-gray-800">Loading rentals...</p>
			{:then rentals}
				<div class="grid grid-cols-2 gap-8">
					{#each rentals as rental (rental.id)}
						<Rental {rental} />
					{/each}
				</div>
			{/await}
		</div>
	</div>
</div>
