<template>
  <BaseBreadcrumb :title="page.title" :breadcrumbs="breadcrumbs"></BaseBreadcrumb>

  <!--本期先去掉-->
  <!-- <UiParentCard class="mb-3">
    <TaskOverview :config="taskDetailConfig" request-url="/finetuning/dashboard" />
  </UiParentCard> -->
  <v-row>
    <v-col>
      <UiParentCard>
        <v-row>
          <!-- <v-col cols="12" lg="3" md="4" sm="6">
            <v-text-field
              density="compact"
              v-model="formData.name"
              label="请输入任务名称"
              hide-details
              clearable
              variant="outlined"
              color="red"
              @keyup.enter="doQueryFirstPage"
              @click:clear="doQueryFirstPage"
            ></v-text-field>
          </v-col> -->
          <v-col cols="12" class="flex justify-between">
            <div style="width: 300px">
              <Select
                @change="doQueryFirstPage"
                label="请选择标注状态"
                :mapDictionary="{ code: 'local_mark_status' }"
                v-model="formData.status"
              >
              </Select>
            </div>

            <ButtonsInForm>
              <v-btn color="primary" @click="onAdd()">创建文本标注任务</v-btn>
              <refresh-button ref="refreshButtonRef" @refresh="doQueryCurrentPage" />
            </ButtonsInForm>
          </v-col>
          <!-- <v-col cols="12" lg="6" md="4" sm="6" class="flex justify-end">

          </v-col> -->

          <v-col cols="12">
            <TableWithPager @query="doTableQuery" ref="refTableWithPager" :infos="state.tableInfos">
              <el-table-column label="任务名称" width="200px" fixed="left">
                <template #default="{ row }">
                  {{ row.name }}
                </template>
              </el-table-column>

              <el-table-column label="数据集名称" width="200px">
                <template #default="{ row }">
                  <el-tooltip content="查看数据集">
                    <a class="link" @click="onViewDataset(row.datasetName)">
                      {{ row.datasetName }}
                    </a>
                  </el-tooltip>
                </template>
              </el-table-column>
              <el-table-column label="数据序列" width="100px">
                <template #default="{ row }">
                  {{ row.dataSequence.join(" ~ ") }}
                </template>
              </el-table-column>
              <el-table-column label="任务类型" width="120px">
                <template #default="{ row }">
                  {{ getLabels([["textannotation_type", row.annotationType]]) }}
                </template>
              </el-table-column>
              <el-table-column width="120px">
                <template #header>完成进度</template>
                <template #default="{ row }">
                  <el-popover placement="right">
                    <div class="space-y-2">
                      <div>
                        已完成：<span class="link font-bold">{{ row.completed }}</span>
                      </div>
                      <div>
                        未完成：<span class="link font-bold">{{ row.total - row.completed - row.abandoned }}</span>
                      </div>
                      <div>
                        丢弃：<span class="link font-bold">{{ row.abandoned }}</span>
                      </div>
                    </div>
                    <template #reference
                      >{{ row.completed }} / {{ row.total - row.completed - row.abandoned }} / {{ row.abandoned }}</template
                    >
                  </el-popover>
                </template>
              </el-table-column>
              <el-table-column label="状态" width="120px">
                <template #default="{ row }">
                  <v-chip label size="small" :class="getStatusClassName(row.status)">{{
                    getLabels([["local_mark_status", row.status]])
                  }}</v-chip>
                </template>
              </el-table-column>
              <el-table-column label="检测状态" width="135px">
                <template #default="{ row }">
                  <template v-if="row.annotationType == 'faq'">
                    <div class="d-flex align-center justify-center">
                      <span :class="getStatusClassName(row.detectionStatus)">{{
                        getLabels([["local_mark_detect_status", row.detectionStatus]])
                      }}</span>
                      <div v-if="['processing', 'completed'].includes(row.detectionStatus)" class="link ml-1" @click="onLog(row)">
                        (日志)
                      </div>
                    </div>
                  </template>
                  <template v-else> -- </template>
                </template>
              </el-table-column>
              <el-table-column label="完成时间" min-width="180px">
                <template #default="{ row }">
                  {{ format.dateFormat(row.completedAt, "YYYY-MM-DD HH:mm:ss") }}
                </template>
              </el-table-column>
              <el-table-column label="创建时间" min-width="180px">
                <template #default="{ row }">
                  {{ format.dateFormat(row.createdAt, "YYYY-MM-DD HH:mm:ss") }}
                </template>
              </el-table-column>
              <el-table-column label="训练集数量" min-width="120px">
                <template #default="{ row }">
                  <el-popover placement="right">
                    <div class="space-y-2">
                      <div>
                        测试集数量：<span class="link font-bold">{{ row.testTotal }}</span>
                      </div>
                      <div>
                        训练集数量：<span class="link font-bold">{{ row.trainTotal }}</span>
                      </div>
                    </div>
                    <template #reference>{{ row.testTotal }} / {{ row.trainTotal }}</template>
                  </el-popover>
                </template>
              </el-table-column>
              <el-table-column label="负责人" width="200px">
                <template #default="{ row }">
                  <span v-copy="row.principal">{{ row.principal }}</span>
                </template>
              </el-table-column>
              <!-- <el-table-column label="操作人" min-width="200px">
                <template #default="{ row }"> {{ row.operator }} </template>
              </el-table-column> -->
              <el-table-column label="操作" min-width="120px" fixed="right">
                <template #default="{ row }">
                  <ButtonsInTable :buttons="getButtons(row)" />
                </template>
              </el-table-column>
            </TableWithPager>
          </v-col>
        </v-row>
      </UiParentCard>
    </v-col>
  </v-row>

  <PaneTextMarkTask ref="refPaneTextMarkTask" @submit="doQueryFirstPage" />

  <ConfirmByInput ref="refConfirmDelete" @submit="doDelete">
    <template #text>
      这是进行一项操作时必须了解的重要信息<br />
      此操作将会<span class="text-primary font-weight-black">删除</span>在进行的标注任务<br />
      数据集ID：<span class="text-primary font-weight-black">{{ state.selectedRow.uuid }}</span
      ><br />你还要继续吗？
    </template>
  </ConfirmByInput>

  <ConfirmByInput ref="refConfirmCancel" @submit="doCancel">
    <template #text>
      这是进行一项操作时必须了解的重要信息<br />
      此操作将会<span class="text-primary font-weight-black">取消</span>尚未开始的标注任务<br />
      数据集ID：<span class="text-primary font-weight-black">{{ state.selectedRow.uuid }}</span
      ><br />你还要继续吗？
    </template>
  </ConfirmByInput>

  <ConfirmByClick ref="refConfirmSplit" @submit="doSplit">
    <template #text>
      <div classs="text-slate-500">
        此操作将新建基于该数据集的测试集和训练集，<span class="text-slate-700 font-bold">原数据集不会改变。</span
        >数据集将会按照以下比例随机切分。
      </div>
      <div>
        <!-- testPercent :{{ state.selectedRow.testPercent }} -->
        <div class="flex justify-between mt-5 mb-2">
          <div>
            训练集 <span class="text-primary">{{ state.selectedRow.testPercent }}%</span>
          </div>
          <div>
            测试集 <span class="text-primary">{{ 100 - state.selectedRow.testPercent }}%</span>
          </div>
        </div>
        <v-slider color="primary" step="1" v-model="state.selectedRow.testPercent"></v-slider>
      </div>
    </template>
  </ConfirmByClick>

  <ConfirmByClick ref="refConfirmDownload" @submit="doDownload">
    <template #text>
      <AlertBlock class="mb-6"> 该操作将直接生成训练数据集，导出的数据将会以jsonl的格式导出。 </AlertBlock>
      <div class="hv-center">
        <div class="flex items-center">
          <div class="mr-4">导出格式：</div>
          <v-chip-group mandatory selected-class="text-primary" v-model="state.formatTypeSelectedIndex">
            <v-chip v-for="(item, index) in options.formatType" :key="index" filter variant="outlined">
              {{ item.label }}
            </v-chip>
          </v-chip-group>
        </div>
      </div>
    </template>
  </ConfirmByClick>

  <DialogLog ref="refDialogLog" :interval="30" @refresh="getLog" />

  <Dialog ref="refReport" style="width: 80%; max-width: 800px">
    <template #title>标注数据检测报告</template>
    <div class="box-report" v-html="state.reportHTML"></div>
  </Dialog>
</template>
<script setup lang="ts">
import { reactive, toRefs, ref, onMounted, inject } from "vue";
import BaseBreadcrumb from "@/components/shared/BaseBreadcrumb.vue";
import UiParentCard from "@/components/shared/UiParentCard.vue";
import PaneTextMarkTask from "./components/PaneTextMarkTask.vue";
import AlertBlock from "@/components/ui/AlertBlock.vue";

import ConfirmByInput from "@/components/business/ConfirmByInput.vue";
import ConfirmByClick from "@/components/business/ConfirmByClick.vue";
import { useMapRemoteStore } from "@/stores";
import { http, format } from "@/utils";
import { useRouter, useRoute } from "vue-router";
import _ from "lodash";
import TaskOverview from "@/components/business/TaskOverview.vue";
import { TypeButtonsInTable } from "@/components/types/components.ts";
import { IconCircleCheckFilled, IconLoader, IconAlarm } from "@tabler/icons-vue";
import { ItfAspectPageState } from "@/types/AspectPageTypes.ts";
import DialogLog from "@/components/ui/log/DialogLog.vue";

const provideAspectPage = inject("provideAspectPage") as ItfAspectPageState;
const { getLabels, loadDictTree } = useMapRemoteStore();
loadDictTree(["textannotation_type"]);

const router = useRouter();
const route = useRoute();

const refPaneTextMarkTask = ref();

const refConfirmDelete = ref();
const refConfirmCancel = ref();
const refTableWithPager = ref();
const refConfirmSplit = ref();
const refReport = ref();
const refConfirmDownload = ref();
const refDialogLog = ref();
const refreshButtonRef = ref();

const state = reactive<{
  [x: string]: any;
}>({
  style: {},
  formData: {
    name: "",
    status: null
  },
  tableInfos: {
    list: [],
    total: ""
  },
  selectedRow: {},
  reportHTML: "",
  formatTypeSelectedIndex: 0,
  options: {
    formatType: [
      { label: "默认", value: "default" },
      { label: "会话", value: "conversation" }
    ]
  }
});
const { style, formData, options } = toRefs(state);

const taskDetailConfig = [
  {
    statusText: "标注中",
    valueText: "个任务",
    key: "waitingJobCount",
    color: "warning",
    bgColor: "lightwarning",
    icon: IconLoader
  },
  {
    statusText: "已完成",
    valueText: "个任务",
    key: "successJobCount",
    color: "success",
    bgColor: "lightsuccess",
    icon: IconCircleCheckFilled
  }
];

const page = ref({ title: "文本标注列表" });
const breadcrumbs = ref([
  {
    text: "任务库",
    disabled: false,
    href: "#"
  },
  {
    text: "文本标注列表",
    disabled: true,
    href: "#"
  }
]);

const onViewDataset = datasetName => {
  router.push(`/sample-library/mgr-datasets/list?name=${datasetName}`);
};

const onView = id => {
  router.push(`/sample-library/text-mark/detail?annotationId=${id}`);
};

const doSplit = async () => {
  const [err, res] = await http.post({
    showSuccess: true,
    url: `/api/mgr/annotation/task/${state.selectedRow.uuid}/split`,
    data: {
      testPercent: state.selectedRow.testPercent / 100
    }
  });
  if (res) {
    refConfirmSplit.value.hide();
    doTableQuery();
  }
};

const getSectionStartHTML = title => {
  return `
   <div class="section">
        <div class="title">${title}</div>
        <div class="content">
   `;
};

const getSectionEndHTML = () => {
  return `
        </div>
      </div>
   `;
};

const genReportHtml = data => {
  let jsonData = null;
  try {
    jsonData = JSON.parse(data);
  } catch (error) {
    console.log(error);
  }
  if (!jsonData) return;
  let { mismatchedIntents, similarIntents } = jsonData;

  let partA = [],
    partB = [];

  if (similarIntents.length) {
    //相似意图
    partA.push(getSectionStartHTML("相似意图"));
    for (let item of similarIntents) {
      for (let itemChild of item.intentPair) {
        partA.push(`<div class="item"> ${itemChild}</div>`);
      }
    }
    partA.push(getSectionEndHTML());
  }
  if (mismatchedIntents.length) {
    //不匹配的意图
    partB.push(getSectionStartHTML("不匹配的意图"));
    for (let item of mismatchedIntents) {
      partB.push(`<div class="item"> `);
      if (item.answer1) {
        partB.push(`<div>问题1：${item.answer1}</div>`);
      }
      if (item.answer2) {
        partB.push(`<div>问题2：${item.answer2}</div>`);
      }
      if (item.intent1) {
        partB.push(`<div>意图1：${item.intent1}</div>`);
      }
      if (item.intent2) {
        partB.push(`<div>意图2：${item.intent2}</div>`);
      }
      if (item.lineNumbers.length) {
        partB.push(`<div>问题点：${item.lineNumbers.join("、")}</div>`);
      }
      partB.push(`</div>`);
    }
    partB.push(getSectionEndHTML());
  }
  let ret = [...partA, ...partB];

  // console.log("ret is ", ret.join(""));
  return ret.join("");
};

const doDownload = async (options = {}) => {
  const res = await http.download({
    ...options,
    url: `/api/mgr/annotation/task/${state.selectedRow.uuid}/export?formatType=${
      state.options.formatType[state.formatTypeSelectedIndex].value
    }`
  });
  if (res) {
    refConfirmDownload.value.hide();
  }
};

const getButtons = (row): Array<TypeButtonsInTable> => {
  let ret: Array<TypeButtonsInTable> = [];
  let { status } = row;

  if (["pending", "processing"].includes(row.status)) {
    ret.push({
      text: "标注",
      click() {
        onView(row.uuid);
      }
    });
  }

  if (["completed"].includes(status)) {
    ret.push({
      text: "导出",
      click() {
        state.selectedRow = row;
        refConfirmDownload.value.show({ width: "500px", title: "导出选择" });
      }
    });

    if (row.testTotal === 0) {
      //大于0 不显示切分
      ret.push({
        text: "切分",
        click() {
          state.selectedRow = row;
          let { testTotal, trainTotal } = row;
          state.selectedRow.testPercent = Math.floor((trainTotal * 100) / (trainTotal + testTotal));
          refConfirmSplit.value.show({
            width: "400px"
          });
        }
      });
    }

    if (row.annotationType == "faq") {
      if (row.detectionStatus === "pending" || row.detectionStatus === "canceled") {
        ret.push({
          text: "检测",
          async click() {
            const [err, res] = await http.post({
              timeout: Number.MAX_SAFE_INTEGER,
              showSuccess: true,
              showLoading: true,
              url: `/api/mgr/annotation/task/${row.uuid}/detect/annotation/async`
            });
            if (res) {
              doQueryCurrentPage();
            }
          }
        });
      }
      if (row.detectionStatus === "processing") {
        ret.push({
          text: "取消检测",
          async click() {
            const [err, res] = await http.post({
              timeout: Number.MAX_SAFE_INTEGER,
              showSuccess: true,
              showLoading: true,
              url: `/api/mgr/annotation/task/${row.uuid}/detect/cancel`
            });
            if (res) {
              doQueryCurrentPage();
            }
          }
        });
      }
      if (row.detectionStatus === "completed") {
        ret.push({
          text: "标注报告",
          click() {
            refReport.value.show({ showActions: false, contentHeight: "calc(100vh - 400px)" });
            state.reportHTML = genReportHtml(row.testReport);
          }
        });
      }
    }
  }
  if (["processing"].includes(status)) {
    ret.push({
      text: "取消",
      color: "error",
      click() {
        onCancel(row);
      }
    });
  }
  if (["completed", "cleaned", "abandoned", "pending"].includes(status)) {
    ret.push({
      text: "删除",
      color: "error",
      click() {
        onDelete(row);
      }
    });
  }
  return ret;
};

const doTableQuery = async (options = {}) => {
  const [err, res] = await http.get({
    url: "/api/mgr/annotation/task/list",
    showLoading: refTableWithPager.value.el,
    data: {
      ...state.formData,
      ...options
    }
  });
  if (res) {
    state.tableInfos.list = res.list || [];
    state.tableInfos.total = res.total;
  } else {
    state.tableInfos.list = [];
    state.tableInfos.total = 0;
  }
  refreshButtonRef.value.start();
};

provideAspectPage.methods.refreshListPage = () => {
  // 重新查询
  doTableQuery();
};

const doQueryFirstPage = () => {
  refTableWithPager.value.query({ page: 1 });
};

const doQueryCurrentPage = () => {
  refTableWithPager.value.query();
};

const onDelete = info => {
  state.selectedRow = info;
  refConfirmDelete.value.show({
    width: "600px",
    confirmText: state.selectedRow.uuid
  });
};

const doDelete = async (options = {}) => {
  const [err, res] = await http.delete({
    ...options,
    showSuccess: true,
    url: `/api/mgr/annotation/task/${state.selectedRow.uuid}/delete`
  });
  if (res) {
    refConfirmDelete.value.hide();
    doTableQuery();
  }
};

const getStatusClassName = status => {
  let map = {
    pending: "undo",
    processing: "doing",
    completed: "done",
    abandoned: "status-gray",
    cleaned: "status-gray",
    canceled: "status-gray"
  };
  return map[status];
};

const onCancel = info => {
  state.selectedRow = info;
  refConfirmCancel.value.show({
    width: "600px",
    confirmText: state.selectedRow.uuid
  });
};

const doCancel = async (options = {}) => {
  const [err, res] = await http.put({
    ...options,
    showSuccess: true,
    url: `/api/mgr/annotation/task/${state.selectedRow.uuid}/clean`
  });
  if (res) {
    refConfirmCancel.value.hide();
    doTableQuery();
  }
};

const onAdd = () => {
  refPaneTextMarkTask.value.show({
    title: "创建文本标注任务",
    operateType: "add"
  });
};

const onLog = async info => {
  state.selectedRow = info;
  refDialogLog.value.show();
};

const getLog = async () => {
  let [err, res] = await http.get({
    url: `/api/mgr/annotation/task/${state.selectedRow.uuid}/detect/annotation/log`
  });
  if (res) {
    const content = Object.keys(res).length === 0 ? "" : res;
    refDialogLog.value.setContent(content);
  }
};

onMounted(() => {
  doTableQuery();
});
</script>
<style lang="scss">
.box-report {
  @apply divide-y-2;
  .section {
    @apply p-2 shadow-sm;
  }
  .title {
    @apply text-gray-700 text-lg font-bold;
  }
  .content {
    @apply pl-4 space-y-4;
    .item {
      @apply p-4 shadow-md;
    }
  }
}
</style>
