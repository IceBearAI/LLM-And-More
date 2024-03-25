<template>
  <!-- <div style="color: red">
      isDone:{{ isDone }} <br />
      state.isReady: {{ state.isReady }}
    </div> -->
  <div
    class="compo-textLog d-flex"
    v-loading="!isDone && !state.isReady"
    element-loading-text="日志加载中..."
    element-loading-background="#333"
  >
    <div v-if="isDone && !modelValue" class="mt-10">暂无日志 ...</div>
    <v-virtual-scroll ref="refBox" :items="state.items" class="scrollbar-auto dark" v-else>
      <template v-slot:default="{ item, index }">
        <span class="index pr-2 font-weight-bold">[{{ index }}]</span> {{ item }}
      </template>
    </v-virtual-scroll>
  </div>
</template>
<script setup lang="ts">
import { ref, watch, reactive, computed, nextTick, onMounted } from "vue";
import $ from "jquery";

interface Props {
  /** 日志内容 */
  modelValue: string;
  /** 是否可调整大小，默认true */
  resizeAble: boolean;
  /** 标记日志是否已完成， 已完成状态且modelValue='' 时，显示 空状态 */
  isDone: boolean;
}

const props = withDefaults(defineProps<Props>(), { modelValue: "", resizeAble: true, isDone: true });

const refBox = ref();

const state = reactive({
  items: [],
  isReady: false
});

const render = () => {
  let { modelValue } = props;
  if (modelValue) {
    state.isReady = false;
    let formatedHTML = modelValue.replace(/\r|\n/g, "あ");
    state.items = formatedHTML.split("あ");

    //滑到底部
    nextTick(() => {
      setTimeout(() => {
        refBox.value.scrollToIndex(Number.MAX_SAFE_INTEGER);
      }, 100);
      setTimeout(() => {
        $(refBox.value.$el).scrollTop(Number.MAX_SAFE_INTEGER);
        state.isReady = true;
      }, 300);
    });
  } else {
    if (props.isDone) {
      state.isReady = true;
    }
  }
};

watch(
  () => props.modelValue,
  newVal => {
    render();
  }
);

onMounted(() => {
  render();
});
</script>
<style lang="scss" scoped>
.compo-textLog {
  padding: 24px 16px;
  font-size: 16px;
  width: 100%;
  height: 100%;
  background: #333;
  color: #fff;
  .index {
    color: gold;
    zoom: 0.8;
  }
}
</style>
