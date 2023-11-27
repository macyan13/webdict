import axios from 'axios';
import authHeader from "@/services/auth-header";

const LAST_USED_TRANSLATION_TAGS_LOCAL_STORAGE_KEY = 'last_used_translation_tags';

class TagService {
    create(tag) {
        return new Promise((resolve, reject) => {
            axios
                .post('/v1/api/tags', tag, {headers: authHeader()})
                .then(response => {
                    resolve(response.data);
                })
                .catch(error => {
                    console.log(error)
                    reject(error);
                });
        });
    }
    get(id) {
        return new Promise((resolve, reject) => {
            axios
                .get('/v1/api/tags/' + id, {headers: authHeader()})
                .then(response => {
                    resolve(response.data);
                })
                .catch(error => {
                    console.log(error)
                    reject(error);
                });
        });
    }
    update(tag) {
        return new Promise((resolve, reject) => {
            axios
                .put('/v1/api/tags/' + tag.id, {
                    name: tag.name,
                }, {headers: authHeader()})
                .then(response => {
                    resolve(response.data);
                })
                .catch(error => {
                    console.log(error)
                    reject(error);
                });
        });
    }
    delete(id) {
        return new Promise((resolve, reject) => {
            axios
                .delete('/v1/api/tags/' + id, {headers: authHeader()})
                .then(response => {
                    resolve(response.data);
                })
                .catch(error => {
                    console.log(error)
                    reject(error);
                });
        });
    }
    getAll() {
        return new Promise((resolve, reject) => {
            axios
                .get('/v1/api/tags', {headers: authHeader()})
                .then(response => {
                    resolve(response.data);
                })
                .catch(error => {
                    console.log(error)
                    reject(error);
                });
        });
    }
    getLastUsedTranslationTagIds() {
        return JSON.parse(localStorage.getItem(LAST_USED_TRANSLATION_TAGS_LOCAL_STORAGE_KEY));
    }
    updateLastUsedTranslationTagIds(tagIds) {
        localStorage.setItem(LAST_USED_TRANSLATION_TAGS_LOCAL_STORAGE_KEY, JSON.stringify(tagIds));
    }
}

export default new TagService();
