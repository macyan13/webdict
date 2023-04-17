<template>
  <div>
    <TagForm :tag="currentTag" :edit-mode="isEditMode" @submit="handleSubmit"/>
    <div v-if="hasError" style="color: red;">There was an error processing your request - {{ errorMessage }}</div>
  </div>
</template>

<script>
import TagForm from '@/components/Tag.vue'
import Tag from '@/models/tag'
import TagService from '@/services/tag.service.js'
import router from "@/router";

export default {
  name: "NewTag",
  components: {
    TagForm
  },
  data() {
    return {
      currentTag: '',
      isEditMode: false,
      hasError: false,
      errorMessage: ''
    }
  },
  methods: {
    async handleSubmit(newTag) {
      try {
        console.log(newTag)
        let tag = new Tag(newTag)
        TagService.create(tag)
            .then((data) => {
              // todo: redirect to tag detail page or tags page
              console.log(data);
              router.push({name: 'Home'});
            })
            .catch((error) => {
              this.hasError = true;
              this.errorMessage = error;
            });
      } catch (error) {
        this.hasError = true;
        this.errorMessage = error;
      }
    },
  },
}
</script>

