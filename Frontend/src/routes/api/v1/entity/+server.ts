import { cookieStore } from '$lib/stores/cookieStore';
import { API_URL } from '$env/static/private';

//@ts-ignore
export async function POST({ request, cookies }) {
    const {
        address,
        category,
        name,
        notes,
    } = await request.json();

    let response = new Response()
    try {
        response = await fetch(
            `${API_URL}v1/entity`,
            {
                method: "POST",
                headers: {
                    Authorization: "Bearer " + cookieStore.get(cookies, "accessToken")
                },
                body: JSON.stringify({
                    address: address,
                    category: category,
                    name: name,
                    notes: notes,
                })
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
