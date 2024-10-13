import { describe, it, expect, vi, beforeEach } from 'vitest';
import '@testing-library/jest-dom';
import { fireEvent, render } from '@testing-library/svelte';
import Page from './+page.svelte';
import { _loginUser } from './+page';

vi.mock('./+page', () => ({
    _loginUser: vi.fn(),
}));

describe('LoginForm', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('FEUT-31: Renders Login Form', () => {
        const result = render(Page);
        expect(result.getByText('Login')).toBeInTheDocument();
        expect(result.getByLabelText('Username:')).toBeInTheDocument();
        expect(result.getByLabelText('Password:')).toBeInTheDocument();
        expect(result.getByRole('button', { name: 'Submit' })).toBeInTheDocument();
        expect(result.getByRole('button', { name: 'Submit' })).toBeDisabled();
    });

    it('FEUT-32: Login Username', async () => {
        const result = render(Page);
        const input = result.getByLabelText('Username:');

        await fireEvent.input(input, { target: { value: '' } });
        await fireEvent.focusOut(input);
        expect(result.getByText('Username is required!')).toBeInTheDocument();

        await fireEvent.input(input, { target: { value: 'testuser' } });
        expect(result.queryByText('Username is required!')).not.toBeInTheDocument();
    });

    it('FEUT-33: Login Password', async () => {
        const result = render(Page);
        const input = result.getByLabelText('Password:');

        await fireEvent.input(input, { target: { value: '' } });
        await fireEvent.focusOut(input);
        expect(result.getByText('Password is required!')).toBeInTheDocument();

        await fireEvent.input(input, { target: { value: 'password123' } });
        expect(result.queryByText('Password is required!')).not.toBeInTheDocument();
    });

    it('FEUT-34: Login Form Validation', async () => {
        const result = render(Page);
        const usernameInput = result.getByLabelText('Username:');
        const passwordInput = result.getByLabelText('Password:');
        const submitButton = result.getByRole('button', { name: 'Submit' });

        await fireEvent.input(usernameInput, { target: { value: 'testuser' } });
        await fireEvent.input(passwordInput, { target: { value: 'password123' } });

        expect(submitButton).not.toBeDisabled();
    });

    it('FEUT-35: Login Successful Submission', async () => {
        vi.mocked(_loginUser).mockResolvedValue([true, '']);
        const result = render(Page);

        const usernameInput = result.getByLabelText('Username:');
        const passwordInput = result.getByLabelText('Password:');
        const submitButton = result.getByRole('button', { name: 'Submit' });

        await fireEvent.input(usernameInput, { target: { value: 'testuser' } });
        await fireEvent.input(passwordInput, { target: { value: 'password123' } });
        await fireEvent.click(submitButton);

        expect(_loginUser).toHaveBeenCalledWith({
            userEmail: 'testuser',
            password: 'password123',
        });
        expect(result.queryByText(/error/i)).not.toBeInTheDocument();
    });

    it('FEUT-36: Login Unsuccessful Submission', async () => {
        vi.mocked(_loginUser).mockResolvedValue([false, 'Invalid credentials']);
        const result = render(Page);

        const usernameInput = result.getByLabelText('Username:');
        const passwordInput = result.getByLabelText('Password:');
        const submitButton = result.getByRole('button', { name: 'Submit' });

        await fireEvent.input(usernameInput, { target: { value: 'testuser' } });
        await fireEvent.input(passwordInput, { target: { value: 'wrongpassword' } });
        await fireEvent.click(submitButton);

        expect(result.getByText('Invalid credentials')).toBeInTheDocument();
    });
});
