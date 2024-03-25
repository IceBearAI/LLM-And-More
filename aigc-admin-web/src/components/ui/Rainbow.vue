<template>
  <div class="compo-rainbow"></div>
</template>
<script setup lang="ts">
import { reactive, toRefs, ref } from "vue";
const state = reactive({
  style: {},
  formData: {}
});
const { style, formData } = toRefs(state);
</script>
<style lang="scss">
@keyframes jumbo {
  from {
    background-position: 50% 50%, 50% 50%;
  }
  to {
    background-position: 350% 50%, 350% 50%;
  }
}
.compo-rainbow {
  --stripes: repeating-linear-gradient(100deg, #ffffff 0%, #ffffff 7%, transparent 10%, transparent 12%, #ffffff 16%);
  --rainbow: repeating-linear-gradient(100deg, #60a5fa 10%, #e879f9 15%, #60a5fa 20%, #5eead4 25%, #60a5fa 30%);

  position: absolute;
  inset: 0;
  opacity: 0.5;
  pointer-events: none;
  background-image: var(--stripes), var(--rainbow);
  filter: blur(10px) invert(100%);
  background-position: 50% 50%, 50% 50%;
  background-size: 300%, 200%;
  mask-image: radial-gradient(ellipse at 100% 0%, black 40%, transparent 75%);
  &:after {
    position: absolute;
    inset: 0;
    content: "";
    background-image: var(--stripes), var(--rainbow);
    background-attachment: fixed;
    background-size: 200%, 200%;
    background-repeat: repeat;
    mix-blend-mode: difference;
    animation: jumbo 60s linear infinite;
  }
}
</style>
