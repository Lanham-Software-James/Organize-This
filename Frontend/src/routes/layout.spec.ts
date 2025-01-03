import { describe, it, expect, vi, beforeEach } from 'vitest';
import { _logoutUser } from './+layout'; // adjust this import to match your file structure
import { goto, invalidateAll } from '$app/navigation';

// Mock the modules
vi.mock('$app/navigation', () => ({
  goto: vi.fn(),
  invalidateAll: vi.fn(),
}));

describe('_logoutUser', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    global.fetch = vi.fn();
    global.console.log = vi.fn();
  });

  it('FEUT-16: Successful Logout', async () => {
    global.fetch = vi.fn().mockResolvedValueOnce({
      ok: true,
    });

    const result = await _logoutUser();

    expect(global.fetch).toHaveBeenCalledWith(
      '/api/v1/token',
      expect.objectContaining({
        method: 'DELETE',
      })
    );
    expect(result).toBe(true);
    expect(invalidateAll).toHaveBeenCalled();
    expect(goto).toHaveBeenCalledWith('/login');
  });

  it('FEUT-17: Unsuccessful Logout', async () => {
    global.fetch = vi.fn().mockResolvedValueOnce({
      ok: false,
    });

    const result = await _logoutUser();

    expect(global.fetch).toHaveBeenCalledWith(
      '/api/v1/token',
      expect.objectContaining({
        method: 'DELETE',
      })
    );
    expect(result).toBe(false);
    expect(invalidateAll).not.toHaveBeenCalled();
    expect(goto).not.toHaveBeenCalled();
  });
});
