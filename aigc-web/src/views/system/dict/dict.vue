<template>
  <BaseBreadcrumb :title="page.title" :breadcrumbs="breadcrumbs"></BaseBreadcrumb>
  <UiParentCard>
    <v-row>
      <v-col cols="12" lg="3" md="4" sm="6">
        <v-text-field
          v-model="searchData.code"
          label="请输入字典编号"
          hide-details
          clearable
          @keyup.enter="doQueryFirstPage"
          @click:clear="doQueryFirstPage"
        ></v-text-field>
      </v-col>
      <v-col cols="12" lg="3" md="4" sm="6">
        <v-text-field
          v-model="searchData.label"
          label="请输入字典名称"
          hide-details
          clearable
          @keyup.enter="doQueryFirstPage"
          @click:clear="doQueryFirstPage"
        ></v-text-field>
      </v-col>
      <v-col cols="12" lg="3" md="4" sm="6">
        <ButtonsInForm>
          <v-btn color="primary" @click="onAdd">添加字典</v-btn>
        </ButtonsInForm>
      </v-col>
      <v-col cols="12">
        <AlertBlock>修改之后将实时生效，请谨慎操作！</AlertBlock>
      </v-col>

      <v-col cols="12">
        <TableWithPager @query="doTableQuery" ref="tableWithPagerRef" :infos="tableInfos">
          <el-table-column label="字典编号" width="200px" prop="code">
            <template #default="{ row }">
              <span v-copy="row.code">{{ row.code }}</span>
            </template>
          </el-table-column>
          <el-table-column label="字典名称" prop="dictLabel" width="150px"></el-table-column>
          <el-table-column label="字典排序" prop="sort" width="100px"></el-table-column>
          <el-table-column label="字典类型" width="100px">
            <template #default="{ row }">
              <span>{{ getLabels([["sys_dict_type", row.dictType]]) }}</span>
            </template>
          </el-table-column>
          <el-table-column label="备注" prop="remark" min-width="200px"></el-table-column>
          <el-table-column label="更新时间" min-width="160px">
            <template #default="{ row }">
              {{ format.dateFormat(row.updatedAt, "YYYY-MM-DD HH:mm:ss") }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="160px" fixed="right">
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
      此操作将会<span class="text-primary">删除</span>正在使用的字典<br />
      字典编号：<span class="text-primary font-weight-black">{{ confirmDelete.currentCode }}</span
      ><br />
      你还要继续吗？
    </template>
  </ConfirmByInput>
  <CreateDictPane ref="createDictPaneRef" @submit="doQueryFirstPage" />
  <ConfigDictPane ref="configDictPaneRef" @refresh="doTableQuery" />
</template>
<script setup lang="ts">
import { reactive, ref, onMounted } from "vue";
import BaseBreadcrumb from "@/components/shared/BaseBreadcrumb.vue";
import UiParentCard from "@/components/shared/UiParentCard.vue";
import AlertBlock from "@/components/ui/AlertBlock.vue";
import CreateDictPane from "./components/CreateDictPane.vue";
import ConfirmByInput from "@/components/business/ConfirmByInput.vue";
import ConfigDictPane from "./components/ConfigDictPane.vue";

import { http, format } from "@/utils";
import { useMapRemoteStore } from "@/stores";

const { getLabels, loadDictTree } = useMapRemoteStore();
loadDictTree(["sys_dict_type"]);

const page = ref({ title: "系统字典" });
const breadcrumbs = ref([
  {
    text: "系统管理",
    disabled: false,
    href: "#"
  },
  {
    text: "系统字典",
    disabled: true,
    href: "#"
  }
]);
const searchData = reactive({
  code: "",
  label: ""
});
const tableInfos = reactive({
  list: [],
  total: 0
});
const createDictPaneRef = ref();
const configDictPaneRef = ref();
const tableWithPagerRef = ref();
const refConfirmDelete = ref();
const confirmDelete = reactive({
  id: null,
  currentCode: ""
});

const getButtons = row => {
  let ret = [];
  // ret.push({
  //   text: "编辑",
  //   color: "info",
  //   click() {
  //     onEdit(row);
  //   }
  // });
  ret.push({
    text: "字典配置",
    color: "info",
    click() {
      onConfig(row);
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
  confirmDelete.currentCode = row.code;
  confirmDelete.id = row.id;
  refConfirmDelete.value.show({
    width: "450px",
    confirmText: confirmDelete.currentCode
  });
};

const doDelete = async (options = {}) => {
  const [err, res] = await http.delete({
    ...options,
    showSuccess: true,
    url: `/sys/dict/${confirmDelete.id}`
  });
  if (res) {
    refConfirmDelete.value.hide();
    doTableQuery();
  }
};

const doTableQuery = async (options = {}) => {
  const [err, res] = await http.get({
    url: "/sys/dict",
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
  createDictPaneRef.value.show({
    title: "添加字典",
    operateType: "add"
  });
};

const onEdit = info => {
  createDictPaneRef.value.show({
    title: "编辑字典",
    infos: info,
    operateType: "edit"
  });
};

const onConfig = row => {
  configDictPaneRef.value.show({
    title: `字典配置（${row.dictLabel}）`,
    id: row.id
  });
};

onMounted(() => {
  doTableQuery();
});
</script>
