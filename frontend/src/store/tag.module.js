import TagService from '../services/tag.service';

export default {
    namespaced: true,
    state: {
        tags: [],
    },
    actions: {
        fetchAll({state, commit}) {
            if (state.tags.length > 0) {
                return Promise.resolve(state.tags);
            }
            return TagService.getAll().then(
                tags => {
                    commit('allTagsSuccess', tags);
                    return Promise.resolve(tags);
                },
                error => {
                    console.log("Error: " + error);
                    return Promise.reject(error);
                }
            );
        },
        clear({commit}) {
            commit('cleanTags');
            return Promise.resolve();
        }
    },
    mutations: {
        allTagsSuccess(state, tags) {
            state.tags = tags;
        },
        cleanTags(state) {
            state.tags = [];
        },
    },
};