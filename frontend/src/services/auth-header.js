const AUTH_CONTEXT_LOCAL_STORAGE_KEY = 'auth_context';

export default function authHeader() {
    let authContext = JSON.parse(localStorage.getItem(AUTH_CONTEXT_LOCAL_STORAGE_KEY));

    if (authContext && authContext.accessToken && authContext.type) {
        return { Authorization: authContext.type + ' ' + authContext.accessToken };
    } else {
        return {};
    }
}