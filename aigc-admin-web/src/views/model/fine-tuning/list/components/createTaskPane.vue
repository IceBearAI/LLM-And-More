<template>
  <Pane ref="refPane" @submit="onSubmit">
    <template #txtLeft v-if="state.estimateTime">
      <p class="orangeTxt">
        预计训练时间：<span>{{ state.estimateTime }}</span>
      </p>
    </template>
    <div class="mx-auto mt-3" style="width: 540px">
      <AlertBlock class="mb-6">微调模型的时间可能会很长，微调完成之后会邮件通知您！</AlertBlock>
      <v-form ref="refForm" class="my-form">
        <v-input :rules="rules.fileId" v-model="formData.fileId" hide-details="auto">
          <template v-if="echoFileSelect">
            <v-chip closable color="info" @click:close="fileChipClose">{{ echoFileSelect.filename }}</v-chip>
          </template>
          <template v-else>
            <v-btn color="info" variant="outlined" @click="openFileDialog">选择文件</v-btn>
          </template>
          <template #prepend> <label class="required">微调文件</label></template>
        </v-input>
        <Select
          placeholder="请选择基础模型"
          :rules="rules.baseModel"
          v-model="formData.baseModel"
          :mapAPI="{ url: '/finetuning/base/model' }"
          hide-details="auto"
          @change="doQueryModal"
        >
          <template #prepend>
            <label class="required">基础模型<Explain>基于基础模型进行微调</Explain></label>
          </template>
        </Select>
        <v-text-field
          type="number"
          placeholder="上限512"
          hide-details="auto"
          :rules="rules.trainEpoch"
          v-model.number="formData.trainEpoch"
          @input="doQueryModal"
        >
          <template #prepend> <label class="required">训练轮次</label></template>
        </v-text-field>
        <v-text-field
          type="text"
          placeholder="请输入微调后模型名称的后缀"
          hide-details="auto"
          clearable
          v-model="formData.suffix"
        >
          <template #prepend>
            <label>后缀 <Explain>微调后模型名称的后缀，通常微调名称为 ft::{模型}:{渠道}:-{随机}:{后缀}</Explain></label></template
          >
        </v-text-field>
        <v-textarea v-model.trim="formData.remark" placeholder="请输入你微调模型的备注" clearable>
          <template #prepend> <label>备注</label></template>
        </v-textarea>
        <v-input>
          <v-expansion-panels ref="refSettingExpansion" v-model="settingExpansion">
            <v-expansion-panel value="advancedSetting" elevation="10">
              <v-expansion-panel-title class="text-h6">高级设置</v-expansion-panel-title>
              <v-expansion-panel-text class="mt-4" eager>
                <v-checkbox label="启用Lora" v-model="formData.lora">
                  <template #prepend>
                    <label>Lora </label>
                  </template>
                </v-checkbox>
                <v-text-field
                  type="number"
                  placeholder="请输入训练批次"
                  :rules="rules.trainBatchSize"
                  v-model.number="formData.trainBatchSize"
                >
                  <template #prepend>
                    <label class="required">训练批次 <Explain>即一次训练所抓取的数据样本数量</Explain></label></template
                  >
                </v-text-field>
                <v-text-field
                  type="number"
                  placeholder="请输入评估批次"
                  :rules="rules.evalBatchSize"
                  v-model.number="formData.evalBatchSize"
                >
                  <template #prepend>
                    <label class="required"
                      >评估批次 <Explain>用于评估的每个 GPU/TPU 核心/CPU 的批量大小</Explain></label
                    ></template
                  >
                </v-text-field>
                <v-text-field
                  type="number"
                  placeholder="请输入梯度累加步数"
                  :rules="rules.accumulationSteps"
                  v-model.number="formData.accumulationSteps"
                >
                  <template #prepend>
                    <label class="required"
                      >梯度累加步数 <Explain>在执行向后/更新传递之前累积梯度的更新步骤数</Explain></label
                    ></template
                  >
                </v-text-field>
                <v-text-field
                  type="number"
                  placeholder="请输入使用GPU数量"
                  :rules="rules.procPerNode"
                  v-model.number="formData.procPerNode"
                  @input="doQueryModal"
                >
                  <template #prepend>
                    <label class="required">使用GPU数量 <Explain>使用GPU数量</Explain></label></template
                  >
                </v-text-field>
                <v-text-field
                  type="number"
                  placeholder="请输入学习率"
                  :rules="rules.learningRate"
                  v-model.number="formData.learningRate"
                  :step="0.00001"
                >
                  <template #prepend>
                    <label class="required">学习率 <Explain>AdamW 优化器的初始学习率</Explain></label></template
                  >
                </v-text-field>
                <v-text-field
                  type="number"
                  placeholder="请输入模型最大长度"
                  :rules="rules.modelMaxLength"
                  v-model.number="formData.modelMaxLength"
                >
                  <template #prepend>
                    <label class="required">模型最大长度 <Explain>模型支持的最大上下文长度</Explain></label></template
                  >
                </v-text-field>
              </v-expansion-panel-text>
            </v-expansion-panel>
          </v-expansion-panels>
          <!-- <template #prepend><label></label></template> -->
        </v-input>
      </v-form>
    </div>
    <Dialog ref="refDialogUpload" attach persistent max-width="1200px" :retain-focus="false">
      <template #title> 选择文件 </template>
      <FileList @selected="handleFileSelect" purpose="fine-tune" />
    </Dialog>
  </Pane>
</template>
<script setup lang="ts">
import { reactive, ref, toRefs, customRef } from "vue";
import _ from "lodash";
import Explain from "@/components/ui/Explain.vue";
import AlertBlock from "@/components/ui/AlertBlock.vue";
import FileList from "./FileList.vue";
import { http } from "@/utils";

const paneConfig = reactive({
  operateType: "add"
});
const formData = reactive({
  fileId: "",
  baseModel: null,
  trainEpoch: 1,
  suffix: "",
  remark: "",
  trainBatchSize: 8,
  evalBatchSize: 32,
  accumulationSteps: 1,
  procPerNode: 2,
  learningRate: 0.00002,
  modelMaxLength: 2048,
  lora: false
});
const state = reactive({
  estimateTime: "",
  readyOnly: true
});
const echoFileSelect = ref(null);

const emits = defineEmits(["submit"]);
const regNum = /^\+?[1-9][0-9]*$/;
const refPane = ref();
const refForm = ref();
const refDialogUpload = ref();
const refSettingExpansion = ref();
const rules = reactive({
  fileId: [v => !!v || "请选择微调文件"],
  baseModel: [v => !!v || "请选择基础模型"],
  trainEpoch: [v => validNumberInput(v, 1, 512, "请输入训练轮次", true)],
  trainBatchSize: [v => validNumberInput(v, 1, 512, "请输入训练批次", true)],
  evalBatchSize: [v => validNumberInput(v, 1, 32, "请输入评估批次", true)],
  accumulationSteps: [v => validNumberInput(v, 1, 32, "请输入梯度累加步数", true)],
  procPerNode: [v => validNumberInput(v, 1, 16, "请输入GPU数量", true)],
  learningRate: [v => validNumberInput(v, 0.000001, 1, "请输入学习率")],
  modelMaxLength: [v => validNumberInput(v, 1, 500000, "请输入模型最大长度", true)]
});
const settingExpansion = ref("");
const {} = toRefs(state);

const validNumberInput = (value, min, max, errorMessage, reg = false) => {
  if (value !== "") {
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
    return errorMessage;
  }
};
const debounceRef = (data, delay = 300) => {
  let timer = null;
  return delay == null
    ? ref(data)
    : customRef((track, trigger) => {
        return {
          get() {
            track();
            return data;
          },
          set(value) {
            if (timer != null) {
              clearTimeout(timer);
              timer = null;
            }
            timer = setTimeout(() => {
              data = value;
              trigger();
            }, delay);
          }
        };
      });
};
debounceRef(formData.trainEpoch);
const openFileDialog = () => {
  refDialogUpload.value.show({
    showActions: false
  });
};

const handleFileSelect = row => {
  refDialogUpload.value.hide();
  echoFileSelect.value = row;
  formData.fileId = row.fileId;
  doQueryModal();
};

const fileChipClose = () => {
  echoFileSelect.value = null;
  formData.fileId = "";
};

const doAdd = async (options = {}) => {
  const [err, res] = await http.post({
    ...options,
    showSuccess: true,
    url: `/finetuning`,
    data: formData
  });
  if (res) {
    refPane.value.hide();
    emits("submit");
  }
};
const getData = async (options = {}) => {
  const [err, res] = await http.post({
    ...options,
    showSuccess: false,
    url: `/api/finetuning/estimate`,
    data: formData
  });
  if (res) {
    state.estimateTime = res.time ? res.time : "";
  }
};

const doEdit = async (options = {}) => {};
const doQueryModal = () => {
  if (
    formData.baseModel &&
    formData.fileId.length > 0 &&
    formData.procPerNode > 0 &&
    regNum.test(String(formData.trainEpoch)) != false &&
    formData.trainEpoch > 0 &&
    regNum.test(String(formData.procPerNode)) != false
  ) {
    getData();
  }
};

const openAdvancedSetting = errors => {
  const settingExpansionEl = refSettingExpansion.value.$el;
  const isOpen = errors.some(item => {
    return settingExpansionEl.querySelector(`#${item.id}`);
  });
  if (isOpen) {
    settingExpansion.value = "advancedSetting";
  }
};

const onSubmit = ({ valid, errors, showLoading }) => {
  if (valid) {
    if (paneConfig.operateType == "add") {
      doAdd({ showLoading });
    } else {
      doEdit({ showLoading });
    }
  } else {
    openAdvancedSetting(errors);
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
      echoFileSelect.value = null;
      formData.baseModel = null;
      formData.fileId = "";
      formData.remark = "";
      formData.suffix = "";
      state.estimateTime = "";
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

:deep() {
  .v-expansion-panel-text__wrapper {
    margin: 0 -20px;
  }
}
.orangeTxt {
  font-size: 14px;
  width: 90%;
  display: flex;
  color: #333;
  align-items: center;
  letter-spacing: 0.1em;
  span {
    font-size: 22px;
    font-weight: 500;
    color: #ff6600;
    border-radius: 18px;
  }
}
</style>
