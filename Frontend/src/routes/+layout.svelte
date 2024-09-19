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
	import { initializeStores, Modal, getModalStore, type ModalComponent, type ModalSettings, Drawer, type DrawerSettings, getDrawerStore, Toast, getToastStore } from '@skeletonlabs/skeleton';
	initializeStores();

	const drawerStore = getDrawerStore();
	function showNavigation() {
		const drawerSettings: DrawerSettings = {
			id: 'navbar',
			padding: 'p-4',
			rounded: 'rounded-xl',
			position: 'right',
			width: 'w-6/12'
		};
		drawerStore.open(drawerSettings);
	}

	import AddNewModal from '$lib/AddNewModal/AddNewModal.svelte';
	const modalStore = getModalStore();
	const modalRegistry: Record<string, ModalComponent> = {
		addNewModal: { ref: AddNewModal },
	};

	function showModal() {
		drawerStore.close();

		const modal: ModalSettings = {
			type: 'component',
			component: 'addNewModal',
			title: "Add New",
			body: "Please complete the form to add a new item, container, shelf, shelving unit, room, or building.",
		};
		modalStore.trigger(modal);
	}

	const toastStore = getToastStore();
</script>

<Drawer>
	{#if $drawerStore.id === 'navbar'}
		<div class="p-2">
			<ul id="pages">
				<p>Overview</p>
				<li>
					<a href="/">
					<button type="button" class="btn bg-initial">
					<span><i class="fa-solid fa-house"></i></span>
					<span>Home</span>
					</button>
					</a>
				</li>
			</ul>
			<ul id="tools" class="pt-4">
				<p>Tools</p>
				<li>
					<button type="button" class="btn bg-initial" on:click={showModal}>
					<span><i class="fa-solid fa-plus"></i></span>
					<span>Add New</span>
					</button>
				</li>
			</ul>
			<ul id ="account" class="pt-4">
				<p>Account</p>
			</ul>
		</div>
	{/if}
</Drawer>

<Modal components={modalRegistry} />

<Toast />

<div class="container p-3 h-full">
	<div class="flex flex-row justify-between items-center pb-2">
		<h1 class="text-3xl">Organize This!</h1>

		<button type="button" class="btn-icon bg-initial" on:click={showNavigation}><i class="fa-solid fa-bars fa-xl"></i></button>
	</div>
	<slot />
</div>
