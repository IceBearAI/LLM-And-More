<template>
  <v-row>
    <v-col class="d-flex justify-end" cols="12">
      <ButtonsInForm>
        <v-btn color="primary" @click="onAdd">上传授权音频</v-btn>
      </ButtonsInForm>
    </v-col>
    <v-col cols="12">
      <TableWithPager @query="doTableQuery" ref="tableWithPagerRef" :infos="tableInfos">
        <el-table-column label="项目ID" prop="projectId" min-width="200px" show-overflow-tooltip></el-table-column>
        <el-table-column label="授权同意ID" prop="consentId" min-width="260px"></el-table-column>
        <el-table-column label="授权音频" min-width="330px">
          <template #default="{ row }">
            <AiAudio :src="row?.audiodataUrl" />
          </template>
        </el-table-column>
        <el-table-column label="备注" prop="description" min-width="200px"></el-table-column>
        <el-table-column label="创建时间" min-width="165px">
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

  <ConfirmByClick ref="refConfirmDelete" @submit="doDelete">
    <template #text>
      <span class="font-weight-bold">这是进行一项操作时必须了解的重要信息</span><br />
      您将要删除 <span class="text-primary">{{ confirmData.consentId }}</span
      ><br />
      确定要继续吗？
    </template>
  </ConfirmByClick>

  <CreateConsentVoicePane ref="createConsentVoiceRef" :speak-name="speakName" @submit="doQueryFirstPage" />
</template>
<script setup lang="ts">
import { reactive, ref, onMounted } from "vue";
import CreateConsentVoicePane from "./CreateConsentVoicePane.vue";

import ConfirmByClick from "@/components/business/ConfirmByClick.vue";
import AiAudio from "@/components/business/AiAudio.vue";
import { http, format } from "@/utils";
import { useRoute } from "vue-router";

const route = useRoute();
const { speakName } = route.query as { speakName: string };

const tableInfos = reactive({
  list: [],
  total: 0
});
const createConsentVoiceRef = ref();
const tableWithPagerRef = ref();
const refConfirmDelete = ref();
const confirmData = reactive({
  consentId: ""
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
  confirmData.consentId = row.consentId;
  refConfirmDelete.value.show({
    width: "450px"
  });
};

const doDelete = async (options = {}) => {
  const [err, res] = await http.delete({
    ...options,
    showSuccess: true,
    url: `/personal/consent/${confirmData.consentId}`
  });
  if (res) {
    refConfirmDelete.value.hide();
    doQueryCurrentPage();
  }
};

const doTableQuery = async (options = {}) => {
  const [err, res] = await http.get({
    url: `/personal/consent`,
    showLoading: tableWithPagerRef.value.el,
    data: {
      speakName,
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
  createConsentVoiceRef.value.show({
    title: "上传授权音频",
    operateType: "add"
  });
};

onMounted(() => {
  doQueryCurrentPage();
});
</script>
