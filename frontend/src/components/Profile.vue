<template>
  <div>
    <b-card title="Edit Profile">
      <flash-message v-if="showFlashMessage" :message="flashMessage"/>
      <b-form @submit.prevent="submitForm">
        <b-form-group
            id="name-group"
            label="Name:"
            label-for="user-name"
            :state="name ? true : false"
            invalid-feedback="required"
        >
          <div style="display: flex; justify-content: center;">
            <b-form-input
                :required=true
                id="user-name"
                v-model="name"
                placeholder="Enter a name..."
                class="w-25"
            ></b-form-input>
          </div>
        </b-form-group>
        <b-form-group
            id="email-group"
            label="Email:"
            label-for="user-email"
            :state="name ? true : false"
            invalid-feedback="required"
        >
          <div style="display: flex; justify-content: center;">
            <b-form-input
                type="email"
                :required=true
                id="user-email"
                v-model="email"
                placeholder="Enter a email..."
                class="w-25"
            ></b-form-input>
          </div>
        </b-form-group>
        <b-form-group
            id="current-password-group"
            label="Current password:"
            label-for="user-current-password"
        >
          <div style="display: flex; justify-content: center;">
            <b-form-input
                type="password"
                id="user-current-password"
                v-model="currentPassword"
                placeholder="Enter current password..."
                class="w-25"
            ></b-form-input>
          </div>
        </b-form-group>
        <b-form-group
            id="new-password-group"
            label="New password:"
            label-for="user-new-password"
            :state="currentPassword && !newPassword ? false : true"
            invalid-feedback="required"
        >
          <div style="display: flex; justify-content: center;">
            <b-form-input
                type="password"
                id="user-new-password"
                v-model="newPassword"
                placeholder="Enter new password..."
                class="w-25"
            ></b-form-input>
          </div>
        </b-form-group>
        <b-form-group
            id="new-password-repeat-group"
            label="Repeat new password:"
            label-for="new-password-repeat"
            :state="currentPassword && !newPasswordRepeat ? false : true"
            invalid-feedback="required"
        >
          <div style="display: flex; justify-content: center;">
            <b-form-input
                type="password"
                id="new-password-repeat"
                v-model="newPasswordRepeat"
                placeholder="Repeat new password..."
                class="w-25"
            ></b-form-input>
          </div>
        </b-form-group>
        <b-form-group
            id="lang-group"
            label="Default language:"
        >
          <div style="display: flex; justify-content: center;">
            <VueMultiselect
                :preselect-first="true"
                :options="langOptions"
                v-model="defaultLang"
                :multiple="false"
                label="name"
                track-by="id"
                placeholder="Pick a language"
                style="width: 25%"
            ></VueMultiselect>
          </div>
        </b-form-group>
        <div v-if="!newPasswordsEqual" style="color: red;">New and Repeated passwords are not identical</div>
        <b-button type="submit" variant="primary">
          Save
        </b-button>
        <div v-if="showLoadSpinner" class="d-flex justify-content-center m-3">
          <b-spinner variant="primary" label="Spinning"></b-spinner>
        </div>
      </b-form>

      <div v-if="hasError" style="color: red;">{{errorMessage}}</div>
    </b-card>
  </div>
</template>

<script>
import VueMultiselect from 'vue-multiselect'
import ProfileService from "@/services/profile.service";
import Profile from "@/models/profile";
import FlashMessage from "@/components/FlashMessage.vue";
import EntityStatusService from "@/services/entity-status.service";

export default {
  name: 'Profile',
  components: {
    VueMultiselect,
    FlashMessage
  },
  data() {
    return {
      id: '',
      name: '',
      email: '',
      currentPassword: '',
      newPassword: '',
      newPasswordRepeat: '',
      showLoadSpinner: false,
      hasError: false,
      errorMessage: '',
      langOptions: [],
      defaultLang: null,
      flashMessage: '',
      showFlashMessage: false,
    }
  },
  mounted() {
    this.showLoadSpinner = true;
    this.loadData();
    this.fetchLangs();
    this.showLoadSpinner = false;
  },
  computed: {
    newPasswordsEqual() {
      return this.newPassword === this.newPasswordRepeat;
    },
  },
  methods: {
    loadData() {
      this.$store.dispatch('profile/fetchProfile')
          .then((data) => {
            this.id = data.id;
            this.name = data.name;
            this.email = data.email;
            this.defaultLang = data.default_lang && data.default_lang.id ? data.default_lang : null;
          })
          .catch((error) => {
            this.hasError = true;
            this.errorMessage = "Can not get user data from server: " + error;
          });
    },
    fetchLangs() {
      this.$store.dispatch('lang/fetchAll')
          .then((langs) =>{
            this.langOptions = langs;
          })
          .catch(() => {
            this.hasError = true;
            this.errorMessage = 'Can not get languages from server :(';
          })
    },
    confirmDelete() {
      this.showConfirmationModal = true;
    },
    validate() {
      if (this.currentPassword) {
        if (!this.newPassword) {
          this.hasError = true;
          this.errorMessage = "If you want to change the password, please fill a new password.";
          return false;
        }

        if (this.newPassword !== this.newPasswordRepeat) {
          this.hasError = true;
          this.errorMessage = "New and Repeated passwords are not identical.";
          return false;
        }
      }
      return true;
    },
    triggerFlashMessage() {
      let status = this.$store.getters["profile/entityStatus"];
      if (status === null) {
        return;
      }

      this.flashMessage = EntityStatusService.getMessageByStatus("Profile", status);
      this.showFlashMessage = true;
      this.$store.dispatch('profile/clearEntityStatus');
      setTimeout(() => {
        this.showFlashMessage = false;
      }, 5000);
    },
    submitForm() {
      this.hasError = false;
      if (!this.validate()) {
        return;
      }

      this.showLoadSpinner = true;
      let languageId = this.defaultLang ? this.defaultLang.id : null;
      ProfileService.update(new Profile(this.id, this.name, this.email, this.currentPassword, this.newPassword, languageId))
          .then(() => {
            this.$store.dispatch('profile/setEntityStatus', EntityStatusService.updated());
            this.triggerFlashMessage();
            this.$store.dispatch('profile/clear');
          })
          .catch((error) => {
            this.hasError = true;
            this.errorMessage = error;
          })
          .finally(() => {
            this.showLoadSpinner = false;
          });
    },
  },
}
</script>