<template>
  <v-alert density="compact" type="info" variant="tonal">
    <template #prepend><IconInfoCircle size="20" class="mr-n2" /></template>
    <div class="text-body-2">
      <slot /></div
  ></v-alert>
</template>
<script setup>
import { reactive, toRefs, ref } from "vue";
import { IconInfoCircle } from "@tabler/icons-vue";

const state = reactive({
  style: {},
  formData: {}
});
const { style, formData } = toRefs(state);
</script>
<style lang="scss"></style>
