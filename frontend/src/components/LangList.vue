<template>
  <b-card :title="title">
    <flash-message v-if="showFlashMessage" :message="flashMessage"/>
    <div class="lang-list" style="display: flex; justify-content: center;">
      <b-list-group style="width: 40%;">
        <b-list-group-item v-for="lang in langs" :key="lang.id">
          <div style="display: flex; justify-content: space-between;">
            <span>{{ lang.name }}</span>
            <span>
              <b-button variant="primary" @click="editLang(lang.id)">Edit</b-button>
              <b-button variant="danger" @click="confirmDelete(lang.id)">Delete</b-button>
            </span>
          </div>
        </b-list-group-item>
      </b-list-group>
    </div>

    <div v-if="hasError" style="color: red;">{{errorMessage}}</div>

    <div v-if="showLoadSpinner" class="d-flex justify-content-center m-3">
      <b-spinner variant="primary" label="Spinning"></b-spinner>
    </div>
    <b-modal v-model="showConfirmationModal" title="Delete Language?" hide-footer hide-backdrop>
      <p>Are you sure you want to delete this language?</p>
      <div class="d-flex justify-content-end">
        <b-button variant="secondary" class="mr-2" @click="deleteCancel">
          Cancel
        </b-button>
        <b-button variant="danger" @click="deleteLang">
          Delete
        </b-button>
      </div>
      <div v-if="showDeleteSpinner" class="d-flex justify-content-center mb-3">
        <b-spinner variant="danger" label="Spinning"></b-spinner>
      </div>
    </b-modal>
  </b-card>
</template>
<script>

import LangService from "@/services/lang.service";
import FlashMessage from "@/components/FlashMessage.vue";
import EntityStatusService from "@/services/entity-status.service";

export default {
  name: 'LangList',
  components: {
    FlashMessage
  },
  props: {
    title: {
      type: String,
      default() {
        return "Your Languages";
      }
    },
  },
  data() {
    return {
      langs: [],
      hasError: false,
      errorMessage: '',
      showConfirmationModal: false,
      idToDelete: null,
      showDeleteSpinner: false,
      showLoadSpinner: true,
      flashMessage: '',
      showFlashMessage: false,
    }
  },
  mounted() {
    this.fetchLangs();
    this.triggerFlashMessage();
  },
  methods: {
    editLang(id) {
      this.$router.push(`/editLang/${id}`)
    },
    confirmDelete(id) {
      this.idToDelete = id;
      this.showConfirmationModal = true;
    },
    deleteLang() {
      this.showDeleteSpinner = true;
      LangService.delete(this.idToDelete)
          .then(() => {
            this.$store.dispatch('lang/setEntityStatus', EntityStatusService.deleted())
            this.triggerFlashMessage();
            this.$store.dispatch('lang/clear');
            this.fetchLangs();
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
    deleteCancel() {
      this.showConfirmationModal = false;
      this.idToDelete = null
    },
    triggerFlashMessage() {
      let status = this.$store.getters["lang/entityStatus"];
      if (status === null) {
        return;
      }

      this.flashMessage = EntityStatusService.getMessageByStatus("Language", status);
      this.showFlashMessage = true;
      this.$store.dispatch('lang/clearEntityStatus');
      setTimeout(() => {
        this.showFlashMessage = false;
      }, 5000);
    },
    fetchLangs() {
      this.showLoadSpinner = true;
      this.$store.dispatch('lang/fetchAll')
          .then((langs) => this.langs = langs)
          .catch(() => {
            this.hasError = true
            this.errorMessage = 'Can not get langs from server :('
          })
          .finally(() => {
            this.showLoadSpinner = false;
          })
    }
  }
};
</script>
