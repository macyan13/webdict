<template>
  <div>
    <Tag :id="$route.params.id"/>
  </div>
</template>

<script>
import Tag from '@/components/Tag.vue'

export default {
  name: "EditTag",
  components: {
    Tag
  },
}
</script>

