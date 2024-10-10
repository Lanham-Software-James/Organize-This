import { describe, it, expect, vi, beforeEach } from 'vitest';
import '@testing-library/jest-dom';
import { fireEvent, render } from '@testing-library/svelte';
import Page from './+page.svelte';
import { _confirmUser } from './+page';

vi.mock('./+page', () => ({
  _confirmUser: vi.fn(),
}));

describe('Home route', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('FEUT-29: Renders Confirm User Page Header', () => {
    const result = render(Page);

    const headerText = result.getByText('Confirm User Account');

    expect(headerText).toBeInTheDocument();
  });

  it('FEUT-30: Renders Confirm User Form', () => {
    const result = render(Page);
    expect(result.getByText('Confirm User Account')).toBeInTheDocument();
    expect(result.getByLabelText('Confirmation Code:')).toBeInTheDocument();
    expect(result.getByRole('button', { name: 'Submit' })).toBeInTheDocument();
  });

  it('FEUT-31: Renders Confirm User Submit Button Disabled', () => {
    const result = render(Page);
    expect(result.getByRole('button', { name: 'Submit' })).toBeDisabled();
  });

  it('FEUT-32: Confirm User Confirmation Code', async () => {
    const result = render(Page);
    const input = result.getByLabelText('Confirmation Code:');
    const submitButton = result.getByRole('button', { name: 'Submit' });

    // Empty input
    await fireEvent.input(input, { target: { value: '' } });
    await fireEvent.focusOut(input);
    expect(result.getByText('Confirmation code is required!')).toBeInTheDocument();
    expect(submitButton).toBeDisabled();

    // Valid input
    await fireEvent.input(input, { target: { value: '123456' } });
    expect(result.queryByText('Confirmation code is required!')).not.toBeInTheDocument();
    expect(submitButton).not.toBeDisabled();
  });

  it('FEUT-33: Confirm User Successful Submission', async () => {
    vi.mocked(_confirmUser).mockResolvedValue([true, '']);
    const result = render(Page);
    const input = result.getByLabelText('Confirmation Code:');
    const submitButton = result.getByRole('button', { name: 'Submit' });

    await fireEvent.input(input, { target: { value: '123456' } });
    await fireEvent.click(submitButton);

    expect(_confirmUser).toHaveBeenCalledWith({ confirmationCode: '123456' });
    expect(result.queryByText(/error/i)).not.toBeInTheDocument();
  });

  it('FEUT-34: Confirm User Unsuccessful Submission', async () => {
    vi.mocked(_confirmUser).mockResolvedValue([false, 'Invalid confirmation code']);
    const result = render(Page);
    const input = result.getByLabelText('Confirmation Code:');
    const submitButton = result.getByRole('button', { name: 'Submit' });

    await fireEvent.input(input, { target: { value: '123456' } });
    await fireEvent.click(submitButton);

    expect(_confirmUser).toHaveBeenCalledWith({ confirmationCode: '123456' });
    expect(result.getByText('Invalid confirmation code')).toBeInTheDocument();
  });

});
