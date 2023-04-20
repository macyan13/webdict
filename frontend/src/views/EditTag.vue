<template>
  <div>
    <TagForm :id="$route.params.id"/>
  </div>
</template>

<script>
import TagForm from '@/components/Tag.vue'

export default {
  name: "EditTag",
  components: {
    TagForm
  },
}
</script>

