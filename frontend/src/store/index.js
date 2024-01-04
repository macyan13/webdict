import Vue from 'vue'
import Vuex from 'vuex'
import auth from './auth.module'
import tag from './tag.module'
import lang from './lang.module'
import profile from './profile.module'
import role from './role.module'
import user from './user.module'
import translation from "./translation.module";
import translationHome from "./translationHome.module";
import translationSearch from "./translationSearch.module";

Vue.use(Vuex)

export default new Vuex.Store({
    modules: {
        auth,
        tag,
        lang,
        profile,
        role,
        user,
        translation,
        translationHome,
        translationSearch
    }
})
