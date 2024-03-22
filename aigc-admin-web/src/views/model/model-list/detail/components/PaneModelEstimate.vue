<template>
  <Pane ref="refPane" class="" @submit="onSubmit">
    <div class="mx-auto mt-10" style="width: 500px">
      <v-form ref="refForm" class="my-form">
        <Select
          density="compact"
          variant="outlined"
          placeholder="请选择评估数据集"
          hide-details="auto"
          :rules="rules.datasetType"
          v-model="formData.datasetType"
          clearable
          :disabled="state.disabledField"
          :mapDictionary="{ code: 'model_eval_dataset_type' }"
          @change="doQueryDataSet"
        >
          <template #prepend>
            <label class="required" style="width: 100%"
              >评估数据集 <Explain>默认是训练数据集，您可以自定义上传验证集</Explain></label
            >
          </template>
        </Select>

        <v-input
          :rules="rules.fileId"
          v-model="formData.fileId"
          hide-details="auto"
          style="position: relative"
          v-if="!showDataSet"
        >
          <template v-if="echoFileSelect">
            <v-chip closable color="info" @click:close="fileChipClose"> {{ echoFileSelect }}</v-chip>
          </template>
          <template v-else>
            <v-btn color="info" variant="outlined"
              >选择文件
              <UploadFile
                show-loading
                label2="上传微调文件"
                purpose="fine-tune-eval"
                v-model="formData.fileId"
                @upload:success="doQueryFirstPage"
                @loading="val => (upLoading = val)"
                style="width: 92px; position: absolute; top: 0; left: -10%; opacity: 0"
              />
            </v-btn>
          </template>
          <template #prepend> <label class="required">上传验证集</label></template>
        </v-input>
        <v-text-field
          v-if="showDataSet"
          density="compact"
          variant="outlined"
          type="number"
          placeholder="请输入评估数量"
          hide-details="auto"
          :rules="rules.evalPercent"
          v-model.number="formData.evalPercent"
          :disabled="state.disabledField"
          style="width: 80%"
        >
          <template #prepend>
            <label class="required" style="width: 100%">评估数量(%) <Explain>数据集数量的百分比，默认100%</Explain></label>
          </template>
        </v-text-field>
        <Select
          density="compact"
          variant="outlined"
          placeholder="请选择评估指标"
          hide-details="auto"
          :rules="rules.metricName"
          v-model="formData.metricName"
          clearable
          :disabled="state.disabledField"
          :mapDictionary="{ code: 'model_eval_metric' }"
        >
          <template #prepend>
            <label class="required">评估指标 <Explain>评估的方式</Explain></label>
          </template>
        </Select>

        <v-textarea v-model.trim="formData.remark" placeholder="请输入">
          <template #prepend> <label>备注</label></template>
        </v-textarea>
      </v-form>
    </div>
  </Pane>
</template>
<script setup>
import { reactive, toRefs, ref, defineProps } from "vue";
import _ from "lodash";
import Explain from "@/components/ui/Explain.vue";
import UploadFile from "@/components/business/UploadFile.vue";
import { http, validator, dataDictionary } from "@/utils";
const props = defineProps({
  modelName: String
});
const state = reactive({
  operateType: "", //add 添加  、 edit 编辑
  disabledField: false,
  maxTokens: 4096,
  formData: {
    datasetType: "",
    evalPercent: 100,
    metricName: "",
    fileId: "",
    remark: ""
  },
  showDataSet: false
});
const { formData, showDataSet } = toRefs(state);
const echoFileSelect = ref(null);
const emits = defineEmits(["submit"]);
const refPane = ref();
const refForm = ref();
const validNumberInput = (value, min, max, errorMessage) => {
  if (value !== "") {
    if (value < min) {
      return `下限 ${min}`;
    } else if (value > max) {
      return `上限 ${max}`;
    } else {
      return true;
    }
  } else {
    return errorMessage;
  }
};

const rules = reactive({
  operateType: [
    value => {
      if (value && value.length > 0) {
        return true;
      } else {
        return "请选择评估数据集";
      }
    }
  ],
  fileId: [v => !!v || "请选择微调文件"],
  datasetType: [v => !!v || "请选择评估数据集"],
  metricName: [v => !!v || "请选择评估指标"],
  evalPercent: [v => validNumberInput(v, 1, 101, "请输入评估数量")]
});

const doAdd = async (options = {}) => {
  state.formData = { ...state.formData };
  const [err, res] = await http.post({
    ...options,
    showSuccess: true,
    url: `/api/models/eval`,
    data: {
      ...state.formData,
      modelName: props.modelName
    }
  });
  if (res) {
    refPane.value.hide();
    emits("submit");
  }
};
const doQueryFirstPage = item => {
  echoFileSelect.value = item.filename;
};
const fileChipClose = () => {
  echoFileSelect.value = null;
  state.formData.fileId = "";
};

const doQueryDataSet = item => {
  state.showDataSet = item == "train" ? true : false;
};

const onSubmit = ({ valid, showLoading }) => {
  if (valid) {
    if (state.operateType == "add") {
      doAdd({ showLoading });
    }
  }
};

defineExpose({
  show({
    title,
    infos = {
      fileId: "",
      datasetType: "train",
      metricName: "",
      remark: "",
      evalPercent: 100
    },
    operateType
  }) {
    refPane.value.show({
      title,
      refForm
    });
    state.formData = _.pick(_.cloneDeep(infos), [
      "id",
      "name",
      "alias",
      "datasetType",
      "email",
      "projectName",
      "serviceName",
      "remark"
    ]);

    state.operateType = operateType;
    if (operateType == "add") {
      //添加
      state.disabledField = false;
    } else {
      //编辑
      state.disabledField = true;
      // state.formData.modelId = infos.model.list.map(item => {
      //   return item.id;
      // });
    }
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
