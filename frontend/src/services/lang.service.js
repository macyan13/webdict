import axios from 'axios';
import authHeader from "@/services/auth-header";

class LangService {
    get(id) {
        return new Promise((resolve, reject) => {
            axios
                .get('/v1/api/langs/' + id, {headers: authHeader()})
                .then(response => {
                    resolve(response.data);
                })
                .catch(error => {
                    console.log(error)
                    reject(error);
                });
        });
    }
    create(lang) {
        return new Promise((resolve, reject) => {
            axios
                .post('/v1/api/langs', lang, {headers: authHeader()})
                .then(response => {
                    resolve(response.data);
                })
                .catch(error => {
                    console.log(error)
                    reject(error);
                });
        });
    }
    update(lang) {
        return new Promise((resolve, reject) => {
            axios
                .put('/v1/api/langs/' + lang.id, {
                    name: lang.name,
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
                .delete('/v1/api/langs/' + id, {headers: authHeader()})
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
                .get('/v1/api/langs', {headers: authHeader()})
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

export default new LangService();
