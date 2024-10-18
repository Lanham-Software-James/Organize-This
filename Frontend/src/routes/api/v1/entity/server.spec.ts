import { describe, it, expect, vi, beforeEach } from 'vitest';
import { POST, PUT } from './+server'; // adjust this import to match your file structure
import { cookieStore } from '$lib/stores/cookieStore';

// Mock the API_URL
vi.mock('$env/static/private', () => ({
  API_URL: 'http://mock-api.com/'
}));

// Mock the cookieStore
vi.mock('$lib/stores/cookieStore', () => ({
  cookieStore: {
    get: vi.fn()
  }
}));

describe('POST function', () => {
  beforeEach(() => {
    // Clear all mocks before each test
    vi.clearAllMocks();

    // Reset fetch mock
    global.fetch = vi.fn();
  });

  it('FEUT-53: Create Entity Sever Request Success', async () => {
    const mockCookies = {};
    const mockRequestBody = {
      address: '123 Test St',
      category: 'Test Category',
      name: 'Test Name',
      notes: 'Test Notes'
    };

    const mockRequest = {
      json: vi.fn().mockResolvedValue(mockRequestBody)
    };

    // Mock cookieStore.get to return a token
    vi.mocked(cookieStore.get).mockReturnValue('mock-token');

    // Mock fetch to return a successful response
    global.fetch = vi.fn().mockResolvedValue(new Response('{"data": "mock data"}', { status: 200 }));

    const response = await POST({ request: mockRequest, cookies: mockCookies });

    // Check if fetch was called with correct arguments
    expect(global.fetch).toHaveBeenCalledWith(
      'http://mock-api.com/v1/entity',
      {
        method: 'POST',
        headers: {
          Authorization: 'Bearer mock-token'
        },
        body: JSON.stringify(mockRequestBody)
      }
    );

    // Check if the response is correctly returned
    expect(response.status).toBe(200);
    expect(await response.json()).toEqual({ data: 'mock data' });
  });

  it('FEUT-54: Create Entity Sever Request Unsuccess', async () => {
    const mockCookies = {};
    const mockRequest = {
      json: vi.fn().mockResolvedValue({})
    };

    // Mock fetch to throw an error
    global.fetch = vi.fn().mockRejectedValue(new Error('Network error'));

    // Mock console.error to check if it's called
    console.error = vi.fn();

    const response = await POST({ request: mockRequest, cookies: mockCookies });

    // Check if console.error was called
    expect(console.error).toHaveBeenCalledWith(new Error('Network error'));

    // Check if an empty response is returned
    expect(response.status).toBe(400);
    expect(await response.text()).toBe('{}');
  });
});

describe('PUT function', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    global.fetch = vi.fn();
  });

  it('FEUT-55: Update Entity Server Request Success', async () => {
    const mockCookies = {};
    const mockRequestBody = {
      id: '456',
      address: '456 Update St',
      category: 'Updated Category',
      name: 'Updated Name',
      notes: 'Updated Notes',
      parentID: '789',
      parentCategory: 'Updated Parent Category'
    };

    const mockRequest = {
      json: vi.fn().mockResolvedValue(mockRequestBody)
    };

    vi.mocked(cookieStore.get).mockReturnValue('mock-token');
    global.fetch = vi.fn().mockResolvedValue(new Response('{"data": "updated mock data"}', { status: 200 }));

    const response = await PUT({ request: mockRequest, cookies: mockCookies });

    expect(global.fetch).toHaveBeenCalledWith(
      'http://mock-api.com/v1/entity',
      {
        method: 'PUT',
        headers: {
          Authorization: 'Bearer mock-token'
        },
        body: JSON.stringify({...mockRequestBody, id: '456'})
      }
    );

    expect(response.status).toBe(200);
    expect(await response.json()).toEqual({ data: 'updated mock data' });
  });

  it('FEUT-56: Update Entity Server Request Unsuccess', async () => {
    const mockCookies = {};
    const mockRequest = {
      json: vi.fn().mockResolvedValue({id: '456'})
    };

    global.fetch = vi.fn().mockRejectedValue(new Error('Network error'));
    console.error = vi.fn();

    const response = await PUT({ request: mockRequest, cookies: mockCookies });

    expect(console.error).toHaveBeenCalledWith(new Error('Network error'));
    expect(response.status).toBe(400);
    expect(await response.text()).toBe(JSON.stringify(new Error('Network error')));
  });
});
