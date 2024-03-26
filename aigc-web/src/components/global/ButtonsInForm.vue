<template>
  <div class="compo-buttonsInForm">
    <slot />
  </div>
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
.compo-buttonsInForm {
  display: flex;
  > * {
    margin-right: 8px;
    &:last-child {
      margin-right: 0px;
    }
  }
}
</style>
