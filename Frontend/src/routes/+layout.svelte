<script lang="ts">
	import '../app.postcss';

	// Highlight JS
	import hljs from 'highlight.js/lib/core';
	import 'highlight.js/styles/github-dark.css';
	import { storeHighlightJs } from '@skeletonlabs/skeleton';
	import xml from 'highlight.js/lib/languages/xml'; // for HTML
	import css from 'highlight.js/lib/languages/css';
	import javascript from 'highlight.js/lib/languages/javascript';
	import typescript from 'highlight.js/lib/languages/typescript';

	hljs.registerLanguage('xml', xml); // for HTML
	hljs.registerLanguage('css', css);
	hljs.registerLanguage('javascript', javascript);
	hljs.registerLanguage('typescript', typescript);
	storeHighlightJs.set(hljs);

	// Floating UI for Popups
	import { computePosition, autoUpdate, flip, shift, offset, arrow } from '@floating-ui/dom';
	import { storePopup } from '@skeletonlabs/skeleton';
	import { AppRail, AppRailTile, AppRailAnchor } from '@skeletonlabs/skeleton';
	import '@fortawesome/fontawesome-free/css/all.min.css';

	storePopup.set({ computePosition, autoUpdate, flip, shift, offset, arrow });
	let currentTile: number = 0;

	let visible = true;
	function toggleVisible() {
		visible = !visible;
	}
</script>

<div class="flex flex-row-reverse h-full">
	{#if visible}
		<AppRail>
			<svelte:fragment slot="lead">
				<AppRailAnchor on:click={toggleVisible}><i class="fa-solid fa-bars"></i></AppRailAnchor>
			</svelte:fragment>
			<!-- --- -->
			<AppRailTile bind:group={currentTile} name="tile-2" value={0} title="tile-2">
				<svelte:fragment slot="lead"><i class="fa-solid fa-plus"></i></svelte:fragment>
				<span>Add New</span>
			</AppRailTile>
			<AppRailTile bind:group={currentTile} name="tile-3" value={1} title="tile-3">
				<svelte:fragment slot="lead"><i class="fa-solid fa-filter"></i></svelte:fragment>
				<span>Filter</span>
			</AppRailTile>
			<!-- --- -->
			<svelte:fragment slot="trail">
				<AppRailAnchor
					href="https://github.com/Lanham-Software-James/Organize-This"
					target="_blank"
					title="Account"><i class="fa-brands fa-github"></i></AppRailAnchor
				>
			</svelte:fragment>
		</AppRail>
	{/if}

	{#if !visible}
		<AppRail height="h-fit">
			<svelte:fragment slot="lead">
				<AppRailAnchor on:click={toggleVisible}><i class="fa-solid fa-bars"></i></AppRailAnchor>
			</svelte:fragment>
		</AppRail>
	{/if}

	<div class="container p-3">
		<h1 class="pb-4 text-3xl">Organize This!</h1>
		<slot />
	</div>
</div>
