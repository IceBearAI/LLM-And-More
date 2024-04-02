<template>
  <BaseBreadcrumb :title="page.title" :breadcrumbs="breadcrumbs"></BaseBreadcrumb>
  <UiParentCard>
    <v-row>
      <v-col cols="12" lg="3" md="4" sm="6">
        <ButtonsInForm>
          <v-btn color="primary" @click="onAdd">创建检测</v-btn>
        </ButtonsInForm>
      </v-col>

      <v-col cols="12">
        <TableWithPager @query="doTableQuery" ref="tableWithPagerRef" :infos="tableInfos">
          <el-table-column label="检测图片" min-width="150px">
            <template #default="{ row }">
              <v-avatar size="80" rounded="md">
                <img :src="row.inputS3Url" alt="检测图片" height="80" />
              </v-avatar>
            </template>
          </el-table-column>
          <el-table-column label="比对图片" min-width="150px">
            <template #default="{ row }">
              <template v-if="row.outputS3Url">
                <v-avatar size="80" rounded="md">
                  <img :src="row.outputS3Url" alt="比对图片" height="80" />
                </v-avatar>
              </template>
              <template v-else> -- </template>
            </template>
          </el-table-column>
          <el-table-column label="人脸个数" prop="faceNum" min-width="160px"></el-table-column>
          <el-table-column label="比对阈值" prop="denoiseStrength" min-width="160px"></el-table-column>
          <el-table-column label="是否同一个人" min-width="120px">
            <template #default="{ row }">
              <ChipBoolean v-model="row.isSame"></ChipBoolean>
            </template>
          </el-table-column>
          <el-table-column label="操作人" prop="operatorEmail" min-width="150px" show-overflow-tooltip></el-table-column>
          <el-table-column label="操作" width="80px" fixed="right">
            <template #default="{ row }">
              <ButtonsInTable :buttons="getButtons(row)" onlyOne />
            </template>
          </el-table-column>
        </TableWithPager>
      </v-col>
    </v-row>
  </UiParentCard>
  <CreateRecognitionPane ref="createRecognitionPaneRef" @submit="doQueryFirstPage" />
</template>
<script setup lang="ts">
import { reactive, ref, onMounted } from "vue";
import BaseBreadcrumb from "@/components/shared/BaseBreadcrumb.vue";
import UiParentCard from "@/components/shared/UiParentCard.vue";
import ChipBoolean from "@/components/ui/ChipBoolean.vue";
import CreateRecognitionPane from "./components/CreateRecognitionPane.vue";

import { http, format } from "@/utils";

const page = ref({ title: "人脸检测" });
const breadcrumbs = ref([
  {
    text: "图像服务",
    disabled: false,
    href: "#"
  },
  {
    text: "人脸检测",
    disabled: true,
    href: "#"
  }
]);
const tableInfos = reactive({
  list: [],
  total: 0
});
const createRecognitionPaneRef = ref();
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
      modelType: "faceRecognition",
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
  createRecognitionPaneRef.value.show({
    title: "创建检测",
    operateType: "add"
  });
};

const onEdit = info => {
  createRecognitionPaneRef.value.show({
    title: "查看",
    infos: info,
    operateType: "edit"
  });
};

onMounted(() => {
  doTableQuery();
});
</script>
