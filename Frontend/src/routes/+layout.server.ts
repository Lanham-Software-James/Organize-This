import { cookieStore } from '$lib/stores/cookieStore';
import type { LayoutServerLoad } from './$types';

//@ts-ignore
export const load = (async ({ cookies }) => {

    const tastyCookie = cookieStore.get(cookies, "refreshToken")
    let cookieExists = false
    if (tastyCookie != undefined) {
        cookieExists = true;
    }

    return {
        cookieExists: cookieExists
    }
}) satisfies LayoutServerLoad;
