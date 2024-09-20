import { PUBLIC_API_URL } from '$env/static/public';

export const createEntity = async (formData: { Category: string; Address: any; Name: any; Notes: any; }) => {
    const response = await fetch(PUBLIC_API_URL + 'v1/entity-management/' + formData.Category, {
        method: 'POST',
        body: JSON.stringify({
            Address: formData.Address,
            Name: formData.Name,
            Notes: formData.Notes
        })
    });

    return response;
}
