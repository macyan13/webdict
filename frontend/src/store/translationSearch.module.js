export default {
    namespaced: true,
    state: {
        pageSize: 20,
        currentPage: 1,
        totalRecords: 0,
        translations: [],
        target: '',
        source: '',
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
        setTarget({ commit }, newTarget) {
            commit("setTarget", newTarget);
        },
        setSource({ commit }, newSource) {
            commit("setSource", newSource);
        },
        resetTranslations({ commit }) {
            commit("setTranslations", []);
            commit("setTotalRecords", 0);
            commit("setCurrentPage", 1);
            commit("setTarget", '');
            commit("setSource", '');
        },
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
        getTarget: function (state) {
            return state.target;
        },
        getSource: function (state) {
            return state.source;
        },
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
        setTarget(state, newTarget) {
            state.target = newTarget;
        },
        setSource(state, newSource) {
            state.source = newSource;
        },
    },
};
