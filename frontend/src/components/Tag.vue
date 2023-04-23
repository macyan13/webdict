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
        <b-button v-if="id" variant="danger" @click="confirmDelete">
          Delete
        </b-button>
      </b-form>
      <b-modal v-model="showConfirmationModal" title="Delete Tag?" hide-footer hide-backdrop>
        <p>Are you sure you want to delete this tag?</p>
        <div class="d-flex justify-content-end">
          <b-button variant="secondary" class="mr-2" @click="showConfirmationModal = false">
            Cancel
          </b-button>
          <b-button variant="danger" @click="deleteTag">
            Delete
          </b-button>
        </div>
      </b-modal>
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
      showConfirmationModal: false,
    }
  },
  mounted() {
    if (this.id) {
      this.loadData();
      this.title = 'Edit Tag'
      this.buttonLabel = 'Save'
    } else {
      this.title = 'Create New Tag'
      this.buttonLabel = 'Create'
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
    confirmDelete() {
      this.showConfirmationModal = true;
    },
    deleteTag() {
      TagService.delete(this.id)
          .then(() => {
            this.$store.dispatch('tag/clear');
            router.push({name: 'Tags'});
          })
          .catch((error) => {
            this.hasError = true;
            this.errorMessage = error;
          });
      this.showConfirmationModal = false;
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
  },
}
</script>