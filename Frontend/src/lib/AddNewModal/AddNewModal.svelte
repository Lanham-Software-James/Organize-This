<script lang="ts">
	import { onMount, type SvelteComponent } from 'svelte';
	import { getModalStore, getToastStore, type ToastSettings } from '@skeletonlabs/skeleton';
	import { createEntity, editEntity, getEntity, getParents, type parentData } from './AddNewModal';

	export let parent: SvelteComponent;
	export let edit: boolean;
	export let id: number;
	export let category: string;

	const modalStore = getModalStore();
	const toastStore = getToastStore();

	const formData = {
		id: 0,
		category: 'item',
		name: '',
		address: '',
		notes: '',
		parent: ''
	};

	const entities = [
		{ value: 'item', display: 'Item' },
		{ value: 'container', display: 'Container' },
		{ value: 'shelf', display: 'Shelf' },
		{ value: 'shelving_unit', display: 'Shelving Unit' },
		{ value: 'room', display: 'Room' },
		{ value: 'building', display: 'Building' }
	];

	// Base Classes
	const cBase = 'card p-4 w-modal shadow-xl space-y-4 max-h-screen overflow-scroll';
	const cHeader = 'text-2xl font-bold';
	const cForm = 'border border-surface-500 p-4 space-y-4 rounded-container-token';

	var isFormInvalid = true && !edit;
	var formError = {
		name: ''
	};
	var formErrorClass = {
		name: ''
	};

	var parents: parentData[] = [];

	onMount(async function () {
		if (edit) {
			var [[parentsMessage, parentsData], [entityMessage, entityData]] = await Promise.all([getParents(category), getEntity(id, category)]);
			if (parentsMessage == 'success') {
				parents = parentsData;
			}

			if (entityMessage == 'success') {
				formData.id = entityData.Entity.ID
				formData.category = category
				formData.name = entityData.Entity.Name
				formData.parent = entityData.Parent.ParentID + '-' + entityData.Parent.ParentCategory
				formData.address = entityData.Address || ''
				formData.notes = entityData.Entity.Notes || ''
			}
		} else {
			var [message, data] = await getParents(formData.category);
			if (message == 'success') {
				parents = data;
			}
		}
	});

	async function updateParents() {
		if(formData.category != 'building') {
			var [message, data] = await getParents(formData.category);
			if (message == 'success') {
				parents = data;
				formData.parent = '0-zero';
			}
		}

	}

	function validateForm() {
		if (formData.name == '') {
			isFormInvalid = true;
			formError.name = 'Name is required!';
			formErrorClass.name = 'input-error';
		} else {
			isFormInvalid = false;
			formError.name = '';
			formErrorClass.name = '';
		}
	}

	async function onFormSubmit() {
		modalStore.close();

		const addEdit = edit ? 'edit' : 'add'

		var message = ''
		if(edit){
			[message, ] = await editEntity(formData);
		} else {
			[message, ] = await createEntity(formData);
		}

		let toastMessage = '';
		let toastBackground = 'variant-filled-secondary';

		if (message == 'success') {
			toastMessage = `Successfully ${addEdit}ed!`;
		} else {
			toastMessage = `There was an issue ${addEdit}ing your item.`;
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
		<form class="modal-form {cForm}">
			<label class="label" for="category">Category:</label>
			<select id="category" class="select"  disabled={edit} bind:value={formData.category} on:change={updateParents}>
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

			{#if formData.category != 'building'}
				<label class="label" for="parents">Parent:</label>
				<select id="parents" class="select" bind:value={formData.parent}>
					{#each parents as parent}
						<option value={parent.ID + '-' + parent.Category}>{parent.Name}</option>
					{/each}
				</select>
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
