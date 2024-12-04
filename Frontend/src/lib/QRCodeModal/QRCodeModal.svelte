<script lang="ts">
	import { onMount, type SvelteComponent } from 'svelte';
	import { getModalStore } from '@skeletonlabs/skeleton';
	import { generateQR } from './QRCodeModal';

	export let parent: SvelteComponent;
	export let id: number;
	export let category: string;
	let url: string = ""

	const modalStore = getModalStore();

	// Base Classes
	const cBase = 'card p-4 w-modal shadow-xl space-y-4 max-h-screen overflow-scroll';
	const cHeader = 'text-2xl font-bold';

	onMount(async function () {
		var message: string = "";
		[message, url] = await generateQR(category, id);
	});

	function downloadFile() {
		const link = document.createElement('a');
		link.href = url;
		link.download = id + "-" + category + ".jpg";
		document.body.appendChild(link);
		link.click();
		document.body.removeChild(link);
	}
</script>

{#if $modalStore[0]}
	<div class="modal-add-entity-form {cBase}">
		<header class={cHeader}>{$modalStore[0].title ?? '(title missing)'}</header>
		<article>{$modalStore[0].body ?? '(body missing)'}</article>

		{#if url != ""}
			<!-- svelte-ignore a11y-missing-attribute -->
			<img class="size-96 m-auto py-2" src={url} />
		{:else}
			<div class="size-96 m-auto py-2 flex items-center justify-center">
				<p class="text-center animate-spin">
					<i class="fa-solid fa-spinner fa-2xl"></i>
				</p>
			</div>
		{/if}

		<!-- prettier-ignore -->
		<footer class="modal-footer flex justify-end items-center">
			<div class="w-5/6 md:w-5/12 flex justify-evenly">
				<button class="btn {parent.buttonNeutral}" on:click={parent.onClose}>{parent.buttonTextCancel}</button>
				<button class="btn {parent.buttonPositive}"on:click={downloadFile}>Download</button>
			</div>

		</footer>
	</div>
{/if}
