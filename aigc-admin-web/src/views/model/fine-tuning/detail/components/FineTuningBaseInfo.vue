<template>
  <v-row class="my-form waterfall">
    <v-col xs="12" sm="6" md="4" lg="3">
      <v-input hide-details>
        <template #prepend> <label>创建人</label></template>
        {{ rawData.trainPublisher }}
      </v-input>
    </v-col>
    <v-col xs="12" sm="6" md="4" lg="3">
      <v-input hide-details>
        <template #prepend> <label>基础模型</label></template>
        <router-link to="" class="link">{{ rawData.baseModel }}</router-link>
      </v-input>
    </v-col>
    <v-col xs="12" sm="6" md="4" lg="3">
      <v-input hide-details>
        <template #prepend> <label>训练模型</label></template>
        {{ rawData.fineTunedModel }}
      </v-input>
    </v-col>
    <v-col xs="12" sm="6" md="4" lg="3">
      <v-input hide-details>
        <template #prepend> <label>是否Lora训练</label></template>
        <ChipBoolean v-model="rawData.lora" />
      </v-input>
    </v-col>
    <v-col xs="12" sm="6" md="4" lg="3">
      <v-input hide-details>
        <template #prepend> <label>使用GPU</label></template>
        {{ rawData.procPerNode }}
      </v-input>
    </v-col>
    <v-col xs="12" sm="6" md="4" lg="3">
      <v-input hide-details>
        <template #prepend> <label>创建时间</label></template>
        {{ format.dateFormat(rawData.createdAt, "YYYY-MM-DD HH:mm:ss") }}
      </v-input>
    </v-col>
    <v-col xs="12" sm="6" md="4" lg="3">
      <v-input hide-details>
        <template #prepend> <label>训练轮次</label></template>
        {{ rawData.trainEpoch }}
      </v-input>
    </v-col>
    <v-col xs="12" sm="6" md="4" lg="3">
      <v-input hide-details>
        <template #prepend> <label>模型后缀</label></template>
        {{ rawData.suffix }}
      </v-input>
    </v-col>
    <v-col xs="12" sm="6" md="4" lg="3">
      <v-input hide-details>
        <template #prepend> <label>训练批次</label></template>
        {{ rawData.trainBatchSize }}次训练,{{ rawData.evalBatchSize }}次评估,{{ rawData.accumulationSteps }} 次梯度累加
      </v-input>
    </v-col>
    <v-col xs="12" sm="6" md="4" lg="3">
      <v-input hide-details>
        <template #prepend> <label>学习率</label></template>
        {{ format.toScientfic(rawData.learningRate) }}
      </v-input>
    </v-col>
    <v-col xs="12" sm="6" md="4" lg="3">
      <v-input hide-details>
        <template #prepend> <label>开始时间</label></template>
        {{ format.dateFormat(rawData.startTrainTime, "YYYY-MM-DD HH:mm:ss") }}
      </v-input>
    </v-col>
    <v-col xs="12" sm="6" md="4" lg="3">
      <v-input hide-details>
        <template #prepend> <label>备注</label></template>
        {{ rawData.remark }}
      </v-input>
    </v-col>
    <v-col xs="12" sm="6" md="4" lg="3">
      <v-input hide-details>
        <template #prepend> <label>文件</label></template>
        <a class="link1 line1" @click="onDownload(rawData)">{{ rawData.fileId }}</a>
      </v-input>
    </v-col>
    <v-col xs="12" sm="6" md="4" lg="3">
      <v-input hide-details>
        <template #prepend> <label>完成时间</label></template>
        {{ format.dateFormat(rawData.finishedAt, "YYYY-MM-DD HH:mm:ss") }}
      </v-input>
    </v-col>
    <v-col xs="12" sm="6" md="4" lg="3">
      <v-input hide-details>
        <template #prepend> <label>模型最大长度</label></template>
        {{ rawData.modelMaxLength }}
      </v-input>
    </v-col>
  </v-row>
</template>
<script setup>
import { reactive, toRefs, ref, inject, computed } from "vue";
import { format, http } from "@/utils";
import ChipBoolean from "@/components/ui/ChipBoolean.vue";

const state = reactive({
  style: {},
  formData: {}
});
const { style, formData } = toRefs(state);

const provideFineTuningDetail = inject("provideFineTuningDetail");
const onDownload = row => {
  // window.open(row.fileUrl);
  http.downloadByUrl({
    fileUrl: row.fileUrl,
    suffixName: "jsonl"
  });
};
const rawData = computed(() => {
  return provideFineTuningDetail.rawData;
});
</script>
<style lang="scss" scoped>
.my-form > * {
  margin-bottom: 0;
}
.link1 {
  color: #539cff;

  cursor: pointer;
  &:hover {
    text-decoration: underline;
  }
}
</style>
