import { goto } from '$app/navigation';
import { PUBLIC_API_URL } from '$env/static/public';

export const _signUpUser = async (formData: { confirmationCode: string; }): Promise<boolean> => {

    let success = false;

    try {
        const response = await fetch(`${PUBLIC_API_URL}api/v1/user`, {
            method: 'PUT',
            body: JSON.stringify({
                confirmationCode: formData.confirmationCode,
            })
        });

        success = response.ok
    } catch (error) {
        console.log(error)
    }

    if (success) {
        goto('/login');
    }

    return success
}
