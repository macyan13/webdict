<template>
  <div>
    <Translation :id="$route.params.id"/>
  </div>
</template>

<script>
import Translation from '@/components/Translation'

export default {
  name: "EditTag",
  components: {
    Translation
  },
}
</script>

