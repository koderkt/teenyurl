import { redirect } from '@sveltejs/kit'
import type { Actions, PageServerLoad } from './$types'
import { PRIVATE_BASE_URL } from '$env/static/private';


export const actions: Actions = {
    default: async ({ cookies }) => {
        
        const cookie = cookies.get("sessionId")
        const response = await fetch(`${PRIVATE_BASE_URL}/signout`, {
            method: 'POST',
            headers: {
                Accept: 'application/json',
                'content-type': 'application/json',
                "Authorization": `Bearer ${cookie}`
            },
        });
        let data = await response.json()
        console.log(data)
        if (response.status <= 299) {
            cookies.set('sessionId', '', {
                path: '/',
                expires: new Date(0),
            });
            redirect(302, '/login');
        }
        return {
            error: data.message
        }
    },
}
