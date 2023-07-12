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

	<title>Roomiesz - Find A Shared Rental</title>
	<meta
		name="description"
		content="Easily explore shared rentals on Facebook groups, customized to your preferred location and budget. Find your ideal living space with our welcoming platform, connecting you to the perfect accommodation."
	/>
</svelte:head>

<Nav />

<slot />
