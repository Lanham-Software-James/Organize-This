import { describe, it, expect, vi, beforeEach } from 'vitest';
import '@testing-library/jest-dom';
import { fireEvent, render } from '@testing-library/svelte';
import Page from './+page.svelte';
import { _signUpUser } from './+page';

vi.mock('./+page', () => ({
  _signUpUser: vi.fn(),
}));

describe('SignUpForm', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('FEUT-15: Renders User Sign Up Page Header', () => {
    const result = render(Page);

    const headerText = result.getByText('Sign Up');

    expect(headerText).toBeInTheDocument();
  });

  it('FEUT-16: Renders Sign Up User Form Correctly', () => {
    const result = render(Page);
    expect(result.getByLabelText('Username:')).toBeInTheDocument();
    expect(result.getByLabelText('Password:')).toBeInTheDocument();
    expect(result.getByLabelText('Confirm Password:')).toBeInTheDocument();
    expect(result.getByLabelText('First Name:')).toBeInTheDocument();
    expect(result.getByLabelText('Last Name:')).toBeInTheDocument();
    expect(result.getByLabelText('Birthday:')).toBeInTheDocument();
    expect(result.getByRole('button', { name: 'Submit' })).toBeInTheDocument();
  });

  it('FEUT-17: Renders Sign Up User Submit Button Disabled', () => {
    const result = render(Page);
    expect(result.getByRole('button', { name: 'Submit' })).toBeDisabled();
  });

  it('FEUT-18: Sign Up User Username Field Validation', async () => {
    const result = render(Page);
    const input = result.getByLabelText('Username:');

    await fireEvent.input(input, { target: { value: '' } });
    await fireEvent.focusOut(input);
    expect(result.getByText('Username is required!')).toBeInTheDocument();

    await fireEvent.input(input, { target: { value: 'testuser' } });
    expect(result.queryByText('Username is required!')).not.toBeInTheDocument();
  });

  it('FEUT-19: Sign Up User Password Field Validation', async () => {
    const result = render(Page);
    const passwordInput = result.getByLabelText('Password:');
    const confirmPasswordInput = result.getByLabelText('Confirm Password:');

    await fireEvent.input(passwordInput, { target: { value: '' } });
    await fireEvent.focusOut(passwordInput);
    expect(result.getByText('Password is required!')).toBeInTheDocument();

    await fireEvent.input(passwordInput, { target: { value: 'password123' } });
    await fireEvent.input(confirmPasswordInput, { target: { value: 'password456' } });
    await fireEvent.focusOut(confirmPasswordInput);
    expect(result.getByText('Passwords must match!')).toBeInTheDocument();

    await fireEvent.input(confirmPasswordInput, { target: { value: 'password123' } });
    expect(result.queryByText('Passwords must match!')).not.toBeInTheDocument();
  });

  it('FEUT-20: Sign Up User First Name Field Validation', async () => {
    const result = render(Page);
    const input = result.getByLabelText('First Name:');

    await fireEvent.input(input, { target: { value: '' } });
    await fireEvent.focusOut(input);
    expect(result.getByText('First name is required!')).toBeInTheDocument();

    await fireEvent.input(input, { target: { value: 'John' } });
    expect(result.queryByText('First name is required!')).not.toBeInTheDocument();
  });

  it('FEUT-21: Sign Up User Last Name Field Validation', async () => {
    const result = render(Page);
    const input = result.getByLabelText('Last Name:');

    await fireEvent.input(input, { target: { value: '' } });
    await fireEvent.focusOut(input);
    expect(result.getByText('Last name is required!')).toBeInTheDocument();

    await fireEvent.input(input, { target: { value: 'Doe' } });
    expect(result.queryByText('Last name is required!')).not.toBeInTheDocument();
  });

  it('FEUT-22: Sign Up User Birthday Field Validation', async () => {
    const result = render(Page);
    const input = result.getByLabelText('Birthday:');

    await fireEvent.input(input, { target: { value: '' } });
    await fireEvent.focusOut(input);
    expect(result.getByText('Birthday is required!')).toBeInTheDocument();

    await fireEvent.input(input, { target: { value: '1990-01-01' } });
    expect(result.queryByText('Birthday is required!')).not.toBeInTheDocument();
  });

  it('FEUT-23: Sign Up User Form Validation', async () => {
    const result = render(Page);
    const usernameInput = result.getByLabelText('Username:');
    const passwordInput = result.getByLabelText('Password:');
    const confirmPasswordInput = result.getByLabelText('Confirm Password:');
    const firstNameInput = result.getByLabelText('First Name:');
    const lastNameInput = result.getByLabelText('Last Name:');
    const birthdayInput = result.getByLabelText('Birthday:');
    const submitButton = result.getByRole('button', { name: 'Submit' });

    await fireEvent.input(usernameInput, { target: { value: 'testuser' } });
    await fireEvent.input(passwordInput, { target: { value: 'password123' } });
    await fireEvent.input(confirmPasswordInput, { target: { value: 'password123' } });
    await fireEvent.input(firstNameInput, { target: { value: 'John' } });
    await fireEvent.input(lastNameInput, { target: { value: 'Doe' } });
    await fireEvent.input(birthdayInput, { target: { value: '1990-01-01' } });

    expect(submitButton).not.toBeDisabled();
  });

  it('FEUT-24: Verify Successful Sign Up User Form Submission', async () => {
    vi.mocked(_signUpUser).mockResolvedValue([true, '']);
    const result = render(Page);

    const usernameInput = result.getByLabelText('Username:');
    const passwordInput = result.getByLabelText('Password:');
    const confirmPasswordInput = result.getByLabelText('Confirm Password:');
    const firstNameInput = result.getByLabelText('First Name:');
    const lastNameInput = result.getByLabelText('Last Name:');
    const birthdayInput = result.getByLabelText('Birthday:');
    const submitButton = result.getByRole('button', { name: 'Submit' });

    await fireEvent.input(usernameInput, { target: { value: 'testuser' } });
    await fireEvent.input(passwordInput, { target: { value: 'password123' } });
    await fireEvent.input(confirmPasswordInput, { target: { value: 'password123' } });
    await fireEvent.input(firstNameInput, { target: { value: 'John' } });
    await fireEvent.input(lastNameInput, { target: { value: 'Doe' } });
    await fireEvent.input(birthdayInput, { target: { value: '1990-01-01' } });

    await fireEvent.click(submitButton);

    expect(_signUpUser).toHaveBeenCalledWith({
      userEmail: 'testuser',
      password: 'password123',
      confirmPassword: 'password123',
      firstName: 'John',
      lastName: 'Doe',
      birthday: '1990-01-01',
    });
    expect(result.queryByText(/error/i)).not.toBeInTheDocument();
  });

  it('FEUT-25: Verify Unsuccessful Sign Up User Form Submission', async () => {
    vi.mocked(_signUpUser).mockResolvedValue([false, 'Registration failed']);
    const result = render(Page);

    const usernameInput = result.getByLabelText('Username:');
    const passwordInput = result.getByLabelText('Password:');
    const confirmPasswordInput = result.getByLabelText('Confirm Password:');
    const firstNameInput = result.getByLabelText('First Name:');
    const lastNameInput = result.getByLabelText('Last Name:');
    const birthdayInput = result.getByLabelText('Birthday:');
    const submitButton = result.getByRole('button', { name: 'Submit' });

    await fireEvent.input(usernameInput, { target: { value: 'testuser' } });
    await fireEvent.input(passwordInput, { target: { value: 'password123' } });
    await fireEvent.input(confirmPasswordInput, { target: { value: 'password123' } });
    await fireEvent.input(firstNameInput, { target: { value: 'John' } });
    await fireEvent.input(lastNameInput, { target: { value: 'Doe' } });
    await fireEvent.input(birthdayInput, { target: { value: '1990-01-01' } });

    await fireEvent.click(submitButton);

    expect(result.getByText('Registration failed')).toBeInTheDocument();
  });
});