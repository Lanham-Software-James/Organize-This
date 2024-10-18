import { describe, it, expect, vi, beforeEach } from 'vitest';
import { GET } from './+server';
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

  it('FEUT-66 Get Parents Server Request Success', async () => {
    const validCategories = ['item', 'container', 'shelf', 'shelving_unit', 'room'];

    for (const category of validCategories) {
      const mockCookies = {};
      const mockResponseData = { data: `mock${category}Data` };

      vi.mocked(cookieStore.get).mockReturnValue('mockAccessToken');

      global.fetch = vi.fn().mockResolvedValue({
        status: 200,
        json: vi.fn().mockResolvedValue(mockResponseData)
      });

      const response = await GET({ params: { category }, cookies: mockCookies });

      expect(global.fetch).toHaveBeenCalledWith(
        `http://mock-api.com/v1/parents/${category}`,
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

  it('FEUT-67: Get Parents Server Request Invalid Category', async () => {
    const mockCookies = {};
    const invalidCategory = 'invalid';

    const response = await GET({ params: { category: invalidCategory }, cookies: mockCookies });

    expect(global.fetch).not.toHaveBeenCalled();
    expect(response).toBeInstanceOf(Response);
    expect(response.status).toBe(200);
    expect(await response.text()).toBe('');
  });

  it('FEUT-68: Get Parents Server Request Unsuccess', async () => {
    const mockCookies = {};
    const category = 'item';

    vi.mocked(cookieStore.get).mockReturnValue('mockAccessToken');

    global.fetch = vi.fn().mockRejectedValue(new Error('Network error'));
    console.error = vi.fn();

    const response = await GET({ params: { category }, cookies: mockCookies });

    expect(console.error).toHaveBeenCalledWith(new Error('Network error'));
    expect(response.status).toBe(400);
    const responseBody = await response.json();
    expect(responseBody).toEqual({});
  });
});
