<script lang="ts">
	import type { Writable } from "svelte/store";

	export let name: string;
	export let values: string[];
	export let store: Writable<string>;

	const storeKey = `STORE:FILTER:${name}`;

	const saved = localStorage.getItem(storeKey);
	if (saved) store.set(saved);

	store.subscribe((value) => localStorage.setItem(storeKey, value));
</script>

<div class="space-y-2">
	<p class="text-gray-800 font-medium text-sm">{name}</p>
	<select class="text-gray-700" bind:value={$store}>
		<option value="">Any</option>
		{#each values as value}
			<option {value}>{value}</option>
		{/each}
	</select>
</div>
