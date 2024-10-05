import { goto, invalidateAll } from '$app/navigation';
import { PUBLIC_API_URL } from '$env/static/public';

export const _loginUser = async (formData: { userEmail: string; password: string; }): Promise<boolean> => {

    let success = false;

    try {
        const response = await fetch(`${PUBLIC_API_URL}api/v1/token`, {
            method: 'POST',
            body: JSON.stringify({
                userEmail: formData.userEmail,
                password: formData.password,
            })
        });

        success = response.ok
    } catch (error) {
        console.log(error)
    }

    if (success) {
        await invalidateAll();
        goto('/');
    }

    return success
}
