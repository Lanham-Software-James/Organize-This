<script lang="ts">
	import { Paginator, type PaginationSettings } from '@skeletonlabs/skeleton';

	const source = [ /* any array of objects */ ];

	let paginationSettings = {
	page: 0,
	limit: 5,
	size: source.length,
	amounts: [10, 25, 50, 100],
} satisfies PaginationSettings;
</script>

<h2 class="pb-4 text-xl">Items</h2>
<div class="flex flex-col justify-between">
	<div class="table-container pb-4">
		<table class="table table-compact table-hover">
			<thead>
				<tr>
					<th class="!py-2">Name</th>
					<th class="!py-2">Location</th>
				</tr>
			</thead>
			<tbody>
					<tr>
						<td> </td>
						<td> </td>
					</tr>
			</tbody>
		</table>
	</div>

	<div>
		<Paginator
		bind:settings={paginationSettings}
		showFirstLastButtons="{false}"
		showPreviousNextButtons="{true}"
		/>
	</div>

</div>

