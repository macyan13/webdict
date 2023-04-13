import axios from 'axios';

const USER_LOCAL_STORAGE_KEY = 'user';

class AuthService {
    login(user) {
        return axios
            .post('/api/auth/signin', {
                    username: user.username,
                    password: user.password
                },
                {withCredentials: true}
            )
            .then(response => {
                console.log(response);
                if (response.data.accessToken) {
                    localStorage.setItem(USER_LOCAL_STORAGE_KEY, JSON.stringify(response.data));
                }
                return response.data;
            });
    }

    refresh() {
        return new Promise((resolve, reject) => {
            axios
                .post('/api/auth/refresh', {}, {withCredentials: true})
                .then(response => {
                    console.log(response);
                    if (response.data && response.data.accessToken) {
                        localStorage.setItem(USER_LOCAL_STORAGE_KEY, JSON.stringify(response.data));
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
        localStorage.removeItem(USER_LOCAL_STORAGE_KEY);
    }

    register(user) {
        return axios.post('/api/auth/signup', {
            username: user.username,
            email: user.email,
            password: user.password
        });
    }

    getUser() {
        return JSON.parse(localStorage.getItem(USER_LOCAL_STORAGE_KEY));
    }

    getRefreshUrl() {
        return '/api/auth/refresh'
    }
}

export default new AuthService();