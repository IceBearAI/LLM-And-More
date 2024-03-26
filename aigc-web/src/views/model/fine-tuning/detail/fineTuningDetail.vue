<template>
  <NavBack backUrl="/model/fine-tuning/list">微调详情</NavBack>
  <UiParentCard class="mt-4">
    <template #header> 任务ID：{{ rawData.jobId }}</template>
    <FineTuningBaseInfo />
  </UiParentCard>
  <UiParentCard class="mt-5">
    <template #header>
      <v-tabs v-model="state.tabIndex" align-tabs="start" color="primary">
        <v-tab :value="1">详情</v-tab>
        <v-tab :value="2">日志</v-tab>
      </v-tabs>
    </template>

    <v-window v-model="state.tabIndex">
      <v-window-item :value="1"> <TabFineTuningDetail /> </v-window-item>
      <v-window-item :value="2">
        <TextLog v-model="rawData.trainLog" style="height: 600px" :idDone="true" />
      </v-window-item>
    </v-window>
  </UiParentCard>
</template>

<script setup>
import { reactive, toRefs, provide } from "vue";
import NavBack from "@/components/business/NavBack.vue";
import { http } from "@/utils";
import { useRoute } from "vue-router";
import UiParentCard from "@/components/shared/UiParentCard.vue";
import FineTuningBaseInfo from "./components/FineTuningBaseInfo.vue";
import TabFineTuningDetail from "./components/TabFineTuningDetail.vue";
import TextLog from "@/components/ui/log/TextLog.vue";

const route = useRoute();
const state = reactive({
  tabIndex: "",
  style: {},
  rawData: {}
});
const { style, rawData } = toRefs(state);

// defineOptions({
//   name: "FineTuningDetail"
// });

provide("provideFineTuningDetail", state);

const getData = async () => {
  let { jobId } = route.query;
  let [err, res] = await http.get({
    showLoading: true,
    url: `/api/finetuning/${jobId}`
  });
  if (res) {
    state.rawData = res;
  }
};

getData();
</script>
<style lang="scss"></style>
