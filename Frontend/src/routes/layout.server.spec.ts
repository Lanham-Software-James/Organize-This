import { describe, it, expect, vi, beforeEach } from 'vitest';
import { load } from './+layout.server'; // adjust this import to match your file structure
import { cookieStore } from '$lib/stores/cookieStore';

// Mock the cookieStore
vi.mock('$lib/stores/cookieStore', () => ({
  cookieStore: {
    get: vi.fn(),
  },
}));

describe('Layout load function', () => {
  let mockCookies: any;

  beforeEach(() => {
    vi.clearAllMocks();
    mockCookies = {
      get: vi.fn(),
      set: vi.fn(),
    };
  });

  it('FEUT-13: Cookie Exists', async () => {
    vi.mocked(cookieStore.get).mockReturnValue('some-refresh-token');

    //@ts-ignore
    const result = await load({ cookies: mockCookies });

    expect(cookieStore.get).toHaveBeenCalledWith(mockCookies, 'refreshToken');
    expect(result).toEqual({ cookieExists: true });
  });

  it('FEUT-14: Cookie Does Not Exist', async () => {
    vi.mocked(cookieStore.get).mockReturnValue(undefined);

    //@ts-ignore
    const result = await load({ cookies: mockCookies });

    expect(cookieStore.get).toHaveBeenCalledWith(mockCookies, 'refreshToken');
    expect(result).toEqual({ cookieExists: false });
  });
});
