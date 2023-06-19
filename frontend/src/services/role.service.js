import axios from 'axios';
import authHeader from "@/services/auth-header";

class RoleService {
    getAll() {
        return new Promise((resolve, reject) => {
            axios
                .get('/v1/api/roles', {headers: authHeader()})
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

export default new RoleService();
