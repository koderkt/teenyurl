<script lang="ts">
	

    
	export let show = false;
	export let onClose = () => {};
	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Enter' || event.key === ' ') {
			onClose();
		}
	}
</script>

{#if show}
	<div class="modal" on:click={onClose} on:keydown={handleKeydown} role="dialog" aria-modal="true">
		<div class="modal-content" on:click|stopPropagation>
			<slot></slot>
		
		</div>
	</div>
{/if}

<style>
	.modal {
		position: fixed;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		background: rgba(0, 0, 0, 0.5);
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.modal-content {
		background: white;
		padding: 20px;
		border-radius: 5px;
		max-width: 500px;
		width: 100%;
	}
</style>
