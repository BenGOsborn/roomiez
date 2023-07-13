<script lang="ts">
	import { z } from "zod";
	import { location, mapReady, radius } from "../../stores";
	import { onMount } from "svelte";

	const STORE_KEY = "STORE:LOCATION";

	const stateSchema = z.object({
		location: z.string(),
		radius: z.number()
	});
	type State = z.infer<typeof stateSchema>;

	const saved = localStorage.getItem(STORE_KEY);
	if (saved) {
		const state = stateSchema.parse(JSON.parse(saved));

		location.set(state.location);
		radius.set(state.radius);
	}

	location.subscribe((value) => localStorage.setItem(STORE_KEY, JSON.stringify({ location: value, radius: $radius } as State)));
	radius.subscribe((value) => localStorage.setItem(STORE_KEY, JSON.stringify({ location: $location, radius: value } as State)));

	let autocomplete: any = null;
	let completions: string[] = [];
	let closeCompletions: boolean = false;

	let timeoutId: number | null = null;

	location.subscribe((value) => {
		if (closeCompletions) {
			completions = [];
			closeCompletions = false;
		} else if (autocomplete && !!value) {
			if (timeoutId) clearTimeout(timeoutId);

			timeoutId = setTimeout(() => {
				autocomplete
					.getPlacePredictions({
						input: value
					})
					.then((value: any) => (completions = value.predictions.slice(0, 3).map((result: any) => result.description)));
			}, 500);
		}
	});

	onMount(() => {
		mapReady.subscribe((value) => {
			if (value && !autocomplete)
				// @ts-ignore
				autocomplete = new google.maps.places.AutocompleteService();
		});
	});
</script>

<div class="bg-white p-4 rounded-md flex flex-col space-y-4 drop-shadow z-50">
	<p class="text-gray-800 font-bold">Location</p>
	<div class="relative">
		<input class="w-full outline-none px-2 py-1 border border-gray-300 rounded-md text-gray-600" placeholder="Location" type="text" bind:value={$location} />
		{#if completions.length > 0 && !!$location}
			<ul class="absolute rounded-md drop-shadow bg-white w-full p-2 space-y-2">
				{#each completions as completion}
					<li>
						<button
							class="text-gray-500 hover:text-gray-700 text-left"
							on:click={() => {
								closeCompletions = true;
								location.set(completion);
							}}>{completion}</button
						>
					</li>
				{/each}
			</ul>
		{/if}
	</div>
	<div class="flex items-center justify-between space-x-2">
		<input class="w-4/5" placeholder="Radius" type="range" min="1" max="10" bind:value={$radius} />
		<span class="font-medium text-gray-600">{$radius} km</span>
	</div>
</div>
