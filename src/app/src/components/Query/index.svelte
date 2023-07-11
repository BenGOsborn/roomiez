<script lang="ts">
	import { PUBLIC_API_ENDPOINT } from "$env/static/public";
	import { getFields } from "$lib/api";
	import Location from "./Location.svelte";
	import Features from "./Features.svelte";
	import Filter from "./Filter.svelte";
	import { age, duration, gender, rentalType, tenant } from "../../stores";
</script>

{#await getFields(PUBLIC_API_ENDPOINT)}
	<p class="text-center font-medium text-gray-800">Loading fields...</p>
{:then fields}
	<div class="space-y-8">
		<Location />
		<div class="bg-white p-4 rounded-md flex flex-col space-y-4 drop-shadow">
			<p class="text-gray-800 font-bold">Filters</p>
			<div class="grid grid-cols-2 gap-4">
				<Filter name="Age" values={fields.age} store={age} />
				<Filter name="Duration" values={fields.duration} store={duration} />
				<Filter name="Gender" values={fields.gender} store={gender} />
				<Filter name="Rental Type" values={fields.rentalType} store={rentalType} />
				<Filter name="Tenant" values={fields.tenant} store={tenant} />
			</div>
		</div>
		<Features features={fields.feature} />
	</div>
{/await}
