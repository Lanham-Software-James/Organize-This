<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import type { getEntityData } from './+page';

	$: message = $page.data.message satisfies string;
	$: entity = $page.data.entity satisfies getEntityData;
    $: parentName = $page.data.parentName satisfies string;

    function navParent() {

        goto(`/${entity.Parent.ParentCategory}/${entity.Parent.ParentID}`)
    }
</script>

<div class="flex flex-row justify-between items-center h-16">
	<h2 class="text-xl capitalize">View Entity</h2>
</div>

<div class="flex flex-col justify-between">
	{#if message == 'success'}
		<table class="table table-compact table-hover w-1/4">
            <tbody>
                <tr>
                    <th class="text-left w-1/4 pl-2">Name:</th>
                    <td class="capitalize">{entity.Entity.Name}</td>
                </tr>
                <tr>
                    <th class="text-left w-1/4 pl-2">Category:</th>
                    <td class="capitalize">{entity.Entity.Category}</td>
                </tr>

                {#if entity.Entity.Category === "building"}
                    <tr>
                        <th class="text-left w-1/4 pl-2">Address:</th>
                        <td class="capitalize">{entity.Address}</td>
                    </tr>
                {:else}
                    <tr>
                        <th class="text-left w-1/4 pl-2">Parent:</th>
                        <td class="capitalize hover:cursor-pointer" on:click={navParent}>{parentName}</td>
                    </tr>
                {/if}

                <tr>
                    <th class="text-left w-1/4 pl-2">Notes:</th>
                    {#if entity.Entity.Notes === undefined || entity.Entity.Notes.length === 0}
                        <td>&nbsp;</td>
                    {:else}
                        <td class="capitalize">{entity.Entity.Notes}</td>
                    {/if}
                </tr>
            </tbody>
        </table>
	{:else}
		<div class="pt-32">
			<p class="text-center animate-spin">
				<i class="fa-solid fa-spinner fa-2xl"></i>
			</p>
		</div>
	{/if}
</div>
