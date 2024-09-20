<script lang="ts">
	import type { SvelteComponent } from 'svelte';
	import { getModalStore, getToastStore, type ToastSettings } from '@skeletonlabs/skeleton';
	import { PUBLIC_API_URL } from '$env/static/public';

	export let parent: SvelteComponent;

	const modalStore = getModalStore();
	const toastStore = getToastStore();

	const formData = {
		Category: 'item',
		Name: '',
		Address: '',
		Notes: ''
	};

	async function onFormSubmit() {
		modalStore.close();

		const response = await fetch(PUBLIC_API_URL + 'v1/entity-management/' + formData.Category, {
			method: 'POST',
			body: JSON.stringify({
				Address: formData.Address,
				Name: formData.Name,
				Notes: formData.Notes
			})
		});

		response.json().then((res) => {
			let toastMessage = '';
			let toastBackground = 'variant-filled-secondary';
			if (res.message == 'success') {
				toastMessage = 'Successfully added!';
			} else {
				toastMessage = 'There was an issue adding your item. ' + res.data;
				toastBackground = 'variant-filled-error';
			}

			const t: ToastSettings = {
				message: toastMessage,
				background: toastBackground,
				timeout: 5000
			};
			toastStore.trigger(t);
		});
	}

	const entities = [
		{ value: 'item', display: 'Item' },
		{ value: 'container', display: 'Container' },
		{ value: 'shelf', display: 'Shelf' },
		{ value: 'unit', display: 'Shevling Unit' },
		{ value: 'room', display: 'Room' },
		{ value: 'building', display: 'Building' }
	];

	// Base Classes
	const cBase = 'card p-4 w-modal shadow-xl space-y-4';
	const cHeader = 'text-2xl font-bold';
	const cForm = 'border border-surface-500 p-4 space-y-4 rounded-container-token';
</script>

{#if $modalStore[0]}
	<div class="modal-example-form {cBase}">
		<header class={cHeader}>{$modalStore[0].title ?? '(title missing)'}</header>
		<article>{$modalStore[0].body ?? '(body missing)'}</article>
		<!-- Enable for debugging: -->
		<form class="modal-form {cForm}">
			<label class="label" for="category">Category:</label>
			<select id="category" class="select" bind:value={formData.Category}>
				{#each entities as entity}
					<option value={entity.value}>{entity.display}</option>
				{/each}
			</select>

			<label for="name" class="label">Name:</label>
			<input
				id="name"
				class="input"
				type="text"
				bind:value={formData.Name}
				placeholder="Enter name..."
			/>

			{#if formData.Category == 'building'}
				<label for="address" class="label">Address:</label>
				<input
					id="address"
					class="input"
					type="tel"
					bind:value={formData.Address}
					placeholder="Enter address..."
				/>
			{/if}

			<label for="notes" class="label">Notes:</label>
			<textarea
				id="notes"
				class="textarea"
				rows="4"
				placeholder="Notes..."
				bind:value={formData.Notes}
			/>
		</form>
		<!-- prettier-ignore -->
		<footer class="modal-footer {parent.regionFooter}">
			<button class="btn {parent.buttonNeutral}" on:click={parent.onClose}>{parent.buttonTextCancel}</button>
			<button class="btn {parent.buttonPositive}" on:click={onFormSubmit}>{parent.buttonTextSubmit}</button>
		</footer>
	</div>
{/if}
