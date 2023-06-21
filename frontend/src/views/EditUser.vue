<template>
  <div>
    <user :id="$route.params.id"/>
  </div>
</template>

<script>
import User from '@/components/User.vue'

export default {
  name: "EditUser",
  components: {
    User
  },
}
</script>

