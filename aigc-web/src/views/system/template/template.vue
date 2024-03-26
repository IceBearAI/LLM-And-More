<template>
  <BaseBreadcrumb :title="page.title" :breadcrumbs="breadcrumbs"></BaseBreadcrumb>
  <UiParentCard>
    <v-row>
      <v-col cols="12" lg="3" md="4" sm="6">
        <v-text-field
          v-model="searchData.name"
          label="请输入模型"
          hide-details
          clearable
          @keyup.enter="doQueryFirstPage"
          @click:clear="doQueryFirstPage"
        ></v-text-field>
      </v-col>
      <v-col cols="12" lg="3" md="4" sm="6">
        <Select
          v-model="searchData.templateType"
          :mapDictionary="{ code: 'template_type' }"
          label="请选择模版类型"
          hide-details
          @change="doQueryFirstPage"
        ></Select>
      </v-col>
      <v-col cols="12" lg="3" md="4" sm="6">
        <ButtonsInForm>
          <v-btn color="primary" @click="onAdd">添加模版</v-btn>
        </ButtonsInForm>
      </v-col>
      <v-col cols="12">
        <AlertBlock>修改之后将实时生效，请谨慎操作！</AlertBlock>
      </v-col>

      <v-col cols="12">
        <TableWithPager @query="doTableQuery" ref="tableWithPagerRef" :infos="tableInfos">
          <el-table-column label="名称" width="200px" show-overflow-tooltip>
            <template #default="{ row }">
              <span v-copy="row.name">{{ row.name }}</span>
            </template>
          </el-table-column>
          <el-table-column label="基础模型" prop="baseModel" width="180px"></el-table-column>
          <el-table-column label="最长上下文" prop="maxTokens" width="110px"></el-table-column>
          <el-table-column label="模版类型" width="100px">
            <template #default="{ row }"> {{ getLabels([["template_type", row.templateType]]) }} </template>
          </el-table-column>
          <el-table-column label="镜像" width="200px" show-overflow-tooltip>
            <template #default="{ row }">
              <span v-copy="row.trainImage">{{ row.trainImage }}</span>
            </template>
          </el-table-column>
          <el-table-column label="状态" width="100px">
            <template #default="{ row }">
              <ChipStatus v-model="row.enabled" />
            </template>
          </el-table-column>
          <el-table-column label="备注" prop="remark" min-width="200px"></el-table-column>
          <el-table-column label="更新时间" min-width="165px">
            <template #default="{ row }">
              {{ format.dateFormat(row.updatedAt, "YYYY-MM-DD HH:mm:ss") }}
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
      您将要删除<span class="text-primary font-weight-black">{{ confirmDelete.name }}</span
      >模版，删除之后将无法使用该模版创建微调任务<br />
      确定要继续吗？
    </template>
  </ConfirmByInput>
  <CreateTemplatePane ref="createTemplatePaneRef" @submit="doQueryFirstPage" />
</template>
<script setup lang="ts">
import { reactive, ref, onMounted } from "vue";
import BaseBreadcrumb from "@/components/shared/BaseBreadcrumb.vue";
import UiParentCard from "@/components/shared/UiParentCard.vue";
import AlertBlock from "@/components/ui/AlertBlock.vue";
import ConfirmByInput from "@/components/business/ConfirmByInput.vue";
import CreateTemplatePane from "./components/CreateTemplatePane.vue";
import ChipStatus from "@/components/ui/ChipStatus.vue";

import { http, format } from "@/utils";
import { useMapRemoteStore } from "@/stores";

const { getLabels } = useMapRemoteStore();

const page = ref({ title: "模版管理" });
const breadcrumbs = ref([
  {
    text: "系统管理",
    disabled: false,
    href: "#"
  },
  {
    text: "模版管理",
    disabled: true,
    href: "#"
  }
]);
const searchData = reactive({
  name: "",
  templateType: null
});
const tableInfos = reactive({
  list: [],
  total: 0
});
const createTemplatePaneRef = ref();
const tableWithPagerRef = ref();
const refConfirmDelete = ref();
const confirmDelete = reactive({
  name: ""
});

const getButtons = row => {
  let ret = [];
  ret.push({
    text: "删除",
    color: "error",
    click() {
      onDelete(row);
    }
  });
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
    url: `/sys/template/${confirmDelete.name}`
  });
  if (res) {
    refConfirmDelete.value.hide();
    doTableQuery();
  }
};

const doTableQuery = async (options = {}) => {
  const [err, res] = await http.get({
    url: "/sys/template",
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
  createTemplatePaneRef.value.show({
    title: "添加模版",
    operateType: "add"
  });
};

const onEdit = info => {
  createTemplatePaneRef.value.show({
    title: "编辑模版",
    infos: info,
    operateType: "edit"
  });
};

onMounted(() => {
  doTableQuery();
});
</script>
