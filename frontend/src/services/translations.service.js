import axios from 'axios';
import authHeader from './auth-header';

class TranslationsService {
    getPhrases() {
        return axios.get("phrases", {headers: authHeader()});
    }

    getPhrasalVerbs() {
        return axios.get("phrasal-verbs", {headers: authHeader()});
    }

    getWords() {
        return axios.get("words", {headers: authHeader()});
    }

    getByType(type) {
        return axios.get('/api/translation/' + type, {headers: authHeader()});
    }
}

export default new TranslationsService();