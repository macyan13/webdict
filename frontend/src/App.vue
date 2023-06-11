<template>
  <div id="app">
    <div id="nav">
      <span v-if="isLoggedIn"><router-link to="/">Home</router-link> | </span>
      <span v-if="isLoggedIn"><router-link to="/newTranslation">Add Translation</router-link> | </span>
      <span v-if="isLoggedIn"><router-link to="/newTag">Add Tag</router-link> | </span>
      <span v-if="isLoggedIn"><router-link to="/tags">Tags</router-link> | </span>
      <span v-if="isLoggedIn"><router-link to="/newLang">Add Language</router-link> | </span>
      <span v-if="isLoggedIn"><router-link to="/langs">Languages</router-link> | </span>
      <span v-if="isLoggedIn"><router-link to="/profile">Profile</router-link> | </span>
      <router-link to="/about">About</router-link> |
      <router-link v-if="!isLoggedIn" to="/login">Login</router-link>
      <a to="" v-if="isLoggedIn" href @click.prevent="logOut">Logout</a>
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
    isLoggedIn() {
      return this.$store.getters['auth/isLoggedIn'];
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
