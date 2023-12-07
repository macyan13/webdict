import TagService from '../services/tag.service';

const lastUsedTranslationTagIds = TagService.getLastUsedTranslationTagIds();

export default {
    namespaced: true,
    state: {
        tags: [],
        lastUsedTranslationTagIds: lastUsedTranslationTagIds ? lastUsedTranslationTagIds : [],
        entityStatus: null,
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
        },
        updateLastUsedTranslationTagIds({commit}, tagIds) {
            commit('updateLastUsedTranslationTagIds', tagIds);
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
        lastUsedTranslationTagIds: function (state) {
            return state.lastUsedTranslationTagIds;
        },
        entityStatus: function (state) {
            return state.entityStatus;
        },
    },
    mutations: {
        allTagsSuccess(state, tags) {
            state.tags = tags;
        },
        cleanTags(state) {
            state.tags = [];
        },
        updateLastUsedTranslationTagIds(state, tagIds) {
            state.lastUsedTranslationTagIds = tagIds;
            TagService.updateLastUsedTranslationTagIds(tagIds);
        },
        setEntityStatus(state, status) {
            state.entityStatus = status;
        },
    },
};