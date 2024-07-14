import type { LayoutServerLoad } from './$types';

export const load = (async ({ cookies }) => {
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

}) satisfies LayoutServerLoad;
