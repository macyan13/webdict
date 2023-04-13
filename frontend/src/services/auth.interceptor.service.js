import axios from 'axios';
import store from '../store';
import router from '../router';
import authService from './auth.service.js';

export default () => {
    axios.interceptors.response.use(
        (response) => {
            return response;
        },
        (error) => {
            if (error.response.status !== 401) {
                // todo: add handling for error page
                return new Promise((resolve, reject) => {
                    router.push({name: 'Error'})
                    reject(error);
                });
            }

            if (error.config.url === authService.getRefreshUrl()) {
                store.dispatch('auth/logout').then(() => router.push({name: 'Login'}));
                return new Promise((resolve, reject) => {
                    reject(error);
                });
            }

            return new Promise((resolve, reject) => {
                authService.refresh()
                    .then((data) => {
                        store.dispatch('auth/refresh', data).then(
                            function () {
                                // New request with new token
                                const config = error.config;
                                config.headers['Authorization'] = `Bearer ${data.accessToken}`;

                                return new Promise((resolve, reject) => {
                                    axios.request(config).then(response => {
                                        resolve(response);
                                    }).catch((error) => {
                                        reject(error);
                                    })
                                });
                            }

                        );
                    })
                    .catch((error) => {
                        console.log(error);
                        reject(error);
                    });
            })
        }
    );
}