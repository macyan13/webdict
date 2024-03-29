<template>
  <div class="container">
    <div v-if="showLoadSpinner" class="d-flex justify-content-center m-3">
      <b-spinner variant="primary" label="Spinning"></b-spinner>
    </div>
    <flash-message v-if="showFlashMessage" :message="flashMessage"/>
    <div class="row">
      <div class="col-md-10">
        <!-- Translation search result -->
        <translation-list :translations="translations" @onDelete="refreshData"></translation-list>
        <b-pagination
            v-model="currentPage"
            :total-rows="totalRecords"
            :per-page="pageSize"
            last-number
            align="center"
            aria-controls="search-results"
            @change="onPaginationChange"
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
                v-model="lang"
                :allow-empty="false"
                :multiple="false"
                :options="langOptions"
                :preselect-first="true"
                :show-labels="false"
                label="name"
                placeholder="Pick a language"
                track-by="id"
            ></VueMultiselect>
          </div>
        </b-form-group>

        <b-form-group
            id="tags-group"
            label="Tags:"
        >
          <div style="display: flex; justify-content: center;">
            <VueMultiselect
                v-model="tags"
                :max="5"
                :multiple="true"
                :options="tagOptions"
                :show-labels="false"
                label="name"
                placeholder="Pick a tag"
                track-by="id"
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
    <div class="mt-5">
      <b-button block variant="info" @click="onRandomTranslationsClick">Show random translations</b-button>
    </div>
    <div class="row">
      <div class="col-md-10">
        <div v-if="showRandomTranslations && lang">
          <translation-list :translations="randomTranslations" @onDelete="refreshData"></translation-list>
        </div>
      </div>

      <div class="col-md-2">
        <div v-if="showRandomTranslations && lang">
          <b-form-group
              id="page-size-group"
              label="Random items amount:"
              label-for="page-input"
          >
            <b-form-select v-model="randomLimit" :options="randomLimitOptions"></b-form-select>
          </b-form-group>
          <b-button variant="primary" @click="fetchRandomTranslations">Refresh</b-button>
        </div>
      </div>
    </div>
    <div v-if="hasError" style="color: red;">{{errorMessage}}</div>
  </div>
</template>

<style src="vue-multiselect/dist/vue-multiselect.min.css"></style>

<script>
import VueMultiselect from 'vue-multiselect'
import TranslationService from "@/services/translation.service";
import SearchParams from "@/models/searchParams";
import RandomParams from "@/models/randomParams";
import TranslationList from "@/components/TranslationList.vue";
import FlashMessage from "@/components/FlashMessage.vue";
import EntityStatusService from "@/services/entity-status.service";

export default {
  name: 'translationHome',
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
      tags: [],
      tagOptions: [],
      hasError: false,
      errorMessage: '',
      showLoadSpinner: false,
      pageSizeOptions: [20, 30, 50, 100],
      showRandomTranslations: false,
      randomTranslations: [],
      randomLimit: 10,
      randomLimitOptions: [10, 15, 20],
      flashMessage: '',
      showFlashMessage: false,
    };
  },
  async created() {
    this.showLoadSpinner = true;
    await this.fetchTags();
    await this.fetchLangs();
    await this.initLang();
    this.initSearchState();
    this.initSearch();
    this.showLoadSpinner = false;

    if (this.$store.getters["translation/entityStatus"] !== null) {
      this.triggerFlashMessage();
    }
  },
  computed: {
    currentPage: {
      get() {
        return this.$store.getters["translationHome/getCurrentPage"];
      },
      set(value) {
        this.$store.dispatch('translationHome/setCurrentPage', value);
      }
    },
    pageSize: {
      get() {
        return this.$store.getters["translationHome/getPageSize"];
      },
      set(value) {
        this.$store.dispatch('translationHome/setPageSize', value);
      }
    },
    totalRecords: {
      get() {
        return this.$store.getters["translationHome/getTotalRecords"];
      },
      set(value) {
        this.$store.dispatch('translationHome/setTotalRecords', value);
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
      this.pageSize = this.$store.getters["translationHome/getPageSize"];
      this.totalRecords = this.$store.getters["translationHome/getTotalRecords"];
      this.currentPage = this.$store.getters["translationHome/getCurrentPage"];
      this.tags = this.$store.getters["translationHome/getTags"];
      this.translations = this.$store.getters["translationHome/getTranslations"];
    },
    initSearch() {
      if (this.lang && this.translations.length === 0) {
        this.search();
      }
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

      if (this.showRandomTranslations) {
        this.fetchRandomTranslations();
      }
    },
    onRandomTranslationsClick() {
      if (!this.showRandomTranslations) {
        if (this.randomTranslations.length === 0) {
          this.fetchRandomTranslations();
        }
        this.showRandomTranslations = true;
        return;
      }
      this.showRandomTranslations = false;
    },
    fetchRandomTranslations() {
      if (!this.lang) {
        this.hasError = true;
        this.errorMessage = 'Please select the language to get random translations';
        return;
      }

      this.showLoadSpinner = true;
      let tagIds = this.tags.map(x => x.id);

      TranslationService.random(new RandomParams(tagIds, this.lang.id, this.randomLimit))
          .then(searchResult => {
            this.randomTranslations = searchResult.translations;
            this.hasError = false;
          })
          .catch((error) => {
            this.hasError = true;
            this.errorMessage = error;
          })
          .finally(() => {
            this.showLoadSpinner = false;
          });
    },
    fetchTags() {
      return this.$store.dispatch('tag/fetchAll')
          .then((tags) => this.tagOptions = tags)
          .catch(() => {
            this.hasError = true;
            this.errorMessage = 'Can not get tags from server :(';
          })
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
      this.$store.dispatch('translationHome/setPageSize', this.pageSize);
      this.$store.dispatch('translationHome/setTotalRecords', this.totalRecords);
      this.$store.dispatch('translationHome/setCurrentPage', this.currentPage);
      this.$store.dispatch('translationHome/setTranslations', this.translations);
      this.$store.dispatch('translationHome/setTags', this.tags);
    },
    onPaginationChange(page) {
      this.currentPage = page;
      this.search();
    },
    search() {
      if (!this.lang) {
        this.hasError = true;
        this.errorMessage = 'Please select the language to perform translation search';
        return;
      }

      this.showLoadSpinner = true;
      let tagIds = this.tags.map(x => x.id);

      TranslationService.search(new SearchParams(tagIds, this.lang.id, this.currentPage, this.pageSize))
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