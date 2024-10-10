import { describe, it, expect, vi, beforeEach } from 'vitest';
import { _confirmUser } from './+page'; // adjust this import to match your file structure
import { goto } from '$app/navigation';


// Mock the modules
vi.mock('$app/navigation', () => ({
  goto: vi.fn(),
}));

vi.mock('$env/static/public', () => ({
  PUBLIC_API_URL: 'http://test-api.com/',
}));

describe('_confirmUser', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    global.fetch = vi.fn();
    global.console.log = vi.fn();
  });

  it('FEUT-35: Successful User Confirmation', async () => {
    global.fetch = vi.fn().mockResolvedValueOnce({
      ok: true,
      json: vi.fn().mockResolvedValueOnce({}),
    });

    const result = await _confirmUser({ confirmationCode: '123456' });

    expect(result).toEqual([true, 'Error']);
    expect(global.fetch).toHaveBeenCalledWith(
      'http://test-api.com/api/v1/user',
      expect.objectContaining({
        method: 'PUT',
        body: JSON.stringify({
          confirmationCode: '123456',
        }),
      })
    );
    expect(goto).toHaveBeenCalledWith('/login');
  });

  it('FEUT-36: Unsuccessful User Confirmation', async () => {
    const errorMessage = 'Invalid confirmation code';
    global.fetch = vi.fn().mockResolvedValueOnce({
      ok: false,
      json: vi.fn().mockResolvedValueOnce({ data: errorMessage }),
    });

    const result = await _confirmUser({ confirmationCode: 'wrongcode' });

    expect(result).toEqual([false, errorMessage]);
    expect(goto).not.toHaveBeenCalled();
  });

  it('FEUT-37: Unsuccessful User Confirmation Network Error', async () => {
    global.fetch = vi.fn().mockRejectedValueOnce(new Error('Network error'));

    const result = await _confirmUser({ confirmationCode: '123456' });

    expect(result).toEqual([false, 'Error']);
    expect(console.log).toHaveBeenCalledWith(new Error('Network error'));
    expect(goto).not.toHaveBeenCalled();
  });
});
