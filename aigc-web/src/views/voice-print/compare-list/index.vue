<template>
  <BaseBreadcrumb :title="page.title" :breadcrumbs="breadcrumbs"></BaseBreadcrumb>
  <UiParentCard>
    <v-row>
      <v-col cols="12" lg="3" md="4" sm="6">
        <v-text-field
          density="compact"
          v-model="searchData.name"
          label="音频名称"
          hide-details
          variant="outlined"
          :clearable="true"
          @keyup.enter="doQueryFirstPage"
          @click:clear="doQueryFirstPage"
        ></v-text-field>
      </v-col>
      <v-col cols="12" lg="3" md="4" sm="6">
        <ButtonsInForm>
          <v-btn color="primary" @click="onAdd">声纹比对</v-btn>
        </ButtonsInForm>
      </v-col>
      <v-col cols="12">
        <TableWithPager @query="doTableQuery" ref="tableWithPagerRef" :infos="tableInfos">
          <el-table-column label="音频1" min-width="320px">
            <template #default="{ row }">
              <AiAudio :src="row?.input1S3Url" />
              <div>{{ row.input1Name }}</div>
            </template>
          </el-table-column>
          <el-table-column label="音频2" min-width="320px">
            <template #default="{ row }">
              <AiAudio :src="row?.input2S3Url" />
              <div>{{ row.input2Name }}</div>
            </template>
          </el-table-column>
          <el-table-column label="模型" prop="modelType" min-width="150px"></el-table-column>
          <el-table-column label="比分" prop="dist" min-width="180px"></el-table-column>
          <el-table-column label="操作人" prop="operatorEmail" min-width="200px"></el-table-column>
          <el-table-column label="创建时间" min-width="165px">
            <template #default="{ row }">
              {{ format.dateFormat(row.createdAt, "YYYY-MM-DD HH:mm:ss") }}
            </template>
          </el-table-column>
          <el-table-column label="操作" min-width="80px" fixed="right">
            <template #default="{ row }">
              <ButtonsInTable :buttons="getButtons(row)" onlyOne />
            </template>
          </el-table-column>
        </TableWithPager>
      </v-col>
    </v-row>
  </UiParentCard>

  <ConfirmByClick ref="refConfirmDelete" @submit="doDelete">
    <template #text> 您确定要删除【{{ selectedRow.input1Name }}】与【{{ selectedRow.input2Name }}】的比对结果吗？ </template>
  </ConfirmByClick>
</template>
<script setup lang="ts">
import { ref, reactive, onMounted, inject } from "vue";
import { useRoute, useRouter } from "vue-router";
import BaseBreadcrumb from "@/components/shared/BaseBreadcrumb.vue";
import UiParentCard from "@/components/shared/UiParentCard.vue";
import AiAudio from "@/components/business/AiAudio.vue";
import ConfirmByClick from "@/components/business/ConfirmByClick.vue";
import { http, format } from "@/utils";
import { type ItfAspectPageState } from "@/types/AspectPageTypes.ts";

defineOptions({
  name: "VoicePrintCompareList"
});
const provideAspectPage = inject("provideAspectPage") as ItfAspectPageState;

const page = ref({ title: "声纹比对列表" });
const breadcrumbs = ref([
  {
    text: "智能声纹",
    disabled: false,
    href: "#"
  },
  {
    text: "声纹比对列表",
    disabled: true,
    href: "#"
  }
]);

const route = useRoute();
const router = useRouter();

// 当前选中的行
const selectedRow = ref();
// 删除确认框组件ref
const refConfirmDelete = ref();
const tableWithPagerRef = ref();
const tableInfos = reactive({
  list: [],
  total: 0
});

const searchData = reactive({
  name: ""
});

const doQueryFirstPage = () => {
  tableWithPagerRef.value.query({ page: 1 });
};

// 清空查询条件
const doClear = () => {
  searchData.name = "";
};

const doTableQuery = async (options = {}) => {
  const [err, res] = await http.get({
    url: "/voice/compare",
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

// 声纹比对
const onAdd = () => {
  router.push("/voice-print/compare-list/compare");
};

// 删除
const onDelete = row => {
  selectedRow.value = row;
  refConfirmDelete.value.show({
    width: "400px"
  });
};

// 执行删除
const doDelete = async () => {
  const [err, res] = await http.delete({
    showSuccess: true,
    url: `/voice/compare/${selectedRow.value.id}`
  });
  if (res) {
    refConfirmDelete.value.hide();
    tableWithPagerRef.value.query();
  }
};

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

provideAspectPage.methods.refreshListPage = () => {
  // 清空查询条件
  doClear();
  // 重新查询
  doQueryFirstPage();
};

onMounted(() => {
  doTableQuery();
});
</script>
