
import { PRIVATE_BASE_URL } from "$env/static/private";
import { redirect } from "@sveltejs/kit";

export const load = async (event) => {
    const sessionId = event.cookies.get('sessionId');

    if (sessionId) {
        redirect(301, '/');
    }
};


export const actions = {
    default: async (event) => {
        const formData = await event.request.formData();
        let email = formData.get("email")
        let password = formData.get("password")
        let userName = formData.get("username")
        const response = await fetch(`${PRIVATE_BASE_URL}/signup`, {
            method: 'POST',
            headers: {
                Accept: 'application/json',
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                user_name: userName,
                email: email,
                password: password
            })
        });

        const data = await response.json();
        if (response.status <= 299) {
            redirect(301, '/login');
        }
        return { error: data.message };


    }
}


