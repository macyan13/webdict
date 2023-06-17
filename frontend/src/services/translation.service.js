import axios from 'axios';
import authHeader from "@/services/auth-header";

class TranslationService {
    create(translation) {
        return new Promise((resolve, reject) => {
            axios
                .post('/v1/api/translations', translation, {headers: authHeader()})
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
                .get('/v1/api/translations/' + id, {headers: authHeader()})
                .then(response => {
                    resolve(response.data);
                })
                .catch(error => {
                    console.log(error)
                    reject(error);
                });
        });
    }
    update(translation) {
        return new Promise((resolve, reject) => {
            axios
                .put('/v1/api/translations/' + translation.id, translation, {headers: authHeader()})
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
                .delete('/v1/api/translations/' + id, {headers: authHeader()})
                .then(response => {
                    resolve(response.data);
                })
                .catch(error => {
                    console.log(error)
                    reject(error);
                });
        });
    }
    search(SearchParams) {
        return new Promise((resolve, reject) => {
            axios
                .get('/v1/api/translations/last', {headers: authHeader(), params: SearchParams})
                .then(response => {
                    resolve(response.data);
                })
                .catch(error => {
                    console.log(error)
                    reject(error);
                });
        });
    }
    random(randomParams) {
        return new Promise((resolve, reject) => {
            axios
                .get('/v1/api/translations/random', {headers: authHeader(), params: randomParams})
                .then(response => {
                    resolve(response.data);
                })
                .catch(error => {
                    console.log(error)
                    reject(error);
                });
        });
    }
}

export default new TranslationService();
