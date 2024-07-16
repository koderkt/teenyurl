import { redirect } from "@sveltejs/kit";
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
        console.log(response.status)
        const res: Link[] = await response.json();
        console.log(res)
        return {
            links: res
        };
    }
    return {};
};
