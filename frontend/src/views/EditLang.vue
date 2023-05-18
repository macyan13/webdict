<template>
  <div>
    <Lang :id="$route.params.id"/>
  </div>
</template>

<script>
import Lang from '@/components/Lang.vue'

export default {
  name: "EditLang",
  components: {
    Lang
  },
}
</script>

