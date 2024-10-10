<script lang="ts">
	import { _confirmUser as ConfirmUser } from './+page';

	const formData = {
		confirmationCode: ''
	};

	var isFormInvalid = true;
	var formError = {
		confirmationCode: '',
		error: ''
	};

	var formErrorClass = {
		confirmationCode: '',
		error: ''
	};

	function validateConfirmationCode() {
		if (formData.confirmationCode == '') {
			isFormInvalid = true;
			formError.confirmationCode = 'Confirmation code is required!';
			formErrorClass.confirmationCode = 'input-error';
		} else {
			isFormInvalid = false;
			formError.confirmationCode = '';
			formErrorClass.confirmationCode = '';
		}
	}

	async function onFormSubmit() {
		const [success, message] = await ConfirmUser(formData);

		if (!success) {
			formError.error = message;
			formErrorClass.error = 'input-error';
		}
	}
</script>

<h2 class="text-center text-xl">Confirm User Account</h2>
<form class="w-3/6 p-8 my-0 mx-auto flex flex-col justify-center">
	<label for="confirmation" class="label">Confirmation Code:</label>
	<input
		id="confirmation"
		class="input {formErrorClass.confirmationCode}"
		type="text"
		bind:value={formData.confirmationCode}
		on:input={validateConfirmationCode}
		on:focusout={validateConfirmationCode}
		placeholder="Enter username..."
	/>
	{#if formError.confirmationCode}
		<p class="text-red-500 !mt-0">{formError.confirmationCode}</p>
	{/if}

	<button class="btn variant-filled my-4 mx-auto" disabled={isFormInvalid} on:click={onFormSubmit}
		>Submit</button
	>
	{#if formError.error}
		<p class="text-red-500 !mt-0">{formError.error}</p>
	{/if}
</form>
