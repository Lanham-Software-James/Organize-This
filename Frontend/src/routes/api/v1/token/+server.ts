import { API_URL } from '$env/static/private';
import { cookieStore } from '$lib/stores/cookieStore';

//@ts-ignore
export async function POST({ request, cookies }) {
    const {
        userEmail,
        password,
    } = await request.json();

    let proxyResponse = new Response()
    try {
        proxyResponse = await fetch(`${API_URL}v1/token`, {
            method: 'POST',
            headers: new Headers({ 'content-type': 'application/json' }),
            body: JSON.stringify({
                userEmail: userEmail,
                password: password,
            })
        });

        if (proxyResponse.status == 200) {
            const data = await proxyResponse.json()

            cookieStore.set(cookies, 'accessToken', data.data.AccessToken, {
                path: '/',
                maxAge: data.data.ExpiresIn,
                httpOnly: true,
            });

            cookieStore.set(cookies, 'idToken', data.data.IdToken, {
                path: '/',
                maxAge: 60 * 60 * 24 * 30, // Expire when refresh token does to ensure we can refresh
                httpOnly: true,
            });

            cookieStore.set(cookies, 'refreshToken', data.data.RefreshToken, {
                path: '/',
                maxAge: 60 * 60 * 24 * 30, // Refresh Tokens expire in 30 days from cognito
                httpOnly: true,
            });

            proxyResponse = new Response()
        }

    } catch (error) {
        console.error(error);
    }
    return proxyResponse
}

//@ts-ignore
export async function DELETE({ cookies }) {
    let proxyResponse = new Response()
    try {
        proxyResponse = await fetch(`${API_URL}v1/token`, {
            method: 'DELETE',
            headers: new Headers({ 'content-type': 'application/json' }),
            body: JSON.stringify({
                refreshToken: cookieStore.get(cookies, "refreshToken"),
            })
        });

        if (proxyResponse.status == 200) {
            const data = await proxyResponse.json()

            cookieStore.delete(cookies, 'accessToken', {
                path: '/',
            });

            cookieStore.delete(cookies, 'idToken', {
                path: '/',
            });

            cookieStore.delete(cookies, 'refreshToken', {
                path: '/',
            });

            proxyResponse = new Response()
        }

    } catch (error) {
        console.error(error);
    }
    return proxyResponse
}
