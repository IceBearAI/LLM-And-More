<template>
  <BaseBreadcrumb :title="page.title" :breadcrumbs="breadcrumbs"></BaseBreadcrumb>
  <UiParentCard>
    <v-row>
      <v-col cols="12" lg="3" md="4" sm="6">
        <ButtonsInForm>
          <v-btn color="primary" @click="onAdd">创建抠图</v-btn>
        </ButtonsInForm>
      </v-col>

      <v-col cols="12">
        <TableWithPager @query="doTableQuery" ref="tableWithPagerRef" :infos="tableInfos">
          <el-table-column label="原图" min-width="150px">
            <template #default="{ row }">
              <v-avatar size="80" rounded="md">
                <img :src="row.inputS3Url" alt="原图" height="80" />
              </v-avatar>
            </template>
          </el-table-column>
          <el-table-column label="抠图后" min-width="150px">
            <template #default="{ row }">
              <v-avatar size="80" rounded="md">
                <img :src="row.outputS3Url" alt="抠图后" height="80" />
              </v-avatar>
            </template>
          </el-table-column>
          <el-table-column label="模型名称" prop="modelName" min-width="160px"></el-table-column>
          <el-table-column label="操作人" prop="operatorEmail" min-width="150px" show-overflow-tooltip></el-table-column>
          <el-table-column label="创建时间" min-width="160px">
            <template #default="{ row }">
              {{ format.dateFormat(row.createdAt, "YYYY-MM-DD HH:mm:ss") }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="80px" fixed="right">
            <template #default="{ row }">
              <ButtonsInTable :buttons="getButtons(row)" onlyOne />
            </template>
          </el-table-column>
        </TableWithPager>
      </v-col>
    </v-row>
  </UiParentCard>
  <CreateMattingPane ref="createMattingPaneRef" @submit="doQueryFirstPage" />
</template>
<script setup lang="ts">
import { reactive, ref, onMounted } from "vue";
import BaseBreadcrumb from "@/components/shared/BaseBreadcrumb.vue";
import UiParentCard from "@/components/shared/UiParentCard.vue";
import CreateMattingPane from "./components/CreateMattingPane.vue";

import { http, format } from "@/utils";

const page = ref({ title: "抠图列表" });
const breadcrumbs = ref([
  {
    text: "图像服务",
    disabled: false,
    href: "#"
  },
  {
    text: "抠图列表",
    disabled: true,
    href: "#"
  }
]);
const tableInfos = reactive({
  list: [],
  total: 0
});
const createMattingPaneRef = ref();
const tableWithPagerRef = ref();

const getButtons = row => {
  let ret = [];
  ret.push({
    text: "查看",
    color: "info",
    click() {
      onEdit(row);
    }
  });
  return ret;
};

const doTableQuery = async (options = {}) => {
  const [err, res] = await http.get({
    url: "/esrgan/list",
    showLoading: tableWithPagerRef.value.el,
    data: {
      modelType: "rmBackground",
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
  createMattingPaneRef.value.show({
    title: "创建抠图",
    operateType: "add"
  });
};

const onEdit = info => {
  createMattingPaneRef.value.show({
    title: "查看",
    infos: info,
    operateType: "edit"
  });
};

onMounted(() => {
  doTableQuery();
});
</script>
