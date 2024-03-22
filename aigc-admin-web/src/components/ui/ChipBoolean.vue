<template>
  <v-chip :color="state.color" label size="small"> {{ state.label }} </v-chip>
</template>
<script setup>
import { reactive, toRefs, ref, watchEffect } from "vue";
const state = reactive({
  label: "",
  color: ""
});

const { style, formData } = toRefs(state);

const props = defineProps({
  modelValue: {
    type: Boolean,
    default() {
      return false;
    }
  },
  mapText: {
    type: Object,
    default() {
      return {
        true: "是",
        false: "否"
      };
    }
  }
});

watchEffect(() => {
  let { modelValue, mapText } = props;
  if (modelValue) {
    state.label = mapText[true];
    state.color = "success";
  } else {
    state.label = mapText[false];
    state.color = "default";
  }
});
</script>
<style lang="scss"></style>
