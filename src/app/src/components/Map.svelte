<script lang="ts">
	import { onMount } from "svelte";
	import { mapReady } from "../stores";
	import type { SearchResult } from "$lib/api";
	import mapStyles from "$lib/mapStyles";

	export let rentals: SearchResult[];
	const coordinates = rentals.reduce<{ coords: { lat: number; lng: number }; id: number }[]>((prev, rental) => {
		if (rental.coordinates) {
			const coordinates = rental.coordinates.replace("POINT(", "").replace(")", "").split(" ");
			const lng = parseFloat(coordinates[0]);
			const lat = parseFloat(coordinates[1]);

			prev.push({ coords: { lat, lng }, id: rental.id });
		}

		return prev;
	}, []);

	let container: any;
	let map: any = null;

	onMount(() =>
		mapReady.subscribe((value) => {
			if (value && !map)
				// @ts-ignore
				map = new google.maps.Map(container, {
					zoom: 10,
					center: { lat: -33.8688, lng: 151.2093 },
					disableDefaultUI: true,
					styles: mapStyles
				});

			coordinates.forEach((coordinate) => {
				// @ts-ignore
				new google.maps.Marker({
					map,
					position: coordinate.coords,
					title: coordinate.id.toString()
				});
			});
		})
	);
</script>

<div class="full-screen" bind:this={container} />

<style>
	.full-screen {
		width: 100%;
		height: 400px;
		border-radius: 0.375rem;
	}
</style>
