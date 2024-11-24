import { goto } from '$app/navigation';

export const _signUpUser = async (formData: { userEmail: string; password: string; firstName: string, lastName: string, birthday: string }): Promise<[boolean, string]> => {

    let success = false;
    let message = "Error";

    try {
        const response = await fetch(`/api/v1/user`, {
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

        if(!success) {
            const data = await response.json()
            message = data.data
        }
    } catch (error) {
        console.error(error)
    }

    if (success) {
        goto('/signup/confirm');
    }

    return [success, message]
}
