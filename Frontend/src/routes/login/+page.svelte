<script lang="ts">
	import { _loginUser as loginUser } from './+page';

	const formData = {
		userEmail: '',
		password: ''
	};

	var isFormInvalid = true;
	var formError = {
		username: '',
		password: '',
		incorrect: ''
	};
	var formErrorClass = {
		username: '',
		password: '',
		incorrect: ''
	};

	function validateUsername() {
		if (formData.userEmail == '') {
			isFormInvalid = true;
			formError.username = 'Username is required!';
			formErrorClass.username = 'input-error';
		} else {
			isFormInvalid = false;
			formError.username = '';
			formErrorClass.username = '';
		}
	}

	function validatePassword() {
		if (formData.password == '') {
			isFormInvalid = true;
			formError.password = 'Password is required!';
			formErrorClass.password = 'input-error';
		} else {
			isFormInvalid = false;
			formError.password = '';
			formErrorClass.password = '';
		}
	}

	async function onFormSubmit() {
		const [success, message] = await loginUser(formData);

		if (!success) {
			formError.incorrect = message;
			formErrorClass.incorrect = 'input-error';
		}
	}
</script>

<h2 class="text-center text-xl">Login</h2>

<form class="w-11/12 md:w-9/12 lg:w-3/6 p-8 my-0 mx-auto flex flex-col justify-center" on:submit|preventDefault={onFormSubmit}>
	<label for="username" class="label">Username:</label>
	<input
		id="username"
		class="input {formErrorClass.username}"
		type="text"
		bind:value={formData.userEmail}
		on:input={validateUsername}
		on:focusout={validateUsername}
		placeholder="Enter username..."
	/>
	{#if formError.username}
		<p class="text-red-500 !mt-0">{formError.username}</p>
	{/if}

	<label for="password" class="label pt-4">Password:</label>
	<input
		id="password"
		class="input {formErrorClass.password}"
		type="password"
		bind:value={formData.password}
		on:input={validatePassword}
		on:focusout={validatePassword}
		placeholder="Enter password..."
	/>
	{#if formError.password}
		<p class="text-red-500 !mt-0">{formError.password}</p>
	{/if}

	<button class="btn variant-filled my-4 mx-auto" disabled={isFormInvalid} on:click={onFormSubmit}
		>Submit</button
	>
	{#if formError.incorrect}
		<p class="text-red-500 !mt-0">{formError.incorrect}</p>
	{/if}
</form>
