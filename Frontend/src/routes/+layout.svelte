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
	import { setContext } from 'svelte';
    import { writable } from 'svelte/store';
	import AddNewModal from '$lib/AddNewModal/AddNewModal.svelte';
	import { _logoutUser as logoutUser } from './+layout';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores'
	import { computePosition, autoUpdate, flip, shift, offset, arrow } from '@floating-ui/dom';
	import { storePopup } from '@skeletonlabs/skeleton';
	import {
		initializeStores,
		Modal,
		getModalStore,
		type ModalComponent,
		type ModalSettings,
		Drawer,
		type DrawerSettings,
		getDrawerStore,
		Toast
	} from '@skeletonlabs/skeleton';
	import '@fortawesome/fontawesome-free/css/all.min.css';

	hljs.registerLanguage('xml', xml); // for HTML
	hljs.registerLanguage('css', css);
	hljs.registerLanguage('javascript', javascript);
	hljs.registerLanguage('typescript', typescript);
	storeHighlightJs.set(hljs);

	storePopup.set({ computePosition, autoUpdate, flip, shift, offset, arrow });

	initializeStores();

	const drawerStore = getDrawerStore();
	function showNavigation() {
		const drawerSettings: DrawerSettings = {
			id: 'navbar',
			padding: 'p-4',
			rounded: 'rounded-xl',
			position: 'right',
			width: 'w-7/12 md:w-4/12 lg:w-2/12'
		};
		drawerStore.open(drawerSettings);
	}

	const refreshPageStore = writable(false);

    setContext('refreshPage', {
        subscribe: refreshPageStore.subscribe,
        refresh: () => refreshPageStore.update(n => !n)
    });

	$: isUserAuthed = $page.data.cookieExists satisfies boolean

	const modalStore = getModalStore();

	function showModal() {
		drawerStore.close();

		const modal: ModalSettings = {
			type: 'component',
			component: {ref: AddNewModal, props: {edit: false}},
			title: 'Add New',
			body: 'Please complete the form to add a new item, container, shelf, shelving unit, room, or building.'
		};
		modalStore.trigger(modal);
	}

	async function logout() {
		drawerStore.close();

		var success = await logoutUser();
	}

	function navigate(route: string) {
		drawerStore.close();
		goto(route)
	}
</script>

<Drawer>
	{#if $drawerStore.id === 'navbar'}
		<div class="p-2">
			{#if isUserAuthed}
				<ul id="pages">
					<p>Overview</p>
					<li>
						<button type="button" class="btn bg-initial" on:click={() => navigate("/")}>
							<span><i class="fa-solid fa-house"></i></span>
							<span>Home</span>
						</button>
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
				<ul id="feedback" class="pt-4">
					<p>Feedback</p>

					<button type="button" class="btn bg-initial" on:click={() => navigate("/feedback")}>
						<span><i class="fa-solid fa-comment"></i></span>
						<span>Leave Feedback</span>
					</button>
				</ul>

				<ul id="account" class="pt-4">
					<p>Account</p>

					<button type="button" class="btn bg-initial" on:click={logout}>
						<span><i class="fa-solid fa-arrow-right-from-bracket"></i></span>
						<span>Log Out</span>
					</button>
				</ul>
			{/if}
			{#if !isUserAuthed}
			<ul id="account" class="pt-4">
				<li>
					<button type="button" class="btn bg-initial" on:click={() => navigate("/login")}>
						<span><i class="fa-solid fa-arrow-right-from-bracket"></i></span>
						<span>Sign In</span>
					</button>
				</li>
				<li>
					<button type="button" class="btn bg-initial" on:click={() => navigate("/signup")}>
						<span><i class="fa-solid fa-user-plus"></i></span>
						<span>Create Account</span>
					</button>
				</li>
			</ul>
			{/if}
		</div>
	{/if}
</Drawer>

<Modal />

<Toast />

<div class="container p-3 m-h-full mx-auto">
	<div class="flex flex-row justify-between items-center">
		<h1 class="text-3xl">WillowSuite Vault</h1>

		<button id="hamburgerMenu" type="button" class="btn-icon bg-initial" on:click={showNavigation}
			><i class="fa-solid fa-bars fa-xl"></i></button
		>
	</div>
	<slot />
</div>
