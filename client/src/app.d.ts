// See https://kit.svelte.dev/docs/types#app
// for information about these interfaces
declare global {
    namespace App {
        // interface Error {}
        // interface Locals {}
        // interface PageData {}
        // interface PageState {}
        // interface Platform {}
    }
}

export { SignUpForm, ShortURLRespone };


interface SignUpForm {
    username: string
    email: string
    password: string
}

interface ShortURLRespone {
    short_url: string
    long_url: string
    link_id: int
}
