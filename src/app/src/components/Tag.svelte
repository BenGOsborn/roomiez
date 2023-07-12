<script lang="ts">
	export let text: string;

	const colors = ["red", "rose", "gray"] as const;
	type Colors = (typeof colors)[number];

	function hash(text: string): Colors {
		let integerRepresentation = 0;
		for (let i = 0; i < text.length; i++) {
			const charCode = text.charCodeAt(i);
			integerRepresentation += charCode;
		}

		return colors[integerRepresentation % colors.length];
	}

	const color = hash(text);
</script>

{#if color === "red"}
	<span class="py-1 px-2 rounded-md bg-red-200 text-red-400 font-medium">{text}</span>
{:else if color === "rose"}
	<span class="py-1 px-2 rounded-md bg-rose-300 text-rose-500 font-medium">{text}</span>
{:else}
	<span class="py-1 px-2 rounded-md bg-gray-300 text-gray-500 font-medium">{text}</span>
{/if}
