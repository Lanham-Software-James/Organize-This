<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { cleanCategory } from '$lib/CleanCategory/CleanCategory';
	import type { getEntityData, getEntityEntity } from './+page';

	$: message = $page.data.message satisfies string;
	$: entity = $page.data.entity satisfies getEntityData;
    $: parentName = $page.data.parentName satisfies string;
    $: children = $page.data.children satisfies getEntityEntity[];

    function navParent() {
        goto(`/${entity.Parent.ParentCategory}/${entity.Parent.ParentID}`)
    }

    function navChild(category: string , id: string) {
        goto(`/${category}/${id}`)
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
                    <td class="capitalize">{cleanCategory(entity.Entity.Category)}</td>
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

        {#if entity.Entity.Category != "item"}
            <h3 class="text-l capitalize pt-6">Children</h3>
            <div class="table-container pt-3 w-1/2">
                <table class="table table-compact table-hover">
                    <thead>
                        <tr>
                            <th id="th-name" class="!py-2">Name</th>
                            <th id="th-category" class="!py-2 invisible md:visible">Category</th>
                        </tr>
                    </thead>
                    <tbody>
                        {#each children as child}
                            <tr on:click={() => navChild(child.Category, child.ID)}>
                                <td class="capitalize">{child.Name}</td>
                                <td class="invisible md:visible capitalize">{cleanCategory(child.Category)}</td>
                            </tr>
                        {/each}
                    </tbody>
                </table>
            </div>
        {/if}
    {/if}

	{#if message != 'success'}
		<div class="pt-32">
			<p class="text-center animate-spin">
				<i class="fa-solid fa-spinner fa-2xl"></i>
			</p>
		</div>
	{/if}
</div>
