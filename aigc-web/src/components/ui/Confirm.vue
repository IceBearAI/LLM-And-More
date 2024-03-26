<template>
  <v-dialog v-model="state.visible" max-width="800" :width="style.width">
    <v-card class="px-2 pt-3">
      <v-card-title class="text-subtitle-1 color-font" v-if="$slots.title">
        <slot name="title"></slot>
      </v-card-title>
      <v-card-text class="text-body-1 color-font-light"><slot name="text"></slot></v-card-text>

      <v-card-actions>
        <v-spacer />

        <template v-if="$slots.buttons">
          <slot name="buttons"></slot>
        </template>
        <template v-else>
          <v-btn size="small" color="secondary" variant="outlined" @click="onClose">取消</v-btn>
          <v-btn size="small" color="primary" variant="flat" @click="onSubmit">确定</v-btn>
        </template>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
<script setup>
import { reactive, toRefs, ref, useSlots } from "vue";
const state = reactive({
  style: {
    width: "auto"
  },
  formData: {},
  visible: false
});
const { style, formData } = toRefs(state);

const $slots = useSlots();

const emits = defineEmits(["close", "submit"]);

const closePane = () => {
  state.visible = false;
};

const onClose = () => {
  closePane();
  emits("close");
};

const onSubmit = () => {
  emits("submit");
};

defineExpose({
  show({ width = "auto" } = {}) {
    state.style.width = width;
    state.visible = true;
  },
  hide() {
    closePane();
  }
});
</script>
<style lang="scss"></style>
