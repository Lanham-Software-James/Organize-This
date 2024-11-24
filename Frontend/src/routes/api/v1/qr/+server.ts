import { cookieStore } from '$lib/stores/cookieStore';
import { API_URL } from '$env/static/private';

//@ts-ignore
export async function POST({ request, cookies }) {
    const {
        id,
        category,
    } = await request.json();

    let response = new Response()
    try {
        var proxyResponse = await fetch(
            `${API_URL}v1/qr`,
            {
                method: "POST",
                headers: {
                    'Authorization': 'Bearer ' + cookieStore.get(cookies, 'accessToken'),
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    id: id + "",
                    category: category,
                })
            }
        );

        const data = await proxyResponse.json()

        response = new Response(JSON.stringify(data), {
            status: proxyResponse.status,
            headers: {
                'Content-Type': 'application/json',
            }
        });
    } catch (error) {
        console.error(error);
        response = new Response(JSON.stringify(error), {
            status: 400,
            headers: {
                'Content-Type': 'application/json'
            }
        });
    }
    return response
}
