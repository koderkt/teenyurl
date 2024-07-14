import { fail, redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { PRIVATE_BASE_URL } from '$env/static/private';
import type { ShortURLRespone } from '../app';

export const load: PageServerLoad = (async ({ cookies, fetch }) => {
    // fetch the current user's todos from the server
    const sessionId = cookies.get("sessionId");
    const data = {
        sessionId: sessionId
    }

    if (sessionId) {
        return data;
    }
    return {
        sessionId: null
    }
}) satisfies PageServerLoad;



export const actions = {
    default: async (event) => {
        const formData = await event.request.formData();
        const longURL = formData.get('longurl');
        console.log(longURL)
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
