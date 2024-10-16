<script lang="ts">
	import { Paginator, type PaginationSettings } from '@skeletonlabs/skeleton';
	import { onMount } from 'svelte';
	import { _getEntities as getEntities, type GetEntitiesData } from './+page';

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

	onMount(async function () {
		[entities, paginationSettings.size] = await getEntities(offset, limit);
	});

	async function limitChange(e: CustomEvent) {
		limit = e.detail;

		[entities, paginationSettings.size] = await getEntities(offset, limit);
	}

	async function pageChange(e: CustomEvent) {
		page = e.detail;
		offset = page * limit;

		[entities, paginationSettings.size] = await getEntities(offset, limit);
	}
</script>

<div class="flex flex-row justify-between items-center pb-2">
	<h2 class="text-xl">All Things</h2>

	<button id="filter" type="button" class="btn-icon btn-icon-sm variant-filled"
		><i class="fa-solid fa-filter fa-xs"></i></button
	>
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
						<tr>
							<td>{entity.Name}</td>
							<td class="hidden md:block">{entity.Category}</td>
							<td>
								{#each [...entity.Parent].reverse() as parent, index}
									<span>{parent.Name}</span>

									{#if index < entity.Parent.length - 1}
										<span>&nbsp;<i class="fa-solid fa-arrow-right"></i>&nbsp;</span>
									{/if}
								{/each}
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
