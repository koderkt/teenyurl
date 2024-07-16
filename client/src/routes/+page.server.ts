import type { PageServerLoad } from './$types';
import { PRIVATE_BASE_URL } from '$env/static/private';
import { redirect } from '@sveltejs/kit';

export const load: PageServerLoad = (async ({ cookies, fetch }) => {
    // fetch the current user's todos from the server
    const sessionId = cookies.get("sessionId");
    const data = {
        sessionId: sessionId
    }
    if (!sessionId) {
        redirect(302, "/login")
    }
    if (sessionId) {
        return data;
    }

}) satisfies PageServerLoad;



export const actions = {
    default: async (event) => {
        const formData = await event.request.formData();
        const longURL = formData.get('longurl');
        const cookie = event.cookies.get("sessionId");

        const response = await fetch(`${PRIVATE_BASE_URL}/links`, {
            method: 'POST',
            headers: {
                Accept: 'application/json',
                'content-type': 'application/json',
                "Authorization": `Bearer ${cookie}`
            },
            body: JSON.stringify({
                long_url: longURL,
            })
        });

        const data = await response.json();
        console.log(data)

        if (response.ok) {
            return {
                longURL: data.long_url,
                shortURL: data.short_url
            }
        }

        return {
            error: data.message
        }
    }
}
