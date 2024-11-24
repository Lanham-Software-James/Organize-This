import { goto, invalidateAll } from '$app/navigation';

export const _loginUser = async (formData: { userEmail: string; password: string; }): Promise<[boolean, string]> => {

    let success = false;
    let message = "Error"

    try {
        const response = await fetch(`/api/v1/token`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                userEmail: formData.userEmail,
                password: formData.password,
            })
        });

        success = response.ok

        if(!success) {
            const data = await response.json()
            message = data.data
        }
    } catch (error) {
        console.error(error)
    }

    if (success) {
        await invalidateAll();
        goto('/');
    }

    return [success, message]
}
