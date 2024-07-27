import { fail, redirect, type Actions } from "@sveltejs/kit";
import type { PageServerLoad, RequestEvent } from "../$types";
import { PRIVATE_BASE_URL } from "$env/static/private";
import type { Link } from "../../app";

export const load: PageServerLoad = async (event: RequestEvent) => {
    const cookie = event.cookies.get("sessionId");
    if (!cookie) {
        throw redirect(302, "/login");
    }
    const response = await fetch(`${PRIVATE_BASE_URL}/links`, {
        method: 'GET',
        headers: {
            Accept: 'application/json',
            "Authorization": `Bearer ${cookie}`

            // 'Content-Type': 'application/json'
        }
    });
    if (response.status === 401) {
        event.cookies.set('sessionId', '', {
            path: '/',
            expires: new Date(0),
        });
        throw redirect(302, '/login');
    }
    if (response.ok) {
        const res: Link[] = await response.json();
        res.reverse();
        return {
            links: res,
            cookie: cookie
        };
    }
    return {};
};

export const actions: Actions = {
    updateLink: async (event) => {
        const formData = await event.request.formData();
        const shortUrl = formData.get('short_url') as string;
        const originalUrl = formData.get('original_url') as string;

        if (!shortUrl || !originalUrl) {
            return fail(400, { error: 'Missing required fields' });
        }

        try {
            const splits = shortUrl.split('/');
            const response = await fetch(`${PRIVATE_BASE_URL}/${splits[splits.length - 1]}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    Authorization: `Bearer ${event.cookies.get("sessionId")}`, // Adjust according to your auth logic
                },
                body: JSON.stringify({ long_url: originalUrl }),
            });

            if (response.ok) {
                return {
                    success: true,
                    message: 'Link updated successfully',
                };
            } else {
                return fail(response.status, { error: 'Failed to update link' });
            }
        } catch (error: any) {
            return fail(500, { error: 'Error updating link: ' + error.message });
        }
    },
    enableDisableLink: async (event) => {
        const formData = await event.request.formData();
        const shortUrl = formData.get('short_url') as string;
        const val = formData.get("isEnabled") as string;
        // if (!linkId || !shortUrl || !originalUrl) {
        //     return fail(400, { error: 'Missing required fields' });
        // }

        try {
            const splits = shortUrl.split('/');
            console.log(val)
            console.log(splits[splits.length - 1])
            const response = await fetch(`${PRIVATE_BASE_URL}/${splits[splits.length - 1]}/${val}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    Authorization: `Bearer ${event.cookies.get("sessionId")}`, // Adjust according to your auth logic
                },
            });

            if (response.ok) {
                return {
                    success: true,
                    message: 'Link updated successfully',
                };
            } else {
                return fail(response.status, { error: 'Failed to update link' });
            }
        } catch (error: any) {
            return fail(500, { error: 'Error updating link: ' + error.message });
        }
    },
};

