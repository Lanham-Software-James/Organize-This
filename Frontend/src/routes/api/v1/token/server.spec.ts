import { describe, it, expect, vi, beforeEach } from 'vitest';
import { POST, DELETE } from './+server';
import { cookieStore } from '$lib/stores/cookieStore';

// Mock the API_URL
vi.mock('$env/static/private', () => ({
  API_URL: 'http://mock-api.com/'
}));

// Mock the cookieStore
vi.mock('$lib/stores/cookieStore', () => ({
  cookieStore: {
    get: vi.fn(),
    set: vi.fn(),
    delete: vi.fn()
  }
}));

describe('POST function', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    global.fetch = vi.fn();
  });

  it('FEUT-43: Login Sever Request Success', async () => {
    const mockRequest = {
      json: vi.fn().mockResolvedValue({ userEmail: 'test@example.com', password: 'password123' })
    };
    const mockCookies = {};
    const mockResponseData = {
      data: {
        AccessToken: 'mockAccessToken',
        IdToken: 'mockIdToken',
        RefreshToken: 'mockRefreshToken',
        ExpiresIn: 3600
      }
    };

    global.fetch = vi.fn().mockResolvedValue({
      status: 200,
      json: vi.fn().mockResolvedValue(mockResponseData)
    });

    const response = await POST({ request: mockRequest, cookies: mockCookies });

    expect(global.fetch).toHaveBeenCalledWith(
      'http://mock-api.com/v1/token',
      expect.objectContaining({
        method: 'POST',
        body: JSON.stringify({ userEmail: 'test@example.com', password: 'password123' })
      })
    );

    expect(cookieStore.set).toHaveBeenCalledTimes(3);
    expect(cookieStore.set).toHaveBeenCalledWith(mockCookies, 'accessToken', 'mockAccessToken', expect.any(Object));
    expect(cookieStore.set).toHaveBeenCalledWith(mockCookies, 'idToken', 'mockIdToken', expect.any(Object));
    expect(cookieStore.set).toHaveBeenCalledWith(mockCookies, 'refreshToken', 'mockRefreshToken', expect.any(Object));

    expect(response.status).toBe(200);
    expect(await response.json()).toEqual(mockResponseData);
  });

  it('FEUT-44: Login Sever Request Unsuccess', async () => {
    const mockRequest = {
      json: vi.fn().mockResolvedValue({})
    };
    const mockCookies = {};

    global.fetch = vi.fn().mockRejectedValue(new Error('Network error'));
    console.error = vi.fn();

    const response = await POST({ request: mockRequest, cookies: mockCookies });

    expect(console.error).toHaveBeenCalledWith(new Error('Network error'));
    expect(response.status).toBe(400);
    expect(await response.text()).toBe('{}');
  });
});

describe('DELETE function', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    global.fetch = vi.fn();
  });

  it('FEUT-45: Logout Sever Request Success', async () => {
    const mockCookies = {};
    const mockResponseData = { message: 'Logged out successfully' };

    vi.mocked(cookieStore.get).mockReturnValue('mockRefreshToken');

    global.fetch = vi.fn().mockResolvedValue({
      status: 200,
      json: vi.fn().mockResolvedValue(mockResponseData)
    });

    const response = await DELETE({ cookies: mockCookies });

    expect(global.fetch).toHaveBeenCalledWith(
      'http://mock-api.com/v1/token',
      expect.objectContaining({
        method: 'DELETE',
        body: JSON.stringify({ refreshToken: 'mockRefreshToken' })
      })
    );

    expect(cookieStore.delete).toHaveBeenCalledTimes(3);
    expect(cookieStore.delete).toHaveBeenCalledWith(mockCookies, 'accessToken', expect.any(Object));
    expect(cookieStore.delete).toHaveBeenCalledWith(mockCookies, 'idToken', expect.any(Object));
    expect(cookieStore.delete).toHaveBeenCalledWith(mockCookies, 'refreshToken', expect.any(Object));

    expect(response.status).toBe(200);
    expect(await response.json()).toEqual(mockResponseData);
  });

  it('FEUT-46: Logout Sever Request Unsuccess', async () => {
    const mockCookies = {};

    global.fetch = vi.fn().mockRejectedValue(new Error('Network error'));
    console.error = vi.fn();

    const response = await DELETE({ cookies: mockCookies });

    expect(console.error).toHaveBeenCalledWith(new Error('Network error'));
    expect(response.status).toBe(400);
    expect(await response.text()).toBe('{}');
  });
});
