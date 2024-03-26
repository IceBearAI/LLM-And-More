<template>
  <a v-if="$slots.trigger" href="javascript: void(0)" @click="handleSelectFile" v-bind="$attrs">
    <slot :loading="loading" name="trigger"></slot>
  </a>
  <input ref="refUpload" :accept="fileType.join(',')" class="hide" type="file" @change="handleClick" />
</template>
<script setup lang="ts">
import { ref } from "vue";
import { toast } from "vue3-toastify";
import { http } from "@/utils";

const fileTypeMap = {
  "application/pdf": "pdf"
};

interface IProps {
  limit?: number;
  fileSize?: number;
  fileType?: string[];
  isSuffixValid?: boolean;
  autoUpload?: boolean;
}
const props = withDefaults(defineProps<IProps>(), {
  limit: 1,
  fileSize: 2,
  fileType: () => [],
  isSuffixValid: false,
  autoUpload: true
});

interface IEmits {
  (e: "before-upload", value: any): void;
  (e: "after-upload", value: any): void;
  (e: "upload-success", value: any): void;
}

const emits = defineEmits<IEmits>();

const loading = ref(false);
const refUpload = ref();

const handleSelectFile = () => {
  refUpload.value.click();
};

const handleClick = e => {
  const files = e.target.files;
  const rawFile = files[0]; // only use files[0]
  if (!rawFile) return;
  upload(rawFile);
};

const upload = async file => {
  refUpload.value.value = null; // fix can't select the same audio
  // 验证格式
  let type = file.type;
  if (props.isSuffixValid) {
    type = `.${file.name.split(".").pop().toLowerCase()}`;
  }
  if (!props.fileType.includes(type)) {
    toast.warning("上传文件不符合所需的格式！");
    return;
  }
  // 验证大小
  if (file.size / 1024 / 1024 > props.fileSize) {
    toast.warning(`上传文件大小不能超过 ${props.fileSize}M！`);
    return;
  }
  emits("before-upload", file);
  if (!props.autoUpload) return;
  loading.value = true;
  const [err, res] = await http.upload({
    url: "/files",
    data: {
      file
    }
  });
  loading.value = false;
  emits("after-upload", { err, res });
};

defineExpose({
  handleSelectFile
});
</script>
<style lang="scss" scoped>
a {
  color: inherit;
}
.hide {
  display: none;
}
</style>
