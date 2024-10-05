import { goto, invalidateAll } from '$app/navigation';
import { PUBLIC_API_URL } from '$env/static/public';

export const _logoutUser = async (): Promise<boolean> => {

    let success = false;

    try {
        const response = await fetch(`${PUBLIC_API_URL}api/v1/token`, {
            method: 'DELETE',
        });

        success = response.ok
    } catch (error) {
        console.log(error)
    }

    if (success) {
        await invalidateAll();
        goto('/login');
    }

    return success
}
