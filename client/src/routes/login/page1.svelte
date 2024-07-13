<script lang="ts">
	import { goto } from '$app/navigation';
	import { PUBLIC_BASE_URL } from '$env/static/public';

	let email: string = '';
	let password: string = '';
	let errorMessage = '';
	let showPassword = false;
	let signIn = async () => {
		const response = await fetch(`${PUBLIC_BASE_URL}/signin`, {
			method: 'POST',
			// credentials: 'include',
			headers: {
				Accept: 'application/json',
				'content-type': 'application/json'
			},
			body: JSON.stringify({
				email: email,
				password: password
			})
		});

		const data = await response.json();
		if (response.status <= 299) {
			console.log(data);
			await goto('/', { noScroll: false, replaceState: true });
		} else {
			errorMessage = data.message;
		}
	};
</script>

<div class="flex min-h-full flex-col justify-center px-6 py-12 lg:px-8">
	<div class="sm:mx-auto sm:w-full sm:max-w-sm">
		<h2 class="mt-10 text-center text-2xl font-bold leading-9 tracking-tight text-gray-900">
			Signin to your account
		</h2>
	</div>

	<div class="mt-10 sm:mx-auto sm:w-full sm:max-w-sm">
		<form on:submit|preventDefault={signIn} class="space-y-6">
			<div>
				<label for="email" class="block text-sm font-medium leading-6 text-gray-900"
					>Email address</label
				>
				<div class="mt-2">
					<input
						id="email"
						name="email"
						type="email"
						autocomplete="email"
						required
						class="block w-full rounded-md border-0 py-1.5 px-2 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-cyan-500 sm:text-sm sm:leading-6"
						bind:value={email}
					/>
				</div>
			</div>

			<div>
				<div class="flex items-center justify-between">
					<label for="password" class="block text-sm font-medium leading-6 text-gray-900"
						>Password</label
					>
				</div>
				<div class="mt-2">
					{#if showPassword}
						<input
							id="password"
							name="password"
							type="text"
							autocomplete="current-password"
							required
							class="block w-full rounded-md border-0 py-1.5 px-2 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-cyan-500 sm:text-sm sm:leading-6"
							bind:value={password}
						/>
					{:else}
						<input
							id="password"
							name="password"
							type="password"
							autocomplete="current-password"
							required
							class="block w-full rounded-md border-0 py-1.5 px-2 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-cyan-500 sm:text-sm sm:leading-6"
							bind:value={password}
						/>
					{/if}
				</div>
				<div class="mt-2">
					<label>
						<input type="checkbox" bind:checked={showPassword} />
						Show Password
					</label>
				</div>
			</div>

			<div>
				<button
					type="submit"
					class="flex w-full justify-center rounded-md bg-cyan-400 px-3 py-1.5 text-sm font-semibold leading-6 text-white shadow-sm hover:bg-cyan-600 duration-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-cyan-500"
					>Sign in</button
				>
			</div>
			{#if errorMessage}
				<!-- Display error message -->
				<div class="mb-4 text-red-600">{errorMessage}</div>
			{/if}
		</form>
	</div>
</div>
