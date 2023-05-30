import axios from 'axios';

const AUTH_CONTEXT_LOCAL_STORAGE_KEY = 'auth_context';

class AuthService {
    login(user) {
        return axios
            .post('/v1/api/auth/signin', {
                    email: user.email,
                    password: user.password
                },
                {withCredentials: true}
            )
            .then(response => {
                if (response.data.accessToken) {
                    localStorage.setItem(AUTH_CONTEXT_LOCAL_STORAGE_KEY, JSON.stringify(response.data));
                }
                return response.data;
            });
    }

    refresh() {
        return new Promise((resolve, reject) => {
            axios
                .post(this.getRefreshUrl(), {}, {withCredentials: true})
                .then(response => {
                    if (response.data && response.data.accessToken) {
                        localStorage.setItem(AUTH_CONTEXT_LOCAL_STORAGE_KEY, JSON.stringify(response.data));
                    }
                    resolve(response.data);
                })
                .catch(error => {
                    console.log(error)
                    reject(error);
                });
        });
    }

    logout() {
        localStorage.removeItem(AUTH_CONTEXT_LOCAL_STORAGE_KEY);
    }

    getAuthContext() {
        return JSON.parse(localStorage.getItem(AUTH_CONTEXT_LOCAL_STORAGE_KEY));
    }

    getRefreshUrl() {
        return '/v1/api/auth/refresh';
    }
}

export default new AuthService();