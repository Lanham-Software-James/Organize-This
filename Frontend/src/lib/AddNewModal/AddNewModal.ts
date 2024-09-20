import { PUBLIC_API_URL } from '$env/static/public';

export const createEntity = async (formData: { category: string; address: any; name: any; notes: any; }) => {
    const response = await fetch(PUBLIC_API_URL + 'v1/entity', {
        method: 'POST',
        body: JSON.stringify({
            address: formData.address,
            category: formData.category,
            name: formData.name,
            notes: formData.notes,
        })
    });

    return response;
}
