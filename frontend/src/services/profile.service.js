import axios from 'axios';
import authHeader from "@/services/auth-header";

class ProfileService {
    get() {
        return new Promise((resolve, reject) => {
            axios
                .get('/v1/api/profile', {headers: authHeader()})
                .then(response => {
                    resolve(response.data);
                })
                .catch(error => {
                    console.log(error)
                    reject(error);
                });
        });
    }
    update(profile) {
        return new Promise((resolve, reject) => {
            axios
                .put('/v1/api/profile', {
                    name: profile.name,
                    email: profile.email,
                    current_password: profile.currentPassword,
                    new_password: profile.newPassword,
                    default_lang_id: profile.defaultLangId,
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
}

export default new ProfileService();
