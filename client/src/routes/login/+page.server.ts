import { applyAction } from '$app/forms';
import { PRIVATE_BASE_URL } from '$env/static/private';
import { redirect } from '@sveltejs/kit';

export const load = async (event) => {
    const sessionId = event.cookies.get("sessionId");

    if (sessionId) {
        redirect(301, '/');
    }
}


export const actions = {
    default: async (event) => {
        const formData = await event.request.formData();
        let email = formData.get("email")
        let password = formData.get("password")

        console.log(formData);
        const response = await fetch(`${PRIVATE_BASE_URL}/signin`, {
            method: 'POST',
            headers: {
                Accept: 'application/json',
                'content-type': 'application/json'
            },
            body: JSON.stringify({
                email: email,
                password: password
            })
        });
        const data = await response.json();
        if (response.ok) {
            const sessionId = response.headers.get("Authorization");
            console.log(sessionId);
            event.cookies.set("sessionId", sessionId?.split("Bearer ")[1] ?? "", {
                path: "/",
            });

            redirect(301, "/")
        }
        return { error: data.message }
    }
}

