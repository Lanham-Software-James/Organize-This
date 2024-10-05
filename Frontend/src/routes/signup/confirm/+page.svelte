<script lang="ts">
	import { _signUpUser as signUpUser } from './+page';

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
		const success = await signUpUser(formData);

		if (!success) {
			formError.error = 'Invalid Confirmation Code!';
			formErrorClass.error = 'input-error';
		}
	}
</script>

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
