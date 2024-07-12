<script lang="ts">
	import type { SignUpForm } from '../../app';
	import type { PageData } from './$types';
	let signUpForm: SignUpForm = {
		email: '',
		password: '',
		username: ''
	};
	let showPassword = false;

	let signup = () => {
		console.log(signUpForm);
	};

	// export let data: PageData;
</script>

<div class="flex min-h-full flex-col justify-center px-6 py-12 lg:px-8">
	<div class="sm:mx-auto sm:w-full sm:max-w-sm">
		<h2 class="mt-10 text-center text-2xl font-bold leading-9 tracking-tight text-gray-900">
			Create your account
		</h2>
	</div>

	<div class="mt-10 sm:mx-auto sm:w-full sm:max-w-sm">
		<form on:submit|preventDefault={signup} class="space-y-6">
			<div>
				<label for="username" class="block text-sm font-medium leading-6 text-gray-900"
					>User name</label
				>
				<div class="mt-2">
					<input
						id="username"
						name="username"
						type="text"
						autocomplete="username"
						required
						class="block w-full rounded-md border-0 py-1.5 px-2 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-cyan-500 sm:text-sm sm:leading-6"
						bind:value={signUpForm.username}
					/>
				</div>
			</div>
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
						bind:value={signUpForm.email}
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
							bind:value={signUpForm.password}
						/>
					{:else}
						<input
							id="password"
							name="password"
							type="password"
							autocomplete="current-password"
							required
							class="block w-full rounded-md border-0 py-1.5 px-2 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-cyan-500 sm:text-sm sm:leading-6"
							bind:value={signUpForm.password}
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
					>Sign Up</button
				>
			</div>
		</form>
	</div>
</div>
