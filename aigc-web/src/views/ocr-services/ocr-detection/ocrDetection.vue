<template>
  <BaseBreadcrumb :title="page.title" :breadcrumbs="breadcrumbs"></BaseBreadcrumb>
  <UiParentCard>
    <v-row>
      <v-col cols="12" lg="3" md="4" sm="6">
        <Select
          v-model="searchData.fileType"
          :mapDictionary="{ code: 'ocr_file_type' }"
          label="请选择文件类型"
          hide-details
          @change="doQueryFirstPage"
        ></Select>
      </v-col>
      <v-col cols="12" lg="3" md="4" sm="6">
        <ButtonsInForm>
          <v-btn color="primary" @click="onAdd">创建</v-btn>
        </ButtonsInForm>
      </v-col>

      <v-col cols="12">
        <TableWithPager @query="doTableQuery" ref="tableWithPagerRef" :infos="tableInfos">
          <el-table-column label="文件" min-width="150px">
            <template #default="{ row }">
              <template v-if="row.fileType === 'image'">
                <v-avatar size="80" rounded="md">
                  <img :src="row.inputS3Url" alt="图片" height="80" />
                </v-avatar>
              </template>
              <template v-else>
                <a class="link" :href="row.inputS3Url" target="_blank">
                  <component class="align-top" :is="map.fileTypes[row.fileType].icon"></component>
                </a>
              </template>
            </template>
          </el-table-column>
          <el-table-column label="文件类型" min-width="100px">
            <template #default="{ row }"> {{ map.fileTypes[row.fileType].text }} </template>
          </el-table-column>
          <el-table-column label="文本内容" min-width="100px">
            <template #default="{ row }">
              <a class="link" href="javascript: void(0)" @click="onLog(row.text)">
                <IconFileText class="align-top" :size="20" />
              </a>
            </template>
          </el-table-column>
          <el-table-column label="可视化图片" min-width="150px">
            <template #default="{ row }">
              <template v-if="row.viewImgs.length === 0"> -- </template>
              <template v-else>
                <el-image
                  style="width: 80px; height: 80px"
                  :src="row.viewImgs[0]"
                  :zoom-rate="1.2"
                  :max-scale="7"
                  :min-scale="0.2"
                  :preview-src-list="row.viewImgs"
                  :initial-index="4"
                  :preview-teleported="true"
                  fit="cover"
                />
              </template>
            </template>
          </el-table-column>
          <el-table-column label="页数" prop="pageNum" min-width="80px"></el-table-column>
          <el-table-column label="接口耗时" prop="duration" min-width="100px"></el-table-column>
          <el-table-column label="操作人" prop="operatorEmail" min-width="150px" show-overflow-tooltip></el-table-column>
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
  </UiParentCard>
  <DialogLog ref="refDialogLog" />
  <ConfirmByClick ref="refConfirmDelete" @submit="doDelete">
    <template #text>确认要删除这条ocr记录吗？</template>
  </ConfirmByClick>
  <CreateDetectionPane ref="createDetectionPaneRef" @submit="doQueryFirstPage" />
</template>
<script setup lang="ts">
import { reactive, ref, onMounted } from "vue";
import BaseBreadcrumb from "@/components/shared/BaseBreadcrumb.vue";
import UiParentCard from "@/components/shared/UiParentCard.vue";
import ConfirmByClick from "@/components/business/ConfirmByClick.vue";
import DialogLog from "@/components/ui/log/DialogLog.vue";
import CreateDetectionPane from "./components/CreateDetectionPane.vue";
import { IconFileText } from "@tabler/icons-vue";

import { http, format, map } from "@/utils";

const page = ref({ title: "OCR检测识别" });
const breadcrumbs = ref([
  {
    text: "OCR服务",
    disabled: false,
    href: "#"
  },
  {
    text: "OCR检测识别",
    disabled: true,
    href: "#"
  }
]);
const searchData = reactive({
  fileType: null
});
const tableInfos = reactive({
  list: [],
  total: 0
});
const refDialogLog = ref();
const createDetectionPaneRef = ref();
const tableWithPagerRef = ref();
const refConfirmDelete = ref();
const confirmData = reactive({
  id: ""
});

const onLog = content => {
  refDialogLog.value.show({
    content,
    title: "文本内容"
  });
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

const onDelete = row => {
  confirmData.id = row.id;
  refConfirmDelete.value.show({
    width: "450px"
  });
};

const doDelete = async (options = {}) => {
  const [err, res] = await http.delete({
    ...options,
    showSuccess: true,
    url: `/ocr/recognition/${confirmData.id}`
  });
  if (res) {
    refConfirmDelete.value.hide();
    doQueryCurrentPage();
  }
};

const doTableQuery = async (options = {}) => {
  const [err, res] = await http.get({
    url: "/ocr/list",
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
  createDetectionPaneRef.value.show({
    title: "创建",
    operateType: "add"
  });
};

onMounted(() => {
  doQueryCurrentPage();
});
</script>
