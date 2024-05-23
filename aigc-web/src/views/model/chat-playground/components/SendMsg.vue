<template>
  <div class="d-flex align-center pa-2">
    <v-btn size="x-small" icon flat @click="handleClear"><IconTrash stroke-width="1.5" :size="20" class="text-error" /></v-btn>
    <v-textarea
      variant="solo"
      hide-details
      v-model="msg"
      color="primary"
      class="shadow-none"
      density="compact"
      placeholder="在此处键入用户查询（按 Shift + Enter 生成新行）"
      no-resize
      @keypress="handleEnter"
      auto-grow
      :rows="1"
      :max-rows="4"
    ></v-textarea>
    <v-btn icon variant="text" @click="handleSend" class="text-medium-emphasis" :disabled="sendBtnDisabled">
      <IconSend :size="20" />
    </v-btn>
  </div>
</template>
<script setup lang="ts">
import { ref, computed } from "vue";
import { IconTrash, IconSend } from "@tabler/icons-vue";

import { useVModel } from "@vueuse/core";

interface IProps {
  modelValue?: string;
}

interface IEmits {
  (e: "update:modelValue", val: string): void;
  (e: "submit"): void;
  (e: "clear"): void;
}

const props = withDefaults(defineProps<IProps>(), {
  modelValue: ""
});
const emit = defineEmits<IEmits>();

const msg = useVModel(props, "modelValue", emit);

const sendBtnDisabled = computed(() => {
  return !msg.value || msg.value.trim() === "";
});

const handleEnter = event => {
  if (event.key === "Enter" && !event.shiftKey) {
    event.preventDefault();
    handleSend();
  }
};
const handleSend = () => {
  if (sendBtnDisabled.value) return;
  emit("submit");
};

const handleClear = () => {
  emit("clear");
};
</script>
<style lang="scss" scoped>
.shadow-none :deep(.v-field--no-label) {
  --v-field-padding-top: -4px;
}
</style>
