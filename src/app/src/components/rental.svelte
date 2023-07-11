<script lang="ts">
	import type { SearchResult } from "$lib/api";
	import Tag from "./tag.svelte";

	export let rental: SearchResult;
</script>

<div class="bg-white p-4 rounded-md flex flex-col space-y-2 drop-shadow">
	<a class="text-blue-700 font-medium underline" href={rental.url} target="_blank" rel="nofollow">View original post</a>
	<div class="space-y-4">
		<p class="text-gray-800 font-bold text-lg">About</p>
		<div class="grid grid-cols-2 gap-2">
			{#if rental.location}
				<p class="text-gray-500"><span class="text-gray-700 font-medium">Location:</span> {rental.location}</p>
			{/if}
			{#if rental.price}
				<p class="text-gray-500"><span class="text-gray-700 font-medium">Price:</span> ${rental.price} AUD / week</p>
			{/if}
			{#if rental.bond}
				<p class="text-gray-500"><span class="text-gray-700 font-medium">Bond:</span> ${rental.bond} AUD</p>
			{/if}
			{#if rental.rentalType}
				<p class="text-gray-500"><span class="text-gray-700 font-medium">Rental type:</span> {rental.rentalType}</p>
			{/if}
			{#if rental.gender}
				<p class="text-gray-500"><span class="text-gray-700 font-medium">Gender preference:</span> {rental.gender}</p>
			{/if}
			{#if rental.age}
				<p class="text-gray-500"><span class="text-gray-700 font-medium">Age preference:</span> {rental.age}</p>
			{/if}
			{#if rental.duration}
				<p class="text-gray-500"><span class="text-gray-700 font-medium">Duration:</span> {rental.duration}</p>
			{/if}
			{#if rental.tenant}
				<p class="text-gray-500"><span class="text-gray-700 font-medium">Tenant:</span> {rental.tenant}</p>
			{/if}
		</div>
	</div>
	{#if rental.features}
		<div class="space-y-4">
			<p class="text-gray-800 font-bold">Features</p>
			<div class="space-x-2 flex items-center">
				{#each rental.features as feature}
					<Tag text={feature} />
				{/each}
			</div>
		</div>
	{/if}
</div>
