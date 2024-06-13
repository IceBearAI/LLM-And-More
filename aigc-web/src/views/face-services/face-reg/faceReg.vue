<template>
  <BaseBreadcrumb :title="page.title" :breadcrumbs="breadcrumbs"></BaseBreadcrumb>
  <UiParentCard>
    <v-row>
      <v-col cols="12" lg="3" md="4" sm="6">
        <v-text-field
          v-model="searchData.userKey"
          label="请输入用户标识"
          hide-details
          clearable
          @keyup.enter="doQueryFirstPage"
          @click:clear="doQueryFirstPage"
        ></v-text-field>
      </v-col>
      <v-col cols="12" lg="3" md="4" sm="6">
        <v-text-field
          v-model="searchData.userType"
          label="请输入用户类型"
          hide-details
          clearable
          @keyup.enter="doQueryFirstPage"
          @click:clear="doQueryFirstPage"
        ></v-text-field>
      </v-col>
      <v-col cols="12" lg="3" md="4" sm="6">
        <v-text-field
          v-model="searchData.userRegion"
          label="请输入用户渠道"
          hide-details
          clearable
          @keyup.enter="doQueryFirstPage"
          @click:clear="doQueryFirstPage"
        ></v-text-field>
      </v-col>
      <v-col cols="12" lg="3" md="4" sm="6">
        <ButtonsInForm>
          <v-btn color="primary" @click="onAdd">注册人脸</v-btn>
          <v-btn color="primary" @click="onSearch">人脸搜索</v-btn>
        </ButtonsInForm>
      </v-col>

      <v-col cols="12">
        <TableWithPager @query="doTableQuery" ref="tableWithPagerRef" :infos="tableInfos">
          <el-table-column label="人脸图片" min-width="150px">
            <template #default="{ row }">
              <v-avatar size="80" rounded="md">
                <img :src="row.inputS3Url" alt="人脸图片" height="80" />
              </v-avatar>
            </template>
          </el-table-column>
          <el-table-column label="向量库uuid" min-width="200px" show-overflow-tooltip>
            <template #default="{ row }">
              <span v-copy="row.uuid">{{ row.uuid }}</span>
            </template>
          </el-table-column>
          <el-table-column label="用户标识" min-width="100px" show-overflow-tooltip>
            <template #default="{ row }">
              <template v-if="row.userKey">
                {{ row.userKey }}
              </template>
              <template v-else> 未指定 </template>
            </template>
          </el-table-column>
          <el-table-column label="用户类型" prop="userType" min-width="100px" show-overflow-tooltip></el-table-column>
          <el-table-column label="用户渠道" prop="userRegion" min-width="100px" show-overflow-tooltip></el-table-column>
          <el-table-column label="py耗时" prop="durationPy" min-width="90px">
            <template #default="{ row }">
              {{ row.durationPy.toFixed(4) }}
            </template>
          </el-table-column>
          <el-table-column label="接口耗时" prop="duration" min-width="90px">
            <template #default="{ row }">
              {{ row.duration.toFixed(4) }}
            </template>
          </el-table-column>
          <el-table-column label="创建时间" min-width="165px">
            <template #default="{ row }">
              {{ format.dateFormat(row.createdAt, "YYYY-MM-DD HH:mm:ss") }}
            </template>
          </el-table-column>
          <el-table-column label="更新时间" min-width="165px">
            <template #default="{ row }">
              {{ format.dateFormat(row.updatedAt, "YYYY-MM-DD HH:mm:ss") }}
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
  <ConfirmByInput ref="refConfirmDelete" @submit="doDelete">
    <template #text>
      您将要删除 <span class="text-primary font-weight-black">{{ confirmData.uuid }}</span> 这条记录<br />
      确定要继续吗？
    </template>
  </ConfirmByInput>
  <CreateRegPane ref="createRegPaneRef" @submit="doQueryFirstPage" />
  <FaceSearchPane ref="faceSearchPaneRef" />
</template>
<script setup lang="ts">
import { reactive, ref, onMounted } from "vue";
import BaseBreadcrumb from "@/components/shared/BaseBreadcrumb.vue";
import UiParentCard from "@/components/shared/UiParentCard.vue";
import ConfirmByInput from "@/components/business/ConfirmByInput.vue";
import CreateRegPane from "./components/CreateRegPane.vue";
import FaceSearchPane from "./components/FaceSearchPane.vue";

import { http, format } from "@/utils";

const page = ref({ title: "人脸注册" });
const breadcrumbs = ref([
  {
    text: "人脸服务",
    disabled: false,
    href: "#"
  },
  {
    text: "人脸注册",
    disabled: true,
    href: "#"
  }
]);
const searchData = reactive({
  userKey: "",
  userType: "",
  userRegion: ""
});
const tableInfos = reactive({
  list: [],
  total: 0
});
const createRegPaneRef = ref();
const faceSearchPaneRef = ref();
const tableWithPagerRef = ref();
const refConfirmDelete = ref();
const confirmData = reactive({
  uuid: ""
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
  return ret;
};

const onDelete = row => {
  confirmData.uuid = row.uuid;
  refConfirmDelete.value.show({
    width: "500px",
    confirmText: confirmData.uuid
  });
};

const doDelete = async (options = {}) => {
  const [err, res] = await http.delete({
    ...options,
    showSuccess: true,
    url: `/esrgan/face/reg/${confirmData.uuid}`
  });
  if (res) {
    refConfirmDelete.value.hide();
    doTableQuery();
  }
};

const doTableQuery = async (options = {}) => {
  const [err, res] = await http.get({
    url: "/esrgan/face/reg",
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
  createRegPaneRef.value.show({
    title: "注册人脸",
    operateType: "add"
  });
};

const onSearch = () => {
  faceSearchPaneRef.value.show({
    title: "人脸搜索"
  });
};

onMounted(() => {
  doTableQuery();
});
</script>
