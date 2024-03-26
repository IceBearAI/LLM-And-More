<template>
  <div class="">
    <v-row>
      <v-col cols="12" class="d-flex justify-space-between align-center" style="height: 76px">
        <div class="d-flex">
          <div style="width: 250px" class="mr-4">
            <Select
              @change="doQueryFirstPage"
              label="请选择评估状态"
              :mapDictionary="{ code: 'model_eval_status' }"
              v-model="searchData.status"
            ></Select>
          </div>
          <div style="width: 250px">
            <Select
              @change="doQueryFirstPage"
              label="请选择指标"
              :mapDictionary="{ code: 'model_evaluate_target_type' }"
              v-model="searchData.evalTargetType"
            ></Select>
          </div>
        </div>
        <el-tooltip
          ref="tooltipRef"
          :visible="state.showTooltip"
          :popper-options="{
            modifiers: [
              {
                name: 'computeStyles',
                options: {
                  adaptive: false,
                  enabled: false
                }
              }
            ]
          }"
          :auto-close="1"
          :virtual-ref="buttonRef"
          virtual-triggering
          popper-class="singleton-tooltip"
        >
          <template #content>
            <span>请先部署模型 </span>
          </template>
        </el-tooltip>
        <ButtonsInForm>
          <v-btn color="primary" @click="onAdd">创建评估</v-btn>
        </ButtonsInForm>
      </v-col>
      <v-col cols="12">
        <TableWithPager @query="doTableQuery" ref="refTableWithPager" :infos="state.tableInfos">
          <el-table-column label="指标" min-width="100px">
            <template #default="{ row }">{{ getLabels([["model_evaluate_target_type", row.evalTargetType]]) }} </template>
          </el-table-column>
          <el-table-column label="评估状态" min-width="120px">
            <template #default="{ row }">
              <el-tooltip :disabled="row.status !== 'failed'" :content="row.statusMsg" placement="top" raw-content>
                <v-chip label size="small" :color="map.statusMap[row.status].color">{{
                  getLabels([["model_eval_status", row.status]])
                }}</v-chip>
              </el-tooltip>
            </template>
          </el-table-column>
          <el-table-column label="数据量" min-width="100px">
            <template #default="{ row }">{{ row.dataSize }}</template>
          </el-table-column>
          <el-table-column label="评估数据集" min-width="120px">
            <template #default="{ row }"> {{ getLabels([["model_evaluate_data_type", row.dataType]]) }}</template>
          </el-table-column>
          <el-table-column label="平均分" min-width="100px">
            <template #default="{ row }">
              <template v-if="row.evalTargetType === 'five' && row.score > 0">
                <el-tooltip content="点击可查看五维图指标详情" placement="top">
                  <span class="link" @click="onViewDetail(row)">{{ row.score }}</span>
                </el-tooltip>
              </template>
              <template v-else>{{ row.score }}</template>
            </template>
          </el-table-column>
          <el-table-column label="备注" min-width="200px" show-overflow-tooltip>
            <template #default="{ row }">{{ row.remark }}</template>
          </el-table-column>
          <el-table-column label="创建时间" min-width="160px">
            <template #default="{ row }">{{ format.dateFormat(row.createdAt, "YYYY-MM-DD HH:mm:ss") }}</template>
          </el-table-column>
          <el-table-column label="操作" width="80px" fixed="right">
            <template #default="{ row }">
              <ButtonsInTable :buttons="getButtons(row)" onlyOne />
            </template>
          </el-table-column>
        </TableWithPager>
      </v-col>
    </v-row>
  </div>
  <ConfirmByInput ref="refConfirmAbort" @submit="doAbort">
    <template #text>
      此操作将会<span class="text-primary">取消</span>正在进行的模型评估<br />
      任务ID：<span class="text-primary font-weight-black">{{ state.currentJobId }}</span
      ><br />
      你还要继续吗？
    </template>
  </ConfirmByInput>
  <ConfirmByClick ref="refConfirmDelete" @submit="doRemove">
    <template #text>确定要删除该评估数据？</template>
  </ConfirmByClick>
  <CreatePerformanceEvalPane ref="createPerformanceEvalPaneRef" @submit="doQueryFirstPage" />
  <ViewEvalDetailPane ref="viewEvalDetailPaneRef" />
</template>
<script setup lang="ts">
import { reactive, toRefs, ref, onMounted } from "vue";
import { http, format, map } from "@/utils";
import { TypeButtonsInTable } from "@/components/types/components.ts";
import CreatePerformanceEvalPane from "./CreatePerformanceEvalPane.vue";
import ViewEvalDetailPane from "./ViewEvalDetailPane.vue";
import ConfirmByInput from "@/components/business/ConfirmByInput.vue";
import ConfirmByClick from "@/components/business/ConfirmByClick.vue";
import { useMapRemoteStore } from "@/stores";
import { useRoute } from "vue-router";
interface Props {
  /** 音频地址 */
  showArrange: string;
  modelTitle: string;
}
const props = defineProps({
  showArrange: String,
  modelTitle: String
});
const route = useRoute();
const { modelName, jobId } = route.query;
const { loadDictTree, getLabels } = useMapRemoteStore();
const refTableWithPager = ref();
const createPerformanceEvalPaneRef = ref();
const viewEvalDetailPaneRef = ref();
const buttonRef = ref();
const refConfirmAbort = ref();
const refConfirmDelete = ref();
const state = reactive({
  style: {},
  formData: {
    modelName: "",
    evalTargetType: ""
  },
  showTooltip: false,
  timer: null,
  tableInfos: {
    list: [],
    total: 0
  },
  currentJobId: ""
});
const searchData = reactive({
  status: null,
  evalTargetType: null
});
const confirmDelete = reactive({
  uuid: null
});
const { formData, showTooltip } = toRefs(state);

// const onAdd = () => {
//   refPaneModelEstimate.value.show({
//     title: "添加评估",
//     operateType: "add"
//   });
// };
const onAdd = () => {
  createPerformanceEvalPaneRef.value.show({
    title: "性能评估",
    operateType: "add"
  });
};

const onViewDetail = row => {
  viewEvalDetailPaneRef.value.show({
    title: "详情",
    id: row.id,
    modelId: jobId
  });
};

const onAbort = row => {
  state.currentJobId = row.uuid;
  refConfirmAbort.value.show({
    width: "450px",
    confirmText: state.currentJobId
  });
};
const remove = data => {
  confirmDelete.uuid = data.uuid;
  refConfirmDelete.value.show({
    width: "450px"
  });
};

const doRemove = async (options = {}) => {
  const [err, res] = await http.delete({
    ...options,
    showSuccess: true,
    url: `/evaluate/delete/${confirmDelete.uuid}`
  });
  if (res) {
    refConfirmDelete.value.hide();
    doQueryCurrentPage();
  }
};

const doAbort = async (options = {}) => {
  const [err, res] = await http.put({
    ...options,
    showSuccess: true,
    url: `/evaluate/cancel/${state.currentJobId}`
  });
  if (res) {
    refConfirmAbort.value.hide();
    doQueryCurrentPage();
  }
};

const doTableQuery = async (options = {}) => {
  await loadDictTree(["model_evaluate_data_type"]);
  const [err, res] = await http.get({
    url: `/evaluate/list`,
    showLoading: refTableWithPager.value.el,
    data: {
      modelId: jobId,
      ...searchData,
      ...options
    }
  });

  if (res) {
    state.tableInfos.total = res.total;
    state.tableInfos.list = res.list || [];
  } else {
    state.tableInfos.list = [];
    state.tableInfos.total = 0;
  }
};
const doQueryFirstPage = () => {
  refTableWithPager.value.query({ page: 1 });
};

const doQueryCurrentPage = () => {
  refTableWithPager.value.query();
};

const getButtons = (row): Array<TypeButtonsInTable> => {
  let ret = [];
  if (row.status == "waiting" || row.status == "running") {
    ret.push({
      text: "取消",
      color: "",
      click() {
        onAbort(row);
      }
    });
  } else {
    ret.push({
      text: "删除",
      color: "error",
      click() {
        remove(row);
      }
    });
  }
  return ret;
};

onMounted(() => {
  doTableQuery();
});
</script>
<style lang="scss">
.v-tooltip__content {
  pointer-events: initial;
}
</style>
