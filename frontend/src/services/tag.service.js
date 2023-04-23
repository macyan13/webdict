import axios from 'axios';
import authHeader from "@/services/auth-header";

class TagService {
    create(tag) {
        return new Promise((resolve, reject) => {
            axios
                .post('/v1/api/tags', {
                    tag: tag.tag,
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
                    tag: tag.tag,
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
}

export default new TagService();
