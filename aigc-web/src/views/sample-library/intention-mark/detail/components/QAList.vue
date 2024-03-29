<template>
  <v-row>
    <v-col cols="12" lg="3" md="4" sm="6">
      <v-text-field
        v-model="searchData.keyword"
        label="请输入关键字"
        hide-details
        clearable
        @keyup.enter="doQueryFirstPage"
        @click:clear="doQueryFirstPage"
      ></v-text-field>
    </v-col>
    <v-col cols="12" lg="3" md="4" sm="6">
      <ButtonsInForm>
        <v-btn color="primary" @click="onAdd">添加问答</v-btn>
      </ButtonsInForm>
    </v-col>
    <v-col cols="12">
      <TableWithPager @query="doTableQuery" ref="tableWithPagerRef" :infos="tableInfos">
        <el-table-column label="ID" min-width="240px">
          <template #default="{ row }">
            <a href="javascript: void(0)" class="link" @click="onEdit(row)">{{ row.documentId }}</a>
          </template>
        </el-table-column>
        <el-table-column label="标准问法" min-width="200px" show-overflow-tooltip>
          <template #default="{ row }">
            <a href="javascript: void(0)" class="link" @click="onEdit(row)">{{ row.question }}</a>
          </template>
        </el-table-column>
        <el-table-column label="答案" prop="output" min-width="200px" show-overflow-tooltip></el-table-column>
        <el-table-column label="意图" prop="intent" min-width="200px" show-overflow-tooltip></el-table-column>
        <el-table-column label="更新时间" min-width="200px">
          <template #default="{ row }">
            {{ format.dateFormat(row.updatedAt, "YYYY-MM-DD HH:mm:ss") }}
          </template>
        </el-table-column>
        <el-table-column label="添加时间" min-width="200px">
          <template #default="{ row }">
            {{ format.dateFormat(row.createdAt, "YYYY-MM-DD HH:mm:ss") }}
          </template>
        </el-table-column>
        <el-table-column label="操作人" prop="operator" min-width="220px"></el-table-column>
        <el-table-column label="操作" width="100px" fixed="right">
          <template #default="{ row }">
            <ButtonsInTable :buttons="getButtons(row)" />
          </template>
        </el-table-column>
      </TableWithPager>
    </v-col>
  </v-row>

  <ConfirmByClick ref="refConfirmDelete" @submit="doDelete">
    <template #text>
      <span class="font-weight-bold">这是进行一项操作时必须了解的重要信息</span><br />
      您将要废弃<span class="text-primary">{{ confirmData.question }}</span
      ><br />
      确定要继续吗？
    </template>
  </ConfirmByClick>

  <CreateQAPane ref="createQAPaneRef" :intentId="intentId" @submit="doQueryFirstPage" />
</template>
<script setup lang="ts">
import { reactive, ref, onMounted } from "vue";
import CreateQAPane from "./CreateQAPane.vue";

import ConfirmByClick from "@/components/business/ConfirmByClick.vue";
import { http, format } from "@/utils";
import { useRoute } from "vue-router";

const route = useRoute();
const { intentId } = route.query as { intentId: string };

const searchData = reactive({
  keyword: ""
});
const tableInfos = reactive({
  list: [],
  total: 0
});
const createQAPaneRef = ref();
const tableWithPagerRef = ref();
const refConfirmDelete = ref();
const confirmData = reactive({
  documentId: "",
  question: ""
});

const getButtons = row => {
  let ret = [];
  ret.push({
    text: "编辑",
    color: "info",
    click() {
      onEdit(row);
    }
  });
  ret.push({
    text: "废弃",
    color: "error",
    click() {
      onDelete(row);
    }
  });
  return ret;
};

const onDelete = row => {
  confirmData.documentId = row.documentId;
  confirmData.question = row.question;
  refConfirmDelete.value.show({
    width: "450px"
  });
};

const doDelete = async (options = {}) => {
  const [err, res] = await http.delete({
    ...options,
    showSuccess: true,
    url: `/intent/${intentId}/document/${confirmData.documentId}/delete`
  });
  if (res) {
    refConfirmDelete.value.hide();
    doQueryCurrentPage();
  }
};

const doTableQuery = async (options = {}) => {
  const [err, res] = await http.get({
    url: `/intent/${intentId}/document/list`,
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
  createQAPaneRef.value.show({
    title: "添加问答",
    operateType: "add"
  });
};

const onEdit = info => {
  createQAPaneRef.value.show({
    title: "编辑问答",
    infos: info,
    operateType: "edit"
  });
};

onMounted(() => {
  doQueryCurrentPage();
});
</script>
