<template>
  <div class="container">
    <div v-if="showLoadSpinner" class="d-flex justify-content-center m-3">
      <b-spinner variant="primary" label="Spinning"></b-spinner>
    </div>
    <div class="row">
      <div class="col-md-10">
        <!-- Translation search result -->
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
          <tr v-for="translation in translations" :key="translation.id">
            <td>{{ translation.source }}</td>
            <td>{{ translation.transcription }}</td>
            <td>{{ translation.target }}</td>
            <td>
              <span v-for="tag in translation.tags" :key="tag.id" class="badge badge-primary">{{ tag.tag }}</span>
            </td>
            <td>
              <button class="btn btn-sm btn-primary" @click="editTranslation(translation.id)">Edit</button>
              <button class="btn btn-sm btn-danger" @click="confirmDelete(translation.id)">Delete</button>
            </td>
          </tr>
          </tbody>
        </table>
        <b-pagination
            v-model="currentPage"
            :total-rows="pageSize * totalPages"
            :per-page="pageSize"
            last-number
            align="center"
            aria-controls="search-results"
        ></b-pagination>
      </div>

      <div class="col-md-2">
        <!-- Language input -->
        <b-form-group
            id="lang-group"
            label="Language:"
            label-for="lang-input"
            :state="langOptions.length > 0 ? true : false"
            invalid-feedback="Please create at least one language"
        >
          <div style="display: flex; justify-content: center;">
            <VueMultiselect
                :preselect-first="true"
                :allow-empty="false"
                :options="langOptions"
                v-model="lang"
                :multiple="false"
                label="name"
                track-by="id"
                deselectLabel=""
                placeholder="Pick a language"
            ></VueMultiselect>
          </div>
        </b-form-group>

        <!-- Tags input -->
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
                :max="5"
                :show-labels="false"
                label="tag"
                track-by="id"
                placeholder="Pick a tag"
            ></VueMultiselect>
          </div>
        </b-form-group>
        <b-button variant="primary" class="mr-2" @click="search">
          Search
        </b-button>
      </div>
    </div>
    <div v-if="hasError" style="color: red;">{{errorMessage}}</div>
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

<style src="vue-multiselect/dist/vue-multiselect.min.css"></style>

<script>
import VueMultiselect from 'vue-multiselect'
import TranslationService from "@/services/translation.service";
import SearchParams from "@/models/searchParams";

export default {
  name: 'translationSearch',
  components: {
    VueMultiselect
  },
  data() {
    return {
      lang: null,
      langOptions: [],
      tags: [],
      tagOptions: [],
      currentPage: 1,
      pageSize: 20,
      totalPages: null,
      hasError: false,
      errorMessage: '',
      showLoadSpinner: false,
      showConfirmationModal: false,
      showDeleteSpinner: false,
      idToDelete: null,
      translationResult: '',
      translations: [
        {
          id: 25,
          source: "we",
          transcription: '[we]',
          target: 'мы',
          tags: [{id: 1, tag: "business"}, {id: 2, tag: "php"}]
        }
      ]
    };
  },
  mounted() {
    this.showLoadSpinner = true;
    this.fetchTags();
    this.fetchLangsAndInitSearch();
    this.showLoadSpinner = false;
  },
  created() {
    this.showLoadSpinner = true;
    this.search();
    this.showLoadSpinner = false;
  },
  methods: {
    fetchLangsAndInitSearch() {
      this.$store.dispatch('lang/fetchAll')
          .then((langs) =>{
            this.langOptions = langs;
            this.lang = langs.length > 0 ? langs[0] : null;
            this.search()
          })
          .catch(() => {
            this.hasError = true;
            this.errorMessage = 'Can not get languages from server :(';
          })
    },
    fetchTags() {
      this.$store.dispatch('tag/fetchAll')
          .then((tags) => this.tagOptions = tags)
          .catch(() => {
            this.hasError = true;
            this.errorMessage = 'Can not get tags from server :(';
          })
    },
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
            this.search();
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
    search() {
      if (!this.lang) {
        this.hasError = true;
        this.errorMessage = 'Please select the language to perform translation search';
        return
      }

      this.showLoadSpinner = true;
      let tagIds = this.tags.map(x => x.id);

      TranslationService.search(new SearchParams(tagIds, this.lang.id, this.currentPage, this.pageSize))
          .then(searchResult => {
            this.translations = searchResult.translations;
            this.totalPages = searchResult.total_pages;
            this.hasError = false;
          })
          .catch((error) => {
            this.hasError = true;
            this.errorMessage = error;
          })
          .finally(() => {
            this.showLoadSpinner = false;
          });
    }
  }
};
</script>