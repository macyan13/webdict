import RoleService from '../services/role.service';

export default {
    namespaced: true,
    state: {
        roles: [],
    },
    actions: {
        fetchAll({state, commit}) {
            if (state.roles.length > 0) {
                return Promise.resolve(state.roles);
            }
            return RoleService.getAll().then(
                response => {
                    commit('allRolesSuccess', response.roles);
                    return Promise.resolve(response.roles);
                },
                error => {
                    console.log("Error: " + error);
                    return Promise.reject(error);
                });
        },
    },
    mutations: {
        allRolesSuccess(state, roles) {
            state.roles = roles;
        },
    },
};
