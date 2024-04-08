<template>
  <NavBack backUrl="/sample-library/intention-mark/list">意图模型标注详情</NavBack>
  <UiParentCard class="mt-4">
    <template #header> 名称：{{ data.name }}</template>
    <!-- <template #action></template> -->
    <DetailBaseInfo :info="data" />
  </UiParentCard>
  <UiParentCard class="mt-5">
    <template #header>
      <v-tabs v-model="tabIndex" align-tabs="start" color="primary">
        <v-tab :value="1">问答</v-tab>
        <v-tab :value="2">标注</v-tab>
      </v-tabs>
    </template>
    <v-window v-model="tabIndex">
      <v-window-item class="pt-2" :value="1">
        <QAList />
      </v-window-item>
      <v-window-item class="pt-2" :value="2">
        <MarkList />
      </v-window-item>
    </v-window>
  </UiParentCard>
</template>
<script setup lang="ts">
import { ref, onMounted } from "vue";
import NavBack from "@/components/business/NavBack.vue";
import { http } from "@/utils";
import { useRoute } from "vue-router";
import UiParentCard from "@/components/shared/UiParentCard.vue";
import DetailBaseInfo from "./components/DetailBaseInfo.vue";
import QAList from "./components/QAList.vue";
import MarkList from "./components/MarkList.vue";

const route = useRoute();
const { intentId } = route.query;

const data = ref<Record<string, any>>({});
const tabIndex = ref("");

const getData = async () => {
  let [err, res] = await http.get({
    showLoading: true,
    url: `/intent/${intentId}`
  });
  if (res) {
    data.value = res;
  }
};

onMounted(() => {
  getData();
});
</script>
