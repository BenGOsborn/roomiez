<script lang="ts">
	import { PUBLIC_API_ENDPOINT } from "$env/static/public";
	import { getFields, getRentals } from "$lib/api";
	import Rental from "../components/rental.svelte";
</script>

{#await getFields(PUBLIC_API_ENDPOINT)}
	<p>Loading fields...</p>
{:then data}
	<p>{JSON.stringify(data)}</p>
{/await}

{#await getRentals(PUBLIC_API_ENDPOINT, { page: 1 })}
	<p>Loading rentals...</p>
{:then rentals}
	{#each rentals as rental (rental.id)}
		<Rental {rental} />
	{/each}
{/await}
