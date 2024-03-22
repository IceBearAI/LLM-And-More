<template>
  <BaseBreadcrumb :title="page.title" :breadcrumbs="breadcrumbs"></BaseBreadcrumb>
  <UiParentCard>
    <v-row>
      <v-col cols="12" lg="3" md="4" sm="6">
        <v-text-field
          v-model="searchData.name"
          label="请输入样本名称"
          hide-details
          clearable
          @keyup.enter="doQueryFirstPage"
          @click:clear="doQueryFirstPage"
        ></v-text-field>
      </v-col>
      <v-col cols="12" lg="3" md="4" sm="6">
        <ButtonsInForm>
          <v-btn color="primary" @click="onAdd">上传样本</v-btn>
        </ButtonsInForm>
      </v-col>
      <v-col cols="12">
        <TableWithPager @query="doTableQuery" ref="tableWithPagerRef" :infos="tableInfos">
          <el-table-column label="样本名称" prop="name" min-width="200px" show-overflow-tooltip></el-table-column>
          <el-table-column label="样本ID" min-width="200px" show-overflow-tooltip>
            <template #default="{ row }">
              <span v-copy="row.uuid">{{ row.uuid }}</span>
            </template>
          </el-table-column>
          <el-table-column label="样本数量" prop="segmentCount" min-width="110px"></el-table-column>
          <el-table-column label="备注" prop="remark" min-width="200px" show-overflow-tooltip></el-table-column>
          <el-table-column label="创建时间" min-width="200px">
            <template #default="{ row }">
              {{ format.dateFormat(row.createdAt, "YYYY-MM-DD HH:mm:ss") }}
            </template>
          </el-table-column>
          <el-table-column label="操作人" prop="creatorEmail" min-width="220px"></el-table-column>
          <el-table-column label="操作" width="80px" fixed="right">
            <template #default="{ row }">
              <ButtonsInTable :buttons="getButtons(row)" onlyOne />
            </template>
          </el-table-column>
        </TableWithPager>
      </v-col>
    </v-row>
  </UiParentCard>

  <ConfirmByInput ref="refConfirmDelete" @submit="doDelete">
    <template #text>
      此操作将会<span class="text-primary font-weight-black">删除</span>该样本数据集<br />
      数据集ID：<span class="text-primary font-weight-black">{{ confirmData.id }}</span
      ><br />
      确定要继续吗？
    </template>
  </ConfirmByInput>

  <CreateSamplePane ref="createSamplePaneRef" @submit="doQueryFirstPage" />
</template>
<script setup lang="ts">
import { reactive, toRefs, ref, onMounted } from "vue";
import BaseBreadcrumb from "@/components/shared/BaseBreadcrumb.vue";
import UiParentCard from "@/components/shared/UiParentCard.vue";
import CreateSamplePane from "./components/CreateSamplePane.vue";

import ConfirmByInput from "@/components/business/ConfirmByInput.vue";
import { http, format } from "@/utils";
import { useRoute } from "vue-router";

const route = useRoute();
const page = ref({ title: "文本样本列表" });
const breadcrumbs = ref([
  {
    text: "样本库",
    disabled: false,
    href: "#"
  },
  {
    text: "文本样本列表",
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
const createSamplePaneRef = ref();
const tableWithPagerRef = ref();
const refConfirmDelete = ref();
const confirmData = reactive({
  id: ""
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
    text: "删除",
    color: "error",
    click() {
      onDelete(row);
    }
  });
  return ret;
};

const onDelete = row => {
  confirmData.id = row.uuid;
  refConfirmDelete.value.show({
    width: "450px",
    confirmText: confirmData.id
  });
};

const doDelete = async (options = {}) => {
  const [err, res] = await http.delete({
    ...options,
    showSuccess: true,
    url: `/mgr/datasets/${confirmData.id}`
  });
  if (res) {
    refConfirmDelete.value.hide();
    doQueryCurrentPage();
  }
};

const doTableQuery = async (options = {}) => {
  const [err, res] = await http.get({
    url: "/mgr/datasets/list",
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
  createSamplePaneRef.value.show({
    title: "上传样本",
    operateType: "add"
  });
};

// const onEdit = info => {
//   createSamplePaneRef.value.show({
//     title: "编辑样本",
//     infos: info,
//     operateType: "edit"
//   });
// };

onMounted(() => {
  searchData.name = route.query.name as string;
  doQueryCurrentPage();
});
</script>
