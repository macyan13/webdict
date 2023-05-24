<template>
  <div class="home">
<!--    <img alt="Vue logo" src="../assets/logo.png">-->
<!--    <div>Hi, you are on home page!</div>-->
    <translation-search></translation-search>
  </div>
</template>

<script>
import TranslationSearch from "@/components/TranslationSearch.vue";

export default {
  name: 'Home',
  components: {TranslationSearch},
}
</script>
