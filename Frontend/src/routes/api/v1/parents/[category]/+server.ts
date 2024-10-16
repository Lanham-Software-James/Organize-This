import { cookieStore } from '$lib/stores/cookieStore';
import { API_URL } from '$env/static/private';

//@ts-ignore
export async function GET({ params, url, cookies }) {
    let proxyResponse = new Response()

    if( params.category == 'item' ||
        params.category == 'container' ||
        params.category == 'shelf' ||
        params.category == 'shelvingunit' ||
        params.category == 'room'
    ){
        try {
            proxyResponse = await fetch(
                `${API_URL}v1/parents/${params.category}`,
                {
                    headers: {
                        Authorization: "Bearer " + cookieStore.get(cookies, "accessToken")
                    }
                }
            );
        } catch (error) {
            console.error(error + "test");
        }
    }

    return proxyResponse
}
