import { API_URL } from '$env/static/private';
import { cookieStore } from '$lib/stores/cookieStore.js';

//@ts-ignore
export async function POST({ request, cookies }) {
    const {
        userEmail,
        password,
        firstName,
        lastName,
        birthday,
    } = await request.json();

    let proxyResponse = new Response()
    try {
        proxyResponse = await fetch(`${API_URL}v1/user`, {
            method: 'POST',
            headers: new Headers({ 'content-type': 'application/json' }),
            body: JSON.stringify({
                userEmail: userEmail,
                password: password,
                firstName: firstName,
                lastName: lastName,
                birthday: birthday,
            })
        });

        if (proxyResponse.status == 200) {
            cookieStore.set(cookies, "userEmail", userEmail, {
                path: '/',
                maxAge: 60 * 60 * 24 * 30,
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
export async function PUT({ request, cookies }) {
    const {
        confirmationCode,
    } = await request.json();

    let proxyResponse = new Response()
    try {
        proxyResponse = await fetch(`${API_URL}v1/user`, {
            method: 'PUT',
            headers: new Headers({ 'content-type': 'application/json' }),
            body: JSON.stringify({
                userEmail: cookieStore.get(cookies, "userEmail"),
                confirmationCode: confirmationCode,
            })
        });

        if (proxyResponse.status == 200) {
            cookieStore.delete(cookies, "userEmail", {
                path: '/',
            })
            proxyResponse = new Response()
        }


    } catch (error) {
        console.error(error);
    }
    return proxyResponse
}
