import { goto, invalidateAll } from '$app/navigation';

export const _logoutUser = async (): Promise<boolean> => {

    let success = false;

    try {
        const response = await fetch(`/api/v1/token`, {
            method: 'DELETE',
        });

        success = response.ok
    } catch (error) {
        console.error(error)
    }

    if (success) {
        await invalidateAll();
        goto('/login');
    }

    return success
}
