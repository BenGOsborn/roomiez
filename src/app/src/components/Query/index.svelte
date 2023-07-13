<script lang="ts">
	import { PUBLIC_API_ENDPOINT } from "$env/static/public";
	import { getFields } from "$lib/api";
	import Location from "./Location.svelte";
	import Features from "./Features.svelte";
	import Price from "./Price.svelte";
	import Filters from "./Filters.svelte";

	let hidden = true;
</script>

{#await getFields(PUBLIC_API_ENDPOINT)}
	<p class="text-center font-medium text-gray-800">Loading fields...</p>
{:then fields}
	<div class="space-y-8 hidden md:flex flex-col">
		<Location />
		<Price />
		<Filters {fields} />
		<Features features={fields.feature} />
	</div>
	<div class="space-y-8 flex md:hidden flex-col">
		<div class="flex items-center justify-between">
			<p class="text-gray-800 font-bold">Search Options</p>
			<button class="py-1 px-2 rounded-md bg-gray-200 text-gray-500 font-medium hover:bg-gray-300" on:click={() => (hidden = !hidden)}>{hidden ? "Show" : "Hide"}</button>
		</div>
		{#if !hidden}
			<Location />
			<Price />
			<Filters {fields} />
			<Features features={fields.feature} />
		{/if}
	</div>
{/await}
