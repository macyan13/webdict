<template>
  <b-card title="Created Users">
    <div class="lang-list" style="display: flex; justify-content: center;">
      <table class="table" id="users" style="width: 75%;">
        <thead>
        <tr>
          <th>Name</th>
          <th>Email</th>
          <th>Role</th>
          <th>Actions</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="user in users" :key="user.id" :id="user.id">
          <td>{{ user.name }}</td>
          <td>{{ user.email }}</td>
          <td>{{ user.role.name }}</td>
          <td>
            <button class="btn btn-sm btn-primary" @click="editUser(user.id)">Edit</button>
            <button class="btn btn-sm btn-danger" @click="confirmDelete(user.id)">Delete</button>
          </td>
        </tr>
        </tbody>
      </table>
    </div>

    <div v-if="hasError" style="color: red;">{{ errorMessage }}</div>

    <div v-if="showLoadSpinner" class="d-flex justify-content-center m-3">
      <b-spinner variant="primary" label="Spinning"></b-spinner>
    </div>
    <b-modal v-model="showConfirmationModal" title="Delete User?" hide-footer hide-backdrop>
      <p>Are you sure you want to delete this user and all user related content?</p>
      <div class="d-flex justify-content-end">
        <b-button variant="secondary" class="mr-2" @click="deleteCancel">
          Cancel
        </b-button>
        <b-button variant="danger" @click="deleteUser">
          Delete
        </b-button>
      </div>
      <div v-if="showDeleteSpinner" class="d-flex justify-content-center mb-3">
        <b-spinner variant="danger" label="Spinning"></b-spinner>
      </div>
    </b-modal>
    <b-modal v-model="showDeleteResults" title="User deletion results" hide-footer hide-backdrop>
      <p>{{deletedCount}} records related to user have been deleted.</p>
      <div class="d-flex justify-content-end">
        <b-button variant="success" class="mr-2" @click="deleteResultsClose" @close="deleteResultsClose">
          Ok
        </b-button>
      </div>
      <div v-if="showDeleteSpinner" class="d-flex justify-content-center mb-3">
        <b-spinner variant="danger" label="Spinning"></b-spinner>
      </div>
    </b-modal>
  </b-card>
</template>
<script>

import UserService from "@/services/user.service";

export default {
  name: 'UserList',
  data() {
    return {
      users: [],
      hasError: false,
      errorMessage: '',
      showConfirmationModal: false,
      idToDelete: null,
      showDeleteSpinner: false,
      showLoadSpinner: true,
      showDeleteResults: false,
      deletedCount: 0,
    }
  },
  mounted() {
    this.fetchUsers();
  },
  methods: {
    editUser(id) {
      this.$router.push(`/editUser/${id}`)
    },
    confirmDelete(id) {
      this.idToDelete = id;
      this.showConfirmationModal = true;
    },
    deleteUser() {
      this.showDeleteSpinner = true;
      UserService.delete(this.idToDelete)
          .then((results) => {
            this.$store.dispatch('user/clear');
            this.fetchUsers();
            this.deletedCount = results.count;
            this.showDeleteResults = true;
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
      this.idToDelete = null;
    },
    deleteResultsClose() {
      this.showDeleteResults = false;
      this.deletedCount = 0;
    },
    fetchUsers() {
      this.showLoadSpinner = true;
      this.$store.dispatch('user/fetchAll')
          .then((users) => this.users = users)
          .catch(() => {
            this.hasError = true
            this.errorMessage = 'Can not get users from server :('
          })
          .finally(() => {
            this.showLoadSpinner = false
          })
    }
  }
};
</script>
