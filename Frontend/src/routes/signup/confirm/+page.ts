import { goto } from '$app/navigation';
import { PUBLIC_API_URL } from '$env/static/public';

export const _confirmUser = async (formData: { confirmationCode: string; }): Promise<[boolean, string]> => {

    let success = false;
    let message = "Error";

    try {
        const response = await fetch(`${PUBLIC_API_URL}api/v1/user`, {
            method: 'PUT',
            body: JSON.stringify({
                confirmationCode: formData.confirmationCode,
            })
        });

        success = response.ok

        if (!success) {
            const data = await response.json()
            message = data.data
        }
    } catch (error) {
        console.log(error)
    }

    if (success) {
        goto('/login');
    }

    return [success, message]
}
