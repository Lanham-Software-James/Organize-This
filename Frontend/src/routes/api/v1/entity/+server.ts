import { cookieStore } from '$lib/stores/cookieStore';
import { API_URL } from '$env/static/private';

//@ts-ignore
export async function POST({ request, cookies }) {
    const {
        address,
        category,
        name,
        notes,
        parentID,
        parentCategory,
    } = await request.json();

    let proxyResponse = new Response()
    try {
        proxyResponse = await fetch(
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
                    parentID: parentID,
                    parentCategory: parentCategory,
                })
            }
        );
    } catch (error) {
        console.error(error);
    }
    return proxyResponse
}
