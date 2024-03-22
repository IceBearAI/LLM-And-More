<template>
  <BaseBreadcrumb :title="page.title" :breadcrumbs="breadcrumbs"></BaseBreadcrumb>
  <UiParentCard>
    <v-row>
      <v-col cols="12" lg="3" md="4" sm="6">
        <Select
          v-model="searchData.modelName"
          :mapDictionary="{ code: 'esrgan_model_type' }"
          label="请选择模型"
          hide-details
          @change="doQueryFirstPage"
        ></Select>
      </v-col>
      <v-col cols="12" lg="3" md="4" sm="6">
        <ButtonsInForm>
          <v-btn color="primary" @click="onAdd">创建超分</v-btn>
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
          <el-table-column label="超分后" min-width="150px">
            <template #default="{ row }">
              <v-avatar size="80" rounded="md">
                <img :src="row.outputS3Url" alt="超分后" height="80" />
              </v-avatar>
            </template>
          </el-table-column>
          <el-table-column label="超分模型" min-width="160px">
            <template #default="{ row }"> {{ getLabels([["esrgan_model_type", row.modelName]]) }} </template>
          </el-table-column>
          <el-table-column label="是否面部增强" min-width="120px">
            <template #default="{ row }">
              <ChipBoolean v-model="row.faceEnhance"></ChipBoolean>
            </template>
          </el-table-column>
          <el-table-column label="降噪强度" min-width="100px">
            <template #default="{ row }">
              {{ row.modelName === "realesr-general-x4v3" ? row.denoiseStrength : "-" }}
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
  <CreateSuperPane ref="createSuperPaneRef" @submit="doQueryFirstPage" />
</template>
<script setup lang="ts">
import { reactive, ref, onMounted } from "vue";
import BaseBreadcrumb from "@/components/shared/BaseBreadcrumb.vue";
import UiParentCard from "@/components/shared/UiParentCard.vue";
import CreateSuperPane from "./components/CreateSuperPane.vue";
import ChipBoolean from "@/components/ui/ChipBoolean.vue";

import { http } from "@/utils";
import { useMapRemoteStore } from "@/stores";

const { getLabels } = useMapRemoteStore();

const page = ref({ title: "超分列表" });
const breadcrumbs = ref([
  {
    text: "图像服务",
    disabled: false,
    href: "#"
  },
  {
    text: "超分列表",
    disabled: true,
    href: "#"
  }
]);
const searchData = reactive({
  modelType: "esrgan",
  modelName: null
});
const tableInfos = reactive({
  list: [],
  total: 0
});
const createSuperPaneRef = ref();
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
  createSuperPaneRef.value.show({
    title: "创建超分",
    operateType: "add"
  });
};

const onEdit = info => {
  createSuperPaneRef.value.show({
    title: "查看",
    infos: info,
    operateType: "edit"
  });
};

onMounted(() => {
  doTableQuery();
});
</script>
