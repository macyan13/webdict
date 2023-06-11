import AuthService from '../services/auth.service';

const authContext = AuthService.getAuthContext();
const initialState = authContext
    ? {loggedIn: true, authContext: authContext}
    : {loggedIn: false, user: null};

export default {
    namespaced: true,
    state: initialState,
    actions: {
        login({commit}, authParams) {
            return AuthService.login(authParams).then(
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
    },
    mutations: {
        loginSuccess(state, user) {
            state.loggedIn = true;
            state.authContext = user;
        },
        loginFailure(state) {
            state.loggedIn = false;
            state.authContext = null;
        },
        logout(state) {
            state.loggedIn = false;
            state.authContext = null;
        }
    },
    getters: {
        isLoggedIn: function (state) {
            return state.loggedIn;
        },
        authContext: function (state) {
            return state.authContext;
        },
    }
};