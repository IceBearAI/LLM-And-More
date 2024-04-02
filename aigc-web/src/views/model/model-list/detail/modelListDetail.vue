<template>
  <NavBack backUrl="/model/model-list/list">模型详情</NavBack>
  <UiParentCard class="mt-4">
    <template #header> 模型：{{ rawData.modelName }}</template>
    <template #action>
      <v-btn
        v-if="rawData.deployStatus == 'success' || rawData.deployStatus == 'running'"
        flat
        color="secondary"
        @click="onUninstall"
        >卸载</v-btn
      >
      <v-btn v-if="rawData.deployStatus == 'failed'" flat color="secondary" @click="onArrange">部署</v-btn>
      <!-- <ButtonsInTable :buttons="getButtons()" style="width: 80px" /> -->
    </template>

    <ModelListDetailBaseInfo />
  </UiParentCard>
  <UiParentCard class="mt-5">
    <template #header>
      <v-tabs v-model="state.tabIndex" align-tabs="start" color="primary">
        <!-- <v-tab :value="1">详情</v-tab> -->
        <v-tab :value="2">当前模型下的评估列表</v-tab>
        <!-- <v-tab :value="3">数据集</v-tab> -->
      </v-tabs>
    </template>

    <v-window v-model="state.tabIndex">
      <!-- <v-window-item :value="1"> <TabFineTuningDetail /> </v-window-item> -->
      <v-window-item :value="2">
        <TabModelEstimate
          :showArrange="rawData.deployStatus"
          :modelTitle="rawData.modelName"
          :providerName="rawData.providerName"
      /></v-window-item>
    </v-window>
  </UiParentCard>
  <ConfirmByClick ref="refConfirmByClick" @submit="onConfirmByClick">
    <template #text> <div v-html="state.confirmByClickInfo.html"></div></template>
  </ConfirmByClick>
</template>
<script setup lang="ts">
import { reactive, toRefs, ref, provide } from "vue";
import NavBack from "@/components/business/NavBack.vue";
import { http } from "@/utils";
import { useRoute } from "vue-router";
import UiParentCard from "@/components/shared/UiParentCard.vue";
import ModelListDetailBaseInfo from "./components/ModelListDetailBaseInfo.vue";
import TabModelEstimate from "./components/TabModelEstimate.vue";
import ConfirmByClick from "@/components/business/ConfirmByClick.vue";
import _ from "lodash";
const route = useRoute();
const refConfirmByClick = ref();
const state = reactive({
  tabIndex: "2",
  style: {},
  formData: {},
  rawData: {
    modelName: "",
    deployStatus: "",
    providerName: ""
  },
  confirmByClickInfo: {
    html: "",
    action: "",
    row: {}
  }
});
const { style, rawData, formData } = toRefs(state);
provide("provideModelListDetail", state);
provide("provideModelListDetailed", state);
const getData = async () => {
  let { jobId } = route.query;
  let [err, res] = await http.get({
    showLoading: true,
    url: `/api/models/${jobId}`
  });
  if (res) {
    state.rawData = res;
  }
};
getData();

const onConfirmByClick = (options = {}) => {
  let { action, row } = state.confirmByClickInfo;
  if (action == "deploy") {
    onDeploy(row, options);
  } else if (action == "undeploy") {
    onUndeploy(row, options);
  }
};
const onDeploy = async (row, options) => {
  let [err, res] = await http.post({
    ...options,
    showSuccess: true,
    url: `/api/models/${row.id}/deploy`,
    data: {
      ...state.formData,
      id: row.id
    }
  });
  if (res) {
    //关闭确认框
    refConfirmByClick.value.hide();
    // doQueryCurrentPage();
  }
};
const onUndeploy = async (row, options) => {
  let [err, res] = await http.post({
    ...options,
    showSuccess: true,
    url: `/api/models/${row.id}/undeploy`,
    data: {
      ...state.formData,
      id: row.id
    }
  });
  if (res) {
    //关闭确认框
    refConfirmByClick.value.hide();
    // doQueryCurrentPage();
  }
};

const onArrange = () => {
  _.assign(state.confirmByClickInfo, {
    html: "确认部署模型 <span class='text-primary mx-1  font-weight-black'>" + state.rawData.modelName + "</span>吗？",
    action: "deploy",
    row: state.rawData
  });
  refConfirmByClick.value.show({ width: "350px" });
};
const onUninstall = () => {
  _.assign(state.confirmByClickInfo, {
    html: "确认卸载模型 <span class='text-primary mx-1  font-weight-black'>" + state.rawData.modelName + "</span>吗？",
    action: "undeploy",
    row: state.rawData
  });
  refConfirmByClick.value.show({ width: "350px" });
};
const getButtons = () => {
  let ret = [];
  if (state.rawData.deployStatus == "success" || state.rawData.deployStatus == "running") {
    ret.push({
      text: "部署",
      click() {
        onArrange();
      }
    });
  } else {
    ret.push({
      text: "卸载",
      color: "error",
      click() {
        onUninstall();
      }
    });
  }

  return ret;
};
</script>
<style lang="scss"></style>
