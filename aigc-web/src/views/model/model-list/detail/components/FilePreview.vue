<template>
  <div class="file-preview min-h-[230px]">
    <template v-if="previewCodeSuffix.includes(type)">
      <CodeMirror :model-value="content" :language="type" disabled />
    </template>
    <template v-if="previewTextSuffix.includes(type)">
      <TextComp class="h-100" :text="content" msg-type="answer" />
    </template>
  </div>
</template>
<script setup lang="ts">
import { computed } from "vue";
import { previewCodeSuffix, previewTextSuffix } from "../../modelList";
import TextComp from "@/views/model/chat-playground/components/Message/Text.vue";
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
