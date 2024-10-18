import { describe, it, expect, vi, beforeEach } from 'vitest';
import { POST, PUT } from './+server';
import { cookieStore } from '$lib/stores/cookieStore.js';

// Mock the API_URL
vi.mock('$env/static/private', () => ({
  API_URL: 'http://mock-api.com/'
}));

// Mock the cookieStore
vi.mock('$lib/stores/cookieStore.js', () => ({
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

  it('FEUT-45: Create User Sever Request Success', async () => {
    const mockRequest = {
      json: vi.fn().mockResolvedValue({
        userEmail: 'test@example.com',
        password: 'password123',
        firstName: 'John',
        lastName: 'Doe',
        birthday: '1990-01-01'
      })
    };
    const mockCookies = {};
    const mockResponseData = { message: 'User created successfully' };

    global.fetch = vi.fn().mockResolvedValue({
      status: 200,
      json: vi.fn().mockResolvedValue(mockResponseData)
    });

    const response = await POST({ request: mockRequest, cookies: mockCookies });

    expect(global.fetch).toHaveBeenCalledWith(
      'http://mock-api.com/v1/user',
      expect.objectContaining({
        method: 'POST',
        body: JSON.stringify({
          userEmail: 'test@example.com',
          password: 'password123',
          firstName: 'John',
          lastName: 'Doe',
          birthday: '1990-01-01'
        })
      })
    );

    expect(cookieStore.set).toHaveBeenCalledWith(mockCookies, 'userEmail', 'test@example.com', expect.any(Object));

    expect(response.status).toBe(200);
    expect(await response.json()).toEqual(mockResponseData);
  });

  it('FEUT-46: Create User Sever Request Unsuccess', async () => {
    const mockRequest = {
      json: vi.fn().mockResolvedValue({})
    };
    const mockCookies = {};

    global.fetch = vi.fn().mockRejectedValue(new Error('Network error'));
    console.error = vi.fn();

    const response = await POST({ request: mockRequest, cookies: mockCookies });

    expect(console.error).toHaveBeenCalledWith(new Error('Network error'));
    expect(response.status).toBe(400);
    expect(await response.text()).toBe("{}");
  });
});

describe('PUT function', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    global.fetch = vi.fn();
  });

  it('FEUT-47: Confirm User Sever Request Success', async () => {
    const mockRequest = {
      json: vi.fn().mockResolvedValue({
        confirmationCode: '123456'
      })
    };
    const mockCookies = {};
    const mockResponseData = { message: 'User confirmed successfully' };

    vi.mocked(cookieStore.get).mockReturnValue('test@example.com');

    global.fetch = vi.fn().mockResolvedValue({
      status: 200,
      json: vi.fn().mockResolvedValue(mockResponseData)
    });

    const response = await PUT({ request: mockRequest, cookies: mockCookies });

    expect(global.fetch).toHaveBeenCalledWith(
      'http://mock-api.com/v1/user',
      expect.objectContaining({
        method: 'PUT',
        body: JSON.stringify({
          userEmail: 'test@example.com',
          confirmationCode: '123456'
        })
      })
    );

    expect(cookieStore.delete).toHaveBeenCalledWith(mockCookies, 'userEmail', expect.any(Object));

    expect(response.status).toBe(200);
    expect(await response.json()).toEqual(mockResponseData);
  });

  it('FEUT-48: Confirm User Sever Request Unsuccess', async () => {
    const mockRequest = {
      json: vi.fn().mockResolvedValue({})
    };
    const mockCookies = {};

    global.fetch = vi.fn().mockRejectedValue(new Error('Network error'));
    console.error = vi.fn();

    const response = await PUT({ request: mockRequest, cookies: mockCookies });

    expect(console.error).toHaveBeenCalledWith(new Error('Network error'));
    expect(response.status).toBe(400);
    expect(await response.text()).toBe("{}");
  });
});
