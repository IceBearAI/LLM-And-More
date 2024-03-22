<template>
  <Pane ref="refPane" class="" @submit="onSubmit">
    <div class="mx-auto mt-10" style="width: 500px">
      <v-form ref="refForm" class="my-form">
        <v-text-field hide-details="auto" v-model="state.modelName" :disabled="true">
          <template #prepend> <label class="required">模型名称</label></template>
        </v-text-field>
        <v-text-field
          type="number"
          hide-details="auto"
          placeholder="请输入使用实例数量"
          :rules="rules.replicas"
          v-model.number="formData.replicas"
          @input="validNumberInput(formData.replicas, 1, 8, '请输入实例数量', true)"
          max="8"
        >
          <template #prepend>
            <label class="required">实例数 <Explain>启动实例数量</Explain></label></template
          >
        </v-text-field>
        <Select
          hide-details="auto"
          placeholder="请选择调度标签"
          :rules="rules.label"
          :mapDictionary="{ code: 'model_deploy_label' }"
          v-model="formData.label"
        >
          <template #prepend>
            <label class="required">调度标签 <Explain>设置实例调度的节点亲和性</Explain></label>
          </template>
        </Select>
        <v-radio-group
          v-model="formData.inferredType"
          inline
          color="primary"
          hide-details="auto"
          @update:modelValue="inferredTypeChange"
        >
          <v-radio label="CPU" value="cpu"></v-radio>
          <v-radio label="GPU" value="gpu"></v-radio>
          <template #prepend>
            <label class="required">推理类型</label>
          </template>
        </v-radio-group>
        <v-text-field
          v-if="formData.inferredType === 'cpu'"
          type="number"
          placeholder="请输入使用CPU数量"
          :rules="rules.cpu"
          v-model.number="formData.cpu"
          @input="validNumberInput(formData.cpu, 1, 80, '请输入使用CPU数量', true)"
          max="80"
          hide-details="auto"
        >
          <template #prepend> <label class="required">CPU数量</label></template>
        </v-text-field>
        <v-text-field
          v-if="formData.inferredType === 'gpu'"
          type="number"
          placeholder="请输入使用GPU数量"
          :rules="rules.gpu"
          v-model.number="formData.gpu"
          @input="validNumberInput(formData.gpu, 1, 8, '请输入使用GPU数量', true)"
          max="8"
          hide-details="auto"
        >
          <template #prepend>
            <label class="required">GPU数量 <Explain>单实例GPU数量</Explain></label></template
          >
        </v-text-field>
        <v-text-field
          v-if="formData.gpu > 1"
          type="number"
          placeholder="请输入GPU内存"
          :rules="rules.maxGpuMemory"
          v-model.number="formData.maxGpuMemory"
          @input="validNumberInput(formData.maxGpuMemory, 1, 80, '请输入GPU内存', true)"
          max="80"
          hide-details="auto"
        >
          <template #prepend>
            <label class="required"
              >GPU内存
              <Explain
                >指定每个 GPU
                用于存储模型权重的最大内存。这允许它为激活分配更多内存，因此您可以使用更长的上下文长度或更大的批量大小</Explain
              ></label
            ></template
          >
          <template #append-inner> GiB </template>
        </v-text-field>
        <Select
          placeholder="请选择k8s集群"
          hide-details="auto"
          :mapDictionary="{ code: 'k8s_cluster' }"
          v-model="formData.k8sCluster"
        >
          <template #prepend>
            <label>k8s集群</label>
          </template>
        </Select>
        <v-checkbox label="开启量化" v-model="state.isQuantify" @change="checkQuantity" hide-details="auto" color="primary">
          <template #prepend>
            <label>量化 <Explain>对模型启动进行量化，节约显存提升推理速度</Explain></label>
          </template>
        </v-checkbox>
        <Select
          v-show="state.isQuantify"
          placeholder="请选择精度"
          :mapDictionary="{ code: 'model_deploy_quantization' }"
          v-model="formData.quantization"
          style="margin-top: -14px"
          hide-details="auto"
        >
          <template #prepend>
            <label> </label>
          </template>
        </Select>

        <v-checkbox style="margin-top: -14px" label="开启VLLM" v-model="formData.vllm" hide-details="auto" color="primary">
          <template #prepend>
            <label>VLLM <Explain>如果想要高吞吐量批量处理，您可以尝试开启，开启后会占用gpu的所有显存容量</Explain></label>
          </template>
        </v-checkbox>
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
const regNum = /^\+?[1-9][0-9]*$/;
const state = reactive<{
  // formData: ItfModel;
  /** 操作类型 add 添加  、 edit 编辑 ，默认add */
  operateType: "";
  id: any;
  [x: string]: any;
}>({
  operateType: "",
  maxTokens: 500000,
  id: "",
  modelName: "",
  isQuantify: false,
  formData: {
    replicas: "",
    label: "",
    cpu: 1,
    gpu: 1,
    maxGpuMemory: "",
    quantization: "",
    vllm: "",
    inferredType: "cpu",
    k8sCluster: null
  }
});
const { formData, id } = toRefs(state);

const emits = defineEmits(["submit"]);

const refPane = ref();
const refForm = ref();
const validNumberInput = (value, min, max, errorMessage, reg = false) => {
  if (value) {
    if (value < min) {
      return `下限 ${min}`;
    } else if (value > max) {
      return `上限 ${max}`;
    } else if (reg && regNum.test(value) != true) {
      return "请输入正整数";
    } else {
      return true;
    }
  } else {
    if (errorMessage) {
      return errorMessage;
    }
  }
  return true;
};
const rules = reactive({
  replicas: [v => validNumberInput(v, 1, 8, "请输入实例数量", true)],
  label: [v => !!v || "请选择调度标签"],
  gpu: [v => validNumberInput(v, 1, 8, "请输入使用GPU数量", true)],
  cpu: [v => validNumberInput(v, 1, 80, "请输入使用CPU数量", true)],
  maxGpuMemory: [v => validNumberInput(v, 1, 80, "请输入GPU内存", true)],
  quantization: [v => !!v || "请选择精度"],
  isQuantify: [
    value => {
      if (value && value != false) {
        return true;
      } else {
        return "请开启量化";
      }
    }
  ],
  vllm: [
    value => {
      if (value && value != false) {
        return true;
      } else {
        return "请开启VLLM";
      }
    }
  ]
});

const checkQuantity = async (options = {}) => {
  state.formData.quantization = state.isQuantify ? "float16" : "";
};

const inferredTypeChange = val => {
  if (val === "cpu") {
    state.formData.gpu = 1;
  } else {
    state.formData.cpu = 1;
  }
};

const onSubmit = async ({ valid, showLoading }) => {
  if (valid) {
    const data = { ...state.formData };
    if (data.inferredType === "cpu") {
      data.gpu = 0;
    } else {
      data.cpu = 0;
    }
    const [err, res] = await http.post({
      ...showLoading,
      showSuccess: true,
      url: `/api/models/${state.id}/deploy`,
      data
    });
    if (res) {
      refPane.value.hide();
      emits("submit");
    }
  }
};

defineExpose({
  show({
    title,
    infos = {
      id: "",
      modelName: ""
    },
    operateType
  }) {
    refPane.value.show({
      title,
      refForm
    });
    state.formData = _.pick(_.cloneDeep(infos), []);
    state.formData.inferredType = "cpu";
    state.formData.cpu = 1;
    state.formData.gpu = 1;
    state.operateType = operateType;
    state.modelName = infos.modelName;
    state.id = infos.id;
  }
});
</script>
<style lang="scss" scoped>
label {
  width: 100px;
  text-align: right;
  .compo-explain {
    margin-top: 1px;
  }
}
</style>
