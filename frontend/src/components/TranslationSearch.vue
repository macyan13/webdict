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
                label="name"
                track-by="id"
                deselectLabel=""
                placeholder="Pick a language"
            ></VueMultiselect>
          </div>
        </b-form-group>

        <b-form-group
            id="tags-group"
            label="Tags:"
        >
          <div style="display: flex; justify-content: center;">
            <VueMultiselect
                :options="tagOptions"
                v-model="tags"
                :multiple="true"
                :max="5"
                :show-labels="false"
                label="name"
                track-by="id"
                placeholder="Pick a tag"
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
      tags: [],
      tagOptions: [],
      currentPage: 1,
      pageSize: 20,
      totalRecords: null,
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
  async mounted() {
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
      this.tags = this.$store.getters["translationSearch/getTags"];
      this.translations = this.$store.getters["translationSearch/getTranslations"];
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
      this.$store.dispatch('translationSearch/setPageSize', this.pageSize);
      this.$store.dispatch('translationSearch/setTotalRecords', this.totalRecords);
      this.$store.dispatch('translationSearch/setCurrentPage', this.currentPage);
      this.$store.dispatch('translationSearch/setTranslations', this.translations);
      this.$store.dispatch('translationSearch/setTags', this.tags);
    },
    search() {
      if (!this.lang) {
        this.hasError = true;
        this.errorMessage = 'Please select the language to perform translation search';
        return;
      }

      this.showLoadSpinner = true;
      let tagIds = this.tags.map(x => x.id);

      TranslationService.search(new SearchParams(tagIds, this.lang.id, this.$store.getters["translationSearch/getCurrentPage"], this.pageSize))
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