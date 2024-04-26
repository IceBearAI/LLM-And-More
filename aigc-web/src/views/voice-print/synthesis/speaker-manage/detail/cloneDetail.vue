<template>
  <NavBack backUrl="/voice-print/synthesis/speaker/list">声音克隆（微软个人版）</NavBack>
  <UiParentCard class="mt-4">
    <template #header> 发声人：{{ `${data.speakCname}（${data.speakName}）` }}</template>
    <!-- <template #action></template> -->
    <CloneDetailBaseInfo :info="data" />
  </UiParentCard>
  <UiParentCard class="mt-5">
    <template #header>
      <v-tabs v-model="tabIndex" align-tabs="start" color="primary">
        <v-tab :value="1">授权语音列表</v-tab>
        <v-tab :value="2">参考音频列表</v-tab>
      </v-tabs>
    </template>
    <v-window v-model="tabIndex">
      <v-window-item class="pt-2" :value="1">
        <ConsentVoice />
      </v-window-item>
      <v-window-item class="pt-2" :value="2">
        <ReferenceVoice />
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
import CloneDetailBaseInfo from "./components/CloneDetailBaseInfo.vue";
import ConsentVoice from "./components/ConsentVoice.vue";
import ReferenceVoice from "./components/ReferenceVoice.vue";

const route = useRoute();
const { speakName } = route.query;

const data = ref<Record<string, any>>({});
const tabIndex = ref("");

const getData = async () => {
  let [err, res] = await http.get({
    showLoading: true,
    url: `/voice/speak/${speakName}`
  });
  if (res) {
    data.value = res;
  }
};

onMounted(() => {
  getData();
});
</script>
