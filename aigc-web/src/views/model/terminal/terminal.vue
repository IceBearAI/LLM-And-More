<template>
  <UiParentCard>
    <v-row>
      <v-col cols="12" lg="3" md="4" sm="6">
        <v-select
          v-model="searchData.containerName"
          @update:modelValue="containerNameChange"
          :clearable="false"
          hide-details
          :items="containerOptions"
          label="请选择容器"
        >
        </v-select>
      </v-col>

      <v-col style="height: calc(100vh - 160px)" cols="12">
        <Terminal class="h-100" v-if="wsUrl" :ws-url="wsUrl" :start-data="startData" />
      </v-col>
    </v-row>
  </UiParentCard>
</template>
<script setup lang="ts">
import { ref, reactive, onMounted, computed } from "vue";
import { useRoute } from "vue-router";
import { http } from "@/utils";
import UiParentCard from "@/components/shared/UiParentCard.vue";
import Terminal from "@/components/business/Terminal.vue";
import { useCustomizerStore } from "@/stores/customizer";

const customizer = useCustomizerStore();

const route = useRoute();
const { modelId } = route.query;
const searchData = reactive({
  containerName: null
});
const containerOptions = ref([]);
const modelData = reactive({
  serviceName: ""
});
const sessionId = ref("");
const startData = ref({});

const wsUrl = computed(() => {
  let domain = "";
  if (window.env === "development") {
    domain = "http://aigc-admin-web.aigc.paas.test";
    // domain = "http://10.21.20.46:8080";
  } else {
    domain = location.protocol + "//" + location.host;
  }
  if (sessionId.value) {
    return `${domain}/api/ws/terminal/console/exec?sessionId=${sessionId.value}`;
  } else {
    return "";
  }
});

const containerNameChange = () => {
  sessionId.value = "";
  getServiceWebsocketToken();
};

const getData = async () => {
  let [err, res] = await http.get({
    showLoading: true,
    url: `/models/${modelId}`
  });
  if (res) {
    containerOptions.value = res.containerNames || [];
    searchData.containerName = containerOptions.value[0];
    modelData.serviceName = res.serviceName;
    getServiceWebsocketToken();
  }
};

const getServiceWebsocketToken = async () => {
  const [err, res] = await http.get({
    showLoading: true,
    url: `/ws/resource/deployment/service/${modelData.serviceName}/container/${searchData.containerName}/token`
  });
  if (res) {
    startData.value = res;
    sessionId.value = res.sessionId;
  }
};

onMounted(() => {
  customizer.SET_MINI_SIDEBAR(true);
  getData();
});
</script>
