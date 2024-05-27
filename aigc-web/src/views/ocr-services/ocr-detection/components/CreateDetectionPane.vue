<template>
  <Pane ref="refPane" @submit="onSubmit">
    <div class="mx-auto mt-3" style="width: 540px">
      <v-form ref="refForm" class="my-form">
        <Select
          placeholder="请选择文件类型"
          :rules="rules.fileType"
          :mapDictionary="{ code: 'ocr_file_type' }"
          v-model="formData.fileType"
          :clearable="false"
        >
          <template #prepend>
            <label class="required">文件类型</label>
          </template>
        </Select>
        <v-input :rules="rules.file" v-model="formData.file" hide-details="auto">
          <template v-if="formData.file">
            <v-chip closable color="info" @click:close="fileChipClose">{{ formData.file.name }}</v-chip>
          </template>
          <template v-else>
            <custom-upload
              :file-type="acceptMap[formData.fileType]"
              @before-upload="handleBeforeUpload"
              :auto-upload="false"
              :file-size="100"
              :disabled="!formData.fileType"
            >
              <template #trigger>
                <v-btn color="info" variant="outlined" :disabled="!formData.fileType">选择文件</v-btn>
              </template>
            </custom-upload>
          </template>
          <template #prepend> <label class="required">文件上传</label></template>
        </v-input>
        <v-switch v-model="formData.viewImage" color="primary" hide-details="auto">
          <template #prepend><label>生成可视化图片</label></template>
        </v-switch>
        <v-text-field
          v-if="formData.fileType === 'pdf'"
          type="number"
          placeholder="请输入pdf页数"
          hide-details="auto"
          v-model.number="formData.pageNum"
        >
          <template #prepend> <label>pdf页数</label></template>
        </v-text-field>
      </v-form>
    </div>
  </Pane>
</template>
<script setup lang="ts">
import { reactive, ref } from "vue";
import CustomUpload from "@/components/business/CustomUpload.vue";
import { http } from "@/utils";

const acceptMap = {
  image: ["image/jpeg", "image/png"],
  pdf: ["application/pdf"]
};

const emits = defineEmits(["submit"]);

const paneConfig = reactive({
  operateType: "add"
});
const formData = reactive({
  fileType: null,
  file: null,
  viewImage: true,
  pageNum: 0
});
const refPane = ref();
const refForm = ref();
const rules = reactive({
  fileType: [v => !!v || "请选择文件类型"],
  file: [v => !!v || "请上传样本文件"]
});

const handleBeforeUpload = (file: File) => {
  formData.file = file;
};

const fileChipClose = () => {
  formData.file = null;
};

const onSubmit = async ({ valid, showLoading }) => {
  if (valid) {
    if (formData.fileType !== "pdf" || !formData.pageNum) {
      formData.pageNum = 0;
    }
    const [err, res] = await http.upload({
      timeout: 30 * 60 * 1000, // 请求超时时间设置为30分钟
      showLoading,
      showSuccess: true,
      url: "/ocr/recognition",
      data: formData
    });

    if (res) {
      refPane.value.hide();
      emits("submit");
    }
  }
};

defineExpose({
  show({ title, operateType, infos }) {
    refPane.value.show({
      title,
      refForm
    });
    paneConfig.operateType = operateType;
    if (paneConfig.operateType === "add") {
      formData.file = null;
      formData.fileType = null;
      formData.viewImage = true;
      formData.pageNum = 0;
    }
  }
});
</script>
<style lang="scss" scoped>
label {
  width: 120px;
  text-align: right;
  .compo-explain {
    margin-top: 1px;
  }
}
</style>
