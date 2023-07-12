<script lang="ts">
	import { PUBLIC_MAPS_KEY } from "$env/static/public";
	import "../app.css";
	import Nav from "../components/Nav.svelte";
	import { mapReady } from "../stores";

	// @ts-ignore
	window.initMap = () => {
		mapReady.set(true);
	};
</script>

<svelte:head>
	<script defer async src={`https://maps.googleapis.com/maps/api/js?key=${PUBLIC_MAPS_KEY}&callback=initMap`}>
	</script>
</svelte:head>

<Nav />

<slot />
