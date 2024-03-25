<template>
  <v-row>
    <v-col cols="12" lg="3" md="4" sm="6">
      <v-text-field
        v-model="searchData.name"
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
        <el-table-column label="ID" min-width="200px" show-overflow-tooltip class-name="link-ellipsis-color">
          <template #default="{ row }">
            <a href="javascript: void(0)" class="link" @click="onEdit(row)">{{ row.uuid }}</a>
          </template>
        </el-table-column>
        <el-table-column label="问题" min-width="200px" show-overflow-tooltip class-name="link-ellipsis-color">
          <template #default="{ row }">
            <a href="javascript: void(0)" class="link" @click="onEdit(row)">{{ row.title }}</a>
          </template>
        </el-table-column>
        <el-table-column label="对话轮次" prop="turns" width="180px"></el-table-column>
        <el-table-column label="更新时间" min-width="165px">
          <template #default="{ row }">
            {{ format.dateFormat(row.updatedAt, "YYYY-MM-DD HH:mm:ss") }}
          </template>
        </el-table-column>
        <el-table-column label="添加时间" min-width="165px">
          <template #default="{ row }">
            {{ format.dateFormat(row.createdAt, "YYYY-MM-DD HH:mm:ss") }}
          </template>
        </el-table-column>
        <el-table-column label="操作人" prop="creatorEmail" min-width="150px" show-overflow-tooltip></el-table-column>
        <el-table-column label="操作" width="80px" fixed="right">
          <template #default="{ row }">
            <ButtonsInTable :onlyOne="true" :buttons="getButtons(row)" />
          </template>
        </el-table-column>
      </TableWithPager>
    </v-col>
  </v-row>
  <ConfirmByClick ref="refConfirmDiscard" @submit="doDiscard">
    <template #text>
      <span class="font-weight-bold">这是进行一项操作时必须了解的重要信息</span><br />
      您将要废弃<span class="text-primary ml-1">{{ currentOperateData.title }}</span
      ><br />
      确定要继续吗？
    </template>
  </ConfirmByClick>
  <CreateQA ref="createQAPaneRef" @submit="doQueryFirstPage" :uuid="uuid" />
</template>
<script setup lang="ts">
import { reactive, ref, onMounted } from "vue";
import ConfirmByClick from "@/components/business/ConfirmByClick.vue";
import CreateQA from "./CreateQA.vue";
import { useRoute } from "vue-router";

import { http, format } from "@/utils";

const route = useRoute();

const { uuid } = route.query as { uuid: string };

const searchData = reactive({
  name: ""
});
const tableInfos = reactive({
  list: [],
  total: 0
});
const createQAPaneRef = ref();
const tableWithPagerRef = ref();
const refConfirmDiscard = ref();
const currentOperateData = reactive({
  title: "",
  id: ""
});

const getButtons = row => {
  let ret = [];
  ret.push({
    text: "废弃",
    color: "error",
    click() {
      onDiscard(row);
    }
  });
  return ret;
};

const onDiscard = row => {
  currentOperateData.title = row.title;
  currentOperateData.id = row.uuid;
  refConfirmDiscard.value.show({
    width: "450px"
  });
};

const doDiscard = async (options = {}) => {
  const [err, res] = await http.delete({
    ...options,
    showSuccess: true,
    url: `/api/datasets/${uuid}/samples/${currentOperateData.id}`
  });
  if (res) {
    refConfirmDiscard.value.hide();
    doTableQuery();
  }
};

const doTableQuery = async (options = {}) => {
  const [err, res] = await http.get({
    url: `/api/datasets/${uuid}/samples`,
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
  createQAPaneRef.value.show({
    title: "添加问答",
    operateType: "add"
  });
};

const onEdit = row => {
  createQAPaneRef.value.show({
    title: "编辑问答",
    infos: {
      uuid: row.uuid,
      messages: row.messages
    },
    operateType: "edit"
  });
};

onMounted(() => {
  doTableQuery();
});
</script>
