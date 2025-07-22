<script lang="ts">
  let email = '';
  let message = '';
  let overallSatisfaction = 0;
  let easeOfUse = 0;
  let featureCompleteness = 0;
  let suggestionType = '';
  let isSubmitting = false;
  let submitStatus = ''; // 'success', 'error', or ''

  var isFormInvalid = true;
  var formError = {
    email: '',
    message: '',
    overallSatisfaction: '',
    suggestionType: '',
    incorrect: ''
  };
  var formErrorClass = {
    email: '',
    message: '',
    overallSatisfaction: '',
    suggestionType: '',
    incorrect: ''
  };

  const satisfactionOptions = [
    { value: 1, label: 'Very Dissatisfied' },
    { value: 2, label: 'Dissatisfied' },
    { value: 3, label: 'Neutral' },
    { value: 4, label: 'Satisfied' },
    { value: 5, label: 'Very Satisfied' }
  ];

  const suggestionTypes = [
    { value: '', label: 'Select suggestion type...' },
    { value: 'feature', label: 'New Feature Request' },
    { value: 'improvement', label: 'Feature Improvement' },
    { value: 'ui', label: 'UI/UX Suggestion' },
    { value: 'other', label: 'Other' }
  ];



  function validateEmail() {
    if (email == '') {
      isFormInvalid = true;
      formError.email = 'Email is required!';
      formErrorClass.email = 'input-error';
    } else {
      isFormInvalid = false;
      formError.email = '';
      formErrorClass.email = '';
    }
  }

  function validateMessage() {
    if (message == '') {
      isFormInvalid = true;
      formError.message = 'Message is required!';
      formErrorClass.message = 'input-error';
    } else {
      isFormInvalid = false;
      formError.message = '';
      formErrorClass.message = '';
    }
  }

  function validateOverallSatisfaction() {
    if (overallSatisfaction == 0) {
      isFormInvalid = true;
      formError.overallSatisfaction = 'Please rate your overall satisfaction!';
      formErrorClass.overallSatisfaction = 'input-error';
    } else {
      isFormInvalid = false;
      formError.overallSatisfaction = '';
      formErrorClass.overallSatisfaction = '';
    }
  }

  function validateSuggestionType() {
    if (suggestionType == '') {
      isFormInvalid = true;
      formError.suggestionType = 'Please select a suggestion type!';
      formErrorClass.suggestionType = 'input-error';
    } else {
      isFormInvalid = false;
      formError.suggestionType = '';
      formErrorClass.suggestionType = '';
    }
  }

  function checkFormValidity() {
    validateEmail();
    validateMessage();
    validateOverallSatisfaction();
    validateSuggestionType();
  }

  async function handleSubmit(event: SubmitEvent) {
    event.preventDefault();
    
    checkFormValidity();
    
    if (isFormInvalid) {
      return;
    }

    isSubmitting = true;
    submitStatus = '';

    try {
      const response = await fetch('https://formspree.io/f/xqalppbq', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          email: email,
          message: message,
          overallSatisfaction: overallSatisfaction,
          easeOfUse: easeOfUse,
          featureCompleteness: featureCompleteness,
          suggestionType: suggestionType
        })
      });

      if (response.ok) {
        submitStatus = 'success';
        email = '';
        message = '';
        overallSatisfaction = 0;
        easeOfUse = 0;
        featureCompleteness = 0;
        suggestionType = '';
        // Reset form validation state
        isFormInvalid = true;
        formError.email = '';
        formError.message = '';
        formError.overallSatisfaction = '';
        formError.suggestionType = '';
        formErrorClass.email = '';
        formErrorClass.message = '';
        formErrorClass.overallSatisfaction = '';
        formErrorClass.suggestionType = '';
      } else {
        submitStatus = 'error';
        formError.incorrect = 'There was an error submitting your feedback. Please try again.';
        formErrorClass.incorrect = 'input-error';
      }
    } catch (error) {
      console.error('Error submitting form:', error);
      submitStatus = 'error';
      formError.incorrect = 'There was an error submitting your feedback. Please try again.';
      formErrorClass.incorrect = 'input-error';
    } finally {
      isSubmitting = false;
    }
  }
</script>

<div class="flex flex-row justify-center items-center h-16">
    <h2 class="text-xl">Leave Feedback</h2>
</div>

<form class="w-11/12 md:w-9/12 lg:w-3/6 p-8 my-0 mx-auto flex flex-col justify-center" on:submit|preventDefault={handleSubmit}>
  <label for="email" class="label">Email:</label>
  <input
    id="email"
    class="input {formErrorClass.email}"
    type="email"
    bind:value={email}
    on:input={validateEmail}
    on:focusout={validateEmail}
    placeholder="Enter your email..."
  />
  {#if formError.email}
    <p class="text-red-500 !mt-0">{formError.email}</p>
  {/if}

  <label for="suggestionType" class="label pt-4">Suggestion Type:</label>
  <select
    id="suggestionType"
    class="select {formErrorClass.suggestionType}"
    bind:value={suggestionType}
    on:change={validateSuggestionType}
    on:focusout={validateSuggestionType}
  >
    {#each suggestionTypes as type}
      <option value={type.value}>{type.label}</option>
    {/each}
  </select>
  {#if formError.suggestionType}
    <p class="text-red-500 !mt-0">{formError.suggestionType}</p>
  {/if}

  <label class="label pt-4">Overall Satisfaction:</label>
  <div class="flex flex-wrap gap-2">
    {#each satisfactionOptions as option}
      <label class="flex items-center space-x-2">
        <input
          type="radio"
          name="overallSatisfaction"
          value={option.value}
          bind:group={overallSatisfaction}
          on:change={validateOverallSatisfaction}
          class="radio"
        />
        <span class="text-sm">{option.value} - {option.label}</span>
      </label>
    {/each}
  </div>
  {#if formError.overallSatisfaction}
    <p class="text-red-500 !mt-0">{formError.overallSatisfaction}</p>
  {/if}

  <label class="label pt-8">Ease of Use (Optional):</label>
  <div class="flex flex-wrap gap-2">
    {#each satisfactionOptions as option}
      <label class="flex items-center space-x-2">
        <input
          type="radio"
          name="easeOfUse"
          value={option.value}
          bind:group={easeOfUse}
          class="radio"
        />
        <span class="text-sm">{option.value} - {option.label}</span>
      </label>
    {/each}
  </div>

  <label class="label pt-8">Feature Completeness (Optional):</label>
  <div class="flex flex-wrap gap-2">
    {#each satisfactionOptions as option}
      <label class="flex items-center space-x-2">
        <input
          type="radio"
          name="featureCompleteness"
          value={option.value}
          bind:group={featureCompleteness}
          class="radio"
        />
        <span class="text-sm">{option.value} - {option.label}</span>
      </label>
    {/each}
  </div>

  <label for="message" class="label pt-4">Detailed Feedback:</label>
  <textarea
    id="message"
    class="textarea {formErrorClass.message}"
    rows="6"
    bind:value={message}
    on:input={validateMessage}
    on:focusout={validateMessage}
    placeholder="Please provide detailed feedback, suggestions, or describe any issues you've encountered..."
  />
  {#if formError.message}
    <p class="text-red-500 !mt-0">{formError.message}</p>
  {/if}

  <button 
    class="btn variant-filled my-4 mx-auto" 
    disabled={isFormInvalid || isSubmitting}
  >
    {isSubmitting ? 'Sending...' : 'Submit Feedback'}
  </button>

  {#if submitStatus === 'success'}
    <p class="text-green-500 !mt-0 text-center">Thank you! Your feedback has been submitted successfully.</p>
  {/if}

  {#if formError.incorrect}
    <p class="text-red-500 !mt-0">{formError.incorrect}</p>
  {/if}
</form>
