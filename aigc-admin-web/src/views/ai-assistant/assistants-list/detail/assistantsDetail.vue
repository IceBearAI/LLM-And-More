<template>
  <NavBack backUrl="/ai-assistant/assistants/list">助手详情</NavBack>
  <UiParentCard class="mt-4">
    <template #header> 名称：{{ data.name }}</template>
    <template #action>
      <div>
        <!-- <v-btn color="primary" @click="onPublish">发布</v-btn> -->
        <v-btn class="ml-4" color="primary" @click="goPlayground">去试试</v-btn>
        <v-btn class="ml-4" color="primary" @click="onEdit">编辑</v-btn>
      </div>
    </template>
    <DetailBaseInfo :info="data" />
  </UiParentCard>
  <UiParentCard class="mt-5">
    <template #header>
      <v-tabs v-model="tabIndex" align-tabs="start" color="primary">
        <v-tab :value="1">指令描述</v-tab>
      </v-tabs>
    </template>
    <v-window v-model="tabIndex">
      <v-window-item :value="1">
        <UiChildCard title="指令描述">
          <div class="whitespace-pre-wrap" v-text="data.instructions"></div>
        </UiChildCard>
        <UiChildCard class="mt-4" title="工具集">
          <ToolsTableInfo />
        </UiChildCard>
      </v-window-item>
    </v-window>
  </UiParentCard>
  <ConfirmByInput ref="refConfirmPublish" @submit="doPublish">
    <template #text>
      您将要发布 <span class="text-primary font-weight-black">{{ data.name }}</span> 助手<br />
      确定要继续吗？
    </template>
  </ConfirmByInput>
  <CreateAssistantPane ref="createAssistantPaneRef" @submit="getData" />
</template>
<script setup lang="ts">
import { ref, onMounted } from "vue";
import CreateAssistantPane from "../components/CreateAssistantPane.vue";
import NavBack from "@/components/business/NavBack.vue";
import { http } from "@/utils";
import { useRoute, useRouter } from "vue-router";
import UiParentCard from "@/components/shared/UiParentCard.vue";
import UiChildCard from "@/components/shared/UiChildCard.vue";
import ConfirmByInput from "@/components/business/ConfirmByInput.vue";
import DetailBaseInfo from "./components/DetailBaseInfo.vue";
import ToolsTableInfo from "./components/ToolsTableInfo.vue";

const route = useRoute();
const router = useRouter();
const { assistantId } = route.query;

const createAssistantPaneRef = ref();
const data = ref<Record<string, any>>({});
const tabIndex = ref("");
const refConfirmPublish = ref();

const getData = async () => {
  let [err, res] = await http.get({
    showLoading: true,
    url: `/assistants/${assistantId}`
  });
  if (res) {
    data.value = res;
  }
};

const goPlayground = () => {
  router.push(`/ai-assistant/assistants/playground?assistantId=${assistantId}`);
};

const onPublish = () => {
  refConfirmPublish.value.show({
    width: "450px",
    confirmText: data.value.name
  });
};

const doPublish = async (options = {}) => {
  const [err, res] = await http.post({
    ...options,
    showSuccess: true,
    url: `/assistants/${assistantId}/publish`
  });
  if (res) {
    refConfirmPublish.value.hide();
  }
};

const onEdit = () => {
  createAssistantPaneRef.value.show({
    title: "编辑助手",
    infos: data.value,
    operateType: "edit"
  });
};

onMounted(() => {
  getData();
});
</script>
<style lang="scss"></style>
