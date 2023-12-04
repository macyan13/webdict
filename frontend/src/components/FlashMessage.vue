<template>
  <div class="flash-message">
    {{ message }}
  </div>
</template>

<style >
@import '../assets/css/components/FlashMessage.css';
</style>

<script>
export default {
  props: {
    message: {
      type: String,
      default: null,
    },
  },
};
</script>