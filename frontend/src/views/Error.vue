<template>
  <div class="error error-screen" @mousemove="move">
    <div class='container'>
      <span class="error-num">5</span>
      <div ref="eye" v-bind:style="rotateStyle" id="eye-left" class="eye"></div>
      <div v-bind:style="rotateStyle" class="eye"></div>

      <p class="sub-text">Oh eyeballs! Something went wrong. We're <span class="italic">looking</span> to see what happened. Please come back later.</p>
    </div>
  </div>
</template>

<style >
@import '../assets/css/views/Error.css';
</style>

<script>
import axios from 'axios';

export default {
  name: 'error',
  data() {
    return {
      rot: 0,
    }
  },
  computed: {
    rotateStyle() {
      console.log(axios.defaults.baseURL);
      return {
        '-webkit-transform': 'rotate(' + this.rot + 'deg)',
        '-moz-transform': 'rotate(' + this.rot + 'deg)',
        '-ms-transform': 'rotate(' + this.rot + 'deg)',
        'transform': 'rotate(' + this.rot + 'deg)'
      }
    },
  },
  methods: {
    move(event) {
      const eye = this.$refs.eye.getBoundingClientRect();
      const x = (eye.left) + (eye.width / 2);
      const y = (eye.top) + (eye.height / 2);
      const rad = Math.atan2(event.pageX - x, event.pageY - y);
      this.rot= (rad * (180 / Math.PI) * -1) + 180;
    }
  },
}
</script>
