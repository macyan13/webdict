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
