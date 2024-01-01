<template>
  <div>
    <b-card :title="title">
      <b-form @submit.prevent="submitForm">
        <b-form-group
            id="tag-group"
            label="Name:"
            label-for="name-input"
            :state="name ? true : false"
            invalid-feedback="required"
        >
          <div style="display: flex; justify-content: center;">
            <b-form-input
                :required=true
                id="name-input"
                v-model="name"
                placeholder="Enter a tag..."
                class="w-25"
            ></b-form-input>
          </div>
        </b-form-group>
        <b-button type="submit" variant="primary">
          {{ buttonLabel }}
        </b-button>
        <b-button v-if="id" variant="danger" @click="confirmDelete">
          Delete
        </b-button>
        <div v-if="showEditSpinner" class="d-flex justify-content-center m-3">
          <b-spinner variant="primary" label="Spinning"></b-spinner>
        </div>
      </b-form>

      <div v-if="hasError" style="color: red;">{{errorMessage}}</div>

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
        <div v-if="showDeleteSpinner" class="d-flex justify-content-center mb-3">
          <b-spinner variant="danger" label="Spinning"></b-spinner>
        </div>
      </b-modal>

    </b-card>
  </div>
</template>

<script>
import TagService from "@/services/tag.service";
import Tag from "@/models/tag";
import router from "@/router";
import EntityStatusService from "@/services/entity-status.service";

export default {
  name: 'Tag',
  props: {
    id: {
      type: String,
      default: null,
    },
  },
  data() {
    return {
      title: '',
      buttonLabel: '',
      name: '',
      showConfirmationModal: false,
      showDeleteSpinner: false,
      showEditSpinner: false,
      hasError: false,
      errorMessage: '',
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
            this.name = data.name;
          })
          .catch((error) => {
            this.hasError = true;
            this.errorMessage = "Can not get tag data from server: " + error;
          });
    },
    confirmDelete() {
      this.showConfirmationModal = true;
    },
    deleteTag() {
      this.showDeleteSpinner = true;
      TagService.delete(this.id)
          .then(() => {
            this.$store.dispatch('tag/setEntityStatus',  EntityStatusService.deleted());
            this.$store.dispatch('tag/clear');
            this.$store.dispatch('translationSearch/resetTags');
            router.push({name: 'Tags'});
          })
          .catch((error) => {
            this.hasError = true;
            this.errorMessage = error;
          })
          .finally(() => {
            this.showDeleteSpinner = false;
            this.showConfirmationModal = false;
          });
    },
    submitForm() {
      this.showEditSpinner = true;
      let method = this.id ? TagService.update : TagService.create;
      method(new Tag(this.name, this.id))
          .then(() => {
            this.$store.dispatch('tag/clear');
            this.$store.dispatch('tag/setEntityStatus',  this.id ? EntityStatusService.updated() : EntityStatusService.created());
            router.push({name: 'Tags'});
          })
          .catch((error) => {
            this.hasError = true;
            this.errorMessage = error;
          })
          .finally(() => {
            this.showEditSpinner = false;
          });
    },
  },
}
</script>