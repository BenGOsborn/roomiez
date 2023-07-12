<script lang="ts">
	import { z } from "zod";
	import { features as featuresStore } from "../../stores";

	export let features: string[];

	const STORE_KEY = "STORE:FEATURES";

	featuresStore.set(
		features.reduce((prev: { [key: string]: boolean }, feature) => {
			prev[feature] = false;
			return prev;
		}, {})
	);

	const saved = localStorage.getItem(STORE_KEY);
	if (saved) {
		const state = z.record(z.boolean()).parse(JSON.parse(saved));

		featuresStore.update((map) => {
			Object.entries(state).forEach(([key, value]) => {
				map[key] = value;
			});

			return map;
		});
	}

	featuresStore.subscribe((value) => localStorage.setItem(STORE_KEY, JSON.stringify(value)));
</script>

<div class="bg-white p-4 rounded-md flex flex-col space-y-4 drop-shadow">
	<p class="text-gray-800 font-bold">Features</p>
	<div class="grid grid-cols-1 xl:grid-cols-2 gap-x-8 gap-y-4">
		{#each features as feature}
			<div class="flex items-center justify-between space-x-2">
				<label class="text-gray-700" for={`feature:${feature}`}>{feature}</label>
				<input
					id={`feature:${feature}`}
					value={feature}
					type="checkbox"
					checked={$featuresStore[feature]}
					on:change={() =>
						featuresStore.update((map) => {
							map[feature] = !map[feature];
							return map;
						})}
				/>
			</div>
		{/each}
	</div>
</div>
