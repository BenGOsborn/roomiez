<script lang="ts">
	import { z } from "zod";
	import { location, radius } from "../../stores";

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

	location.subscribe((loc) => localStorage.setItem(STORE_KEY, JSON.stringify({ location: loc, radius: $radius } as State)));
	radius.subscribe((rad) => localStorage.setItem(STORE_KEY, JSON.stringify({ location: $location, radius: rad } as State)));
</script>

<div class="bg-white p-4 rounded-md flex flex-col space-y-4 drop-shadow">
	<p class="text-gray-800 font-bold">Location</p>
	<input class="outline-none px-2 py-1 border border-gray-300 rounded-md text-gray-600" placeholder="Location" type="text" bind:value={$location} />
	<div class="flex items-center justify-between space-x-2">
		<input class="w-4/5" placeholder="Radius" type="range" min="1" max="10" bind:value={$radius} />
		<span class="font-medium text-gray-600">{$radius} km</span>
	</div>
</div>
