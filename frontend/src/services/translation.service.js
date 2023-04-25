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
}

export default new TranslationService();
