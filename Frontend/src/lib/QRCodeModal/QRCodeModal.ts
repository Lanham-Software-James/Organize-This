export const generateQR = async (category: string, id: number): Promise<[string, string]> => {
    let message: string = ""
    let url: string = ""

    const response = await fetch(`/api/v1/qr`, {
        method: 'POST',
        headers:{
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            id: id,
            category: category,
        })
    });

    const data = await response.json()

    message = data.message
    if (message == "success") {
        url = data.data
    }

    return [message, url];
}
