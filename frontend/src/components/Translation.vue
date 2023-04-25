<template>
  <div>
    <b-card :title="title">
      <b-form @submit.prevent="submitForm">
        <b-form-group
            id="source-group"
            label="Source text:"
            label-for="source-input"
        >
          <div style="display: flex; justify-content: center;">
            <b-form-input
                :required=true
                id="source-input"
                v-model="source"
                placeholder="Enter source text..."
                class="w-25"
            ></b-form-input>
          </div>
        </b-form-group>

        <b-form-group
            id="transcription-group"
            label="Transcription:"
            label-for="transcription-input"
        >
          <div style="display: flex; justify-content: center;">
            <b-form-input
                class="w-25"
                :required=true
                id="transcription-input"
                v-model="transcription"
                placeholder="Enter transcription..."
            ></b-form-input>
          </div>
        </b-form-group>

        <b-form-group
            id="target-group"
            label="Target text:"
            label-for="target-input"
        >
          <div style="display: flex; justify-content: center;">
            <b-form-input
                :required=true
                id="target-input"
                v-model="target"
                placeholder="Enter target text..."
                class="w-25"
            ></b-form-input>
          </div>
        </b-form-group>

        <b-form-group
            id="tags-group"
            label="Tags:"
            label-for="tags-input"
        >
          <div style="display: flex; justify-content: center;">
            <VueMultiselect
                :options="tagOptions"
                v-model="tags"
                :multiple="true"
                :show-labels="false"
                label="tag"
                track-by="id"
                placeholder="Pick a tag"
                style="width: 25%"
            ></VueMultiselect>
          </div>
        </b-form-group>

        <b-form-group
            id="example-group"
            label="Example:"
            label-for="example-input"
        >
          <div style="display: flex; justify-content: center;">
            <b-form-textarea
                :required=true
                id="example-input"
                v-model="example"
                placeholder="Enter example..."
                style="width: 40%"
            ></b-form-textarea>
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

      <b-modal v-model="showConfirmationModal" title="Delete Translation?" hide-footer hide-backdrop>
        <p>Are you sure you want to delete this Translation?</p>
        <div class="d-flex justify-content-end">
          <b-button variant="secondary" class="mr-2" @click="showConfirmationModal = false">
            Cancel
          </b-button>
          <b-button variant="danger" @click="deleteTranslation">
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

<style src="vue-multiselect/dist/vue-multiselect.min.css"></style>

<script>
import TranslationService from "@/services/translation.service";
import Translation from "@/models/translation";
import router from "@/router";
import VueMultiselect from 'vue-multiselect'

export default {
  name: 'Translation',
  components: {
    VueMultiselect
  },
  props: {
    id: {
      type: String,
      default: null,
    },
  },
  data() {
    return {
      source: '',
      transcription: '',
      target: '',
      example: '',
      tags: [],
      tagOptions: ["tag1", "tag2", "tag3"],
      buttonLabel: '',
      title: '',
      showConfirmationModal: false,
      showDeleteSpinner: false,
      showEditSpinner: false,
    }
  },
  mounted() {
    if (this.id) {
      this.loadData();
      this.title = 'Edit Translation';
      this.buttonLabel = 'Save';
    } else {
      this.title = 'Create New Translation';
      this.buttonLabel = 'Create';
    }
    this.fetchTags();
  },
  methods: {
    fetchTags() {
      this.$store.dispatch('tag/fetchAll')
          .then((tags) => this.tagOptions = tags)
          .catch(() => {
            this.hasError = true
            this.errorMessage = 'Can not get tags from server :('
          })
    },
    loadData() {
      // todo: load data from server
    },
    confirmDelete() {
      this.showConfirmationModal = true;
    },
    deleteTranslation() {
      // this.showDeleteSpinner = true;
      // TagService.delete(this.id)
      //     .then(() => {
      //       this.$store.dispatch('tag/clear');
      //       router.push({name: 'Tags'});
      //     })
      //     .catch((error) => {
      //       this.hasError = true;
      //       this.errorMessage = error;
      //     })
      //     .finally(() => {
      //       this.showDeleteSpinner = false;
      //       this.showConfirmationModal = false;
      //     });
    },
    submitForm() {
      this.showEditSpinner = true;
      // let method = this.id ? TranslationService.update : TranslationService.create;
      let method = TranslationService.create;
      let tagIds = this.tags.map((tag) => tag.id);
      method(new Translation(this.source, this.transcription, this.target, this.example, tagIds))
          .then(() => {
            // this.$store.dispatch('tag/clear');
            router.push({name: 'Home'});
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
