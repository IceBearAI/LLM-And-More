<template>
  <Pane ref="refPane" @submit="onSubmit">
    <div class="mx-auto mt-3" style="width: 540px">
      <v-form ref="refForm" class="my-form">
        <Select
          placeholder="请选择评测指标"
          :rules="rules.evalTargetType"
          :mapDictionary="{ code: 'model_evaluate_target_type' }"
          v-model="formData.evalTargetType"
        >
          <template #prepend>
            <label class="required">评测指标</label>
          </template>
        </Select>
        <v-input v-if="formData.evalTargetType !== 'five'" :rules="rules.fileId" v-model="formData.fileId" hide-details="auto">
          <template v-if="echoFileName">
            <v-chip closable color="info" @click:close="fileChipClose">{{ echoFileName }}</v-chip>
          </template>
          <template v-else>
            <custom-upload :file-type="['.txt', '.json', '.jsonl']" isSuffixValid @after-upload="handleAfterUpload">
              <template #trigger>
                <v-btn color="info" variant="outlined">上传文件</v-btn>
              </template>
            </custom-upload>
          </template>
          <template #prepend> <label class="required">待评测数据集</label></template>
        </v-input>
        <v-text-field
          type="number"
          placeholder="请输入最大输出序列长度"
          hide-details="auto"
          :rules="rules.maxLength"
          v-model.number="formData.maxLength"
        >
          <template #prepend> <label class="required">最大输出序列长度</label></template>
        </v-text-field>
        <v-text-field
          type="number"
          placeholder="请输入单卡Batch大小"
          hide-details="auto"
          :rules="rules.batchSize"
          v-model.number="formData.batchSize"
        >
          <template #prepend> <label class="required">单卡Batch大小</label></template>
        </v-text-field>
        <v-textarea v-model.trim="formData.remark" placeholder="请输入你微调模型的备注" clearable>
          <template #prepend> <label>备注</label></template>
        </v-textarea>
        <v-input>
          <v-expansion-panels>
            <v-expansion-panel elevation="10">
              <v-expansion-panel-title class="text-h6">高级</v-expansion-panel-title>
              <v-expansion-panel-text class="mt-4" eager>
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
                <v-radio-group v-model="formData.inferredType" inline color="primary" @update:modelValue="inferredTypeChange">
                  <v-radio label="CPU" value="cpu"></v-radio>
                  <v-radio label="GPU" value="gpu"></v-radio>
                  <template #prepend>
                    <label>推理类型</label>
                  </template>
                </v-radio-group>
                <v-text-field
                  v-if="formData.inferredType === 'cpu'"
                  type="number"
                  placeholder="请输入使用CPU数量"
                  :rules="rules.cpu"
                  v-model.number="formData.cpu"
                >
                  <template #prepend> <label class="required">CPU数量</label></template>
                </v-text-field>
                <v-text-field
                  v-if="formData.inferredType === 'gpu'"
                  type="number"
                  placeholder="请输入使用GPU数量"
                  :rules="rules.gpu"
                  v-model.number="formData.gpu"
                >
                  <template #prepend> <label class="required">GPU数量</label></template>
                </v-text-field>
                <v-text-field
                  v-if="formData.gpu > 1"
                  type="number"
                  placeholder="请输入GPU内存"
                  :rules="rules.maxGpuMemory"
                  v-model.number="formData.maxGpuMemory"
                  max="80"
                  hide-details="auto"
                >
                  <template #prepend>
                    <label
                      >GPU内存
                      <Explain
                        >指定每个 GPU
                        用于存储模型权重的最大内存。这允许它为激活分配更多内存，因此您可以使用更长的上下文长度或更大的批量大小</Explain
                      ></label
                    ></template
                  >
                  <template #append-inner> GiB </template>
                </v-text-field>
              </v-expansion-panel-text>
            </v-expansion-panel>
          </v-expansion-panels>
        </v-input>
        <!-- <v-input>
          <span class="text-primary">预计评估时间：1小时28分钟</span>
          <template #prepend><label></label></template>
        </v-input> -->
      </v-form>
    </div>
  </Pane>
</template>
<script setup lang="ts">
import { reactive, ref } from "vue";
import CustomUpload from "@/components/business/CustomUpload.vue";
import { http, validator } from "@/utils";
import { useRoute } from "vue-router";
import Explain from "@/components/ui/Explain.vue";

const route = useRoute();

const paneConfig = reactive({
  operateType: "add"
});
const echoFileName = ref("");
const formData = reactive({
  modelId: Number(route.query.jobId),
  fileId: "",
  evalTargetType: null,
  maxLength: 512,
  batchSize: 32,
  remark: "",
  label: null,
  k8sCluster: null,
  inferredType: "",
  cpu: 0,
  gpu: 0,
  maxGpuMemory: ""
});

const emits = defineEmits(["submit"]);
const refPane = ref();
const refForm = ref();
const rules = reactive({
  fileId: [v => !!v || "请选择微调文件"],
  evalTargetType: [v => !!v || "请选择评测指标"],
  maxLength: [v => !!v || "请输入最大输出序列长度"],
  batchSize: [v => !!v || "请输入单卡Batch大小"],
  cpu: [v => validator.validNumberInput(v, 0, 100, "请输入使用CPU数量", true)],
  gpu: [v => validator.validNumberInput(v, 0, 20, "请输入使用GPU数量", true)],
  maxGpuMemory: [v => validator.validNumberInput(v, 1, 80, "", true)]
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

const inferredTypeChange = () => {
  if (formData.cpu !== 1) {
    formData.cpu = 1;
  }
  if (formData.gpu !== 1) {
    formData.gpu = 1;
  }
};

const onSubmit = async ({ valid, errors, showLoading }) => {
  if (valid) {
    const data = { ...formData };
    if (formData.inferredType === "cpu") {
      data.gpu = 0;
    } else {
      data.cpu = 0;
    }
    if (data.gpu < 2 || !data.maxGpuMemory) {
      delete data.maxGpuMemory;
    }
    const [err, res] = await http.post({
      showLoading,
      showSuccess: true,
      url: "/evaluate/create",
      data
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
    if (operateType == "add") {
      fileChipClose(); // init the file
      formData.evalTargetType = null;
      formData.maxLength = 512;
      formData.batchSize = 32;
      formData.remark = "";
      formData.label = null;
      formData.k8sCluster = null;
      formData.inferredType = "";
      formData.cpu = 0;
      formData.gpu = 0;
      formData.maxGpuMemory = "";
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
