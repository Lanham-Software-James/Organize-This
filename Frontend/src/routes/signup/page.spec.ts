import { describe, it, expect, vi, beforeEach } from 'vitest';
import { _signUpUser } from './+page'; // adjust this import to match your file structure
import { goto } from '$app/navigation';

// Mock the modules
vi.mock('$app/navigation', () => ({
  goto: vi.fn(),
}));

vi.mock('$env/static/public', () => ({
  PUBLIC_API_URL: 'http://test-api.com/',
}));

describe('_signUpUser', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    global.fetch = vi.fn();
    global.console.log = vi.fn();
  });

  it('FEUT-26: Successful Sign Up', async () => {
    global.fetch = vi.fn().mockResolvedValueOnce({
      ok: true,
      json: vi.fn().mockResolvedValueOnce({}),
    });

    const formData = {
      userEmail: 'test@example.com',
      password: 'password123',
      firstName: 'John',
      lastName: 'Doe',
      birthday: '1990-01-01'
    };

    const result = await _signUpUser(formData);

    expect(result).toEqual([true, 'Error']);
    expect(global.fetch).toHaveBeenCalledWith(
      'http://test-api.com/api/v1/user',
      expect.objectContaining({
        method: 'POST',
        body: JSON.stringify(formData),
      })
    );
    expect(goto).toHaveBeenCalledWith('/signup/confirm');
  });

  it('FEUT-27: Unsuccessful Sign Up', async () => {
    const errorMessage = 'Email already exists';
    global.fetch = vi.fn().mockResolvedValueOnce({
      ok: false,
      json: vi.fn().mockResolvedValueOnce({ data: errorMessage }),
    });

    const formData = {
      userEmail: 'existing@example.com',
      password: 'password123',
      firstName: 'Jane',
      lastName: 'Doe',
      birthday: '1995-05-05'
    };

    const result = await _signUpUser(formData);

    expect(result).toEqual([false, errorMessage]);
    expect(goto).not.toHaveBeenCalled();
  });

  it('FEUT-28: Unsuccessful Sign Up Network Error', async () => {
    global.fetch = vi.fn().mockRejectedValueOnce(new Error('Network error'));

    const formData = {
      userEmail: 'test@example.com',
      password: 'password123',
      firstName: 'John',
      lastName: 'Doe',
      birthday: '1990-01-01'
    };

    const result = await _signUpUser(formData);

    expect(result).toEqual([false, 'Error']);
    expect(console.log).toHaveBeenCalledWith(new Error('Network error'));
    expect(goto).not.toHaveBeenCalled();
  });
});