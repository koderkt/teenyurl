import type { RequestEvent } from "../routes/$types";

export function handleLoginRedirect(event: RequestEvent) {
    const redirectTo = event.url.pathname + event.url.search

    return `/login?redirectTo=${redirectTo}`
}
