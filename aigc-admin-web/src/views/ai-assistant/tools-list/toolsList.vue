<template>
  <BaseBreadcrumb :title="page.title" :breadcrumbs="breadcrumbs"></BaseBreadcrumb>
  <UiParentCard>
    <v-row>
      <v-col cols="12" lg="3" md="4" sm="6">
        <v-text-field
          v-model="searchData.name"
          label="请输入工具名称"
          hide-details
          clearable
          @keyup.enter="doQueryFirstPage"
          @click:clear="doQueryFirstPage"
        ></v-text-field>
      </v-col>
      <v-col cols="12" lg="3" md="4" sm="6">
        <ButtonsInForm>
          <v-btn color="primary" @click="onAdd">创建工具</v-btn>
        </ButtonsInForm>
      </v-col>
      <v-col cols="12">
        <AlertBlock>修改之后将实时生效，请谨慎操作！</AlertBlock>
      </v-col>

      <v-col cols="12">
        <TableWithPager @query="doTableQuery" ref="tableWithPagerRef" :infos="tableInfos">
          <el-table-column label="工具ID" min-width="220px">
            <template #default="{ row }">
              <a href="javascript: void(0)" class="link" @click="onEdit(row)">{{ row.toolId }}</a>
            </template>
          </el-table-column>
          <el-table-column label="工具名称" min-width="200px" show-overflow-tooltip class-name="link-ellipsis-color">
            <template #default="{ row }">
              <a href="javascript: void(0)" class="link" @click="onEdit(row)">{{ row.name }}</a>
            </template>
          </el-table-column>
          <el-table-column label="工具描述" prop="description" min-width="300px" show-overflow-tooltip></el-table-column>
          <el-table-column label="工具类型" width="100px">
            <template #default="{ row }">
              {{ getLabels([["assistant_tool_type", row.toolType]]) }}
            </template>
          </el-table-column>
          <el-table-column label="被助手使用" prop="assistant" min-width="110px">
            <template #default="{ row }">
              <template v-if="row.assistants.length > 0">
                <el-tooltip :content="getAssistantContent(row.assistants)" placement="top" raw-content>
                  {{ row.assistants.length }}
                </el-tooltip>
              </template>
              <template v-else> 0 </template>
            </template>
          </el-table-column>
          <el-table-column label="备注" prop="remark" min-width="200px"></el-table-column>
          <el-table-column label="更新时间" min-width="165px">
            <template #default="{ row }">
              {{ format.dateFormat(row.updatedAt, "YYYY-MM-DD HH:mm:ss") }}
            </template>
          </el-table-column>
          <el-table-column label="操作人" prop="operator" min-width="180px" show-overflow-tooltip></el-table-column>
          <el-table-column label="操作" width="120px" fixed="right">
            <template #default="{ row }">
              <ButtonsInTable :buttons="getButtons(row)" />
            </template>
          </el-table-column>
        </TableWithPager>
      </v-col>
    </v-row>
  </UiParentCard>
  <ConfirmByInput ref="refConfirmDelete" @submit="doDelete">
    <template #text>
      您将要删除 <span class="text-primary font-weight-black">{{ confirmDelete.name }}</span> 工具<br />
      确定要继续吗？
    </template>
  </ConfirmByInput>
  <CreateToolPane ref="createToolPaneRef" @submit="doQueryFirstPage" />
</template>
<script setup lang="ts">
import { reactive, ref, onMounted } from "vue";
import BaseBreadcrumb from "@/components/shared/BaseBreadcrumb.vue";
import UiParentCard from "@/components/shared/UiParentCard.vue";
import AlertBlock from "@/components/ui/AlertBlock.vue";
import ConfirmByInput from "@/components/business/ConfirmByInput.vue";
import CreateToolPane from "./components/CreateToolPane.vue";

import { http, format } from "@/utils";
import { useMapRemoteStore } from "@/stores";

const { getLabels, loadDictTree } = useMapRemoteStore();
loadDictTree(["assistant_tool_type"]);

const page = ref({ title: "工具列表" });
const breadcrumbs = ref([
  {
    text: "AI助手",
    disabled: false,
    href: "#"
  },
  {
    text: "工具列表",
    disabled: true,
    href: "#"
  }
]);
const searchData = reactive({
  name: ""
});
const tableInfos = reactive({
  list: [],
  total: 0
});
const createToolPaneRef = ref();
const tableWithPagerRef = ref();
const refConfirmDelete = ref();
const confirmDelete = reactive({
  name: "",
  toolId: ""
});

const getButtons = row => {
  let ret = [];
  ret.push({
    text: "编辑",
    color: "info",
    click() {
      onEdit(row);
    }
  });
  ret.push({
    text: "删除",
    color: "error",
    click() {
      onDelete(row);
    }
  });
  return ret;
};

const onDelete = row => {
  confirmDelete.toolId = row.toolId;
  confirmDelete.name = row.name;
  refConfirmDelete.value.show({
    width: "450px",
    confirmText: confirmDelete.name
  });
};

const doDelete = async (options = {}) => {
  const [err, res] = await http.delete({
    ...options,
    showSuccess: true,
    url: `/tools/${confirmDelete.toolId}`
  });
  if (res) {
    refConfirmDelete.value.hide();
    doTableQuery();
  }
};

const doTableQuery = async (options = {}) => {
  const [err, res] = await http.get({
    url: "/tools/list",
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
};

const doQueryFirstPage = () => {
  tableWithPagerRef.value.query({ page: 1 });
};

const onAdd = () => {
  createToolPaneRef.value.show({
    title: "创建工具",
    operateType: "add"
  });
};

const onEdit = info => {
  createToolPaneRef.value.show({
    title: "编辑工具",
    infos: info,
    operateType: "edit"
  });
};

const getAssistantContent = list => {
  return list.map(item => item.name).join("<br/>");
};

onMounted(() => {
  doTableQuery();
});
</script>
