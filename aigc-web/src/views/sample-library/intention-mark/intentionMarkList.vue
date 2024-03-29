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
          <v-btn color="primary" @click="onAdd">创建样本</v-btn>
        </ButtonsInForm>
      </v-col>
      <v-col cols="12">
        <TableWithPager @query="doTableQuery" ref="tableWithPagerRef" :infos="tableInfos">
          <el-table-column label="ID" min-width="240px">
            <template #default="{ row }">
              <a href="javascript: void(0)" class="link" @click="onViewDetail(row.intentId)">{{ row.intentId }}</a>
            </template>
          </el-table-column>
          <el-table-column label="意图标注名称" prop="name" min-width="200px" show-overflow-tooltip></el-table-column>
          <el-table-column label="样本数量" prop="documentCount" min-width="110px"></el-table-column>
          <el-table-column label="备注" prop="remark" min-width="200px" show-overflow-tooltip></el-table-column>
          <el-table-column label="修改时间" min-width="200px">
            <template #default="{ row }">
              {{ format.dateFormat(row.updatedAt, "YYYY-MM-DD HH:mm:ss") }}
            </template>
          </el-table-column>
          <el-table-column label="创建时间" min-width="200px">
            <template #default="{ row }">
              {{ format.dateFormat(row.createdAt, "YYYY-MM-DD HH:mm:ss") }}
            </template>
          </el-table-column>
          <el-table-column label="操作人" prop="operator" min-width="220px"></el-table-column>
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
      此操作将会<span class="text-primary font-weight-black">删除</span>该意图标注<br />
      意图标注ID：<span class="text-primary font-weight-black">{{ confirmData.intentId }}</span
      ><br />
      确定要继续吗？
    </template>
  </ConfirmByInput>

  <CreateIntentionMarkPane ref="createIntentionMarkPaneRef" @submit="doQueryFirstPage" />
</template>
<script setup lang="ts">
import { reactive, toRefs, ref, onMounted } from "vue";
import BaseBreadcrumb from "@/components/shared/BaseBreadcrumb.vue";
import UiParentCard from "@/components/shared/UiParentCard.vue";
import CreateIntentionMarkPane from "./components/CreateIntentionMarkPane.vue";

import ConfirmByInput from "@/components/business/ConfirmByInput.vue";
import { http, format } from "@/utils";
import { useRouter } from "vue-router";

const router = useRouter();

const page = ref({ title: "意图模型标注列表" });
const breadcrumbs = ref([
  {
    text: "样本库",
    disabled: false,
    href: "#"
  },
  {
    text: "意图模型标注列表",
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
const createIntentionMarkPaneRef = ref();
const tableWithPagerRef = ref();
const refConfirmDelete = ref();
const confirmData = reactive({
  intentId: ""
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
  confirmData.intentId = row.intentId;
  refConfirmDelete.value.show({
    width: "480px",
    confirmText: confirmData.intentId
  });
};

const doDelete = async (options = {}) => {
  const [err, res] = await http.delete({
    ...options,
    showSuccess: true,
    url: `/intent/${confirmData.intentId}/delete`
  });
  if (res) {
    refConfirmDelete.value.hide();
    doQueryCurrentPage();
  }
};

const doTableQuery = async (options = {}) => {
  const [err, res] = await http.get({
    url: "/intent/list",
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
  createIntentionMarkPaneRef.value.show({
    title: "创建意图标注",
    operateType: "add"
  });
};

const onViewDetail = id => {
  router.push(`/sample-library/intention-mark/detail?intentId=${id}`);
};

onMounted(() => {
  doQueryCurrentPage();
});
</script>
