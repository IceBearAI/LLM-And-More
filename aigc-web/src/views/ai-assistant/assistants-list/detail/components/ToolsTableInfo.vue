<template>
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
        <v-btn color="primary" @click="onAdd">添加工具</v-btn>
      </ButtonsInForm>
    </v-col>
    <v-col cols="12">
      <TableWithPager @query="doTableQuery" ref="tableWithPagerRef" :infos="tableInfos">
        <el-table-column label="工具ID" prop="toolId" min-width="220px"></el-table-column>
        <el-table-column label="工具名称" prop="name" min-width="200px" show-overflow-tooltip></el-table-column>
        <el-table-column label="工具描述" prop="description" min-width="300px" show-overflow-tooltip></el-table-column>
        <el-table-column label="工具类型" width="100px">
          <template #default="{ row }">
            {{ getLabels([["assistant_tool_type", row.toolType]]) }}
          </template>
        </el-table-column>
        <el-table-column label="备注" prop="remark" min-width="200px"></el-table-column>
        <el-table-column label="操作" width="80px" fixed="right">
          <template #default="{ row }">
            <ButtonsInTable :buttons="getButtons(row)" only-one />
          </template>
        </el-table-column>
      </TableWithPager>
    </v-col>
  </v-row>
  <ConfirmByInput ref="refConfirmDelete" @submit="doDelete">
    <template #text>
      您将要删除 <span class="text-primary font-weight-black">{{ confirmDelete.name }}</span> 工具<br />
      确定要继续吗？
    </template>
  </ConfirmByInput>
  <Dialog ref="refDialogTools" persistent max-width="1200px" :retain-focus="false">
    <template #title> 选择工具 </template>
    <ToolsList @selected="handleToolsSelect" />
  </Dialog>
</template>
<script setup lang="ts">
import { reactive, ref, onMounted } from "vue";
import ConfirmByInput from "@/components/business/ConfirmByInput.vue";
import ToolsList from "./ToolsList.vue";

import { http } from "@/utils";
import { useMapRemoteStore } from "@/stores";
import { useRoute } from "vue-router";

const { getLabels, loadDictTree } = useMapRemoteStore();
loadDictTree(["assistant_tool_type"]);

const route = useRoute();
const { assistantId } = route.query;

const searchData = reactive({
  name: ""
});
const tableInfos = reactive({
  list: [],
  total: 0
});
const tableWithPagerRef = ref();
const refConfirmDelete = ref();
const confirmDelete = reactive({
  name: "",
  toolId: ""
});
const refDialogTools = ref();

const getButtons = row => {
  let ret = [];
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
    url: `/assistants/${assistantId}/tools/${confirmDelete.toolId}`
  });
  if (res) {
    refConfirmDelete.value.hide();
    doTableQuery();
  }
};

const doTableQuery = async (options = {}) => {
  const [err, res] = await http.get({
    url: `/assistants/${assistantId}/tools`,
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
  refDialogTools.value.show({
    showActions: false
  });
};

const handleToolsSelect = async row => {
  const [err, res] = await http.post({
    showSuccess: true,
    url: `/assistants/${assistantId}/tools`,
    data: {
      toolId: row.toolId
    }
  });
  if (res) {
    refDialogTools.value.hide();
    doTableQuery();
  }
};

onMounted(() => {
  doTableQuery();
});
</script>
