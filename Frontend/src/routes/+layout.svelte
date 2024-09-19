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
	storePopup.set({ computePosition, autoUpdate, flip, shift, offset, arrow });

	// Font Awesome
	import '@fortawesome/fontawesome-free/css/all.min.css';

	// Modal and drawers
	import { initializeStores, Drawer, type DrawerSettings } from '@skeletonlabs/skeleton';
	initializeStores();

	import { getDrawerStore } from "@skeletonlabs/skeleton";

	const drawerStore = getDrawerStore();

	function showNavigation() {
		const drawerSettings: DrawerSettings = {
			padding: 'p-4',
			rounded: 'rounded-xl',
			position: 'right',
			width: 'w-6/12'
		};
		drawerStore.open(drawerSettings);
	}
</script>

<Drawer />

<div class="container p-3 h-full">
	<div class="flex flex-row justify-between items-center pb-2">
		<h1 class="text-3xl">Organize This!</h1>

		<button type="button" class="btn-icon variant-filled" on:click={showNavigation}><i class="fa-solid fa-bars fa-sm"></i></button>
	</div>
	<slot />
</div>
