import Vue from 'vue'
import Vuex from 'vuex'
import auth from './auth.module'
import tag from './tag.module'
import lang from './lang.module'

Vue.use(Vuex)

export default new Vuex.Store({
    modules: {
        auth,
        tag,
        lang,
    }
})
