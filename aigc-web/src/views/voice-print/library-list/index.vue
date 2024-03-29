<template>
  <BaseBreadcrumb :title="page.title" :breadcrumbs="breadcrumbs"></BaseBreadcrumb>
  <UiParentCard>
    <v-row>
      <v-col cols="12" lg="3" md="4" sm="6">
        <v-text-field
          density="compact"
          v-model="searchData.userName"
          label="用户姓名"
          hide-details
          variant="outlined"
          :clearable="true"
          @keyup.enter="doQueryFirstPage"
          @click:clear="doQueryFirstPage"
        ></v-text-field>
      </v-col>
      <v-col cols="12" lg="3" md="4" sm="6">
        <v-text-field
          density="compact"
          v-model="searchData.userKey"
          label="用户标识"
          hide-details
          variant="outlined"
          :clearable="true"
          @keyup.enter="doQueryFirstPage"
          @click:clear="doQueryFirstPage"
        ></v-text-field>
      </v-col>
      <v-col cols="12" lg="3" md="4" sm="6">
        <ButtonsInForm>
          <v-btn color="primary" @click="onAdd">声纹注册</v-btn>
          <v-btn color="secondary" @click="onQuery">声纹查询</v-btn>
        </ButtonsInForm>
      </v-col>
      <v-col cols="12">
        <TableWithPager @query="doTableQuery" ref="tableWithPagerRef" :infos="tableInfos">
          <el-table-column label="用户姓名" prop="userName" min-width="200px"></el-table-column>
          <el-table-column label="用户标识" min-width="200px">
            <template #default="{ row }">
              <span v-copy="row.userKey">{{ row.userKey }}</span>
            </template>
          </el-table-column>
          <el-table-column label="音频文件" min-width="300px">
            <template #default="{ row }">
              <AiAudio :src="row?.s3Url" />
            </template>
          </el-table-column>
        </TableWithPager>
      </v-col>
    </v-row>
  </UiParentCard>
</template>
<script setup lang="ts">
import { ref, reactive, onMounted, inject } from "vue";
import { useRouter } from "vue-router";
import BaseBreadcrumb from "@/components/shared/BaseBreadcrumb.vue";
import UiParentCard from "@/components/shared/UiParentCard.vue";
import AiAudio from "@/components/business/AiAudio.vue";
import { http } from "@/utils";
import { type ItfAspectPageState } from "@/types/AspectPageTypes.ts";

defineOptions({
  name: "VoicePrintLibraryList"
});

const provideAspectPage = inject("provideAspectPage") as ItfAspectPageState;

const page = ref({ title: "声纹库列表" });
const breadcrumbs = ref([
  {
    text: "智能声纹",
    disabled: false,
    href: "#"
  },
  {
    text: "声纹库列表",
    disabled: true,
    href: "#"
  }
]);

const router = useRouter();

const tableWithPagerRef = ref();
const tableInfos = reactive({
  list: [],
  total: 0
});

const searchData = reactive({
  userName: "",
  userKey: ""
});

const doQueryFirstPage = () => {
  tableWithPagerRef.value.query({ page: 1 });
};

// 清空查询条件
const doClear = () => {
  searchData.userName = "";
  searchData.userKey = "";
};

const doTableQuery = async (options = {}) => {
  const [err, res] = await http.get({
    url: "/voice/list",
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

// 声纹注册
const onAdd = () => {
  router.push("/voice-print/library-list/register");
};

// 声纹查询
const onQuery = () => {
  router.push("/voice-print/library-list/search");
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
