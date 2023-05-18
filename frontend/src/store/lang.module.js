import LangService from '../services/lang.service';

export default {
    namespaced: true,
    state: {
        langs: [],
    },
    actions: {
        fetchAll({state, commit}) {
            if (state.langs.length > 0) {
                return Promise.resolve(state.langs);
            }
            return LangService.getAll().then(
                langs => {
                    commit('allLangSuccess', langs);
                    return Promise.resolve(langs);
                },
                error => {
                    console.log("Error: " + error);
                    return Promise.reject(error);
                }
            );
        },
        clear({commit}) {
            commit('cleanLangs');
            return Promise.resolve();
        }
    },
    mutations: {
        allLangSuccess(state, langs) {
            state.langs = langs;
        },
        cleanLangs(state) {
            state.langs = [];
        },
    },
};