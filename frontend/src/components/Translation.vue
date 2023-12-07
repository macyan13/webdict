<template>
  <div>
    <b-card :title="title">
      <b-form @submit.prevent="submitForm">
        <div v-if="hasError" style="color: red;">{{ errorMessage }}</div>
        <flash-message v-if="showFlashMessage" :message="flashMessage"/>
        <div class="row">
          <div class="col-md-1"/>
          <div class="col-md-7">
            <b-form-group
                id="source-group"
                label="Source text:"
                label-for="source-input"
                :state="source ? true : false"
                invalid-feedback="required"
            >
              <div style="display: flex; justify-content: center;">
                <b-form-input
                    :required=true
                    id="source-input"
                    v-model="source"
                    placeholder="Enter source text..."
                    style="width: 40%;"
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
                    style="width: 40%;"
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
                :state="target ? true : false"
                invalid-feedback="required"
            >
              <div style="display: flex; justify-content: center;">
                <mavon-editor
                    v-model="target"
                    placeholder="Enter target text..."
                    language="en"
                    style="width: 75%;"
                    :toolbars="markdownOption">
                </mavon-editor>
              </div>
            </b-form-group>

            <b-form-group
                id="example-group"
                label="Example:"
                label-for="example-input"
                class="markdown-editor"
            >
              <div style="display: flex; justify-content: center;">
                <mavon-editor
                    v-model="example"
                    placeholder="Enter example..."
                    language="en"
                    style="width: 75%;"
                    :toolbars="markdownOption">
                </mavon-editor>
              </div>
            </b-form-group>
          </div>

          <div class="col-md-3">
            <b-form-group
                id="lang-group"
                label="Language:"
                label-for="lang-input"
                :state="lang ? true : false"
                invalid-feedback="required"
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
                    style="width: 75%"
                ></VueMultiselect>
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
                    :max="5"
                    :show-labels="false"
                    label="name"
                    track-by="id"
                    placeholder="Pick a tag"
                    style="width: 85%"
                ></VueMultiselect>
              </div>
            </b-form-group>

            <b-form-group
                v-if="id"
                id="createdAt-group"
                label="Created:"
                label-for="createdAt-input"
            >
              <div class="d-inline-flex">
                <b-form-input
                    :disabled="true"
                    id="createdAt-input"
                    v-model="createdAtFormatted"
                    class="flex-grow-1"
                ></b-form-input>
              </div>
            </b-form-group>

            <b-button type="submit" variant="primary">
              {{ buttonLabel }}
            </b-button>
            <b-button v-if="id" variant="danger" @click="confirmDelete">
              Delete
            </b-button>
          </div>
          <div class="col-md-1"/>
        </div>

        <div v-if="showEditSpinner" class="d-flex justify-content-center m-3">
          <b-spinner variant="primary" label="Spinning"></b-spinner>
        </div>
      </b-form>

      <b-modal v-model="showConfirmationModal" title="Delete Translation?" hide-footer hide-backdrop>
        <p>Are you sure you want to delete this translation?</p>
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
<style src="mavon-editor/dist/css/index.css"></style>

<script>
import TranslationService from "@/services/translation.service";
import Translation from "@/models/translation";
import router from "@/router";
import VueMultiselect from 'vue-multiselect'
import {mavonEditor} from 'mavon-editor'
import FlashMessage from "@/components/FlashMessage";
import EntityStatusService from "@/services/entity-status.service";

export default {
  name: 'Translation',
  components: {
    VueMultiselect,
    mavonEditor,
    FlashMessage
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
      tagOptions: [],
      buttonLabel: '',
      title: '',
      lang: null,
      langOptions: [],
      showConfirmationModal: false,
      showDeleteSpinner: false,
      showEditSpinner: false,
      createdAt: null,
      createdAtFormatted: null,
      hasError: false,
      errorMessage: '',
      flashMessage: '',
      showFlashMessage: false,
      markdownOption: {
        bold: true,
        italic: true,
        header: true,
        underline: true,
        strikethrough: true,
        mark: true,
        quote: true,
        ol: true,
        ul: true,
        table: true,
        undo: true,
        redo: true,
        trash: true,
        navigation: true,
        subfield: true,
        preview: true,
        help: true,
      }
    }
  },
  mounted() {
    this.fetchTags();
    this.fetchLangs();

    if (this.id) {
      this.loadData();
      this.title = 'Edit Translation';
      this.buttonLabel = 'Save';
    } else {
      this.title = 'Create New Translation';
      this.buttonLabel = 'Create';
    }

    this.triggerFlashMessage();
  },
  watch: {
    tagOptions: function (newVal) {
      let lastUsedTags = this.$store.getters["tag/lastUsedTranslationTagIds"];
      if (!this.id && newVal.length > 0 && this.tags.length === 0 && lastUsedTags.length > 0) {
        this.tags = newVal.filter((tag) => lastUsedTags.includes(tag.id));
      }
    },
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
    fetchLangs() {
      this.$store.dispatch('lang/fetchAll')
          .then((langs) => {
            this.langOptions = langs
            if (langs.length > 0) {
              this.$store.dispatch('profile/fetchProfile')
                  .then((profile) => {
                    if (profile.default_lang && profile.default_lang.id) {
                      this.lang = profile.default_lang;
                    } else {
                      this.lang = langs[0];
                    }
                  })
                  .catch((error) => {
                    this.hasError = true;
                    this.errorMessage = "Can not get user data from server: " + error;
                  });
            }
          })
          .catch(() => {
            this.hasError = true
            this.errorMessage = 'Can not get languages from server :('
          })
    },
    loadData() {
      TranslationService.get(this.id)
          .then((translation) => {
            this.source = translation.source;
            this.transcription = translation.transcription;
            this.target = translation.target;
            this.example = translation.example;
            this.tags = translation.tags;
            this.lang = translation.lang;
            this.createdAt = translation.created_at;
            const dateObj = new Date(translation.created_at);
            this.createdAtFormatted = dateObj.toLocaleString();
          })
          .catch((error) => {
            this.hasError = true;
            this.errorMessage = "Can not get translation data from server: " + error;
          });
    },
    confirmDelete() {
      this.showConfirmationModal = true;
    },
    deleteTranslation() {
      this.showDeleteSpinner = true;
      TranslationService.delete(this.id)
          .then(() => {
            this.$store.dispatch('translation/setEntityStatus', EntityStatusService.deleted());
            router.push({name: 'Home'});
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
    validate() {
      if (!this.lang) {
        this.hasError = true;
        this.errorMessage = "Please select the language";
        return false;
      }

      if (!this.target) {
        this.hasError = true;
        this.errorMessage = "Please fill the target text";
        return false;
      }
      return true;
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
    submitForm() {
      this.hasError = false;
      if (!this.validate()) {
        return;
      }

      this.showEditSpinner = true;
      let method = this.id ? TranslationService.update : TranslationService.create;
      let tagIds = this.tags.map((tag) => tag.id);
      this.$store.dispatch('tag/updateLastUsedTranslationTagIds', tagIds);
      method(new Translation(this.id, this.source, this.transcription, this.target, this.example, tagIds, this.lang.id))
          .then((data) => {
            if (!this.id) {
              this.$store.dispatch('translation/setEntityStatus', EntityStatusService.created());
              router.push("/editTranslation/" + data.id);
            } else {
              this.$store.dispatch('translation/setEntityStatus', EntityStatusService.updated());
              this.triggerFlashMessage();
            }
          })
          .catch((error) => {
            console.log(error);
            this.hasError = true;
            this.errorMessage = 'Can not handle request:' + error;
          })
          .finally(() => {
            this.showEditSpinner = false;
          });
    },
  },
}
</script>