<template>
  <BaseBreadcrumb :title="page.title" :breadcrumbs="breadcrumbs"></BaseBreadcrumb>
  <v-row>
    <v-col>
      <UiParentCard>
        <v-row>
          <v-col cols="12" lg="3" md="4" sm="6">
            <v-text-field
              density="compact"
              v-model="formData.modelName"
              label="搜索模型名称"
              hide-details
              clearable
              variant="outlined"
              color="red"
              @keyup.enter="doQueryFirstPage"
              @click:clear="doQueryFirstPage"
            ></v-text-field>
          </v-col>
          <!-- <v-col cols="12" lg="2" md="2" sm="6">
            <Select
              @change="doQueryFirstPage"
              label="是否私有"
              :mapDictionary="{ code: 'boolean' }"
              v-model="formData.isPrivate"
            ></Select>
          </v-col> -->
          <v-col cols="12" lg="2" md="4" sm="6">
            <Select
              @change="doQueryFirstPage"
              label="请选择状态"
              :mapDictionary="{ code: 'local_enabled' }"
              v-model="formData.enabled"
            ></Select>
          </v-col>
          <v-col cols="12" lg="2" md="4" sm="6">
            <Select
              @change="doQueryFirstPage"
              label="请选择模型类型"
              :mapDictionary="{ code: 'model_type' }"
              v-model="formData.modelType"
            ></Select>
          </v-col>

          <v-col cols="12" lg="2" md="4" sm="6">
            <Select
              @change="doQueryFirstPage"
              label="请选择供应"
              :mapDictionary="{ code: 'model_provider_name' }"
              v-model="formData.providerName"
            ></Select>
          </v-col>
          <v-col cols="12" lg="3" md="4" sm="6">
            <ButtonsInForm>
              <v-btn color="primary" @click="onChat">聊天操场</v-btn>
              <v-btn color="primary" @click="onAdd">添加模型</v-btn>
              <refresh-button ref="refreshButtonRef" @refresh="doQueryCurrentPage" />
            </ButtonsInForm>
          </v-col>

          <v-col cols="12">
            <AlertBlock> 修改之后将实时生效，请谨慎操作！ </AlertBlock>
          </v-col>
          <v-col cols="12">
            <TableWithPager @query="doTableQuery" ref="refTableWithPager" :infos="state.tableInfos">
              <el-table-column label="模型" width="200px" fixed="left">
                <template #default="{ row }">
                  <el-tooltip content="查看详情" placement="top">
                    <span class="link" @click="onView(row)">{{ row.modelName }}</span>
                  </el-tooltip>
                </template>
              </el-table-column>
              <el-table-column label="最长上下文" width="110px">
                <template #default="{ row }">
                  <span>{{ row.maxTokens }}</span>
                </template>
              </el-table-column>
              <el-table-column label="模型类型" width="110px">
                <template #default="{ row }">
                  <span>{{ getLabels([["model_type", row.modelType]]) }}</span>
                </template>
              </el-table-column>
              <el-table-column label="是否启用" min-width="100px">
                <template #default="{ row }">
                  <ChipStatus v-model="row.enabled"></ChipStatus>
                </template>
              </el-table-column>
              <!-- <el-table-column label="是否私有" min-width="100px">
                <template #default="{ row }">
                  <ChipBoolean v-model="row.isPrivate"></ChipBoolean>
                </template>
              </el-table-column> -->
              <el-table-column label="是否微调" min-width="100px">
                <template #default="{ row }">
                  <ChipBoolean v-model="row.isFineTuning"></ChipBoolean>
                </template>
              </el-table-column>
              <el-table-column label="部署" min-width="150px">
                <template #default="{ row }">
                  <div v-if="row.deployStatus == 'pending'" style="width: 80px" class="mx-auto mb-4">
                    <v-progress-linear height="24" striped indeterminate color="#4eb879" rounded class="mx-auto">
                      <span class="text-body-2 text-lightText">部署中</span>
                    </v-progress-linear>
                  </div>
                  <!-- <div v-else>
                    {{ row.deployStatus }}
                  </div> -->

                  <div v-if="row.enabled && row.modelType === 'text-generation'">
                    <v-btn @click="onChat(row)" prepend-icon="mdi-chat" size="small" variant="outlined" color="info">
                      <template v-slot:prepend>
                        <IconMessages :size="16" />
                      </template>
                      对话
                    </v-btn>
                  </div>
                </template>
              </el-table-column>
              <el-table-column label="备注" min-width="200px" show-overflow-tooltip>
                <template #default="{ row }"> {{ row.remark }} </template>
              </el-table-column>
              <el-table-column label="更新时间" min-width="200px">
                <template #default="{ row }">
                  {{ format.dateFormat(row.updatedAt, "YYYY-MM-DD HH:mm:ss") }}
                </template>
              </el-table-column>
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

  <PaneModel ref="refPaneModel" @submit="doQueryFirstPage" />
  <ArrangeModel ref="refArrangeModel" @submit="doQueryFirstPage" />

  <ConfirmByInput ref="refConfirmDelete" @submit="doDelete">
    <template #text>
      您将要删除 <span class="text-primary font-weight-black">{{ state.selectedInfo.modelName }}</span> 模型，删除之后，<br />
      若有使用该模型的应用场景都将取消授权。 <br />确定要继续吗？
    </template>
  </ConfirmByInput>

  <ConfirmByClick ref="refConfirmByClick" @submit="onConfirmByClick">
    <template #text> <div v-html="state.confirmByClickInfo.html"></div></template>
  </ConfirmByClick>
</template>
<script setup lang="ts">
import { reactive, toRefs, ref, onMounted } from "vue";
import BaseBreadcrumb from "@/components/shared/BaseBreadcrumb.vue";
import UiParentCard from "@/components/shared/UiParentCard.vue";
import AlertBlock from "@/components/ui/AlertBlock.vue";
import ChipBoolean from "@/components/ui/ChipBoolean.vue";
import ChipStatus from "@/components/ui/ChipStatus.vue";

import PaneModel from "./components/PaneModel.vue";
import ArrangeModel from "./components/ArrangeModel.vue";
import ConfirmByInput from "@/components/business/ConfirmByInput.vue";
import ConfirmByClick from "@/components/business/ConfirmByClick.vue";

import { http, format } from "@/utils";
import { useRouter, useRoute } from "vue-router";
import { IconMessages } from "@tabler/icons-vue";
import _ from "lodash";

import { type ItfModel } from "./types/modelList.ts";
import { TypeButtonsInTable } from "@/components/types/components.ts";
import { useMapRemoteStore } from "@/stores";

const { loadDictTree, getLabels } = useMapRemoteStore();

const router = useRouter();
const route = useRoute();

const state = reactive<{
  /** 所选行，编辑操作 */
  selectedInfo: Partial<ItfModel>;
  [x: string]: any;
}>({
  style: {},
  formData: {
    modelName: "",
    // isPrivate: null,
    enabled: null,
    providerName: null,
    modelType: null
  },
  selectedInfo: {
    id: "",
    modelName: ""
  },
  tableInfos: {
    list: [],
    total: ""
  },
  confirmByClickInfo: {
    html: "",
    action: "",
    row: {}
  }
});
const { style, formData } = toRefs(state);

const refPaneModel = ref();
const refArrangeModel = ref();
const refConfirmDelete = ref();
const refTableWithPager = ref();
const refConfirmByClick = ref();
const refreshButtonRef = ref();

const page = ref({ title: "模型列表" });
const breadcrumbs = ref([
  {
    text: "模型管理",
    disabled: false,
    href: "#"
  },
  {
    text: "模型列表",
    disabled: true,
    href: "#"
  }
]);

const onChat = ({ id, modelName }: ItfModel = {} as ItfModel) => {
  if (modelName) {
    router.push("/model/chat-playground?modelName=" + modelName + "&modelId=" + id);
  } else {
    router.push("/model/chat-playground");
  }
};

const onView = ({ id, modelName }) => {
  router.push(`/model/model-list/detail?jobId=${id}&modelName=` + modelName);
};

const getButtons = (row): Array<TypeButtonsInTable> => {
  let ret: Array<TypeButtonsInTable> = [];
  for (let item of row.operation) {
    if (item == "edit") {
      ret.push({
        text: "编辑",
        color: "info",
        click() {
          onEdit(row);
        }
      });
    } else if (item == "delete") {
      ret.push({
        text: "删除",
        color: "error",
        click() {
          onDelete(row);
        }
      });
    } else if (item == "deploy") {
      ret.push({
        text: "部署",
        color: "info",
        click() {
          refArrangeModel.value.show({
            title: "部署模型",
            infos: row
          });
          // _.assign(state.confirmByClickInfo, {
          //   html: "确认部署模型 <span class='text-primary mx-1  font-weight-black'>" + row.modelName + "</span>吗？",
          //   action: "deploy",
          //   row: row
          // });
          // refConfirmByClick.value.show({ width: "350px" });
        }
      });
    } else if (item == "undeploy") {
      ret.push({
        text: "卸载",
        color: "info",
        click() {
          _.assign(state.confirmByClickInfo, {
            html: "确认卸载模型 <span class='text-primary mx-1  font-weight-black'>" + row.modelName + "</span>吗？",
            action: "undeploy",
            row: row
          });
          refConfirmByClick.value.show({ width: "350px" });
        }
      });
    }
  }
  return ret;
};

const onConfirmByClick = (options = {}) => {
  let { action, row } = state.confirmByClickInfo;
  if (action == "deploy") {
    onDeploy(row, options);
  } else if (action == "undeploy") {
    onUndeploy(row, options);
  }
};

const doTableQuery = async (options = {}) => {
  await loadDictTree(["model_type"]);
  const [err, res] = await http.get({
    url: "/models",
    showLoading: refTableWithPager.value.el,
    data: {
      ...state.formData,
      ...options
    }
  });
  if (res) {
    state.tableInfos.list = res.list;
    state.tableInfos.total = res.total;
  } else {
    state.tableInfos.list = [];
    state.tableInfos.total = 0;
  }
  refreshButtonRef.value.start();
};

const doQueryFirstPage = () => {
  refTableWithPager.value.query({ page: 1 });
};

const doQueryCurrentPage = () => {
  refTableWithPager.value.query();
};

const onDelete = info => {
  state.selectedInfo = info;
  refConfirmDelete.value.show({
    width: "400px",
    confirmText: state.selectedInfo.modelName
  });
};

const doDelete = async (options = {}) => {
  const [err, res] = await http.delete({
    ...options,
    showSuccess: true,
    url: `/api/models/${state.selectedInfo.id}`
  });
  if (res) {
    refConfirmDelete.value.hide();
    doQueryFirstPage();
  }
};

const onAdd = () => {
  refPaneModel.value.show({
    title: "添加模型",
    operateType: "add"
  });
};

const onEdit = info => {
  refPaneModel.value.show({
    title: "编辑模型",
    infos: info,
    operateType: "edit"
  });
};

const onDeploy = async (row, options) => {
  let [err, res] = await http.post({
    ...options,
    showSuccess: true,
    url: `/api/models/${row.id}/deploy`
  });
  if (res) {
    //关闭确认框
    refConfirmByClick.value.hide();
    doQueryCurrentPage();
  }
};

const onUndeploy = async (row, options) => {
  let [err, res] = await http.post({
    ...options,
    showSuccess: true,
    url: `/api/models/${row.id}/undeploy`
  });
  if (res) {
    //关闭确认框
    refConfirmByClick.value.hide();
    doQueryCurrentPage();
  }
};

onMounted(() => {
  let { modelName } = route.query;
  if (modelName) {
    state.formData.modelName = modelName;
  }
  doQueryFirstPage();
});
</script>
<style lang="scss"></style>
./types/modelList.type.ts ./types/modelList.ts
