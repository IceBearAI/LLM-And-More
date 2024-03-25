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
          <v-btn color="primary" @click="onAdd">创建降噪</v-btn>
        </ButtonsInForm>
      </v-col>
      <v-col cols="12">
        <TableWithPager @query="doTableQuery" ref="tableWithPagerRef" :infos="tableInfos">
          <el-table-column label="降噪前" min-width="320px">
            <template #default="{ row }">
              <AiAudio :src="row?.inputS3Url" />
              <div>{{ row.inputName }}</div>
            </template>
          </el-table-column>
          <el-table-column label="降噪后" min-width="320px">
            <template #default="{ row }">
              <AiAudio :src="row?.outS3Url" />
              <div>{{ row.outName }}</div>
            </template>
          </el-table-column>
          <el-table-column label="采样率" min-width="150px">
            <template #default="{ row }">
              <span>{{ getLabels([["denoise_sample_rate", row.sampleRate]]) }}</span>
            </template>
          </el-table-column>
          <el-table-column label="更新时间" min-width="160px">
            <template #default="{ row }">
              {{ format.dateFormat(row.updatedAt, "YYYY-MM-DD HH:mm:ss") }}
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
    <template #text> 您确定要删除【{{ selectedRow.inputName }}】的降噪结果吗？ </template>
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
import { useMapRemoteStore } from "@/stores";
import { ItfAspectPageState } from "@/types/AspectPageTypes.ts";

const provideAspectPage = inject("provideAspectPage") as ItfAspectPageState;
const page = ref({ title: "降噪列表" });
const breadcrumbs = ref([
  {
    text: "语音服务",
    disabled: false,
    href: "#"
  },
  {
    text: "降噪列表",
    disabled: true,
    href: "#"
  }
]);

const { loadDictTree, getLabels } = useMapRemoteStore();

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

provideAspectPage.methods.refreshListPage = () => {
  // 清空查询条件
  doClear();
  // 重新查询
  doQueryFirstPage();
};

const doTableQuery = async (options = {}) => {
  const [err, res] = await http.get({
    url: "/voice/denoise",
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

// 创建降噪
const onAdd = () => {
  router.push("/voice-print/denoise-list/denoise");
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
    url: `/voice/denoise/${selectedRow.value.id}`
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

onMounted(async () => {
  await loadDictTree(["denoise_sample_rate"]);
  doTableQuery();
});
</script>
