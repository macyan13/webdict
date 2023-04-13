import Vue from 'vue'
import Vuex from 'vuex'
import App from './App.vue'
import router from './router'
import store from './store'
import axios from 'axios'
import axiosInterceptor from './services/auth.interceptor.service';
import BootstrapVue from 'bootstrap-vue'

import "bootstrap/dist/css/bootstrap.min.css"

Vue.config.productionTip = false
Vue.prototype.$http = axios;
axiosInterceptor();

Vue.use(BootstrapVue)
Vue.use(Vuex)

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
