<script lang="ts">
	import { email as emailStore } from "../stores";

	const SUBSCRIPTION_ID = "subscription";
	const STORE_KEY = "STORE:EMAIL";

	let email: string;

	const saved = localStorage.getItem(STORE_KEY);
	if (saved) {
		email = saved;
		emailStore.set(email);
	}

	emailStore.subscribe((value) => localStorage.setItem(STORE_KEY, value));
</script>

<nav class="p-6 bg-red-400">
	<div class="mx-auto w-4/5 flex items-center justify-between">
		<h1 class="text-white font-medium uppercase text-lg">Roomiez</h1>
		<ul class="flex space-x-4">
			<li>
				<div class="flex flex-col space-y-2">
					<label for={SUBSCRIPTION_ID} class="text-white font-medium text-sm">Get personalized rentals delivered to your inbox every week.</label>
					<div class="flex justify-between space-x-2">
						<input
							id={SUBSCRIPTION_ID}
							class="w-4/5 outline-none px-2 py-1 border border-gray-300 rounded-md text-gray-600"
							placeholder="johndoe@mail.com"
							type="email"
							bind:value={email}
						/>
						<button class="w-1/5 py-1 px-2 rounded-md bg-gray-100 text-gray-500 font-medium hover:bg-gray-200" on:click={() => emailStore.set(email)}>Join</button>
					</div>
				</div>
			</li>
		</ul>
	</div>
</nav>
