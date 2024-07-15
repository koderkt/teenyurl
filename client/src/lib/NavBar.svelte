<script lang="ts">
	export let data;
	import { goto } from '$app/navigation';
	let isLoggedIn = false;
	const logout = async () => {
		await fetch('http://localhost:3030/api/logout', {
			method: 'POST',
			credentials: 'include',
			headers: {
				Accept: 'application/json',
				'content-type': 'application/json'
			}
		});
		// user.update(val => val = null);
		await goto('/', { noScroll: false, replaceState: true });
	};
</script>

<nav class="flex mt-2 ml-3 mr-3 justify-between">
	<a href="/" class="font-sans text-2xl font-[750] p-2">TEENYURL</a>
	{#if data.sessionId == null}
		<div class="flex item-center font-bold">
			<div class="py-4">
				<a class="p-2 px-4 text-gray-700" href="/login">Login</a>
			</div>
			<div class=" py-4">
				<a class="p-2 px-4 bg-cyan-400 rounded-3xl text-white" href="/signup">Sign Up</a>
			</div>
		</div>
	{:else}
		<div class="flex font-sans items-center font-bold">
			<div class="py-4">
				<a class="p-2 px-4 text-gray-700" href="/links">My Links</a>
			</div>
			<form action="/logout" method="POST">
				<div class="py-4">
					<button class="p-2 px-4 bg-cyan-400 rounded-3xl text-white">Logout</button>
				</div>
			</form>
		</div>
	{/if}
</nav>
