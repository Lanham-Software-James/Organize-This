import { cookieStore } from '$lib/stores/cookieStore';
import { API_URL } from '$env/static/private';

//@ts-ignore
export async function GET({ url, cookies }) {
    const offset = url.searchParams.get('offset');
    const limit = url.searchParams.get('limit');
    const search = url.searchParams.get('search');
    const filter = url.searchParams.get('filter');

    let response = new Response()

    try {
        response = await fetch(
            `${API_URL}v1/entities?offset=${offset}&limit=${limit}&search=${search}&filter=${filter}`,
            {
                headers: {
                    Authorization: "Bearer " + cookieStore.get(cookies, "accessToken")
                }
            }
        );
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
