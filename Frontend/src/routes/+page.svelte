<script lang="ts">
	import {
		Paginator,
		type PaginationSettings,
		getModalStore,
		popup,
		type ModalSettings,
		type PopupSettings
	} from '@skeletonlabs/skeleton';
	import { getContext, onMount, tick } from 'svelte';
	import { _getEntities as getEntities, type GetEntitiesData } from './+page';
	import AddNewModal from '$lib/AddNewModal/AddNewModal.svelte';
	import { slide } from 'svelte/transition';
	import QrCodeModal from '$lib/QRCodeModal/QRCodeModal.svelte';
	import { cleanCategory } from '$lib/CleanCategory/CleanCategory';
	import { goto } from '$app/navigation';
	import { onDestroy } from 'svelte';

	let entities: GetEntitiesData[] = [];
	let offset = 0;
	let limit = 15;
	let page = 0;
	let paginationSettings = {
		page: page,
		limit: limit,
		size: 0,
		amounts: [5, 10, 15, 20, 25]
	} satisfies PaginationSettings;

	const refreshPage = getContext('refreshPage');

	let width: number;
	let parentMax = 0;

	$: {
		if (typeof width !== 'undefined') {
			if (width <= 640) {
				parentMax = 2;
			} else if (width <= 768) {
				parentMax = 5;
			} else {
				parentMax = 15;
			}
		}
	}

	onMount(async function () {
		loadData();

		//@ts-ignore
		const unsubscribe = refreshPage.subscribe(() => {
			loadData();
		});

		return unsubscribe;
	});

	onDestroy(() => {
		if (searchTimeout) clearTimeout(searchTimeout);
	});

	async function loadData() {
		[entities, paginationSettings.size] = await getEntities(offset, limit, searchString, filters);
	}

	async function limitChange(e: CustomEvent) {
		limit = e.detail;

		[entities, paginationSettings.size] = await getEntities(offset, limit, searchString, filters);
	}

	async function pageChange(e: CustomEvent) {
		page = e.detail;
		offset = page * limit;

		[entities, paginationSettings.size] = await getEntities(offset, limit, searchString, filters);
	}

	async function searchFilter() {
		[entities, paginationSettings.size] = await getEntities(offset, limit, searchString, filters);
	}

	function triggerSearchWithDebounce() {
		// Clear existing timeout
		if (searchTimeout) {
			clearTimeout(searchTimeout);
		}

		if (searchString.length > 0) {
			// Debounce for 400ms
			searchTimeout = setTimeout(() => {
				searchFilter();
			}, 400);
		} else if (searchString.length === 0) {
			// If search is cleared, search immediately
			searchFilter();
		}
	}

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape') {
			searchString = '';
			searchVisible = false;
			searchFilter(); // Clear the search results
		}
	}

	const modalStore = getModalStore();

	async function rowClick(event: MouseEvent, id: number, category: string) {
		const target = event.target as HTMLElement;
		if (id == 0) {
			return;
		}

		if (
			(target.tagName === 'TD' && target.classList.contains('qr')) ||
			(target.tagName === 'I' && target.classList.contains('fa-qrcode'))
		) {
			// Display QR Modal
			const modal: ModalSettings = {
				type: 'component',
				component: {
					ref: QrCodeModal,
					props: {
						id: id,
						category: category
					}
				},
				title: 'QR Code',
				body: `Please scan or download the QR code for this ${cleanCategory(category)}`
			};
			modalStore.trigger(modal);
		} else if (
			(target.tagName === 'TD' && target.classList.contains('info')) ||
			(target.tagName === 'I' && target.classList.contains('fa-circle-info'))
		) {
			// Navigate to details page
			goto(`${category}/${id}`);
		} else {
			// Display Edit Modal
			const modal: ModalSettings = {
				type: 'component',
				component: {
					ref: AddNewModal,
					props: {
						edit: true,
						id: id,
						category: category,
						displayCategory: cleanCategory(category)
					}
				},
				title: 'Edit Entity',
				body: 'Please complete the form to edit an existing item, container, shelf, shelving unit, room, or building.'
			};
			modalStore.trigger(modal);
		}
	}

	let filters: { [key: string]: boolean } = {
		building: true,
		room: true,
		shelving_unit: true,
		shelf: true,
		container: true,
		item: true
	};

	let searchString = '';
	let searchTimeout: ReturnType<typeof setTimeout> | undefined = undefined;

	let searchVisible = false;
	let searchInput: HTMLInputElement | null = null;

	// Reactive statement to handle search when visibility changes
	$: if (searchVisible && searchString.length >= 0) {
		triggerSearchWithDebounce();
	}

	// Reactive statement to focus the input when it becomes visible
	$: if (searchVisible && searchInput) {
		// Use setTimeout to ensure the input is rendered before focusing
		tick().then(() => searchInput?.focus());
	}

	// Reactive statement to handle filter changes
	$: if (filters) {
		searchFilter();
	}

	const popupFeatured: PopupSettings = {
		event: 'click',
		target: 'popupFeatured',
		placement: 'bottom',
		closeQuery: '.btn-close-popup'
	};

	function toggleSearch() {
		searchVisible = !searchVisible;
	}

	function resetFilters() {
		filters = {
			building: true,
			room: true,
			shelving_unit: true,
			shelf: true,
			container: true,
			item: true
		};
	}

	function deselectAllFilters() {
		filters = {
			building: false,
			room: false,
			shelving_unit: false,
			shelf: false,
			container: false,
			item: false
		};
	}


</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<!-- svelte-ignore a11y-no-static-element-interactions -->
<div class="card p-4 w-72 shadow-xl z-50" data-popup="popupFeatured" on:click|stopPropagation>
	<label class="flex items-center space-x-2">
		<input class="checkbox" type="checkbox" bind:checked={filters['building']} />
		<p>Buildings</p>
	</label>

	<label class="flex items-center space-x-2">
		<input class="checkbox" type="checkbox" bind:checked={filters['room']} />
		<p>Rooms</p>
	</label>

	<label class="flex items-center space-x-2">
		<input class="checkbox" type="checkbox" bind:checked={filters['shelving_unit']} />
		<p>Shelving Units</p>
	</label>

	<label class="flex items-center space-x-2">
		<input class="checkbox" type="checkbox" bind:checked={filters['shelf']} />
		<p>Shelves</p>
	</label>

	<label class="flex items-center space-x-2">
		<input class="checkbox" type="checkbox" bind:checked={filters['container']} />
		<p>Containers</p>
	</label>

	<label class="flex items-center space-x-2">
		<input class="checkbox" type="checkbox" bind:checked={filters['item']} />
		<p>Items</p>
	</label>

	<div class="flex justify-center gap-2 mt-4">
		<button 
			type="button" 
			class="btn btn-sm variant-soft" 
			on:mousedown={(e) => { e.preventDefault(); e.stopPropagation(); resetFilters(); }}
		>
			Reset Filters
		</button>
		<button 
			type="button" 
			class="btn btn-sm variant-soft" 
			on:mousedown={(e) => { e.preventDefault(); e.stopPropagation(); deselectAllFilters(); }}
		>
			Deselect All
		</button>
	</div>
</div>

<div class="flex flex-row justify-between items-center h-16">
	<h2 class="text-xl">My Vault</h2>

	<div class="flex flex-row justify-end items-center w-2/3 md:5/12 lg:w-1/4">
		{#if searchVisible}
			<input
				bind:this={searchInput}
				class="input w-2/3"
				transition:slide={{ duration: 300, axis: 'x' }}
				type="text"
				bind:value={searchString}
				placeholder="search"
				on:keydown={handleKeydown}
				aria-label="Search"
			/>
		{/if}

		<div class="flex flex-row justify-evenly items-center w-1/3">
			<button
				id="filter"
				type="button"
				class="btn-icon btn-icon-sm variant-filled"
				on:click={toggleSearch}
			>
				<i class="fa-solid fa-magnifying-glass fa-s"></i>
			</button>

			<button
				id="filter"
				type="button"
				class="btn-icon btn-icon-sm variant-filled"
				use:popup={popupFeatured}
			>
				<i class="fa-solid fa-filter fa-xs"></i>
			</button>
		</div>
	</div>
</div>

<div class="flex flex-col justify-between">
	{#if entities.length > 0}
		<div class="table-container pb-4" bind:clientWidth={width}>
			<table class="table table-compact table-hover">
				<thead>
					<tr>
						<th id="th-name" class="!py-2">Name</th>
						<th id="th-category" class="!py-2 hidden md:table-cell">Category</th>
						<th id="th-location" class="!py-2">Location</th>
						<th id="th-notes" class="!py-2 hidden lg:table-cell">Notes</th>
						<th>&nbsp;</th>
						<th>&nbsp;</th>
					</tr>
				</thead>
				<tbody>
					{#each entities as entity}
						<tr on:click={(event) => rowClick(event, entity.ID, entity.Category)}>
							<td class="capitalize">{entity.Name}</td>
							<td class="hidden md:table-cell capitalize">{cleanCategory(entity.Category)}</td>
							<td>
								{#if entity.Category == 'building'}
									{entity.Address}
								{:else}
									{#each [...entity.Parent].slice(0, parentMax).reverse() as parent, index}
										<span class="capitalize block sm:inline">{parent.Name}</span>

										{#if index < entity.Parent.length - 1 && index < parentMax - 1}
											<span class="hidden sm:inline">&nbsp;<i class="fa-solid fa-arrow-right"></i>&nbsp;</span>
											<span class="sm:hidden"><br /></span>
										{/if}
									{/each}
								{/if}
							</td>
							<td class="hidden lg:table-cell">{entity.Notes ?? ''}</td>
							<td class="info">
								<button type="button">
									{#if entity.ID != 0}
										<i class="fa-solid fa-circle-info"></i>
									{:else}
										&nbsp;
									{/if}
								</button>
							</td>
							<td class="qr">
								<button type="button">
									{#if entity.ID != 0}
										<i class="fa-solid fa-qrcode"></i>
									{:else}
										&nbsp;
									{/if}
								</button>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
	{#if entities.length == 0}
		<div class="pt-32">
			<p class="text-center animate-spin">
				<i class="fa-solid fa-spinner fa-2xl"></i>
			</p>
		</div>
	{/if}
	<div class="fixed bottom-0 left-1/2 -translate-x-1/2">
		<Paginator
			bind:settings={paginationSettings}
			showFirstLastButtons={false}
			showPreviousNextButtons={true}
			on:amount={limitChange}
			on:page={pageChange}
		/>
	</div>
</div>
