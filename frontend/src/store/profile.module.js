import ProfileService from '../services/profile.service';

export default {
    namespaced: true,
    state: {
        user: null,
    },
    actions: {
        fetchProfile({state, commit}) {
            if (state.user) {
                return Promise.resolve(state.user);
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
            commit('profileUser');
            return Promise.resolve();
        }
    },
    mutations: {
        profileSuccess(state, user) {
            state.user = user;
        },
        profileUser(state) {
            state.user = null;
        },
    },
};