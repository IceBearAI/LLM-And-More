<template>
  <Pane ref="refPane" class="" @submit="onSubmit">
    <div class="mx-auto mt-10" style="width: 500px">
      <!-- formData :{{ formData }} -->
      <v-form ref="refForm" class="my-form">
        <v-text-field
          density="compact"
          variant="outlined"
          type="text"
          placeholder="只允许字母、数字、“-” 、“.” 和 “:” "
          hide-details="auto"
          clearable
          :rules="rules.modelName"
          v-model="formData.modelName"
          :disabled="state.disabledField"
        >
          <template #prepend> <label class="required">模型名称</label></template>
        </v-text-field>
        <v-text-field
          density="compact"
          variant="outlined"
          type="number"
          :placeholder="'上限 ' + format.commaString(state.maxTokens)"
          hide-details="auto"
          :rules="rules.maxTokens"
          v-model.number="formData.maxTokens"
        >
          <template #prepend> <label class="required">最长上下文</label></template>
        </v-text-field>
        <Select
          placeholder="请选择供应商"
          :rules="rules.providerName"
          :mapDictionary="{ code: 'model_provider_name' }"
          v-model="formData.providerName"
          :disabled="state.disabledField"
        >
          <template #prepend>
            <label class="required">供应 <Explain>供应商指的是外部服务提供，自己有服务请选择Local</Explain></label>
          </template>
        </Select>
        <Select
          placeholder="请选择模型类型"
          :rules="rules.modelType"
          :mapDictionary="{ code: 'model_type' }"
          v-model="formData.modelType"
          :disabled="state.disabledField"
          @update:modelValue="onModelTypeChange"
        >
          <template #prepend>
            <label class="required">模型类型</label>
          </template>
        </Select>
        <!-- <v-switch v-model="formData.isPrivate" :disabled="state.disabledField" color="primary" hide-details="auto">
          <template #prepend>
            <label>私有 <Explain>私有模型，表示部署在本地服务器，如果是请打开</Explain></label>
          </template>
        </v-switch> -->
        <v-switch v-model="formData.enabled" color="primary" hide-details="auto">
          <template #prepend>
            <label>启用 <Explain>模型是否可被使用，如果关闭则无法使用该模型</Explain></label></template
          >
        </v-switch>

        <v-textarea v-model.trim="formData.remark" placeholder="请输入">
          <template #prepend> <label>备注</label></template>
        </v-textarea>
        <v-input>
          <v-expansion-panels>
            <v-expansion-panel elevation="10">
              <v-expansion-panel-title class="text-h6">高级</v-expansion-panel-title>
              <v-expansion-panel-text class="mt-4" eager>
                <Select
                  placeholder="请选择基础模型"
                  v-model="formData.baseModelName"
                  :mapAPI="{
                    url: '/channels/models',
                    data: {
                      providerName: 'LocalAI',
                      modelType: 'text-generation',
                      baseModelName: 'null'
                    },
                    labelField: 'modelName',
                    valueField: 'modelName'
                  }"
                  :hide-details="false"
                >
                  <template #prepend>
                    <label>基础模型</label>
                  </template>
                </Select>
                <v-switch v-model="formData.isFineTuning" :disabled="state.disabledField" color="primary">
                  <template #prepend>
                    <label>微调 <Explain>是否是微调模型</Explain></label>
                  </template>
                </v-switch>
                <v-text-field
                  density="compact"
                  variant="outlined"
                  type="number"
                  placeholder="请输入参数量"
                  v-model.number="formData.parameters"
                  style="width: 60%"
                  :disabled="state.disabledField"
                >
                  <template #prepend>
                    <label>参数量(B) <Explain>单位是B，代理10亿</Explain></label>
                  </template>
                </v-text-field>
                <template v-if="formData.modelType === 'digitalhuman'">
                  <v-text-field
                    type="number"
                    placeholder="请输入并行/实例数"
                    :rules="rules.replicas"
                    v-model.number="formData.replicas"
                  >
                    <template #prepend> <label>并行/实例数</label></template>
                  </v-text-field>
                  <Select
                    placeholder="请选择调度标签"
                    :hide-details="false"
                    :mapDictionary="{ code: 'model_deploy_label' }"
                    v-model="formData.label"
                  >
                    <template #prepend>
                      <label>调度标签</label>
                    </template>
                  </Select>
                  <Select
                    placeholder="请选择k8s集群"
                    :hide-details="false"
                    :mapDictionary="{ code: 'k8s_cluster' }"
                    v-model="formData.k8sCluster"
                  >
                    <template #prepend>
                      <label>k8s集群</label>
                    </template>
                  </Select>
                  <v-radio-group v-model="formData.inferredType" inline color="primary">
                    <v-radio label="CPU" value="cpu"></v-radio>
                    <v-radio label="GPU" value="gpu"></v-radio>
                    <template #prepend>
                      <label>推理类型</label>
                    </template>
                  </v-radio-group>
                  <v-text-field type="number" placeholder="请输入使用CPU数量" :rules="rules.cpu" v-model.number="formData.cpu">
                    <template #prepend> <label>CPU数量</label></template>
                  </v-text-field>
                  <v-text-field type="number" placeholder="请输入使用GPU数量" :rules="rules.gpu" v-model.number="formData.gpu">
                    <template #prepend> <label>GPU数量</label></template>
                  </v-text-field>
                  <v-text-field type="number" placeholder="请输入内存" :rules="rules.memory" v-model.number="formData.memory">
                    <template #prepend> <label>内存</label></template>
                    <template #append-inner>G</template>
                  </v-text-field>
                </template>
              </v-expansion-panel-text>
            </v-expansion-panel>
          </v-expansion-panels>
        </v-input>
      </v-form>
    </div>
  </Pane>
</template>
<script setup lang="ts">
import { reactive, toRefs, ref, nextTick } from "vue";
import _ from "lodash";
import Explain from "@/components/ui/Explain.vue";
import { http, validator, format } from "@/utils";
import { ItfModel } from "../types/modelList.ts";

const state = reactive<{
  formData: ItfModel;
  /** 操作类型 add 添加  、 edit 编辑 ，默认add */
  operateType: "add" | "edit";
  [x: string]: any;
}>({
  operateType: "add",
  disabledField: false,
  maxTokens: 500000,
  formData: {
    id: "",
    modelName: "",
    maxTokens: 0,
    // isPrivate: false,
    isFineTuning: false,
    enabled: true,
    remark: "",
    parameters: 0,
    providerName: "",
    modelType: null,
    baseModelName: null,
    replicas: 0,
    label: null,
    k8sCluster: null,
    inferredType: "",
    cpu: 0,
    gpu: 0,
    memory: 0
  }
});
const { formData } = toRefs(state);

const emits = defineEmits(["submit"]);

const refPane = ref();
const refForm = ref();
const rules = reactive({
  modelName: [
    value => {
      return validator.isModelName({ value, required: true });
    }
  ],
  maxTokens: [
    value => {
      if (value !== "") {
        if (value < 1) {
          return "下限 1";
        } else if (value > state.maxTokens) {
          return "上限 " + format.commaString(state.maxTokens);
        } else {
          return true;
        }
      } else {
        return "请输入模型最长上下文";
      }
    }
  ],
  providerName: [
    value => {
      if (value && value.length > 0) {
        return true;
      } else {
        return "请选择供应";
      }
    }
  ],
  modelType: [v => !!v || "请选择模型类型"],
  replicas: [v => validator.validNumberInput(v, 0, 50, "", true)],
  cpu: [v => validator.validNumberInput(v, 0, 100, "", true)],
  gpu: [v => validator.validNumberInput(v, 0, 20, "", true)],
  memory: [v => validator.validNumberInput(v, 0, 100, "", true)]
});

const onModelTypeChange = val => {
  if (val !== "digitalhuman") {
    state.formData.replicas = 0;
    state.formData.label = null;
    state.formData.k8sCluster = null;
    state.formData.inferredType = "";
    state.formData.cpu = 0;
    state.formData.gpu = 0;
    state.formData.memory = 0;
  }
};

const doAdd = async (options = {}) => {
  const [err, res] = await http.post({
    ...options,
    showSuccess: true,
    url: `/models`,
    data: {
      ...state.formData
    }
  });
  if (res) {
    refPane.value.hide();
    emits("submit");
  }
};

const doEdit = async (options = {}) => {
  const [err, res] = await http.put({
    ...options,
    showSuccess: true,
    url: `/models/${state.formData.id}`,
    data: {
      ...state.formData
    }
  });
  if (res) {
    refPane.value.hide();
    emits("submit");
  }
};

const onSubmit = ({ valid, showLoading }) => {
  if (valid) {
    if (state.operateType == "add") {
      doAdd({ showLoading });
    } else {
      doEdit({ showLoading });
    }
  }
};

defineExpose({
  show({
    title,
    infos = {
      id: "",
      modelName: "",
      maxTokens: 0,
      // isPrivate: false,
      isFineTuning: false,
      enabled: true,
      remark: "",
      providerName: "",
      parameters: 0,
      modelType: null,
      baseModelName: null,
      replicas: 0,
      label: null,
      k8sCluster: null,
      inferredType: "",
      cpu: 0,
      gpu: 0,
      memory: 0
    },
    operateType
  }) {
    refPane.value.show({
      title,
      refForm
    });
    state.formData = _.pick(_.cloneDeep(infos), [
      "id",
      "modelName",
      "maxTokens",
      // "isPrivate",
      "isFineTuning",
      "enabled",
      "remark",
      "providerName",
      "parameters",
      "modelType",
      "baseModelName",
      "replicas",
      "label",
      "k8sCluster",
      "inferredType",
      "cpu",
      "gpu",
      "memory"
    ]);
    state.operateType = operateType;
    if (operateType == "add") {
      //添加
      state.disabledField = false;
    } else {
      //编辑
      state.disabledField = true;
    }
  }
});
</script>
<style lang="scss" scoped>
label {
  width: 90px;
  text-align: right;
  .compo-explain {
    margin-top: 1px;
  }
}
</style>
