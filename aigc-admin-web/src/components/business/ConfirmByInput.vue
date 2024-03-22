<template>
  <Confirm ref="refConfirm" class="compo-ConfirmByInput">
    <template #title>
      <slot name="title" v-if="$slots.title"></slot>
      <span v-else class="font-weight-bold">这是进行一项操作时必须了解的重要信息</span>
    </template>
    <template #text>
      <slot name="text"></slot>
      <div class="mt-4">
        <v-text-field
          v-model="state.text"
          density="compact"
          clearable
          hide-details="auto"
          :placeholder="'请输入 “' + state.confirmText + '”'"
        />
      </div>
    </template>
    <template #buttons>
      <v-btn size="small" color="secondary" variant="outlined" @click="onClose">取消</v-btn>
      <AiBtn
        id="btnConfirmByInput"
        size="small"
        color="primary"
        variant="flat"
        @click="onSubmit"
        :disabled="state.confirmText != (state.text || '').trim()"
        >确定</AiBtn
      >
    </template>
  </Confirm>
</template>
<script setup lang="ts">
import { reactive, toRefs, ref } from "vue";
import Confirm from "@/components/ui/Confirm.vue";
import { useSlots } from "vue";

type TypeParams = {
  /** 文本输入，删除关键字 */
  confirmText: string;
  /** 弹窗宽度， 单位 px */
  width?: string;
};

const $slots = useSlots();
const refConfirm = ref();
const state = reactive({
  title: "",
  text: "",
  confirmText: ""
});

const emits = defineEmits(["close", "submit"]);

const onClose = () => {
  refConfirm.value.hide();
  emits("close");
};

const onSubmit = () => {
  emits("submit", { showLoading: "btn#btnConfirmByInput" });
};

defineExpose({
  /** 显示输入确认框 */
  show({ confirmText, width = "400px" }: TypeParams = { confirmText: "" }) {
    state.confirmText = confirmText;
    state.text = "";
    refConfirm.value.show({
      width
    });
  },
  /** 隐藏输入确认框 */
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
