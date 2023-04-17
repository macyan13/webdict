<template>
  <div>
    <b-card :title="title">
      <b-form @submit.prevent="submitForm">
        <b-form-group
            id="tag-group"
            label="Tag:"
            label-for="tag-input"
        >
          <div style="display: flex; justify-content: center;">
            <b-form-input
                :required=true
                id="tag-input"
                v-model="newTag"
                :placeholder="tagPlaceholder"
                style="width: 40%;"
            ></b-form-input>
          </div>
        </b-form-group>
        <b-button type="submit" variant="primary">
          {{ buttonLabel }}
        </b-button>
      </b-form>
    </b-card>
    <div v-if="tag" class="mt-3">
      <b-card title="Tag Preview">
        <p>{{ tag }}</p>
      </b-card>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Tag',
  props: {
    tag: String,
    editMode: Boolean,
  },
  data() {
    return {
      title: '',
      buttonLabel: '',
      newTag: '',
      tagPlaceholder: 'Enter a tag...',
    }
  },
  mounted() {
    if (this.editMode) {
      this.title = 'Edit Tag'
      this.buttonLabel = 'Save Changes'
      this.newTag = this.tag
    } else {
      this.title = 'Create New Tag'
      this.buttonLabel = 'Submit'
    }
  },
  methods: {
    submitForm() {
      this.$emit('submit', this.newTag)
    },
  },
}
</script>