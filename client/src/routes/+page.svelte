<script lang="ts">
	
	export let form;
	let longURL = form?.longURL || '';
	let shortURL = form?.shortURL || '';
</script>

<div>
	<!-- TOP BAR -->
	<div class="mt-2 ml-3 mr-3">
		<!-- Parent container for body -->

		<div class="flex flex-col justify-between items-center py-20 md:flex-row md:ml-40">
			<div class="p-28 px-10">
				<span class="font-sans text-2xl font-semibold"> Paste your log URL here:</span>
				<form action="/" method="POST">
					<div class="flex flex-col">
						<div class="flex items-center border-2 border-grey-800 rounded-3xl mt-8">
							<input
								class="outline-none flex-grow p-2 rounded-3xl"
								name="longurl"
								placeholder="Enter you link here...."
								bind:value={longURL}
							/>
							<div>
								<button class="px-8 py-2 text-white font-sans font-bold bg-cyan-400 rounded-3xl"
									>Shorten
								</button>
							</div>
						</div>
					</div>
				</form>

				{#if form?.error}
					<div class="p-2 text-red-600">{form?.error}</div>
				{/if}
				<div class="flex items-center mt-8">
					<div class="ml-2">Your Short URL:</div>
					<div class="flex items-center border-2 ml-4 border-grey-800 rounded-3xl">
						<input
							class="ml-3 p-1 outline-none"
							placeholder="your short link"
							bind:value={shortURL}
						/>
						<button class="px-6 p-2 rounded-3xl bg-gray-500 font-bold text-white">Copy</button>
					</div>
				</div>
			</div>

			<div class="px-5 ">
				<img class="" src="src/assets/images/illustration-working.svg" alt="" />
			</div>
		</div>
	</div>
</div>
