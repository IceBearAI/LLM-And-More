<template>
  <Dialog @update:modelValue="updateModelValue" ref="refDialogLog" max-width="1200px" min-height="400px" class="compo-dialogLog">
    <template #title>
      <slot name="title" v-if="$slots.title"></slot>
      <span v-else>{{ state.title }}</span>
    </template>
    <TextLog v-model="state.content" :isDone="state.isDone" :resizeAble="false"> </TextLog>
  </Dialog>
</template>
<script setup lang="ts">
import { reactive, toRefs, ref } from "vue";
import TextLog from "./TextLog.vue";
import _ from "lodash";
import { useTimeoutFn } from "@vueuse/core";

interface IProps {
  interval?: number;
}

interface IEmits {
  (e: "refresh"): void;
}

const props = withDefaults(defineProps<IProps>(), {
  interval: 0 // 100/10s, 300/30s, 600/1min
});

const emits = defineEmits<IEmits>();

const refDialogLog = ref();
const intervalFn = ref(null);

type TypePropsOutter = {
  /** 日志内容 */
  content: string;
  /** 弹窗标题，默认 '日志' */
  title?: string;
};

type TypePropsInner = {
  /** 标记日志是否已完成 */
  isDone: boolean;
};

const state = reactive<Required<TypePropsOutter & TypePropsInner>>({
  content: "",
  title: "",
  isDone: false
});

const start = () => {
  if (!intervalFn.value) {
    // 初始化 intervalFn
    intervalFn.value = useTimeoutFn(() => {
      triggerRefresh();
    }, props.interval * 1000);
  } else {
    if (!intervalFn.value.isPending) {
      intervalFn.value.start();
    }
  }
};

const triggerRefresh = () => {
  state.isDone = false;
  emits("refresh");
};

const updateModelValue = val => {
  if (!val) {
    if (intervalFn.value) {
      intervalFn.value.stop();
    }
  }
};

defineExpose({
  show({ content = "", title = "日志" } = {} as TypePropsOutter) {
    triggerRefresh();
    refDialogLog.value.show({ showActions: false, contentHeight: "calc(100vh - 400px)" });
    _.assign(state, {
      content,
      title
    });
    if (content) {
      this.setContent(content);
    }
  },
  setContent(content: TypePropsOutter["content"]) {
    state.content = content;
    state.isDone = true;
    if (props.interval !== 0) {
      start();
    }
  }
});
</script>
<style lang="scss">
.compo-dialogLog {
  .v-card-text {
    overflow: hidden !important;
    padding: 0 !important;
    background: #333;
  }
}
</style>
