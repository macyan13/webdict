<template>
  <div id="app">
    <div id="nav">
      <span v-if="this.$store.getters['auth/isLoggedIn']"><router-link to="/">Home</router-link> | </span>
      <span v-if="this.$store.getters['auth/isLoggedIn']"><router-link to="/newTag">Create Tag</router-link> | </span>
      <router-link to="/about">About</router-link> |
      <router-link v-if="!currentUser" to="/login">Login</router-link>
      <a to="" v-if="this.$store.getters['auth/isLoggedIn']" href @click.prevent="logOut">Logout</a>
    </div>
    <router-view/>
  </div>
</template>

<style>
@import './assets/css/app.css';
</style>

<script>
export default {
  computed: {
    currentUser() {
      return this.$store.state.auth.user;
    },
  },
  methods: {
    logOut() {
      this.$store.dispatch('auth/logout');
      this.$router.push({name: 'Login'});
    }
  }
};
</script>
