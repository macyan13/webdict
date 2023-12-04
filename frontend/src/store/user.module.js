import UserService from '../services/user.service';

export default {
    namespaced: true,
    state: {
        users: [],
        entityStatus: null,
    },
    actions: {
        fetchAll({state, commit}) {
            if (state.users.length > 0) {
                return Promise.resolve(state.users);
            }
            return UserService.getAll().then(
                users => {
                    commit('allUsersSuccess', users);
                    return Promise.resolve(users);
                },
                error => {
                    console.log("Error: " + error);
                    return Promise.reject(error);
                }
            );
        },
        clear({commit}) {
            commit('cleanUsers');
            return Promise.resolve();
        },
        setEntityStatus({commit}, status) {
            commit('setEntityStatus', status);
            return Promise.resolve();
        },
        clearEntityStatus({commit}) {
            commit('setEntityStatus', null);
            return Promise.resolve();
        }
    },
    getters: {
        entityStatus: function (state) {
            return state.entityStatus;
        },
    },
    mutations: {
        allUsersSuccess(state, users) {
            state.users = users;
        },
        cleanUsers(state) {
            state.users = [];
        },
        setEntityStatus(state, status) {
            state.entityStatus = status;
        },
    },
};