<template>
  <BaseBreadcrumb :title="page.title" :breadcrumbs="breadcrumbs"></BaseBreadcrumb>
  <UiParentCard class="mb-3">
    <TaskOverview :config="taskDetailConfig" request-url="/finetuning/dashboard" />
  </UiParentCard>
  <UiParentCard>
    <v-row>
      <v-col cols="12" lg="3" md="4" sm="6">
        <v-text-field
          density="compact"
          v-model="searchData.fineTunedModel"
          label="请输入关键字"
          hide-details
          clearable
          variant="outlined"
          @keyup.enter="doQueryFirstPage"
          @click:clear="doQueryFirstPage"
        ></v-text-field>
      </v-col>
      <v-col cols="12" lg="3" md="4" sm="6">
        <Select
          @change="doQueryFirstPage"
          label="请选择状态"
          :mapDictionary="{ code: 'local_trainStatus' }"
          v-model="searchData.trainStatus"
        ></Select>
      </v-col>
      <v-col cols="12" lg="3" md="4" sm="6">
        <ButtonsInForm>
          <v-btn color="primary" @click="onAdd">创建微调任务</v-btn>
          <refresh-button ref="refreshButtonRef" @refresh="doQueryCurrentPage" />
        </ButtonsInForm>
      </v-col>
      <v-col cols="12">
        <AlertBlock> GPU资源紧张，暂时只支持同时进行一个微调任务！ </AlertBlock>
      </v-col>

      <v-col cols="12">
        <TableWithPager @query="doTableQuery" ref="tableWithPagerRef" :infos="tableInfos">
          <el-table-column label="任务ID" width="180px">
            <template #default="{ row }">
              <el-tooltip content="查看详情" placement="top">
                <span class="link" @click="onView(row)">{{ row.jobId }}</span>
              </el-tooltip>
            </template>
          </el-table-column>
          <el-table-column label="基础模型" width="120px">
            <template #default="{ row }">
              <el-tooltip effect="dark" content="进入聊天操场" placement="top">
                <router-link class="link" :to="{ path: '/model/chat-playground', query: { modelName: row.baseModel } }">{{
                  row.baseModel
                }}</router-link>
              </el-tooltip>
            </template>
          </el-table-column>
          <el-table-column label="训练轮次" prop="trainEpoch" min-width="90px"></el-table-column>
          <el-table-column label="Lora">
            <template #default="{ row }"> <ChipBoolean v-model="row.lora" /> </template>
          </el-table-column>
          <el-table-column label="训练状态" min-width="150px">
            <template #default="{ row }">
              <TableTrainStatus :item="row" @open:log="onLog(row)" />
            </template>
          </el-table-column>
          <el-table-column label="训练监控" min-width="100px">
            <template #default="{ row }">
              <v-row v-if="row.diagnosis" dense class="text-error">
                <v-col v-if="row.diagnosis.overfitting === 'High'" cols="12">过拟合</v-col>
                <v-col v-if="row.diagnosis.underfitting === 'High'" cols="12">欠拟合</v-col>
                <v-col v-if="row.diagnosis.catastrophicForgetting === 'High'" cols="12">灾难性遗忘</v-col>
              </v-row>
            </template>
          </el-table-column>
          <el-table-column label="训练时长" prop="trainDuration" min-width="120px"></el-table-column>
          <el-table-column label="训练进度" min-width="90px">
            <template #default="{ row }">
              <span>{{ format.toPercent(row.process, 2) }}</span>
            </template>
          </el-table-column>
          <el-table-column label="模型名称" width="110px">
            <template #default="{ row }">
              <template v-if="row.trainStatus === 'success'">
                <router-link class="link" :to="{ path: '/model/model-list', query: { modelName: row.fineTunedModel } }">{{
                  row.fineTunedModel
                }}</router-link>
              </template>
              <template v-else> -- </template>
            </template>
          </el-table-column>
          <el-table-column label="备注" min-width="200px" show-overflow-tooltip>
            <template #default="{ row }"> {{ row.remark }} </template>
          </el-table-column>
          <el-table-column label="完成时间" min-width="160px">
            <template #default="{ row }">
              {{ row.finishedAt ? format.dateFormat(row.finishedAt, "YYYY-MM-DD HH:mm:ss") : "--" }}
            </template>
          </el-table-column>
          <el-table-column label="创建时间" min-width="160px">
            <template #default="{ row }">
              {{ format.dateFormat(row.createdAt, "YYYY-MM-DD HH:mm:ss") }}
            </template>
          </el-table-column>
          <el-table-column label="操作人" prop="trainPublisher" min-width="150px" show-overflow-tooltip></el-table-column>
          <el-table-column label="操作" min-width="120px" fixed="right">
            <template #default="{ row }">
              <ButtonsInTable :buttons="getButtons(row)" />
            </template>
          </el-table-column>
        </TableWithPager>
      </v-col>
    </v-row>
  </UiParentCard>
  <DialogLog ref="refDialogLog" :interval="20" @refresh="getLog" />
  <ConfirmByInput ref="refConfirmAbort" @submit="doAbort">
    <template #text>
      此操作将会<span class="text-primary">取消</span>正在进行的训练任务<br />
      任务ID：<span class="text-primary font-weight-black">{{ confirmAbort.currentJobId }}</span
      ><br />
      你还要继续吗？
    </template>
  </ConfirmByInput>

  <ConfirmByClick ref="refConfirmDelete" @submit="doDelete">
    <template #text>确定删除此任务？</template>
  </ConfirmByClick>
  <CreateTaskPane ref="createTaskPaneRef" @submit="doQueryFirstPage" />
</template>
<script setup>
import { reactive, ref, onMounted } from "vue";
import BaseBreadcrumb from "@/components/shared/BaseBreadcrumb.vue";
import UiParentCard from "@/components/shared/UiParentCard.vue";
import AlertBlock from "@/components/ui/AlertBlock.vue";
import ConfirmByInput from "@/components/business/ConfirmByInput.vue";
import ConfirmByClick from "@/components/business/ConfirmByClick.vue";
import ChipBoolean from "@/components/ui/ChipBoolean.vue";
import CreateTaskPane from "./components/createTaskPane.vue";
import TaskOverview from "@/components/business/TaskOverview.vue";
import TableTrainStatus from "./components/TableTrainStatus.vue";
import DialogLog from "@/components/ui/log/DialogLog.vue";
import { useRouter } from "vue-router";

import { http, format, url } from "@/utils";

import { IconCircleCheckFilled, IconLoader, IconAlarm } from "@tabler/icons-vue";

const taskDetailConfig = [
  {
    statusText: "等待中",
    valueText: "个任务",
    key: "waitingJobCount",
    color: "warning",
    bgColor: "lightwarning",
    icon: IconLoader
  },
  {
    statusText: "已完成",
    valueText: "个微调任务",
    key: "successJobCount",
    color: "success",
    bgColor: "lightsuccess",
    icon: IconCircleCheckFilled
  },
  {
    statusText: "总微调时间",
    valueText: "",
    key: "totalDurationCount",
    color: "secondary",
    bgColor: "lightprimary",
    icon: IconAlarm
  }
];

const router = useRouter();

// defineOptions({
//   name: "FineTuningList"
// });

const page = ref({ title: "微调任务" });
const breadcrumbs = ref([
  {
    text: "模型微调",
    disabled: false,
    href: "#"
  },
  {
    text: "微调任务",
    disabled: true,
    href: "#"
  }
]);
const searchData = reactive({
  fineTunedModel: "",
  trainStatus: null
});
const tableInfos = reactive({
  list: [],
  total: 0
});
const createTaskPaneRef = ref();
const tableWithPagerRef = ref();
const refDialogLog = ref();
const refConfirmAbort = ref();
const confirmAbort = reactive({
  currentJobId: ""
});
const refConfirmDelete = ref();
const confirmDeleteId = ref("");
const currentJobId = ref("");
const refreshButtonRef = ref();

const getButtons = row => {
  let ret = [];
  const status = row.trainStatus;
  ret.push({
    text: "查看",
    click() {
      onView(row);
    }
  });
  if (status === "running" || status === "waiting") {
    ret.push({
      text: "取消",
      color: "error",
      click() {
        onAbort(row);
      }
    });
  } else if (status === "failed" || status === "cancel") {
    ret.push({
      text: "删除",
      color: "error",
      click() {
        onDelete(row);
      }
    });
  }
  if (status === "running") {
    ret.push({
      text: "终端",
      color: "info",
      click() {
        url.onNewPage(`/model/terminal?resourceType=ft-job&serviceName=${row.fineTunedModel}`);
      }
    });
  }
  return ret;
};

const doTableQuery = async (options = {}) => {
  const [err, res] = await http.get({
    url: "/finetuning",
    showLoading: tableWithPagerRef.value.el,
    data: {
      ...searchData,
      ...options
    }
  });
  if (res) {
    tableInfos.list = res.list || [];
    tableInfos.total = res.total;
  } else {
    tableInfos.list = [];
    tableInfos.total = 0;
  }
  refreshButtonRef.value.start();
};

const doQueryFirstPage = () => {
  tableWithPagerRef.value.query({ page: 1 });
};

const doQueryCurrentPage = () => {
  tableWithPagerRef.value.query();
};

const onLog = async ({ jobId }) => {
  currentJobId.value = jobId;
  refDialogLog.value.show();
};

const getLog = async () => {
  let [err, res] = await http.get({
    url: `/api/finetuning/${currentJobId.value}`
  });
  if (res) {
    refDialogLog.value.setContent(res.trainLog);
  }
};

const onAbort = row => {
  confirmAbort.currentJobId = row.jobId;
  refConfirmAbort.value.show({
    width: "450px",
    confirmText: confirmAbort.currentJobId
  });
};

const doAbort = async (options = {}) => {
  const [err, res] = await http.delete({
    ...options,
    showSuccess: true,
    url: `/finetuning/${confirmAbort.currentJobId}/cancel`
  });
  if (res) {
    refConfirmAbort.value.hide();
    tableWithPagerRef.value.query();
  }
};

const onView = ({ jobId }) => {
  router.push(`/model/fine-tuning/detail?jobId=${jobId}`);
};

const onDelete = row => {
  confirmDeleteId.value = row.jobId;
  refConfirmDelete.value.show({
    width: "400px"
  });
};

const doDelete = async (options = {}) => {
  const [err, res] = await http.delete({
    ...options,
    showSuccess: true,
    url: `/api/finetuning/${confirmDeleteId.value}`
  });
  if (res) {
    refConfirmDelete.value.hide();
    doQueryFirstPage();
  }
};

const onAdd = () => {
  createTaskPaneRef.value.show({
    title: "创建微调任务",
    operateType: "add"
  });
};

onMounted(() => {
  doTableQuery();
});
</script>
<style lang="scss"></style>
