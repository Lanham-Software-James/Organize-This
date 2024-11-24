import { describe, it, expect, vi, beforeEach } from 'vitest';
import { DELETE, GET } from './+server';
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

describe('GET function', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    global.fetch = vi.fn();
  });

  it('FEUT-69: Get Entity Server Request Success', async () => {
    const validCategories = ['item', 'container', 'shelf', 'shelving_unit', 'room'];

    for (const category of validCategories) {
      const mockCookies = {};
      const mockId = '123';
      const mockResponseData = { data: `mock${category}Data` };

      vi.mocked(cookieStore.get).mockReturnValue('mockAccessToken');

      global.fetch = vi.fn().mockResolvedValue({
        status: 200,
        json: vi.fn().mockResolvedValue(mockResponseData)
      });

      const response = await GET({ params: { category, id: mockId }, cookies: mockCookies });

      expect(global.fetch).toHaveBeenCalledWith(
        `http://mock-api.com/v1/entity/${category}/${mockId}`,
        expect.objectContaining({
          headers: {
            Authorization: 'Bearer mockAccessToken'
          }
        })
      );

      expect(response.status).toBe(200);
      expect(await response.json()).toEqual(mockResponseData);

      vi.clearAllMocks();
    }
  });

  it('FEUT-70: Get Entity Server Request Invalid Category', async () => {
    const mockCookies = {};
    const invalidCategory = 'invalid';
    const mockId = '123';

    const response = await GET({ params: { category: invalidCategory, id: mockId }, cookies: mockCookies });

    expect(global.fetch).not.toHaveBeenCalled();
    expect(response).toBeInstanceOf(Response);
    expect(response.status).toBe(200);
    expect(await response.text()).toBe('');
  });

  it('FEUT-71: Get Entity Server Request Bad Request', async () => {
    const mockCookies = {};
    const category = 'item';
    const mockId = '123';

    vi.mocked(cookieStore.get).mockReturnValue('mockAccessToken');

    global.fetch = vi.fn().mockRejectedValue(new Error('Network error'));
    console.error = vi.fn();

    const response = await GET({ params: { category, id: mockId }, cookies: mockCookies });

    expect(console.error).toHaveBeenCalledWith(new Error('Network error'));
    expect(response.status).toBe(400);
    const responseBody = await response.json();
    expect(responseBody).toEqual({});
  });
});

describe('DELETE function', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    global.fetch = vi.fn();
  });

  it('FEUT-72: Delete Entity Server Request Success', async () => {
    const validCategories = ['item', 'container', 'shelf', 'shelving_unit', 'room', 'building'];

    for (const category of validCategories) {
      const mockCookies = {};
      const mockId = '123';
      const mockResponseData = { message: 'Entity deleted successfully' };

      vi.mocked(cookieStore.get).mockReturnValue('mockAccessToken');

      global.fetch = vi.fn().mockResolvedValue({
        status: 200,
        json: vi.fn().mockResolvedValue(mockResponseData)
      });

      const response = await DELETE({ params: { category, id: mockId }, cookies: mockCookies });

      expect(global.fetch).toHaveBeenCalledWith(
        `http://mock-api.com/v1/entity/${category}/${mockId}`,
        expect.objectContaining({
          method: 'DELETE',
          headers: {
            'Authorization': 'Bearer mockAccessToken',
          }
        })
      );

      expect(response.status).toBe(200);
      expect(await response.json()).toEqual(mockResponseData);

      vi.clearAllMocks();
    }
  });

  it('FEUT-73: Delete Entity Server Request Invalid Category', async () => {
    const mockCookies = {};
    const invalidCategory = 'invalid';
    const mockId = '123';

    const response = await DELETE({ params: { category: invalidCategory, id: mockId }, cookies: mockCookies });

    expect(global.fetch).not.toHaveBeenCalled();
    expect(response).toBeInstanceOf(Response);
    expect(response.status).toBe(200);
    expect(await response.text()).toBe('');
  });

  it('FEUT-74: Delete Entity Server Bad Request', async () => {
    const mockCookies = {};
    const category = 'item';
    const mockId = '123';

    vi.mocked(cookieStore.get).mockReturnValue('mockAccessToken');

    global.fetch = vi.fn().mockRejectedValue(new Error('Network error'));
    console.error = vi.fn();

    const response = await DELETE({ params: { category, id: mockId }, cookies: mockCookies });

    expect(console.error).toHaveBeenCalledWith(new Error('Network error'));
    expect(response.status).toBe(400);
    const responseBody = await response.json();
    expect(responseBody).toEqual({});
  });
});
