<template>
  <Pane ref="refPane" @submit="onSubmit">
    <div class="mx-auto mt-3" style="width: 540px">
      <v-form ref="refForm" class="my-form">
        <v-input :rules="rules.file" v-model="formData.file" hide-details="auto">
          <template v-if="formData.file">
            <v-chip closable color="info" @click:close="fileChipClose">{{ formData.file.name }}</v-chip>
          </template>
          <template v-else>
            <custom-upload
              :file-type="['.txt', '.jsonl']"
              isSuffixValid
              @before-upload="handleBeforeUpload"
              :auto-upload="false"
              :file-size="10"
            >
              <template #trigger>
                <v-btn color="info" variant="outlined">选择文件</v-btn>
              </template>
            </custom-upload>
          </template>
          <template #prepend>
            <label class="required">样本文件 <Explain>支持扩展名：.txt、.jsonl</Explain></label></template
          >
        </v-input>
        <v-text-field
          type="text"
          placeholder="请输入样本名称"
          hide-details="auto"
          clearable
          :rules="rules.name"
          v-model="formData.name"
        >
          <template #prepend><label class="required">样本名称</label></template>
        </v-text-field>
        <v-text-field type="text" placeholder="请输入分隔符" hide-details="auto" clearable v-model="formData.splitType">
          <template #prepend><label>分隔符</label></template>
        </v-text-field>
        <v-text-field
          type="number"
          placeholder="请输入分割值"
          :rules="rules.splitMax"
          hide-details="auto"
          v-model.number="formData.splitMax"
        >
          <template #prepend> <label>分割值</label></template>
        </v-text-field>
        <v-textarea v-model.trim="formData.remark" placeholder="备注" clearable>
          <template #prepend> <label>备注</label></template>
        </v-textarea>
      </v-form>
      <div class="d-flex justify-center textPrimary">
        点击 <a class="link mx-1" href="/assets/file/mgr-datasets.txt" download="demo.txt">下载</a> 数据集模板
      </div>
    </div>
  </Pane>
</template>
<script setup lang="ts">
import { reactive, ref } from "vue";
import CustomUpload from "@/components/business/CustomUpload.vue";
import Explain from "@/components/ui/Explain.vue";
import { http, validator } from "@/utils";

const paneConfig = reactive({
  operateType: "add"
});
const formData = reactive({
  file: null,
  name: "",
  splitType: "",
  splitMax: "",
  remark: ""
});

const emits = defineEmits(["submit"]);
const refPane = ref();
const refForm = ref();
const rules = reactive({
  file: [v => !!v || "请上传样本文件"],
  name: [v => !!v || "请输入样本名称"],
  splitMax: [v => validator.validNumberInput(v, 0, 2000, "", true)]
});

const handleBeforeUpload = (file: File) => {
  formData.file = file;
};

const fileChipClose = () => {
  formData.file = null;
};

const onSubmit = async ({ valid, showLoading }) => {
  if (valid) {
    const requestConfig = {
      url: "",
      method: ""
    };
    if (paneConfig.operateType == "add") {
      requestConfig.url = "/mgr/datasets/create";
      requestConfig.method = "upload";
    } else {
      // requestConfig.url = `/sys/template/${formData.name}`;
      // requestConfig.method = "put";
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
      formData.file = null;
      formData.name = "";
      formData.splitType = "";
      formData.splitMax = "";
      formData.remark = "";
    } else {
      console.log("edit");
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
