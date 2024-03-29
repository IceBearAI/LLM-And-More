<template>
  <v-dialog v-model="state.visible" v-bind="$attrs" scrollable>
    <v-card ref="contentRef">
      <v-card-title class="d-flex align-center justify-space-between text-subtitle-1 color-font">
        <slot name="title"></slot>
        <v-btn
          v-if="state.showClose"
          @click="onClose"
          :icon="IconX"
          flat
          class="mr-n3"
          title="关闭弹窗"
          width="34px"
          height="34px"
        >
        </v-btn>
      </v-card-title>
      <v-divider></v-divider>
      <v-card-text :style="{ height: state.contentHeight }"><slot name="default"></slot></v-card-text>
      <template v-if="state.showActions">
        <v-divider></v-divider>
        <v-card-actions class="pa-4">
          <v-spacer />
          <template v-if="$slots.buttons">
            <slot name="buttons"></slot>
          </template>
          <template v-else>
            <v-btn size="small" color="secondary" variant="outlined" @click="onClose">取消</v-btn>
            <v-btn size="small" color="primary" variant="flat" @click="onSubmit">确定</v-btn>
          </template>
        </v-card-actions>
      </template>
    </v-card>
  </v-dialog>
</template>
<script setup lang="ts">
import { reactive, ref, watch } from "vue";
import { IconX } from "@tabler/icons-vue";
import _ from "lodash";

/** 外部传入属性 */
type TypePropsOutter = {
  /** 弹窗内容区高度，默认 'auto' */
  contentHeight?: string;
  /** 是否现在弹窗底部默认 按钮 */
  showActions?: boolean;
  /** 是否现在弹窗右上角 关闭按钮 */
  showClose?: boolean;
};
/** 内部属性 */
type TypePropsInner = {
  /** 是否显示弹窗 true显示、false隐藏 */
  visible: boolean;
};

interface IEmits {
  (e: "update:modelValue", val: boolean): void;
  (e: "close"): void;
  (e: "submit"): void;
}

const emits = defineEmits<IEmits>();

const state = reactive<TypePropsInner & TypePropsOutter>({
  visible: false,
  contentHeight: "auto",
  showActions: true,
  showClose: true
});
const contentRef = ref();

const closeDialog = () => {
  state.visible = false;
};

const onClose = () => {
  closeDialog();
  emits("close");
};

const onSubmit = () => {
  emits("submit");
};

watch(
  () => state.visible,
  val => {
    emits("update:modelValue", val);
  }
);

defineExpose({
  show({ contentHeight = "auto", showActions = true, showClose = true } = {} as TypePropsOutter) {
    _.assign(state, {
      contentHeight,
      showActions,
      showClose
    });
    state.visible = true;
  },
  hide() {
    closeDialog();
  },
  getContentEl() {
    return contentRef.value.$el;
  }
});
</script>
<style lang="scss" scoped>
// 修复el-table tooltip位置偏移的问题
:deep(.el-table) {
  transform: scale(1);
}
</style>
