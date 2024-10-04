import { PUBLIC_API_URL } from '$env/static/public';

export const createEntity = async (formData: { category: string; address: any; name: any; notes: any; }): Promise<[string, number]> => {
    let message: string = ""
    let id: number = 0


    const response = await fetch(`${PUBLIC_API_URL}api/v1/entity`, {
        method: 'POST',
        body: JSON.stringify({
            address: formData.address,
            category: formData.category,
            name: formData.name,
            notes: formData.notes,
        })
    });

    const data = await response.json()

    message = data.message
    if (message == "success") {
        id = data.data
    }

    return [message, id];
}
