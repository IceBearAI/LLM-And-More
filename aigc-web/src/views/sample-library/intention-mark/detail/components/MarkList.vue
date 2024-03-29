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
        <custom-upload :file-type="['text/csv']" @after-upload="handleAfterUpload">
          <template v-slot:trigger="{ loading }">
            <v-btn :loading="loading || importLoading" color="primary">导入</v-btn>
          </template>
        </custom-upload>
      </ButtonsInForm>
    </v-col>
    <v-col cols="12">
      <TableWithPager @query="doTableQuery" ref="tableWithPagerRef" :infos="tableInfos">
        <el-table-column label="问题" prop="question" min-width="200px" show-overflow-tooltip></el-table-column>
        <el-table-column label="答案" prop="output" min-width="200px" show-overflow-tooltip></el-table-column>
        <el-table-column label="意图" min-width="200px">
          <template #default="{ row }">
            <div class="d-flex align-center">
              <div class="line1 flex-1">{{ row.intent }}</div>
              <IconPencil
                @click="editIntention(row.annotationId)"
                stroke-width="1.5"
                :size="20"
                class="text-primary cursor-pointer"
              />
            </div>
          </template>
        </el-table-column>
        <el-table-column label="更新时间" min-width="200px">
          <template #default="{ row }">{{ format.dateFormat(row.updatedAt, "YYYY-MM-DD HH:mm:ss") }} </template>
        </el-table-column>
        <el-table-column label="操作人" prop="operator" min-width="220px"></el-table-column>
      </TableWithPager>
    </v-col>
  </v-row>
  <Dialog ref="selectIntentionRef" persistent max-width="1200px" :retain-focus="false">
    <template #title> 选择意图 </template>
    <SelectIntentionList :intentId="intentId" @on-select="handleSelectIntent" />
    <!-- <ToolsList /> -->
  </Dialog>
</template>
<script setup lang="ts">
import { reactive, ref, onMounted } from "vue";
import { http, format } from "@/utils";
import { IconPencil } from "@tabler/icons-vue";
import SelectIntentionList from "./SelectIntentionList.vue";
import CustomUpload from "@/components/business/CustomUpload.vue";
import { useRoute } from "vue-router";
import ToolsList from "@/views/ai-assistant/assistants-list/detail/components/ToolsList.vue";

const route = useRoute();
const { intentId } = route.query as { intentId: string };

const searchData = reactive({
  keyword: ""
});
const tableInfos = reactive({
  list: [],
  total: 0
});
const tableWithPagerRef = ref();
const selectIntentionRef = ref();
const importLoading = ref(false);
const dialogConfig = reactive({
  annotationId: ""
});

const doTableQuery = async (options = {}) => {
  const [err, res] = await http.get({
    url: `/intent/${intentId}/annotation/list`,
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

const editIntention = id => {
  dialogConfig.annotationId = id;
  selectIntentionRef.value.show({
    showActions: false
  });
};

const handleSelectIntent = async row => {
  const [err, res] = await http.post({
    showLoading: selectIntentionRef.value.getContentEl(),
    showSuccess: true,
    url: `/intent/${intentId}/annotation/${dialogConfig.annotationId}/mark`,
    data: {
      documentId: row.documentId
    }
  });
  if (res) {
    selectIntentionRef.value.hide();
    doQueryCurrentPage();
  }
};

const handleAfterUpload = async ({ res: fileRes }) => {
  if (fileRes) {
    importLoading.value = true;
    const [err, res] = await http.post({
      showSuccess: true,
      url: `/intent/${intentId}/annotation/create`,
      data: {
        fileId: fileRes.fileId
      }
    });

    if (res) {
      doQueryFirstPage();
    }
    importLoading.value = false;
  }
};

onMounted(() => {
  doQueryCurrentPage();
});
</script>
