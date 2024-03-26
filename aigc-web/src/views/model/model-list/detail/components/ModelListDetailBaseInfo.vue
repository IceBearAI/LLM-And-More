<template>
  <v-row class="my-form waterfall">
    <v-col xs="12" sm="6" md="4" lg="3">
      <v-input hide-details>
        <template #prepend> <label>供应</label></template>
        {{ rawData.providerName }}
      </v-input>
    </v-col>
    <v-col xs="12" sm="6" md="4" lg="3">
      <v-input hide-details>
        <template #prepend> <label>类型</label></template>
        {{ getLabels([["model_type", rawData.modelType]]) }}
      </v-input>
    </v-col>
    <v-col xs="12" sm="6" md="4" lg="3">
      <v-input hide-details>
        <template #prepend> <label>最长上下文</label></template>
        {{ rawData.maxTokens }}
      </v-input>
    </v-col>
    <v-col xs="12" sm="6" md="4" lg="3" v-if="rawData.isFineTuning">
      <v-input hide-details>
        <template #prepend> <label>微调</label></template>
        <el-tooltip content="微调详情" placement="top">
          <router-link :to="'/model/fine-tuning/detail?jobId=' + rawData.jobId" class="link">{{
            rawData.isFineTuning ? "是" : ""
          }}</router-link>
        </el-tooltip>
      </v-input>
    </v-col>
    <v-col xs="12" sm="6" md="4" lg="3">
      <v-input hide-details>
        <template #prepend> <label>参数量</label></template>
        <span style="color: #539bff">{{ rawData.parameters }}B</span>
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
        <template #prepend> <label>备注</label></template>
        {{ rawData.remark }}
      </v-input>
    </v-col>
  </v-row>
</template>
<script setup>
import { reactive, toRefs, ref, inject, computed, onMounted } from "vue";
import { format } from "@/utils";
import { useMapRemoteStore } from "@/stores";

const { loadDictTree, getLabels } = useMapRemoteStore();
const provideModelListDetail = inject("provideModelListDetail");

const rawData = computed(() => {
  return provideModelListDetail.rawData;
});

const init = async () => {
  await loadDictTree(["model_type"]);
};

onMounted(() => {
  init();
});
</script>
<style lang="scss" scoped>
.my-form > * {
  margin-bottom: 0;
}
</style>
