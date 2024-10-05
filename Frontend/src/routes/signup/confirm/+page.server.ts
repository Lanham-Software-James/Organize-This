import { cookieStore } from '$lib/stores/cookieStore';
import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load = (async ({ cookies }) => {
    const tastyCookie = cookieStore.get(cookies, "refreshToken")

    if (tastyCookie != undefined) {
        redirect(302, '/')
    }

    return {}
}) satisfies PageServerLoad;
