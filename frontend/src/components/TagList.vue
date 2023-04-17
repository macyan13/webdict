<template>
  <b-card :title="title">
    <div class="tag-list" style="display: flex; justify-content: center;">
      <b-list-group style="width: 40%;">
        <b-list-group-item v-for="tag in tags" :key="tag.id">
          {{ tag.tag }}
        </b-list-group-item>
      </b-list-group>
    </div>
    <div v-if="hasError" style="color: red;">There was an error processing your request - {{ errorMessage }}</div>
  </b-card>
</template>

<script>

export default {
  name: 'TagList',
  props: {
    title: {
      type: String,
      default() {
        return "Your tags";
      }
    },
  },
  data() {
    return {
      tags: [],
      hasError: false,
      errorMessage: ''
    }
  },
  mounted() {
    this.$store.dispatch('tag/fetchAll')
        .then((tags) => this.tags = tags)
        .catch(() => {
          this.hasError = true
          this.errorMessage = 'Can not get tags from server :('
        })
  },
};
</script>
