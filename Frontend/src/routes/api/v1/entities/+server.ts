import { cookieStore } from '$lib/stores/cookieStore';
import { API_URL } from '$env/static/private';

//@ts-ignore
export async function GET({ url, cookies }) {
    const offset = url.searchParams.get('offset');
    const limit = url.searchParams.get('limit');

    let proxyResponse = new Response()

    try {
        proxyResponse = await fetch(
            `${API_URL}v1/entities?offset=${offset}&limit=${limit}`,
            {
                headers: {
                    Authorization: "Bearer " + cookieStore.get(cookies, "accessToken")
                }
            }
        );
    } catch (error) {
        console.error(error + "test");
    }
    return proxyResponse

}
