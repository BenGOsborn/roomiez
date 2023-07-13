<script lang="ts">
	import { PUBLIC_GA_ID, PUBLIC_MAPS_KEY } from "$env/static/public";
	import "../app.css";
	import Nav from "../components/Nav.svelte";
	import { mapReady } from "../stores";

	// @ts-ignore
	window.initMap = () => {
		mapReady.set(true);
	};

	// @ts-ignore
	window.dataLayer = window.dataLayer || [];

	function gtag() {
		// @ts-ignore
		dataLayer.push(arguments);
	}

	// @ts-ignore
	gtag("js", new Date());
	// @ts-ignore
	gtag("config", PUBLIC_GA_ID);
</script>

<svelte:head>
	<script defer async src={`https://maps.googleapis.com/maps/api/js?key=${PUBLIC_MAPS_KEY}&libraries=places&callback=initMap`}>
	</script>
	<script async src={`https://www.googletagmanager.com/gtag/js?id=${PUBLIC_GA_ID}`}></script>

	<title>Roomiesz - Find A Shared Rental</title>
	<meta
		name="description"
		content="Easily explore shared rentals on Facebook groups, customized to your preferred location and budget. Find your ideal living space with our welcoming platform, connecting you to the perfect accommodation."
	/>
</svelte:head>

<Nav />

<slot />
