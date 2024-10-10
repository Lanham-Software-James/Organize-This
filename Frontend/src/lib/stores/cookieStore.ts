import { writable } from 'svelte/store';
import type { Cookies } from '@sveltejs/kit';

function createCookieStore() {
  const { subscribe, set, update } = writable({});

  return {
    subscribe,
    set: (cookies: Cookies, name: string, value: string, options?: any) => {
      cookies.set(name, value, options);
      update(store => ({ ...store, [name]: value }));
    },
    get: (cookies: Cookies, name: string) => {
      const value = cookies.get(name);
      update(store => ({ ...store, [name]: value }));
      return value;
    },
    delete: (cookies: Cookies, name: string, options?: any) => {
      cookies.delete(name, options);
      update(store => {
        // @ts-ignore
        const { [name]: _, ...rest } = store;
        return rest;
      });
    }
  };
}

export const cookieStore = createCookieStore();
