<template>
  <v-chip :color="state.color" label size="small">
    <component :is="state.icon" size="15" class="mr-1" />{{ state.label }}
  </v-chip>
</template>
<script setup>
import { reactive, toRefs, ref, watchEffect } from "vue";
import { IconCircle, IconClockHour4 } from "@tabler/icons-vue";
const state = reactive({
  label: "",
  color: "",
  icon: ""
});

const { style, formData } = toRefs(state);

const props = defineProps({
  modelValue: {
    type: Boolean,
    default() {
      return false;
    }
  }
});

watchEffect(() => {
  let { modelValue } = props;
  if (modelValue) {
    state.label = "启用";
    state.color = "success";
    state.icon = IconCircle;
  } else {
    state.label = "停用";
    state.color = "default";
    state.icon = IconClockHour4;
  }
});
</script>
<style lang="scss"></style>
