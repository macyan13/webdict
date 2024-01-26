<template>
  <div class="container">
    <div v-if="showLoadSpinner" class="d-flex justify-content-center m-3">
      <b-spinner variant="primary" label="Spinning"></b-spinner>
    </div>
    <flash-message v-if="showFlashMessage" :message="flashMessage"/>
    <div class="row">
      <div class="col-md-10">
        <b-form-group
            id="search-keyword-group"
            label="Search keyword:"
            label-for="source-input"
            :state="!!source || !!target"
            invalid-feedback="required"
        >
          <div style="display: flex; justify-content: center;">
            <b-form-input
                id="source-input"
                v-model="source"
                :disabled="target.length > 0"
                placeholder="Enter a part of source text..."
                style="width: 40%;"
            ></b-form-input>
            <span></span>
            <b-form-input
                id="target-input"
                v-model="target"
                :disabled="source.length > 0"
                placeholder="Enter a part of target text..."
                style="width: 40%;"
            ></b-form-input>
          </div>
        </b-form-group>
        <!-- Translation search result -->
        <translation-list :translations="translations" @onDelete="refreshData"></translation-list>
        <b-pagination
            v-model="currentPage"
            :total-rows="totalRecords"
            :per-page="pageSize"
            last-number
            align="center"
            aria-controls="search-results"
            @change="search"
        ></b-pagination>
      </div>

      <div class="col-md-2">
        <b-form-group
            id="lang-group"
            label="Language:"
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
                :show-labels="false"
                label="name"
                track-by="id"
                deselectLabel=""
                placeholder="Pick a language"
            ></VueMultiselect>
          </div>
        </b-form-group>

        <b-form-group
            id="page-size-group"
            label="Translation per page:"
            label-for="page-input"
        >
          <b-form-select v-model="pageSize" :options="pageSizeOptions"></b-form-select>
        </b-form-group>
        <b-button variant="primary" class="mr-2" @click="search">
          Search
        </b-button>
      </div>
    </div>
    <div v-if="hasError" style="color: red;">{{ errorMessage }}</div>
  </div>
</template>

<style src="vue-multiselect/dist/vue-multiselect.min.css"></style>

<script>
import VueMultiselect from 'vue-multiselect'
import TranslationService from "@/services/translation.service";
import SearchParams from "@/models/searchParams";
import TranslationList from "@/components/TranslationList.vue";
import FlashMessage from "@/components/FlashMessage.vue";
import EntityStatusService from "@/services/entity-status.service";

export default {
  name: 'translationSearch',
  components: {
    TranslationList,
    VueMultiselect,
    FlashMessage,
  },
  data() {
    return {
      lang: null,
      langOptions: [],
      translations: [],
      source: '',
      target: '',
      hasError: false,
      errorMessage: '',
      showLoadSpinner: false,
      pageSizeOptions: [20, 30, 50, 100],
      flashMessage: '',
      showFlashMessage: false,
    };
  },
  async created() {
    this.showLoadSpinner = true;
    await this.fetchLangs();
    await this.initLang();
    this.initSearchState();
    this.showLoadSpinner = false;

    if (this.$store.getters["translation/entityStatus"] !== null) {
      this.triggerFlashMessage();
    }
  },
  computed: {
    currentPage: {
      get() {
        return this.$store.getters["translationSearch/getCurrentPage"];
      },
      set(value) {
        this.$store.dispatch('translationSearch/setCurrentPage', value);
      }
    },
    pageSize: {
      get() {
        return this.$store.getters["translationSearch/getPageSize"];
      },
      set(value) {
        this.$store.dispatch('translationSearch/setPageSize', value);
      }
    },
    totalRecords: {
      get() {
        return this.$store.getters["translationSearch/getTotalRecords"];
      },
      set(value) {
        this.$store.dispatch('translationSearch/setTotalRecords', value);
      }
    },
  },
  methods: {
    fetchLangs() {
      return this.$store.dispatch('lang/fetchAll')
          .then((langs) => {
            this.langOptions = langs;
          })
          .catch(() => {
            this.hasError = true;
            this.errorMessage = 'Can not get languages from server :(';
          })
    },
    initSearchState() {
      this.pageSize = this.$store.getters["translationSearch/getPageSize"];
      this.totalRecords = this.$store.getters["translationSearch/getTotalRecords"];
      this.currentPage = this.$store.getters["translationSearch/getCurrentPage"];
      this.translations = this.$store.getters["translationSearch/getTranslations"];
      this.target = this.$store.getters["translationSearch/getTarget"];
      this.source = this.$store.getters["translationSearch/getSource"];
    },

    async initLang() {
      if (this.langOptions.length > 0) {
        return this.$store.dispatch('profile/fetchProfile')
            .then((profile) => {
              if (profile.default_lang && profile.default_lang.id) {
                this.lang = profile.default_lang;
              } else {
                this.lang = this.langOptions[0];
              }
            })
            .catch((error) => {
              this.hasError = true;
              this.errorMessage = "Can not get user data from server: " + error;
            });
      }
      return Promise.resolve();
    },
    refreshData() {
      this.search();
    },
    triggerFlashMessage() {
      let status = this.$store.getters["translation/entityStatus"];
      if (status == null) {
        return;
      }
      this.flashMessage = EntityStatusService.getMessageByStatus("Translation", status);
      this.showFlashMessage = true;
      this.$store.dispatch('translation/clearEntityStatus');
      setTimeout(() => {
        this.showFlashMessage = false;
      }, 5000);
    },
    commitSearchState() {
      this.$store.dispatch('translationSearch/setPageSize', this.pageSize);
      this.$store.dispatch('translationSearch/setTotalRecords', this.totalRecords);
      this.$store.dispatch('translationSearch/setCurrentPage', this.currentPage);
      this.$store.dispatch('translationSearch/setTranslations', this.translations);
      this.$store.dispatch('translationSearch/setTarget', this.target);
      this.$store.dispatch('translationSearch/setSource', this.source);
    },
    onPaginationChange(page) {
      this.currentPage = page;
      this.search();
    },
    validate() {
      if (!this.lang) {
        this.hasError = true;
        this.errorMessage = 'Please select the language to perform translation search';
        return false;
      }

      if (!this.source && !this.target) {
        this.hasError = true;
        this.errorMessage = 'Please enter a part of source or target text to perform translation search';
        return false;
      }

      if (this.source && this.target) {
        this.hasError = true;
        this.errorMessage = 'Please enter only a part of source or target text to perform translation search';
        return false;
      }

      return true;
    },
    search() {
      if (!this.validate()) {
        return;
      }

      this.showLoadSpinner = true;

      TranslationService.search(new SearchParams([], this.lang.id, this.currentPage, this.pageSize, this.source, this.target))
          .then(searchResult => {
            this.translations = searchResult.translations;
            this.totalRecords = searchResult.total_records;
            this.hasError = false;
            this.commitSearchState();
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