<script lang="ts">
	import type { Fields } from "$lib/api";
	import Filter from "./Filter.svelte";
	import { age, duration, gender, rentalType, tenant } from "../../stores";

	export let fields: Fields;
</script>

<div class="bg-white p-4 rounded-md flex flex-col space-y-4 drop-shadow">
	<p class="text-gray-800 font-bold">Filters</p>
	<div class="grid grid-cols-1 xl:grid-cols-2 gap-4">
		<Filter name="Age" values={fields.age} store={age} />
		<Filter name="Duration" values={fields.duration} store={duration} />
		<Filter name="Gender" values={fields.gender} store={gender} />
		<Filter name="Rental Type" values={fields.rentalType} store={rentalType} />
		<Filter name="Tenant" values={fields.tenant} store={tenant} />
	</div>
</div>
