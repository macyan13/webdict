import ProfileService from '../services/profile.service';

export default {
    namespaced: true,
    state: {
        profile: null,
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
        }
    },
    mutations: {
        profileSuccess(state, profile) {
            state.profile = profile;
        },
        clearProfile(state) {
            state.profile = null;
        },
    },
    getters: {
        isAdmin: function (state) {
            return state.profile && state.profile.role.is_admin;
        },
    }
};