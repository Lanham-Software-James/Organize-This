import { goto } from '$app/navigation';

export const _confirmUser = async (formData: { confirmationCode: string; }): Promise<[boolean, string]> => {

    let success = false;
    let message = "Error";

    try {
        const response = await fetch(`/api/v1/user`, {
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
        console.error(error)
    }

    if (success) {
        goto('/login');
    }

    return [success, message]
}
