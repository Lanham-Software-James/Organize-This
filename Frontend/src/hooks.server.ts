import { sequence } from '@sveltejs/kit/hooks';
import { API_URL } from '$env/static/private';
import { cookieStore } from '$lib/stores/cookieStore';
import { redirect } from '@sveltejs/kit';

const public_paths = [
    '/login',
    '/signup',
];

export function isPathAllowed(path: string) {
    return public_paths.some(allowedPath =>
        path === allowedPath || path.startsWith(allowedPath + '/')
    );
}

//@ts-ignore
export async function validatePath({ event, resolve }) {
    const url = new URL(event.request.url);
    const is_path_valid = isPathAllowed(url.pathname)
    const is_not_api = !url.pathname.includes("api")
    const is_user_authed = cookieStore.get(event.cookies, "refreshToken") != undefined;


    if (!is_user_authed && !is_path_valid && is_not_api) {
        redirect(302, '/login')
    }
    else if (is_user_authed && is_path_valid && is_not_api) {
        redirect(302, '/')
    }
    return resolve(event);
}

//@ts-ignore
export async function refreshTokenHandle({ event, resolve }) {
    const accessToken = cookieStore.get(event.cookies, "accessToken");
    const idToken = cookieStore.get(event.cookies, "idToken");
    const refreshToken = cookieStore.get(event.cookies, "refreshToken");
    const url = new URL(event.request.url);
    const is_api = url.pathname.includes("api")

    if (!accessToken && refreshToken && idToken) {
        try {
            const refreshResponse = await fetch(`${API_URL}v1/token`, {
                method: 'PUT',
                body: JSON.stringify({
                    refreshToken: refreshToken,
                    idToken: idToken,
                })
            });

            if (refreshResponse.ok) {
                const data = await refreshResponse.json();
                const newAccessToken = data.data.AccessToken;
                const newIdToken = data.data.IdToken;

                // Resolve the event with the new auth header
                const response = await resolve(event, {
                    //@ts-ignore
                    transformPageChunk: ({ html }) => html,
                    //@ts-ignore
                    filterSerializedResponseHeaders: (name) => true,
                });

                // Create a new response with updated headers and cookies
                const newResponse = new Response(response.body, response);
                newResponse.headers.set('Authorization', `Bearer ${newAccessToken}`);
                newResponse.headers.append('Set-Cookie', `accessToken=${newAccessToken}; Path=/; HttpOnly; Max-Age=${data.data.ExpiresIn}`);
                newResponse.headers.append('Set-Cookie', `idToken=${newIdToken}; Path=/; HttpOnly; Max-Age=${60 * 60 * 24 * 30}`);

                return newResponse;
            }
        } catch (error) {
            console.error('Error refreshing token:', error);
        }
    }

    return resolve(event);
}

export const handle = sequence(validatePath, refreshTokenHandle);
