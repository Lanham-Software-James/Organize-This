import { goto } from '$app/navigation';
import { PUBLIC_API_URL } from '$env/static/public';

export const _signUpUser = async (formData: { userEmail: string; password: string; firstName: string, lastName: string, birthday: string }): Promise<boolean> => {

    let success = false;

    try {
        const response = await fetch(`${PUBLIC_API_URL}api/v1/user`, {
            method: 'POST',
            body: JSON.stringify({
                userEmail: formData.userEmail,
                password: formData.password,
                firstName: formData.firstName,
                lastName: formData.lastName,
                birthday: formData.birthday,
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
