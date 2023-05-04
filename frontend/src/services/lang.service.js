import axios from 'axios';
import authHeader from "@/services/auth-header";

class LangService {
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
