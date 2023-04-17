import AuthService from '../services/auth.service';

const user = AuthService.getUser();
const initialState = user
    ? {loggedIn: true, user: user}
    : {loggedIn: false, user: null};

export default {
    namespaced: true,
    state: initialState,
    actions: {
        login({commit}, user) {
            return AuthService.login(user).then(
                user => {
                    commit('loginSuccess', user);
                    return Promise.resolve(user);
                },
                error => {
                    console.log("Error: " + error);
                    commit('loginFailure');
                    return Promise.reject(error);
                }
            );
        },
        refresh({commit}, user) {
            commit('loginSuccess', user);
        },
        logout({commit}) {
            AuthService.logout();
            commit('logout');
        },
        // register({commit}, user) {
        //     return AuthService.register(user).then(
        //         response => {
        //             commit('registerSuccess');
        //             return Promise.resolve(response.data);
        //         },
        //         error => {
        //             commit('registerFailure');
        //             return Promise.reject(error);
        //         }
        //     );
        // }
    },
    mutations: {
        loginSuccess(state, user) {
            state.loggedIn = true;
            state.user = user;
        },
        loginFailure(state) {
            state.loggedIn = false;
            state.user = null;
        },
        logout(state) {
            state.loggedIn = false;
            state.user = null;
        }
    },
    getters: {
        isLoggedIn: function (state) {
            return state.loggedIn;
        },
        user: state => state.user,
    }
};