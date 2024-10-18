import { describe, it, expect, vi, beforeEach } from 'vitest';
import { sequence } from '@sveltejs/kit/hooks';
import { API_URL } from '$env/static/private';
import { cookieStore } from '$lib/stores/cookieStore';
import { redirect } from '@sveltejs/kit';
import { handle, validatePath, refreshTokenHandle, isPathAllowed } from './hooks.server'; // adjust this import to match your file structure

// Mock dependencies
vi.mock('@sveltejs/kit/hooks', () => ({
  sequence: vi.fn((fn1, fn2) => ({ fn1, fn2 })),
}));

vi.mock('$env/static/private', () => ({
  API_URL: 'http://test-api.com/',
}));

vi.mock('$lib/stores/cookieStore', () => ({
  cookieStore: {
    get: vi.fn(),
  },
}));

vi.mock('@sveltejs/kit', () => ({
  redirect: vi.fn(),
}));

describe('Hook functions', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('isPathAllowed', () => {
    it('FEUT-59: Allowed Paths', () => {
      expect(isPathAllowed('/login')).toBe(true);
      expect(isPathAllowed('/signup')).toBe(true);
      expect(isPathAllowed('/signup/confirm')).toBe(true);
    });

    it('FEUT-60: Not Allowed Paths', () => {
      expect(isPathAllowed('/')).toBe(false);
      expect(isPathAllowed('/dashboard')).toBe(false);
    });
  });

  describe('validatePath', () => {
    it('FEUT-61: Redirect Unauthed User', async () => {
      const mockEvent = {
        request: { url: 'http://example.com/dashboard' },
        cookies: {},
      };
      const mockResolve = vi.fn();

      vi.mocked(cookieStore.get).mockReturnValue(undefined);

      await validatePath({ event: mockEvent, resolve: mockResolve });

      expect(redirect).toHaveBeenCalledWith(302, '/login');
    });

    it('FEUT-62: Redirect Authed User', async () => {
      const mockEvent = {
        request: { url: 'http://example.com/login' },
        cookies: {},
      };
      const mockResolve = vi.fn();

      vi.mocked(cookieStore.get).mockReturnValue('some-token');

      await validatePath({ event: mockEvent, resolve: mockResolve });

      expect(redirect).toHaveBeenCalledWith(302, '/');
    });
  });

  describe('refreshTokenHandle', () => {
    it('FEUT-63: Refresh Token Valid', async () => {
      const mockEvent = {
        request: { url: 'http://example.com/api/something' },
        cookies: {},
      };
      const mockResolve = vi.fn().mockResolvedValue({
        body: 'test body',
        headers: new Headers(),
      });

      vi.mocked(cookieStore.get)
        .mockReturnValueOnce(undefined) // accessToken
        .mockReturnValueOnce('id-token') // idToken
        .mockReturnValueOnce('refresh-token'); // refreshToken


      global.fetch = vi.fn().mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({
          data: {
            AccessToken: 'new-access-token',
            IdToken: 'new-id-token',
            ExpiresIn: 3600,
          },
        }),
      });

      const result = await refreshTokenHandle({ event: mockEvent, resolve: mockResolve });

      expect(global.fetch).toHaveBeenCalledWith(
        'http://test-api.com/v1/token',
        expect.objectContaining({
          method: 'PUT',
          body: JSON.stringify({
            refreshToken: 'refresh-token',
            idToken: 'id-token',
          }),
        })
      );

      expect(result.headers.get('Authorization')).toBe('Bearer new-access-token');
      expect(result.headers.get('Set-Cookie')).toContain('accessToken=new-access-token');
    });

    it('FEUT-64: Do Not Refresh', async () => {
      const mockEvent = {
        request: { url: 'http://example.com/api/something' },
        cookies: {},
      };
      const mockResolve = vi.fn();

      vi.mocked(cookieStore.get).mockReturnValue('existing-token');

      await refreshTokenHandle({ event: mockEvent, resolve: mockResolve });

      expect(mockResolve).toHaveBeenCalledWith(mockEvent);
    });
  });

  describe('handle', () => {
    it('FEUT-65: Validate Path and Handle Token Refresh', () => {
      expect(handle).toEqual({
        fn1: validatePath,
        fn2: refreshTokenHandle,
      });
    });
  });
});
