<template>
  <div class="vocabularies-section-wrapper">
    <p>{{type}} vocabulary ({{total}} items total)</p>
    <table>
      <thead>
      <tr>
        <th>{{type}}</th>
        <th>Transcription</th>
        <th>Translation</th>
      </tr>
      </thead>
    </table>
  </div>
</template>

<script>
import TranslationService from '../services/translations.service';

export default {
  name: 'VocabularySection',
  props: {
    type: String,
    requestType: String,
    srcCount: String,
  },
  data() {
    return {
      transactions: null,
      total: 82,
      errors: [],
    }
  },
  created() {
    this.getData()
  },
  methods: {
    getData() {
      TranslationService.getByType(this.requestType)
          .then(response => {
            console.log('Section got the data');
            console.log(response);
            if (response && response.data) {
              this.response = response.data
            }
          })
          .catch(e => {
            this.errors.push(e)
            console.log('error in section');
            console.log(this.errors);
          });
    }
  }
}
</script>
