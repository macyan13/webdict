<template>
  <div class="home">
    <img alt="Vue logo" src="../assets/logo.png">
    <Section type="Main vocabulary" request-type="phrases"></Section>
<!--    <Section type="Phrasal verbs" src-data="/api/v1/translation/phrasal-verbs" src-count="src-data"></Section>-->
<!--    <Section type="Useful phrases" src-data="/api/v1/translation/words" src-count="src-data"></Section>-->
  </div>
</template>

<script>
import Section from '@/components/VocabularySection.vue'

export default {
  name: 'Home',
  components: {
    Section
  }
}
</script>
