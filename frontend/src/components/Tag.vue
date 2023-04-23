<template>
  <div>
    <b-card :title="title">
      <b-form @submit.prevent="submitForm">
        <b-form-group
            id="tag-group"
            label="Tag:"
            label-for="tag-input"
        >
          <div style="display: flex; justify-content: center;">
            <b-form-input
                :required=true
                id="tag-input"
                v-model="tag"
                :placeholder="tagPlaceholder"
                style="width: 40%;"
            ></b-form-input>
          </div>
        </b-form-group>
        <b-button type="submit" variant="primary">
          {{ buttonLabel }}
        </b-button>
      </b-form>
    </b-card>
  </div>
</template>

<script>
import TagService from "@/services/tag.service";
import Tag from "@/models/tag";
import router from "@/router";

export default {
  name: 'Tag',
  props: {
    id: {
      type: String,
      default: null,
    },
    initTag: {
      type: String,
      default: ''
    },
  },
  data() {
    return {
      title: '',
      buttonLabel: '',
      tag: '',
      tagPlaceholder: 'Enter a tag...',
    }
  },
  mounted() {
    if (this.id) {
      this.loadData();
      this.title = 'Edit Tag'
      this.buttonLabel = 'Save Changes'
    } else {
      this.title = 'Create New Tag'
      this.buttonLabel = 'Submit'
    }
  },
  methods: {
    loadData() {
      TagService.get(this.id)
          .then((data) => {
            this.tag = data.tag;
          })
          .catch((error) => {
            this.hasError = true;
            this.errorMessage = error;
          });
    },
    submitForm() {
      let method = this.id ? TagService.update : TagService.create;
      method(new Tag(this.tag, this.id))
          .then(() => {
            this.$store.dispatch('tag/clear');
            router.push({name: 'Tags'});
          })
          .catch((error) => {
            this.hasError = true;
            this.errorMessage = error;
          });
    },
    // updateTag() {
    //   TagService.update(Tag(this.tag, this.id))
    //       .then(() => {
    //         this.$store.dispatch('tag/clear');
    //         router.push({name: 'Tags'});
    //       })
    //       .catch((error) => {
    //         this.hasError = true;
    //         this.errorMessage = error;
    //       });
    // },
    // createTag() {
    //   TagService.create(Tag(this.tag, this.id))
    //       .then(() => {
    //         this.$store.dispatch('tag/clear');
    //         router.push({name: 'Tags'});
    //       })
    //       .catch((error) => {
    //         this.hasError = true;
    //         this.errorMessage = error;
    //       });
    // }
  },
}
</script>