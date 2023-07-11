<script lang="ts">
	export let features: string[];
</script>

<div class="bg-white p-4 rounded-md flex flex-col space-y-4 drop-shadow">
	<p class="text-gray-800 font-bold">Features</p>
	<div class="grid grid-cols-2 gap-x-8 gap-y-4">
		{#each features as feature}
			<div class="flex items-center justify-between space-x-2">
				<span class="text-gray-700">{feature}</span>
				<input value={feature} type="checkbox" />
			</div>
		{/each}
	</div>
</div>
