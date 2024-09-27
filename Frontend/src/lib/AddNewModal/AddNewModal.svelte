<script lang="ts">
	import type { SvelteComponent } from 'svelte';
	import { getModalStore, getToastStore, type ToastSettings } from '@skeletonlabs/skeleton';
	import { createEntity } from './AddNewModal';

	export let parent: SvelteComponent;

	const modalStore = getModalStore();
	const toastStore = getToastStore();

	const formData = {
		category: 'item',
		name: '',
		address: '',
		notes: ''
	};

	const entities = [
		{ value: 'item', display: 'Item' },
		{ value: 'container', display: 'Container' },
		{ value: 'shelf', display: 'Shelf' },
		{ value: 'shelvingunit', display: 'Shevling Unit' },
		{ value: 'room', display: 'Room' },
		{ value: 'building', display: 'Building' }
	];

	// Base Classes
	const cBase = 'card p-4 w-modal shadow-xl space-y-4 max-h-screen overflow-scroll';
	const cHeader = 'text-2xl font-bold';
	const cForm = 'border border-surface-500 p-4 space-y-4 rounded-container-token';

	var isFormInvalid = true;
	var formError = {
		name: ''
	};
	var formErrorClass = {
		name: ''
	};

	function validateForm() {
		if(formData.name == '') {
			isFormInvalid = true;
			formError.name = 'Name is required!'
			formErrorClass.name = 'input-error'
		}
		else {
			isFormInvalid = false;
			formError.name = ''
			formErrorClass.name = ''
		}

	}

	async function onFormSubmit() {
		modalStore.close();

		const [message, _] = await createEntity(formData)


		let toastMessage = '';
		let toastBackground = 'variant-filled-secondary';

		if (message == 'success') {
			toastMessage = 'Successfully added!';
		} else {
			toastMessage = 'There was an issue adding your item.'
			toastBackground = 'variant-filled-error';
		}

		const t: ToastSettings = {
			message: toastMessage,
			background: toastBackground,
			timeout: 5000
		};

		toastStore.trigger(t);
	}
</script>

{#if $modalStore[0]}
	<div class="modal-add-entity-form {cBase}">
		<header class={cHeader}>{$modalStore[0].title ?? '(title missing)'}</header>
		<article>{$modalStore[0].body ?? '(body missing)'}</article>
		<!-- Enable for debugging: -->
		<form class="modal-form {cForm}">
			<label class="label" for="category">Category:</label>
			<select id="category" class="select" bind:value={formData.category}>
				{#each entities as entity}
					<option value={entity.value}>{entity.display}</option>
				{/each}
			</select>

			<label for="name" class="label">Name:</label>
			<input
				id="name"
				class="input {formErrorClass.name}"
				type="text"
				bind:value={formData.name}
				on:input={validateForm}
				on:focusout={validateForm}
				placeholder="Enter name..."
			/>
			{#if formError.name}
				<p class="text-red-500 !mt-0">{formError.name}</p>
			{/if}

			{#if formData.category == 'building'}
				<label for="address" class="label">Address:</label>
				<input
					id="address"
					class="input"
					type="tel"
					bind:value={formData.address}
					placeholder="Enter address..."
				/>
			{/if}

			<label for="notes" class="label">Notes:</label>
			<textarea
				id="notes"
				class="textarea"
				rows="4"
				placeholder="Notes..."
				bind:value={formData.notes}
			/>
		</form>
		<!-- prettier-ignore -->
		<footer class="modal-footer {parent.regionFooter}">
			<button class="btn {parent.buttonNeutral}" on:click={parent.onClose}>{parent.buttonTextCancel}</button>
			<button class="btn {parent.buttonPositive}" disabled={isFormInvalid} on:click={onFormSubmit}>{parent.buttonTextSubmit}</button>
		</footer>
	</div>
{/if}
