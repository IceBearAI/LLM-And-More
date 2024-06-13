<template>
  <BaseBreadcrumb :title="page.title" :breadcrumbs="breadcrumbs"></BaseBreadcrumb>
  <UiParentCard>
    <v-row>
      <v-col cols="12" lg="3" md="4" sm="6">
        <v-text-field
          v-model="searchData.name"
          label="请输入名称"
          hide-details
          clearable
          @keyup.enter="doQueryFirstPage"
          @click:clear="doQueryFirstPage"
        ></v-text-field>
      </v-col>
      <v-col cols="12" lg="3" md="4" sm="6">
        <ButtonsInForm>
          <v-btn color="primary" @click="onAdd">添加租户</v-btn>
        </ButtonsInForm>
      </v-col>
      <v-col cols="12">
        <AlertBlock>修改之后将实时生效，请谨慎操作！</AlertBlock>
      </v-col>

      <v-col cols="12">
        <TableWithPager @query="doTableQuery" ref="tableWithPagerRef" :infos="tableInfos">
          <el-table-column label="名称" prop="name" min-width="200px" show-overflow-tooltip></el-table-column>
          <el-table-column label="联系人邮箱" prop="contactEmail" min-width="200px"></el-table-column>
          <el-table-column label="模型" min-width="200px" show-overflow-tooltip>
            <template #default="{ row }">
              {{ row.modelNames ? row.modelNames.join("，") : "--" }}
            </template>
          </el-table-column>
          <el-table-column label="创建时间" min-width="165px">
            <template #default="{ row }">
              {{ format.dateFormat(row.createdAt, "YYYY-MM-DD HH:mm:ss") }}
            </template>
          </el-table-column>
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
      您将要删除 <span class="text-primary font-weight-black">{{ confirmData.name }}</span> 租户<br />
      确定要继续吗？
    </template>
  </ConfirmByInput>
  <CreateTenantPane ref="createTenantPaneRef" @submit="doQueryFirstPage" />
</template>
<script setup lang="ts">
import { reactive, ref, onMounted } from "vue";
import BaseBreadcrumb from "@/components/shared/BaseBreadcrumb.vue";
import UiParentCard from "@/components/shared/UiParentCard.vue";
import AlertBlock from "@/components/ui/AlertBlock.vue";
import ConfirmByInput from "@/components/business/ConfirmByInput.vue";
import CreateTenantPane from "./components/CreateTenantPane.vue";

import { http, format } from "@/utils";

const page = ref({ title: "租户列表" });
const breadcrumbs = ref([
  {
    text: "系统管理",
    disabled: false,
    href: "#"
  },
  {
    text: "租户列表",
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
const createTenantPaneRef = ref();
const tableWithPagerRef = ref();
const refConfirmDelete = ref();
const confirmData = reactive({
  name: "",
  id: ""
});

const getButtons = row => {
  let ret = [];
  if (row.publicTenantId !== "5f9b3b3d-9b9c-4e1a-8e1a-5a4b4b4b4b4b") {
    ret.push({
      text: "删除",
      color: "error",
      click() {
        onDelete(row);
      }
    });
  }
  ret.push({
    text: "编辑",
    color: "info",
    click() {
      onEdit(row);
    }
  });
  return ret;
};

const onDelete = row => {
  confirmData.name = row.name;
  confirmData.id = row.id;
  refConfirmDelete.value.show({
    width: "450px",
    confirmText: confirmData.name
  });
};

const doDelete = async (options = {}) => {
  const [err, res] = await http.delete({
    ...options,
    showSuccess: true,
    url: `/tenants/${confirmData.id}`
  });
  if (res) {
    refConfirmDelete.value.hide();
    doQueryCurrentPage();
  }
};

const doTableQuery = async (options = {}) => {
  const [err, res] = await http.get({
    url: "/tenants",
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

const doQueryCurrentPage = () => {
  tableWithPagerRef.value.query();
};

const onAdd = () => {
  createTenantPaneRef.value.show({
    title: "添加租户",
    operateType: "add"
  });
};

const onEdit = info => {
  createTenantPaneRef.value.show({
    title: "编辑租户",
    infos: info,
    operateType: "edit"
  });
};

onMounted(() => {
  doTableQuery();
});
</script>
