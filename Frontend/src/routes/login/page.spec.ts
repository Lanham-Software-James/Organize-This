import { describe, it, expect, vi, beforeEach } from 'vitest';
import { _loginUser } from './+page'; // adjust this import to match your file structure
import { goto, invalidateAll } from '$app/navigation';

// Mock the modules
vi.mock('$app/navigation', () => ({
  goto: vi.fn(),
  invalidateAll: vi.fn(),
}));

vi.mock('$env/static/public', () => ({
  PUBLIC_API_URL: 'http://test-api.com/',
}));

describe('_loginUser', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    global.fetch = vi.fn();
    global.console.log = vi.fn();
  });

  it('FEUT-43: Successful Login', async () => {
    global.fetch = vi.fn().mockResolvedValueOnce({
      ok: true,
      json: vi.fn().mockResolvedValueOnce({}),
    });

    const result = await _loginUser({ userEmail: 'test@example.com', password: 'password123' });

    expect(result).toEqual([true, 'Error']);
    expect(global.fetch).toHaveBeenCalledWith(
      'http://test-api.com/api/v1/token',
      expect.objectContaining({
        method: 'POST',
        body: JSON.stringify({
          userEmail: 'test@example.com',
          password: 'password123',
        }),
      })
    );
    expect(invalidateAll).toHaveBeenCalled();
    expect(goto).toHaveBeenCalledWith('/');
  });

  it('FEUT-44: Unsuccessful Login', async () => {
    const errorMessage = 'Invalid credentials';
    global.fetch = vi.fn().mockResolvedValueOnce({
      ok: false,
      json: vi.fn().mockResolvedValueOnce({ data: errorMessage }),
    });

    const result = await _loginUser({ userEmail: 'test@example.com', password: 'wrongpassword' });

    expect(result).toEqual([false, errorMessage]);
    expect(invalidateAll).not.toHaveBeenCalled();
    expect(goto).not.toHaveBeenCalled();
  });

  it('FEUT-45: Unsuccessful Login Network Error', async () => {
    global.fetch = vi.fn().mockRejectedValueOnce(new Error('Network error'));

    const result = await _loginUser({ userEmail: 'test@example.com', password: 'password123' });

    expect(result).toEqual([false, 'Error']);
    expect(console.log).toHaveBeenCalledWith(new Error('Network error'));
    expect(invalidateAll).not.toHaveBeenCalled();
    expect(goto).not.toHaveBeenCalled();
  });
});
