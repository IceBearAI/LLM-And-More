<template>
  <Pane ref="refPane" @submit="onSubmit">
    <div class="mx-auto mt-3" style="width: 540px">
      <v-form ref="refForm" class="my-form">
        <v-input :rules="rules.fileId" v-model="formData.fileId" hide-details="auto">
          <template v-if="echoFileName">
            <v-chip closable color="info" @click:close="fileChipClose">{{ echoFileName }}</v-chip>
          </template>
          <template v-else>
            <custom-upload :file-type="['text/csv']" @after-upload="handleAfterUpload">
              <template #trigger>
                <v-btn color="info" variant="outlined">上传文件</v-btn>
              </template>
            </custom-upload>
          </template>
          <template #prepend> <label class="required">样本文件</label></template>
        </v-input>
        <v-text-field
          type="text"
          placeholder="请输入中文、数字、字母、-、_ "
          hide-details="auto"
          clearable
          :rules="rules.name"
          v-model="formData.name"
        >
          <template #prepend><label class="required">别名</label></template>
        </v-text-field>
        <v-textarea v-model.trim="formData.remark" placeholder="备注" clearable>
          <template #prepend> <label>备注</label></template>
        </v-textarea>
      </v-form>
      <div class="d-flex justify-center textPrimary">
        点击 <a class="link mx-1" href="/assets/file/intent-template.csv" download="intent-template.csv">下载</a> 数据集模板
      </div>
    </div>
  </Pane>
</template>
<script setup lang="ts">
import { reactive, ref } from "vue";
import CustomUpload from "@/components/business/CustomUpload.vue";
import { http, validator } from "@/utils";

const emits = defineEmits(["submit"]);

const paneConfig = reactive({
  operateType: "add"
});
const formData = reactive({
  fileId: "",
  name: "",
  remark: ""
});
const echoFileName = ref("");

const refPane = ref();
const refForm = ref();
const rules = reactive({
  fileId: [v => !!v || "请输入样本名称"],
  name: [v => validator.isName({ value: v, required: true, errorValid: "请输入中文、数字、字母、-、_" })]
});

const handleAfterUpload = ({ res }) => {
  if (res) {
    formData.fileId = res.fileId;
    echoFileName.value = res.filename;
  }
};

const fileChipClose = () => {
  echoFileName.value = "";
  formData.fileId = "";
};

const onSubmit = async ({ valid, showLoading }) => {
  if (valid) {
    const requestConfig = {
      url: "",
      method: ""
    };
    if (paneConfig.operateType == "add") {
      requestConfig.url = "/intent/create";
      requestConfig.method = "post";
    }
    const [err, res] = await http[requestConfig.method]({
      showLoading,
      showSuccess: true,
      url: requestConfig.url,
      data: formData
    });

    if (res) {
      refPane.value.hide();
      emits("submit");
    }
  }
};

defineExpose({
  show({ title, operateType }) {
    refPane.value.show({
      title,
      refForm
    });
    paneConfig.operateType = operateType;
    if (paneConfig.operateType === "add") {
      fileChipClose();
      formData.name = "";
      formData.remark = "";
    }
  }
});
</script>
<style lang="scss" scoped>
label {
  width: 130px;
  text-align: right;
  .compo-explain {
    margin-top: 1px;
  }
}
</style>
