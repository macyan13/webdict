import axios from 'axios';
import authHeader from "@/services/auth-header";

class UserService {
    create(user) {
        return new Promise((resolve, reject) => {
            axios
                .post('/v1/api/users', user, {headers: authHeader()})
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
                .get('/v1/api/users', {headers: authHeader()})
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
                .get('/v1/api/users/' + id, {headers: authHeader()})
                .then(response => {
                    resolve(response.data);
                })
                .catch(error => {
                    console.log(error)
                    reject(error);
                });
        });
    }
    update(user) {
        return new Promise((resolve, reject) => {
            axios
                .put('/v1/api/users/' + user.id, {
                    name: user.name,
                    email: user.email,
                    password: user.password,
                    role: user.role,
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

export default new UserService();
