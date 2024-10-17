import { API_URL } from '$env/static/private';
import { cookieStore } from '$lib/stores/cookieStore';

//@ts-ignore
export async function POST({ request, cookies }) {
    const {
        userEmail,
        password,
    } = await request.json();

    let response = new Response()
    try {
        const proxyResponse = await fetch(`${API_URL}v1/token`, {
            method: 'POST',
            headers: new Headers({ 'content-type': 'application/json' }),
            body: JSON.stringify({
                userEmail: userEmail,
                password: password,
            })
        });

        const data = await proxyResponse.json()

        response = new Response(JSON.stringify(data), {
            status: proxyResponse.status,
            headers: {
                'Content-Type': 'application/json'
            }
        });

        if (proxyResponse.status == 200) {
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
        }

    } catch (error) {
        console.error(error);
        response = new Response(JSON.stringify(error),{
            status: 400,
            headers: {
                'Content-Type': 'application/json'
            }
        });
    }
    return response
}

//@ts-ignore
export async function DELETE({ cookies }) {
    let response = new Response()

    try {
        const response = await fetch(`${API_URL}v1/token`, {
            method: 'DELETE',
            headers: new Headers({ 'content-type': 'application/json' }),
            body: JSON.stringify({
                refreshToken: cookieStore.get(cookies, "refreshToken"),
            })
        });

        if (response.status == 200) {

            cookieStore.delete(cookies, 'accessToken', {
                path: '/',
            });

            cookieStore.delete(cookies, 'idToken', {
                path: '/',
            });

            cookieStore.delete(cookies, 'refreshToken', {
                path: '/',
            });
        }

    } catch (error) {
        console.error(error);
        response = new Response(JSON.stringify(error),{
            status: 400,
            headers: {
                'Content-Type': 'application/json'
            }
        });
    }
    return response
}
