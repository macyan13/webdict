<template>
  <div>
    <table class="table" id="search-results">
      <thead>
      <tr>
        <th>Source</th>
        <th>Transcription</th>
        <th>Target</th>
        <th>Tags</th>
        <th>Actions</th>
      </tr>
      </thead>
      <tbody>
      <tr v-for="translation in translations" :key="translation.id" :id="translation.id">
        <td>{{ translation.source }}</td>
        <td>{{ translation.transcription }}</td>
        <td><vue-markdown>{{translation.target}}</vue-markdown></td>
        <td>
          <span v-for="tag in translation.tags" :key="tag.id" class="badge badge-primary">{{ tag.name }}</span>
        </td>
        <td>
          <button class="btn btn-sm btn-primary" @click="editTranslation(translation.id)">Edit</button>
          <button class="btn btn-sm btn-danger" @click="confirmDelete(translation.id)">Delete</button>
        </td>
        <b-popover :target="translation.id" triggers="hover" placement="top">
          <template #title>Usage example</template>
          <vue-markdown>{{translation.example}}</vue-markdown>
        </b-popover>
      </tr>
      </tbody>
    </table>
    <b-modal v-model="showConfirmationModal" title="Delete Translation?" hide-footer hide-backdrop>
      <p>Are you sure you want to delete this translation?</p>
      <div class="d-flex justify-content-end">
        <b-button variant="secondary" class="mr-2" @click="deleteCancel">
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
  </div>
</template>
<script>

import VueMarkdown from "vue-markdown";
import TranslationService from "@/services/translation.service";

export default {
  name: 'TranslationList',
  components: {
    VueMarkdown
  },
  props: {
    translations: {
      type: Array,
      default() {
        return [];
      }
    },
  },
  data() {
    return {
      showConfirmationModal: false,
      showDeleteSpinner: false,
      idToDelete: null,
    };
  },
  methods: {
    editTranslation(id) {
      this.$router.push(`/editTranslation/${id}`)
    },
    confirmDelete(id) {
      this.idToDelete = id;
      this.showConfirmationModal = true;
    },
    deleteCancel() {
      this.showConfirmationModal = false;
      this.idToDelete = null
    },
    deleteTranslation() {
      this.showDeleteSpinner = true;
      TranslationService.delete(this.idToDelete)
          .then(() => {
            this.$emit('onDelete');
          })
          .catch((error) => {
            this.hasError = true;
            this.errorMessage = error;
          })
          .finally(() => {
            this.showDeleteSpinner = false;
            this.showConfirmationModal = false;
            this.idToDelete = null;
          });
    },
  }
};
</script>
