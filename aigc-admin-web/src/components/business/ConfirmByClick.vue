<template>
  <Confirm ref="refConfirm" class="compo-ConfirmByInput">
    <template #title>
      <slot name="title" />
    </template>
    <template #text>
      <slot name="text"></slot>
    </template>
    <template #buttons>
      <v-btn size="small" color="secondary" variant="outlined" @click="onClose">取消</v-btn>
      <AiBtn id="btnConfirmByClick" size="small" color="primary" variant="flat" @click="onSubmit">确定</AiBtn>
    </template>
  </Confirm>
</template>
<script setup lang="ts">
import { reactive, toRefs, ref } from "vue";
import Confirm from "@/components/ui/Confirm.vue";
import { useSlots } from "vue";

type TypeParams = {
  /** 弹窗宽度， 单位 px */
  width?: string;
};

const $slots = useSlots();
const refConfirm = ref();

const emits = defineEmits(["close", "submit"]);

const onClose = () => {
  refConfirm.value.hide();
  emits("close");
};

const onSubmit = () => {
  emits("submit", { showLoading: "btn#btnConfirmByClick" });
};

defineExpose({
  /** 显示点击确认框 */
  show({ width }: TypeParams = {}) {
    refConfirm.value.show({
      width
    });
  },
  /** 隐藏点击确认框 */
  hide() {
    onClose();
  }
});
</script>
<style lang="scss">
.compo-ConfirmByInput {
  .v-card-text {
    padding-top: 5px !important;
  }
}
</style>
