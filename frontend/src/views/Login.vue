<template>
  <div class="login-wrapper">
    <div class="already-logged-in" v-if="!isLoggedIn()">
      <div v-if="loginError">
        <h1>
          <b-badge variant="danger">You don't have rights here, mate :D</b-badge>
        </h1>
        <h5>Seams that you don't have access rights... </h5>
      </div>
      <div class="unprotected">
        <h1>
          <b-badge variant="info">Please login to get access!</b-badge>
        </h1>
        <h5  v-if="!loginError">You're not logged in - so you don't see much here. Try to log in:</h5>

        <form @submit.prevent="callLogin()">
          <input type="email" placeholder="email" v-model="email" required="required">
          <input type="password" placeholder="password" v-model="password" required="required">
          <b-btn variant="success" type="submit">Login</b-btn>
        </form>
      </div>
    </div>
    <div class="login" v-if="isLoggedIn()">
      <h1>
        <b-badge variant="success">You are already logged in. Please visit Home or Profile page</b-badge>
      </h1>
    </div>
  </div>
</template>

<script>
import AuthParams from '../models/authParams';

export default {
  name: 'login',

  data() {
    return {
      loginError: false,
      email: '',
      password: '',
      error: []
    }
  },
  methods: {
    isLoggedIn() {
      return this.$store.getters['auth/isLoggedIn']
    },
    callLogin() {
      let authParams = new AuthParams(this.email, this.password)
      this.$store.dispatch('auth/login', authParams)
          .then(() => {
            this.$store.dispatch('profile/fetchProfile');
            this.$router.push({name: 'Home'});
          })
          .catch(error => {
            this.loginError = true
            this.error.push(error)
          })
    }
  }
}

</script>