<script lang="ts">
	import { Paginator, type PaginationSettings, getModalStore, popup, type ModalSettings, type PopupSettings } from '@skeletonlabs/skeleton';
	import { getContext, onMount } from 'svelte';
	import { _getEntities as getEntities, type GetEntitiesData } from './+page';
	import AddNewModal from '$lib/AddNewModal/AddNewModal.svelte';
	import { slide } from 'svelte/transition';

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

	onMount(async function () {
		loadData();

		//@ts-ignore
        const unsubscribe = refreshPage.subscribe(() => {
            loadData();
        });

        return unsubscribe;
	});

	async function loadData() {
		[entities, paginationSettings.size] = await getEntities(offset, limit);
	}

	async function limitChange(e: CustomEvent) {
		limit = e.detail;

		[entities, paginationSettings.size] = await getEntities(offset, limit);
	}

	async function pageChange(e: CustomEvent) {
		page = e.detail;
		offset = page * limit;

		[entities, paginationSettings.size] = await getEntities(offset, limit);
	}

	async function searchFilter() {
		console.log("\n\nSearch: ")
		console.log(search)
		console.log("\nFilters: ")
		console.log(filters)
	}

	const modalStore = getModalStore();

	async function editEntity(id: number, category: string) {
		const modal: ModalSettings = {
			type: 'component',
			component: {ref: AddNewModal, props: {edit: true, id: id, category: category, displayCategory: cleanCategory(category)}},
			title: 'Edit Entity',
			body: 'Please complete the form to edit an existing item, container, shelf, shelving unit, room, or building.'
		};
		modalStore.trigger(modal);
	}

	function cleanCategory(category: string): string {
		var cleanedCategory = category;
		if(cleanedCategory == "shelving_unit"){
			cleanedCategory = "shelving unit"
		}

		return cleanedCategory
	}

	let filters: {[key: string]: boolean} = {
		'building': true,
		'room': true,
		'shelving_unit': true,
		'shelf': true,
		'container': true,
		'item': true,
	}

	let search = ''

	let isVisible = false;

	const popupFeatured: PopupSettings = {
		event: 'click',
		target: 'popupFeatured',
		placement: 'bottom',
		state: (e: Record<string, boolean>) => popUpOpenClose(e)
	};

	function popUpOpenClose(e: Record<string, boolean>) {
		if(!e.state) {
			searchFilter()
		}
	}

	function toggleSearch() {
		isVisible = !isVisible;
	}
</script>

<div class="card p-4 w-72 shadow-xl" data-popup="popupFeatured">
	<label class="flex items-center space-x-2">
		<input class="checkbox" type="checkbox" bind:checked={filters['building']}/>
		<p>Buildings</p>
	</label>

	<label class="flex items-center space-x-2">
		<input class="checkbox" type="checkbox" bind:checked={filters['room']}/>
		<p>Rooms</p>
	</label>

	<label class="flex items-center space-x-2">
		<input class="checkbox" type="checkbox" bind:checked={filters['shelving_unit']}/>
		<p>Shelving Units</p>
	</label>

	<label class="flex items-center space-x-2">
		<input class="checkbox" type="checkbox" bind:checked={filters['shelf']}/>
		<p>Shelves</p>
	</label>

	<label class="flex items-center space-x-2">
		<input class="checkbox" type="checkbox" bind:checked={filters['container']}/>
		<p>Containers</p>
	</label>

	<label class="flex items-center space-x-2">
		<input class="checkbox" type="checkbox" bind:checked={filters['item']}/>
		<p>Items</p>
	</label>
</div>

<div class="flex flex-row justify-between items-center h-16">
	<h2 class="text-xl">All Things</h2>

	<div class="flex flex-row justify-end items-center w-3/12">
		{#if isVisible}
			<input
				class="input w-3/4"
				transition:slide={{ duration: 300, axis: 'x' }}
				type="text"
				bind:value={search}
				placeholder="search"
				on:blur = {searchFilter}
			/>
		{/if}

		<div class="flex flex-row justify-evenly items-center w-1/4">
			<button id="filter" type="button" class="btn-icon btn-icon-sm variant-filled" on:click={toggleSearch}>
				<i class="fa-solid fa-magnifying-glass fa-s"></i>
			</button>

			<button id="filter" type="button" class="btn-icon btn-icon-sm variant-filled" use:popup={popupFeatured}>
				<i class="fa-solid fa-filter fa-xs"></i>
			</button>
		</div>
	</div>
</div>

<div class="flex flex-col justify-between">
	{#if entities.length > 0}
		<div class="table-container pb-4">
			<table class="table table-compact table-hover">
				<thead>
					<tr>
						<th id="th-name" class="!py-2">Name</th>
						<th id="th-category" class="!py-2 hidden md:block">Category</th>
						<th id="th-location" class="!py-2">Location</th>
						<th id="th-notes" class="!py-2 hidden lg:block">Notes</th>
					</tr>
				</thead>
				<tbody>
					{#each entities as entity}
						<tr on:click={() => editEntity(entity.ID, entity.Category)}>
							<td class="capitalize">{entity.Name}</td>
							<td class="hidden md:block capitalize">{cleanCategory(entity.Category)}</td>
							<td>
								{#if entity.Category == 'building'}
									{entity.Address}
								{:else}
									{#each [...entity.Parent].reverse() as parent, index}
										<span class="capitalize">{parent.Name}</span>

										{#if index < entity.Parent.length - 1}
											<span>&nbsp;<i class="fa-solid fa-arrow-right"></i>&nbsp;</span>
										{/if}
									{/each}
								{/if}
							</td>
							<td class="hidden lg:block">{entity.Notes ?? ""}</td>
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
