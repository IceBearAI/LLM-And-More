<template>
  <div class="file-preview">
    <template v-if="previewCodeSuffix.includes(type)">
      <CodeMirror :model-value="content" :language="fileLanguageMap[type] || type" disabled />
    </template>
    <template v-if="previewMarkdownSuffix.includes(type)">
      <TextComp class="min-h-[230px]" :text="content" msg-type="answer" />
    </template>
  </div>
</template>
<script setup lang="ts">
import { computed } from "vue";
import { previewCodeSuffix, previewMarkdownSuffix } from "../../modelList";
import TextComp from "@/views/model/chat-playground/components/Message/Text.vue";
import { fileLanguageMap } from "@/utils/map";
interface IProps {
  fileType?: string;
  content: string;
  fileInfo: Record<string, any>;
}
const props = withDefaults(defineProps<IProps>(), {
  fileType: "",
  content: "",
  fileInfo: () => ({})
});

const type = computed(() => {
  return props.fileType || props.fileInfo.name.split(".").pop();
});
</script>
<style>
.txt {
  width: 100%;
  min-height: 200px;
  line-height: 16px;
  padding: 10px;
  background-color: #000;
  border-radius: 6px;
}
textarea[readonly] {
  color: #fff;
}
</style>
