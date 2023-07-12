<script lang="ts">
	import { page } from "../stores";
</script>

<div class="rounded-md drop-shadow p-4 bg-white flex items-center justify-evenly space-x-2">
	<button class="text-gray-500 font-medium hover:text-gray-600" on:click={() => page.update((prev) => Math.max(prev - 1, 1))}>Prev</button>
	<p class="text-gray-700 font-medium">{$page.toLocaleString()}</p>
	<button class="text-gray-500 font-medium hover:text-gray-600" on:click={() => page.update((prev) => prev + 1)}>Next</button>
</div>
