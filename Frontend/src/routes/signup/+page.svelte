<script lang="ts">
	import { _signUpUser as signUpUser } from './+page';

	const formData = {
		userEmail: '',
		password: '',
		confirmPassword: '',
		firstName: '',
		lastName: '',
		birthday: ''
	};

	var isFormInvalid = true;
	var formError = {
		username: '',
		password: '',
		firstName: '',
		lastName: '',
		birthday: '',
		incorrect: ''
	};
	var formErrorClass = {
		username: '',
		password: '',
		firstName: '',
		lastName: '',
		birthday: '',
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
		} else if (formData.password != formData.confirmPassword) {
			isFormInvalid = true;
			formError.password = 'Passwords must match!';
			formErrorClass.password = 'input-error';
		} else {
			isFormInvalid = false;
			formError.password = '';
			formErrorClass.password = '';
		}
	}

	function validateFirstName() {
		if (formData.firstName == '') {
			isFormInvalid = true;
			formError.firstName = 'First name is required!';
			formErrorClass.firstName = 'input-error';
		} else {
			isFormInvalid = false;
			formError.firstName = '';
			formErrorClass.firstName = '';
		}
	}

	function validateLastName() {
		if (formData.lastName == '') {
			isFormInvalid = true;
			formError.lastName = 'Last name is required!';
			formErrorClass.lastName = 'input-error';
		} else {
			isFormInvalid = false;
			formError.lastName = '';
			formErrorClass.lastName = '';
		}
	}

	function validateBirthday() {
		if (formData.birthday == '') {
			isFormInvalid = true;
			formError.birthday = 'Birthday is required!';
			formErrorClass.birthday = 'input-error';
		} else {
			isFormInvalid = false;
			formError.birthday = '';
			formErrorClass.birthday = '';
		}
	}

	async function onFormSubmit() {
		const [success, message] = await signUpUser(formData);

		if (!success) {
			formError.incorrect = message;
			formErrorClass.incorrect = 'input-error';
		}
	}
</script>

<h2 class="text-center text-xl">Sign Up</h2>

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
	<label for="confirm-password" class="label pt-4">Confirm Password:</label>
	<input
		id="confirm-password"
		class="input {formErrorClass.password}"
		type="password"
		bind:value={formData.confirmPassword}
		on:input={validatePassword}
		on:focusout={validatePassword}
		placeholder="Confirm password..."
	/>
	{#if formError.password}
		<p class="text-red-500 !mt-0">{formError.password}</p>
	{/if}

	<label for="firstname" class="label pt-4">First Name:</label>
	<input
		id="firstname"
		class="input {formErrorClass.firstName}"
		type="text"
		bind:value={formData.firstName}
		on:input={validateFirstName}
		on:focusout={validateFirstName}
		placeholder="Enter first name..."
	/>
	{#if formError.firstName}
		<p class="text-red-500 !mt-0">{formError.firstName}</p>
	{/if}

	<label for="lastname" class="label pt-4">Last Name:</label>
	<input
		id="lastname"
		class="input {formErrorClass.lastName}"
		type="text"
		bind:value={formData.lastName}
		on:input={validateLastName}
		on:focusout={validateLastName}
		placeholder="Enter last name..."
	/>
	{#if formError.lastName}
		<p class="text-red-500 !mt-0">{formError.lastName}</p>
	{/if}

	<label for="birthday" class="label pt-4">Birthday:</label>
	<input
		id="birthday"
		class="input {formErrorClass.birthday}"
		type="date"
		bind:value={formData.birthday}
		on:input={validateBirthday}
		on:focusout={validateBirthday}
		placeholder="Enter birthday..."
	/>
	{#if formError.birthday}
		<p class="text-red-500 !mt-0">{formError.birthday}</p>
	{/if}

	<button class="btn variant-filled my-4 mx-auto" disabled={isFormInvalid} on:click={onFormSubmit}
		>Submit</button
	>
	{#if formError.incorrect}
		<p class="text-red-500 !mt-0">{formError.incorrect}</p>
	{/if}
</form>
