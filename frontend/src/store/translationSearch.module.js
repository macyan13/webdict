export default {
    namespaced: true,
    state: {
        pageSize: 20,
        currentPage: 1,
        totalRecords: 0,
        translations: [],
        tags: [],
    },
    actions: {
        setPageSize({ commit }, newSize) {
            commit("setPageSize", newSize);
        },
        setCurrentPage({ commit }, newPage) {
            commit("setCurrentPage", newPage);
        },
        setTotalRecords({ commit }, newTotal) {
            commit("setTotalRecords", newTotal);
        },
        setTranslations({ commit }, newTranslations) {
            commit("setTranslations", newTranslations);
        },
        setTags({ commit }, newTags) {
            commit("setTags", newTags);
        },
        resetTranslations({ commit }) {
            commit("setTranslations", []);
            commit("setTotalRecords", 0);
            commit("setCurrentPage", 1);
        },
        resetTags({ commit }) {
            commit("setTags", []);
            commit("setTranslations", []);
            commit("setTotalRecords", 0);
            commit("setCurrentPage", 1);
        }
    },
    getters: {
        getPageSize: function (state) {
            return state.pageSize;
        },
        getCurrentPage: function (state) {
            return state.currentPage;
        },
        getTotalRecords: function (state) {
            return state.totalRecords;
        },
        getTranslations: function (state) {
            return state.translations;
        },
        getTags: function (state) {
            return state.tags;
        }
    },
    mutations: {
        setPageSize(state, newSize) {
            state.pageSize = newSize;
        },
        setCurrentPage(state, newPage) {
            state.currentPage = newPage;
        },
        setTotalRecords(state, newTotal) {
            state.totalRecords = newTotal;
        },
        setTranslations(state, newTranslations) {
            state.translations = newTranslations;
        },
        setTags(state, newTags) {
            state.tags = newTags;
        }
    },
};
