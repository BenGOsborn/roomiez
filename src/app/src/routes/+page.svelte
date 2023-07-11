<script lang="ts" async="true">
	import { PUBLIC_API_ENDPOINT } from "$env/static/public";
	import { getFields, getRentals } from "$lib/api";

	let x = 0;

    const fieldsPromise = getFields(PUBLIC_API_ENDPOINT)
    const rentalsPromise = getRentals(PUBLIC_API_ENDPOINT, {page: 1})
</script>

<h1>Welcome to SvelteKit</h1>
<p>{x}</p>
<p>{PUBLIC_API_ENDPOINT}</p>

{#await fieldsPromise}
    <p>Loading fields...</p>
{:then data} 
    <p>{JSON.stringify(data)}</p> 
{/await}

{#await rentalsPromise}
    <p>Loading rentals...</p> 
{:then data} 
    <p>{JSON.stringify(data)}</p> 
{/await}

<button on:click={() => (x = x + 1)}>Increment</button>
