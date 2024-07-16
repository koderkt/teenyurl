<script lang="ts">
	import type { PageData } from './$types';
		export let data: PageData;
</script>

<div>
	{#if data.links}
		<div>
			<div class="p-10 pb-0 flex justify-between md:pl-48 md:pr-48">
				<div class="p-2 text-center font-sans font-bold">SHORT URL</div>
				<div class="flex items-center">
					<div class="p-2 text-center font-sans font-bold">CLICKS</div>
					<div class="p-2 text-center font-sans font-bold">OPTIONS</div>
				</div>
			</div>
		</div>

		<div>
			{#each data.links as link}
				<div
					class="link-item pl-10 pr-10 pt-2 mt-3 flex justify-between items-center md:pl-48 md:pr-48"
				>
					<div
						class="p-3 bg-gray-200 border flex-grow flex justify-between items-center rounded-lg"
					>
						<div class="truncate">
							{link.short_url}
						</div>
						<img src="src/assets/icons/copy.svg" alt="copy" class="w-6 h-6 cursor-pointer" />
					</div>
					<div class="flex items-center font-mono ml-4">
						<div
							class="block p-3 px-7 mr-2 ml-10 bg-gray-200 text-center min-w-[3rem] border rounded-lg"
						>
							{link.clicks}
						</div>
						<div class="p-3 bg-gray-200 border rounded-lg">
							<img
								src="src/assets/icons/option.svg"
								alt="settings"
								class="w-6 h-6 cursor-pointer"
							/>
						</div>
					</div>
				</div>
			{/each}
		</div>
	{:else}
		<p>No links available.</p>
	{/if}
</div>
