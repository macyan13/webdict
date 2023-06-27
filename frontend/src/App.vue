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
      <span v-if="isLoggedIn && isAdmin"><router-link to="/newUser">Add User</router-link> | </span>
      <span v-if="isLoggedIn && isAdmin"><router-link to="/users">Users</router-link> | </span>
      <router-link to="/about">About</router-link> |
      <router-link v-if="!isLoggedIn" to="/login">Login</router-link>
      <a to="" v-if="isLoggedIn" href @click.prevent="logOut">Logout</a>
    </div>
    <router-view/>
    <div v-if="hasError" style="color: red;">{{errorMessage}}</div>
  </div>
</template>

<style>
@import './assets/css/app.css';
</style>

<script>
export default {
  data() {
    return {
      hasError: false,
      errorMessage: "",
    }
  },
  mounted() {
    this.fetchProfile();
  },
  computed: {
    isLoggedIn() {
      return this.$store.getters['auth/isLoggedIn'];
    },
    isAdmin() {
      return this.$store.getters['profile/isAdmin'];
    }
  },
  methods: {
    fetchProfile() {
      if (!this.isLoggedIn) {
        return;
      }
      this.$store.dispatch('profile/fetchProfile');
    },
    logOut() {
      this.$store.dispatch('auth/logout');

      if (this.$router.currentRoute && this.$router.currentRoute.name !== 'Login') {
        this.$router.push({name: 'Login'});
      }
    }
  }
};
</script>
