export default {
    namespaced: true,
    state: {
        entityStatus: null,
    },
    actions: {
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
        setEntityStatus(state, status) {
            state.entityStatus = status;
        },
    },
};
