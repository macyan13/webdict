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
            console.log(error)
            if (error.response.status === 400) {
                console.log(error.response)
                return new Promise((resolve, reject) => {
                    reject(error.response.data);
                });
            }

            if (error.response.status !== 401) {
                return new Promise((resolve, reject) => {
                    router.push({name: 'Error'})
                    reject(error);
                });
            }

            if (error.config.url === authService.getRefreshUrl() ) {
                store.dispatch('auth/logout').then(() => router.push({name: 'Login'}));
                return new Promise((resolve, reject) => {
                    reject(error);
                });
            }

            if (error.config.url === authService.getAuthUrl() ) {
                store.dispatch('auth/logout').then(() => {});
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
                                config.headers['Authorization'] = `${data.type} ${data.accessToken}`;

                                // return new Promise((resolve, reject) => {
                                    axios.request(config).then(response => {
                                        resolve(response);
                                    }).catch((error) => {
                                        store.dispatch('auth/logout').then(() => router.push({name: 'Login'}));
                                        reject(error);
                                    })
                                // });
                            }

                        );
                    })
                    .catch((error) => {
                        store.dispatch('auth/logout').then(() => router.push({name: 'Login'}));
                        console.log(error);
                        reject(error);
                    });
            })
        }
    );
}