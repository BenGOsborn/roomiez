<script lang="ts">
	import { env } from "$env/dynamic/public";

	let x = 0;
</script>

<h1>Welcome to SvelteKit</h1>
<p>{x}</p>
<p>{env.PUBLIC_API_ENDPOINT}</p>
<button on:click={() => (x = x + 1)}>Increment</button>
