<template>
  <BaseBreadcrumb :title="page.title" :breadcrumbs="breadcrumbs"></BaseBreadcrumb>
  <UiParentCard>
    <v-row>
      <v-col cols="12" lg="3" md="4" sm="6">
        <v-text-field
          v-model="searchData.name"
          label="请输入助手名称"
          hide-details
          clearable
          @keyup.enter="doQueryFirstPage"
          @click:clear="doQueryFirstPage"
        ></v-text-field>
      </v-col>
      <v-col cols="12" lg="3" md="4" sm="6">
        <ButtonsInForm>
          <v-btn color="primary" @click="onAdd">创建助手</v-btn>
        </ButtonsInForm>
      </v-col>
      <v-col cols="12">
        <AlertBlock>修改之后将实时生效，请谨慎操作！</AlertBlock>
      </v-col>

      <v-col cols="12">
        <TableWithPager @query="doTableQuery" ref="tableWithPagerRef" :infos="tableInfos">
          <el-table-column label="助手ID" min-width="240px">
            <template #default="{ row }">
              <a href="javascript: void(0)" class="link" @click="onView(row.assistantId)">{{ row.assistantId }}</a>
            </template>
          </el-table-column>
          <el-table-column label="助手名称" min-width="150px" show-overflow-tooltip class-name="link-ellipsis-color">
            <template #default="{ row }">
              <a href="javascript: void(0)" class="link" @click="onView(row.assistantId)">{{ row.name }}</a>
            </template>
          </el-table-column>
          <el-table-column label="工具数量" min-width="100px">
            <template #default="{ row }">
              {{ row.tools ? row.tools.length : 0 }}
            </template>
          </el-table-column>
          <el-table-column label="模型" prop="modelName" min-width="200px" show-overflow-tooltip></el-table-column>
          <el-table-column label="备注" prop="remark" min-width="200px"></el-table-column>
          <el-table-column label="更新时间" min-width="165px">
            <template #default="{ row }">
              {{ format.dateFormat(row.updatedAt, "YYYY-MM-DD HH:mm:ss") }}
            </template>
          </el-table-column>
          <el-table-column label="操作人" prop="operator" min-width="150px" show-overflow-tooltip></el-table-column>
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
      此操作将会<span class="text-primary font-weight-black">删除</span>该个人助手，删除之后将无法使用<br />
      助手ID：<span class="text-primary font-weight-black">{{ confirmDelete.assistantId }}</span
      ><br />
      确定要继续吗？
    </template>
  </ConfirmByInput>
  <CreateAssistantPane ref="createAssistantPaneRef" @submit="doQueryFirstPage" />
</template>
<script setup lang="ts">
import { reactive, ref, onMounted } from "vue";
import BaseBreadcrumb from "@/components/shared/BaseBreadcrumb.vue";
import UiParentCard from "@/components/shared/UiParentCard.vue";
import AlertBlock from "@/components/ui/AlertBlock.vue";
import ConfirmByInput from "@/components/business/ConfirmByInput.vue";
import CreateAssistantPane from "./components/CreateAssistantPane.vue";

import { http, format } from "@/utils";
import { useRouter } from "vue-router";

const router = useRouter();

const page = ref({ title: "助手列表" });
const breadcrumbs = ref([
  {
    text: "AI助手",
    disabled: false,
    href: "#"
  },
  {
    text: "助手列表",
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
const createAssistantPaneRef = ref();
const tableWithPagerRef = ref();
const refConfirmDelete = ref();
const confirmDelete = reactive({
  assistantId: ""
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
  confirmDelete.assistantId = row.assistantId;
  refConfirmDelete.value.show({
    width: "550px",
    confirmText: confirmDelete.assistantId
  });
};

const doDelete = async (options = {}) => {
  const [err, res] = await http.delete({
    ...options,
    showSuccess: true,
    url: `/assistants/${confirmDelete.assistantId}`
  });
  if (res) {
    refConfirmDelete.value.hide();
    doTableQuery();
  }
};

const doTableQuery = async (options = {}) => {
  const [err, res] = await http.get({
    url: "/assistants/list",
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
  createAssistantPaneRef.value.show({
    title: "创建助手",
    operateType: "add"
  });
};

const onEdit = info => {
  createAssistantPaneRef.value.show({
    title: "编辑助手",
    infos: info,
    operateType: "edit"
  });
};

const onView = id => {
  router.push(`/ai-assistant/assistants/detail?assistantId=${id}`);
};

onMounted(() => {
  doTableQuery();
});
</script>
