<template>
  <Dialog ref="refDialogLog" max-width="1200px" min-height="400px" class="compo-dialogLog">
    <template #title> {{ state.title }} </template>
    <TextLog v-model="state.content" :isDone="state.isDone" :resizeAble="false"> </TextLog>
  </Dialog>
</template>
<script setup lang="ts">
import { reactive, toRefs, ref } from "vue";
import TextLog from "./TextLog.vue";
import _ from "lodash";
import { ElLoading } from "element-plus";
import $ from "jquery";

const refDialogLog = ref();

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

defineExpose({
  show({ content = "", title = "日志" } = {} as TypePropsOutter) {
    state.isDone = false;
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
