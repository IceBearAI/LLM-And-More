<template>
  <Pane ref="refPane" @submit="onSubmit">
    <div class="mx-auto mt-3" style="width: 750px">
      <v-form ref="refForm" class="my-form">
        <v-text-field
          type="text"
          placeholder="只允许字母、数字、“-” 、“.”"
          hide-details="auto"
          clearable
          :rules="rules.name"
          v-model="formData.name"
          :disabled="isEdit"
        >
          <template #prepend><label class="required">模版名称</label></template>
        </v-text-field>
        <Select
          placeholder="请选择基础模型"
          :rules="rules.baseModel"
          v-model="formData.baseModel"
          :mapAPI="{
            url: '/models',
            data: { pageSize: -1, isPrivate: true, enabled: true },
            labelField: 'modelName',
            valueField: 'modelName'
          }"
          hide-details="auto"
        >
          <template #prepend>
            <label class="required">基础模型</label>
          </template>
        </Select>
        <Select
          placeholder="请选择模版类型"
          :rules="rules.templateType"
          :mapDictionary="{ code: 'template_type' }"
          v-model="formData.templateType"
        >
          <template #prepend>
            <label class="required">模版类型</label>
          </template>
        </Select>
        <v-text-field
          type="text"
          placeholder="请输入镜像完整地址"
          hide-details="auto"
          clearable
          :rules="rules.trainImage"
          v-model="formData.trainImage"
        >
          <template #prepend
            ><label class="required">训练镜像 <Explain>请提前将Docker镜像制作好并上传到镜像仓库</Explain></label></template
          >
        </v-text-field>
        <v-text-field
          type="text"
          placeholder="请输入模型在容器的绝对路径"
          hide-details="auto"
          clearable
          :rules="rules.baseModelPath"
          v-model="formData.baseModelPath"
        >
          <template #prepend
            ><label class="required">基础模型路径 <Explain>请输入模型所存储的路径</Explain></label></template
          >
        </v-text-field>
        <v-text-field type="text" placeholder="/data/ft-model/" hide-details="auto" clearable v-model="formData.outputDir">
          <template #prepend
            ><label>输出目录 <Explain>模型训练所保存的目录</Explain></label></template
          >
        </v-text-field>
        <v-text-field
          type="text"
          placeholder="/app/train.sh"
          hide-details="auto"
          clearable
          :rules="rules.scriptFile"
          v-model="formData.scriptFile"
        >
          <template #prepend
            ><label class="required">脚本文件 <Explain>训练脚本文件</Explain></label></template
          >
        </v-text-field>
        <v-input hide-details="auto" :rules="rules.content" v-model="formData.content" :center-affix="false">
          <CodeMirror v-model="formData.content" language="shell" placeholder="请输入脚本模版" />
          <template #prepend>
            <label class="required">脚本模版内容 <Explain>脚本模版，通常为启动训练脚本的Shell</Explain></label></template
          >
        </v-input>
        <v-switch v-model="formData.enabled" color="primary" hide-details="auto">
          <template #prepend><label class="required">启用状态</label></template>
        </v-switch>
        <v-textarea v-model.trim="formData.remark" placeholder="请输入备注">
          <template #prepend> <label>备注</label></template>
        </v-textarea>
      </v-form>
    </div>
  </Pane>
</template>
<script setup lang="ts">
import { reactive, ref, computed } from "vue";
import Explain from "@/components/ui/Explain.vue";
import { http } from "@/utils";

interface IFormData {
  name: string;
  baseModel: string | null;
  templateType: string | null;
  scriptFile: string;
  content: string;
  trainImage: string;
  outputDir: string;
  baseModelPath: string;
  enabled: boolean;
  remark: string;
}
const initFormData = {
  name: "",
  baseModel: null,
  templateType: null,
  scriptFile: "",
  content: "",
  trainImage: "",
  outputDir: "/data/ft-model/",
  baseModelPath: "",
  enabled: false,
  remark: ""
};

const emits = defineEmits(["submit"]);

const paneConfig = reactive({
  operateType: "add"
});
const formData = ref<IFormData>({ ...initFormData });
const refPane = ref();
const refForm = ref();
const rules = reactive({
  name: [v => /^[a-zA-Z0-9-.]+$/.test(v) || "只允许字母、数字、“-” 、“.”"],
  baseModel: [v => !!v || "请选择基础模型"],
  templateType: [v => !!v || "请选择模版类型"],
  scriptFile: [v => !!v || "请输入脚本文件"],
  trainImage: [v => !!v || "请输入镜像完整地址"],
  baseModelPath: [v => !!v || "请输入基础模型路径"],
  content: [v => !!v || "请输入脚本模版"]
});

const isEdit = computed(() => {
  return paneConfig.operateType === "edit";
});

const validNumberInput = (value, min, max, errorMessage, reg = false) => {
  if (value !== "" && !errorMessage) {
    if (value < min) {
      return `下限 ${min}`;
    } else if (value > max) {
      return `上限 ${max}`;
    } else if (reg && /^\+?[1-9][0-9]*$/.test(value) != true) {
      return "请输入正整数";
    } else {
      return true;
    }
  } else {
    return errorMessage;
  }
};

const onSubmit = async ({ valid, showLoading }) => {
  if (valid) {
    const requestConfig = {
      url: "",
      method: ""
    };
    if (paneConfig.operateType == "add") {
      requestConfig.url = "/sys/template";
      requestConfig.method = "post";
    } else {
      requestConfig.url = `/sys/template/${formData.value.name}`;
      requestConfig.method = "put";
    }
    const [err, res] = await http[requestConfig.method]({
      showLoading,
      showSuccess: true,
      url: requestConfig.url,
      data: formData.value
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
      formData.value = { ...initFormData };
    } else {
      formData.value = { ...infos };
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
