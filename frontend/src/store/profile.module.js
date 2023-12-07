import ProfileService from '../services/profile.service';

export default {
    namespaced: true,
    state: {
        profile: null,
        entityStatus: null,
    },
    actions: {
        fetchProfile({state, commit}) {
            if (state.profile) {
                return Promise.resolve(state.profile);
            }
            return ProfileService.get().then(
                profile => {
                    commit('profileSuccess', profile);
                    return Promise.resolve(profile);
                },
                error => {
                    console.log("Error: " + error);
                    return Promise.reject(error);
                }
            );
        },
        clear({commit}) {
            commit('clearProfile');
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
    mutations: {
        profileSuccess(state, profile) {
            state.profile = profile;
        },
        clearProfile(state) {
            state.profile = null;
        },
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
    getters: {
        isAdmin: function (state) {
            return state.profile && state.profile.role.is_admin;
        },
        entityStatus: function (state) {
            return state.entityStatus;
        },
    }
};