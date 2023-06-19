<template>
  <div>
    <b-card :title="title">
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
            id="password-group"
            label="Password:"
            label-for="password"
        >
          <div style="display: flex; justify-content: center;">
            <b-form-input
                type="password"
                id="password"
                v-model="password"
                placeholder="Enter password..."
                class="w-25"
            ></b-form-input>
          </div>
        </b-form-group>

        <b-form-group
            id="role-group"
            label="Role:"
            label-for="lang-input"
            :state="role ? true : false"
            invalid-feedback="required"
        >
          <div style="display: flex; justify-content: center;">
            <VueMultiselect
                :allow-empty="false"
                :options="roleOptions"
                v-model="role"
                :multiple="false"
                label="name"
                track-by="id"
                deselectLabel=""
                placeholder="Pick a role"
                style="width: 25%"
            ></VueMultiselect>
          </div>
        </b-form-group>

        <b-button type="submit" variant="primary">
          {{ buttonLabel }}
        </b-button>

<!--        <b-button v-if="id" variant="danger" @click="confirmDelete">-->
<!--          Delete-->
<!--        </b-button>-->

        <div v-if="showEditSpinner" class="d-flex justify-content-center m-3">
          <b-spinner variant="primary" label="Spinning"></b-spinner>
        </div>
      </b-form>

      <div v-if="hasError" style="color: red;">{{errorMessage}}</div>

<!--      <b-modal v-model="showConfirmationModal" title="Delete Tag?" hide-footer hide-backdrop>-->
<!--        <p>Are you sure you want to delete this tag?</p>-->
<!--        <div class="d-flex justify-content-end">-->
<!--          <b-button variant="secondary" class="mr-2" @click="showConfirmationModal = false">-->
<!--            Cancel-->
<!--          </b-button>-->
<!--          <b-button variant="danger" @click="deleteTag">-->
<!--            Delete-->
<!--          </b-button>-->
<!--        </div>-->
<!--        <div v-if="showDeleteSpinner" class="d-flex justify-content-center mb-3">-->
<!--          <b-spinner variant="danger" label="Spinning"></b-spinner>-->
<!--        </div>-->
<!--      </b-modal>-->

    </b-card>
  </div>
</template>

<script>
import UserService from "@/services/user.service";
import User from "@/models/user";
import router from "@/router";
import VueMultiselect from 'vue-multiselect'

export default {
  name: 'User',
  components: {
    VueMultiselect,
  },
  props: {
    id: {
      type: String,
      default: null,
    },
  },
  data() {
    return {
      title: '',
      name: '',
      email: '',
      password: '',
      buttonLabel: '',
      roleOptions: [],
      role: null,
      showConfirmationModal: false,
      showDeleteSpinner: false,
      showEditSpinner: false,
      hasError: false,
      errorMessage: '',
    }
  },
  mounted() {
    this.fetchRoles();

    if (this.id) {
      this.loadData();
      this.title = 'Edit User'
      this.buttonLabel = 'Save'
    } else {
      this.title = 'Create New User'
      this.buttonLabel = 'Create'
    }
  },
  methods: {
    loadData() {
      UserService.get(this.id)
          .then((data) => {
            this.tag = data.tag;
          })
          .catch((error) => {
            this.hasError = true;
            this.errorMessage = "Can not get tag data from server: " + error;
          });
    },
    confirmDelete() {
      this.showConfirmationModal = true;
    },
    fetchRoles() {
      this.$store.dispatch('role/fetchAll')
          .then((roles) => this.roleOptions = roles)
          .catch(() => {
            this.hasError = true
            this.errorMessage = 'Can not get roles from server :('
          })
    },
    // deleteTag() {
    //   this.showDeleteSpinner = true;
    //   TagService.delete(this.id)
    //       .then(() => {
    //         this.$store.dispatch('tag/clear');
    //         router.push({name: 'Tags'});
    //       })
    //       .catch((error) => {
    //         this.hasError = true;
    //         this.errorMessage = error;
    //       })
    //       .finally(() => {
    //         this.showDeleteSpinner = false;
    //         this.showConfirmationModal = false;
    //       });
    // },
    submitForm() {
      this.showEditSpinner = true;
      let method = this.id ? UserService.update : UserService.create;
      method(new User(this.id, this.name, this.email, this.password, this.role.id))
          .then(() => {
            // todo: push to user list
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