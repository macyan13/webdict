<template>
  <div>
    <b-card :title="title">
      <b-form @submit.prevent="submitForm">
        <b-form-group
            id="name-group"
            label="Name:"
            label-for="lang-input"
        >
          <div style="display: flex; justify-content: center;">
            <b-form-input
                :required=true
                id="lang-input"
                v-model="name"
                placeholder="Enter a language..."
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

      <b-modal v-model="showConfirmationModal" title="Delete Language?" hide-footer hide-backdrop>
        <p>Are you sure you want to delete this language?</p>
        <div class="d-flex justify-content-end">
          <b-button variant="secondary" class="mr-2" @click="showConfirmationModal = false">
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
  </div>
</template>

<script>
import LangService from "@/services/lang.service";
import router from "@/router";
import Lang from "@/models/lang";

export default {
  name: 'Lang',
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
      this.title = 'Edit Language'
      this.buttonLabel = 'Save'
    } else {
      this.title = 'Create New Language'
      this.buttonLabel = 'Create'
    }
  },
  methods: {
    loadData() {
      LangService.get(this.id)
          .then((data) => {
            this.name = data.name;
          })
          .catch((error) => {
            this.hasError = true;
            this.errorMessage = "Can not get lang data from server: " + error;
          });
    },
    confirmDelete() {
      this.showConfirmationModal = true;
    },
    deleteLang() {
      this.showDeleteSpinner = true;
      LangService.delete(this.id)
          .then(() => {
            this.$store.dispatch('lang/clear');
            router.push({name: 'Langs'});
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
      let method = this.id ? LangService.update : LangService.create;
      method(new Lang(this.name, this.id))
          .then(() => {
            this.$store.dispatch('lang/clear');
            router.push({name: 'Langs'});
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