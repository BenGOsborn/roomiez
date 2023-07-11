<script lang="ts">
	import { PUBLIC_API_ENDPOINT } from "$env/static/public";
	import { getFields } from "$lib/api";
</script>

{#await getFields(PUBLIC_API_ENDPOINT)}
	<p class="text-center font-medium text-gray-800">Loading fields...</p>
{:then fields}
	<div class="space-y-8">
		<div class="bg-white p-4 rounded-md flex flex-col space-y-4 drop-shadow">
			<p class="text-gray-800 font-bold">Location</p>
			<input class="outline-none px-2 py-1 border border-gray-300 rounded-md text-gray-600" placeholder="Location" type="text" />
			<input placeholder="Radius" type="range" min="0" max="10" />
		</div>
		<div class="bg-white p-4 rounded-md flex flex-col space-y-4 drop-shadow">
			<p class="text-gray-800 font-bold">Filters</p>
			<div class="grid grid-cols-2 gap-4">
				<div class="space-y-2">
					<p class="text-gray-800 font-medium text-sm">Age</p>
					<select>
						<option>Any</option>
						{#each fields.age as age}
							<option value={age}>{age}</option>
						{/each}
					</select>
				</div>
				<div class="space-y-2">
					<p class="text-gray-800 font-medium text-sm">Duration</p>
					<select>
						<option>Any</option>
						{#each fields.duration as duration}
							<option value={duration}>{duration}</option>
						{/each}
					</select>
				</div>
				<div class="space-y-2">
					<p class="text-gray-800 font-medium text-sm">Gender</p>
					<select>
						<option>Any</option>
						{#each fields.gender as gender}
							<option value={gender}>{gender}</option>
						{/each}
					</select>
				</div>
				<div class="space-y-2">
					<p class="text-gray-800 font-medium text-sm">Rental Type</p>
					<select>
						<option>Any</option>
						{#each fields.rentalType as rentalType}
							<option value={rentalType}>{rentalType}</option>
						{/each}
					</select>
				</div>
				<div class="space-y-2">
					<p class="text-gray-800 font-medium text-sm">Tenant</p>
					<select>
						<option>Any</option>
						{#each fields.tenant as tenant}
							<option value={tenant}>{tenant}</option>
						{/each}
					</select>
				</div>
			</div>
		</div>
		<div class="bg-white p-4 rounded-md flex flex-col space-y-4 drop-shadow">
			<p class="text-gray-800 font-bold">Features</p>
			<div class="grid grid-cols-2 gap-x-8 gap-y-4">
				{#each fields.feature as feature}
					<div class="flex items-center justify-between space-x-2">
						<span>{feature}</span>
						<input value={feature} type="checkbox" />
					</div>
				{/each}
			</div>
		</div>
	</div>
{/await}
