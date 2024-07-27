<script lang="ts">
	import Modal from '$lib/Modal.svelte';
	import type { Link } from '../../app';
	import type { PageData } from './$types';
	import { writable } from 'svelte/store';

	export let data: PageData;
	let showModal = writable(false);
	let editLink: Link | null = null;
	let editLinkOriginalUrl = writable('');

	function openModal(link: Link) {
		editLink = link;
		editLinkOriginalUrl.set(link.original_url);
		showModal.set(true);
	}

	function closeModal() {
		showModal.set(false);
		editLink = null;
	}

	let openLinks = writable(new Map<string, boolean>());

	function toggleOptions(link: Link) {
		openLinks.update((current) => {
			const newMap = new Map(current);
			if (newMap.has(link.short_url)) {
				newMap.set(link.short_url, !newMap.get(link.short_url));
			} else {
				newMap.set(link.short_url, true);
			}
			return newMap;
		});
	}

	function handleToggle(link: Link, event: Event) {
		event.preventDefault();

		const form = event.target as HTMLFormElement;
		const isEnabledInput = form.querySelector('input[name="isEnabled"]') as HTMLInputElement;

		isEnabledInput.value = (!link.is_enabled).toString();


        
		form.submit();
	}
</script>

{#if data.links}
	<div>
		<div class="p-10 pb-0 flex justify-between md:pl-48 md:pr-44">
			<div class="p-2 text-center font-sans font-bold">SHORT URL</div>
			<div class="p-2 text-center font-sans font-bold">ORIGINAL URL</div>
			<div class="flex items-center">
				<div class="p-2 pr-0 text-center font-sans font-bold">CLICKS</div>
				<div class="p-2 pr-0 text-center font-sans font-bold">OPTIONS</div>
			</div>
		</div>
	</div>

	<div>
		{#each data.links as link}
			<div
				class="link-item pl-10 pr-10 pt-2 mt-3 flex flex-row justify-between items-center md:pl-48 md:pr-48"
			>
				<div
					class="p-3 bg-gray-100 border flex-grow flex justify-between items-center rounded-lg mr-2 w-[45%] max-w-[45%] min-w-[45%] overflow-hidden text-ellipsis whitespace-nowrap"
				>
					<div class="truncate">{link.short_url}</div>
					<img src="src/assets/icons/copy.svg" alt="copy" class="w-6 h-6 cursor-pointer" />
				</div>

				<div
					class="p-3 bg-gray-100 border flex-grow flex justify-between items-center rounded-lg w-[45%] max-w-[45%] min-w-[45%] overflow-hidden text-ellipsis whitespace-nowrap"
				>
					<div class="truncate">{link.original_url}</div>
					<img src="src/assets/icons/copy.svg" alt="copy" class="w-6 h-6 cursor-pointer" />
				</div>

				<div class="flex items-center font-mono ml-4 mt-2 md:mt-0">
					<div
						class="block p-3 px-5 mr-2 ml-10 bg-gray-100 text-center min-w-[3rem] border rounded-lg font-serif"
					>
						{link.clicks}
					</div>
					<button
						type="button"
						class="p-3 bg-gray-100 border rounded-lg flex-shrink-0"
						on:click={() => toggleOptions(link)}
						on:keydown={(event) => {
							if (event.key === 'Enter' || event.key === ' ') {
								event.preventDefault();
								toggleOptions(link);
							}
						}}
						aria-expanded={$openLinks.get(link.short_url)}
						aria-label="Options"
					>
						{#if $openLinks.get(link.short_url)}
							<svg
								height="20px"
								aria-hidden="true"
								focusable="false"
								data-prefix="fas"
								data-icon="window-close"
								class="w-6 h-6 cursor-pointer"
								role="img"
								xmlns="http://www.w3.org/2000/svg"
								viewBox="0 0 512 512"
							>
								<path
									fill="black"
									d="M464 32H48C21.5 32 0 53.5 0 80v352c0 26.5 21.5 48 48 48h416c26.5 0 48-21.5 48-48V80c0-26.5-21.5-48-48-48zm-83.6 290.5c4.8 4.8 4.8 12.6 0 17.4l-40.5 40.5c-4.8 4.8-12.6 4.8-17.4 0L256 313.3l-66.5 67.1c-4.8 4.8-12.6 4.8-17.4 0l-40.5-40.5c-4.8-4.8-4.8-12.6 0-17.4l67.1-66.5-67.1-66.5c-4.8-4.8-4.8-12.6 0-17.4l40.5-40.5c4.8-4.8 12.6-4.8 17.4 0l66.5 67.1 66.5-67.1c4.8-4.8 12.6-4.8 17.4 0l40.5 40.5c4.8 4.8 4.8 12.6 0 17.4L313.3 256l67.1 66.5z"
								></path>
							</svg>
						{:else}
							<img
								src="src/assets/icons/option.svg"
								alt="settings"
								class="w-6 h-6 cursor-pointer"
							/>
						{/if}
					</button>
				</div>
			</div>

			{#if $openLinks.get(link.short_url)}
				<div
					class="link-item pl-10 pr-10 pt-2 mt-3 flex flex-col md:flex-row md:justify-between items-center md:pl-56 md:pr-56 font-sans"
				>
					<!-- Your options go here -->
					<button
						class="bg-black text-white py-3 px-4 rounded-md mr-2 mb-2 flex items-center"
						on:click={() => openModal(link)}
					>
						<img class="cursor-pointer" src="src/assets/icons/edit.svg" alt="QR Code" />
						<span class="pl-4">Edit Link</span>
					</button>
					<button class="bg-black text-white py-3 px-4 rounded-md mr-2 mb-2">Analytics</button>
					<!-- <button class="bg-black text-white py-3 px-4 rounded-md mr-2 mb-2 flex items-center">
						<img class="cursor-pointer" src="src/assets/icons/qr.svg" alt="QR Code" />
						<span class="pl-2">QR Code</span>
					</button> -->
					<form
						method="POST"
						action="?/enableDisableLink"
						on:submit|preventDefault={(event) => handleToggle(link, event)}
					>
						<input type="hidden" name="short_url" value={link?.short_url} />

						<input type="hidden" name="isEnabled" value={link.is_enabled} />
						{#if link.is_enabled}
							<button
								class="bg-black font-bold text-red-700 py-2 px-4 rounded-md mr-2 mb-2"
								type="submit"
							>
								Disable Link
							</button>
						{:else}
							<button
								class="bg-black font-bold text-green-500 py-2 px-4 rounded-md mr-2 mb-2"
								type="submit"
							>
								Enable Link
							</button>
						{/if}
					</form>
				</div>
			{/if}
		{/each}
	</div>
{:else}
	<p>No links available.</p>
{/if}

{#if $showModal}
	<Modal show={$showModal} onClose={closeModal}>
		<form method="POST" action="?/updateLink" >
			<div class="flex flex-row">
				<input type="hidden" name="link_id" value={editLink?.id} />
				<input type="hidden" name="short_url" value={editLink?.short_url} />
				<input
					class="ml-3 p-1 outline-none"
					name="original_url"
					bind:value={$editLinkOriginalUrl}
				/>
				<div>
					<button class="bg-blue-500 text-white py-2 px-4 rounded-md mt-2" type="submit" on:submit|preventDefault>
						Update
					</button>
					<button on:click={closeModal} role="dialog" aria-modal="true"> Close </button>
				</div>
			</div>
		</form>
	</Modal>
{/if}
