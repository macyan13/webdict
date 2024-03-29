<template>
  <b-card :title="title">
    <flash-message v-if="showFlashMessage" :message="flashMessage"/>
    <div class="tag-list" style="display: flex; justify-content: center;">
      <b-list-group style="width: 40%;">
        <b-list-group-item v-for="tag in tags" :key="tag.id">
          <div style="display: flex; justify-content: space-between;">
            <span>{{ tag.name }}</span>
            <span>
              <b-button variant="primary" @click="editTag(tag.id)">Edit</b-button>
              <b-button variant="danger" @click="confirmDelete(tag.id)">Delete</b-button>
            </span>
          </div>
        </b-list-group-item>
      </b-list-group>
    </div>

    <div v-if="hasError" style="color: red;">{{errorMessage}}</div>

    <div v-if="showLoadSpinner" class="d-flex justify-content-center m-3">
      <b-spinner variant="primary" label="Spinning"></b-spinner>
    </div>
    <b-modal v-model="showConfirmationModal" title="Delete Tag?" hide-footer hide-backdrop>
      <p>Are you sure you want to delete this tag?</p>
      <div class="d-flex justify-content-end">
        <b-button variant="secondary" class="mr-2" @click="deleteCancel">
          Cancel
        </b-button>
        <b-button variant="danger" @click="deleteTag">
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

import TagService from "@/services/tag.service";
import FlashMessage from "@/components/FlashMessage.vue";
import EntityStatusService from "@/services/entity-status.service";

export default {
  name: 'TagList',
  components: {
    FlashMessage
  },
  props: {
    title: {
      type: String,
      default() {
        return "Your tags";
      }
    },
  },
  data() {
    return {
      tags: [],
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
    this.fetchTags();
    this.triggerFlashMessage();
  },
  methods: {
    editTag(id) {
      this.$router.push(`/editTag/${id}`)
    },
    confirmDelete(id) {
      this.idToDelete = id;
      this.showConfirmationModal = true;
    },
    deleteTag() {
      this.showDeleteSpinner = true;
      TagService.delete(this.idToDelete)
          .then(() => {
            this.$store.dispatch('tag/setEntityStatus', EntityStatusService.deleted())
            this.triggerFlashMessage();
            this.$store.dispatch('tag/clear');
            this.$store.dispatch('translationHome/resetTags');
            this.fetchTags();
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
      let status = this.$store.getters["tag/entityStatus"];
      if (status === null) {
        return;
      }

      this.flashMessage = EntityStatusService.getMessageByStatus("Tag", status);
      this.showFlashMessage = true;
      this.$store.dispatch('tag/clearEntityStatus');
      setTimeout(() => {
        this.showFlashMessage = false;
      }, 5000);
    },
    fetchTags() {
      this.showLoadSpinner = true;
      this.$store.dispatch('tag/fetchAll')
          .then((tags) => this.tags = tags)
          .catch(() => {
            this.hasError = true
            this.errorMessage = 'Can not get tags from server :('
          })
          .finally(() => {
            this.showLoadSpinner = false;
          });
    }
  }
};
</script>
