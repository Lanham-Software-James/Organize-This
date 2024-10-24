import { cookieStore } from '$lib/stores/cookieStore';
import { API_URL } from '$env/static/private';

//@ts-ignore
export async function GET({ params, cookies }) {
    let response = new Response()

    if( params.category == 'item' ||
        params.category == 'container' ||
        params.category == 'shelf' ||
        params.category == 'shelving_unit' ||
        params.category == 'room'
    ){
        try {
            response = await fetch(
                `${API_URL}v1/parents/${params.category}`,
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
    }

    return response
}
