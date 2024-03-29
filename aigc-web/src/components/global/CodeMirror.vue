<template>
  <div class="codemirror-wrapper">
    <Codemirror v-model="codeValue" :style="codeStyle" :extensions="extensions" :tabSize="2" v-bind="$attrs" />
  </div>
</template>
<script setup lang="ts">
import type { CSSProperties } from "vue";
import { computed } from "vue";
import { Codemirror } from "vue-codemirror";
import { oneDark } from "@codemirror/theme-one-dark";
import { useVModel } from "@vueuse/core";
import { languages, languagesLint } from "./codeMirrorLanguage";
import { linter } from "@codemirror/lint";

interface Props {
  modelValue?: string;
  codeStyle?: CSSProperties; // 代码样式
  dark?: boolean; // 是否暗黑主题
  language?: string; // 语言
}

interface IEmits {
  (e: "update:modelValue", val: string): void;
}
const props = withDefaults(defineProps<Props>(), {
  modelValue: "",
  codeStyle: () => {
    return {};
  },
  dark: true,
  language: "json"
});

const emits = defineEmits<IEmits>();

const codeValue = useVModel(props, "modelValue", emits);

const extensions = computed(() => {
  const result = [];
  if (props.dark) {
    result.push(oneDark);
  }
  if (languagesLint[props.language]) {
    result.push(linter(languagesLint[props.language]));
  }
  if (languages[props.language]) {
    result.push(languages[props.language]());
  }
  return result;
});
</script>
<style lang="scss" scoped>
.codemirror-wrapper {
  width: 100%;
  min-height: 120px;
  max-height: 660px;
}
:deep(.cm-editor) {
  border-radius: 8px;
  outline: none;
  border: 1px solid transparent;
  height: 100%;
  .cm-scroller {
    border-radius: 8px;
  }
}
</style>
