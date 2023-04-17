import axios from 'axios';
import authHeader from "@/services/auth-header";

class TagService {
    create(tag) {
        return axios
            .post('/v1/api/tags', {
                    tag: tag.tag,
                },
                {headers: authHeader()}
            )
            .then(response => {
                    console.log(response);
                    // todo: cleaning store on create or update
                    return response.data;
                }
            );
    }
}

export default new TagService();
