<template>
  <v-row>
    <v-col cols="12" class="d-flex justify-space-between align-center">
      <FileBreadcrumb />
    </v-col>
    <v-col cols="12">
      <TableWithPager
        v-if="currentShowType === 'tree'"
        @query="doTableQuery"
        ref="refTableWithPager"
        :infos="tableInfos"
        :showPager="false"
        :border="false"
        header-cell-class-name=""
        cell-class-name=""
      >
        <el-table-column label="文件名" min-width="220px">
          <template #default="{ row }">
            <a class="inline-flex items-center hover:underline cursor-pointer overflow-hidden" @click="fileClick(row)">
              <component :is="row.isDir ? IconFolder : IconFile" class="mr-2 shrink-0" :size="18" />
              <span>{{ row.name }}</span>
            </a>
          </template>
        </el-table-column>
        <el-table-column label="大小" min-width="80px">
          <template #default="{ row }">
            {{ format.getFileSize(row.size) }}
          </template>
        </el-table-column>
        <el-table-column label="更新时间" min-width="165px">
          <template #default="{ row }">
            {{ format.dateFormat(row.updatedAt, "YYYY-MM-DD HH:mm:ss") }}
          </template>
        </el-table-column>
      </TableWithPager>
      <FilePreview
        ref="filePreviewRef"
        v-else-if="currentShowType === 'file'"
        :content="fileInfos.fileContent"
        :file-info="fileInfos.fileInfo"
      />
    </v-col>
  </v-row>
</template>
<script setup lang="ts">
import { reactive, ref, onMounted, computed, watch } from "vue";
import { http, format, map } from "@/utils";
import { useRoute, useRouter } from "vue-router";
import { IconFile, IconFolder } from "@tabler/icons-vue";
import FileBreadcrumb from "./FileBreadcrumb.vue";
import FilePreview from "./FilePreview.vue";
import { toast } from "vue3-toastify";
import { previewCodeSuffix, previewMarkdownSuffix } from "../../modelList";

const route = useRoute();
const router = useRouter();
const { modelName } = route.query;
const refTableWithPager = ref();
const tableInfos = reactive({
  list: [],
  total: 0
});
const fileInfos = reactive({
  fileContent: "",
  fileInfo: {}
});
const currentShowType = ref("tree");
const filePreviewRef = ref();

const currentFilePath = computed(() => {
  return route.query.filePath || "";
});

const doTableQuery = async (options = {}) => {
  let showLoading = true;
  if (currentShowType.value === "tree") {
    showLoading = refTableWithPager.value.el;
  } else if (currentShowType.value === "file") {
    showLoading = filePreviewRef.value.$el;
  }
  const [err, res] = await http.get({
    url: `/models/${modelName}/tree`,
    showLoading,
    data: {
      path: currentFilePath.value
    }
  });

  if (res) {
    currentShowType.value = res.object;
    if (currentShowType.value === "tree") {
      tableInfos.list = res.tree || [];
    } else if (currentShowType.value === "file") {
      fileInfos.fileInfo = res.fileInfo;
      fileInfos.fileContent = res.fileContent;
    }
  } else {
    tableInfos.list = [];
  }
};

const fileClick = row => {
  const query = route.query;
  const fileType = row.name.split(".").pop();
  const filePath = `${currentFilePath.value}/${row.name}`;
  const previewSuffix = [...previewCodeSuffix, ...previewMarkdownSuffix];
  if (row.isDir || previewSuffix.includes(fileType)) {
    router.push({
      path: "/model/model-list/detail",
      query: {
        ...query,
        filePath
      }
    });
  } else {
    toast.warning("该文件类型暂不支持预览！");
  }
};

watch(currentFilePath, val => {
  doTableQuery();
});

onMounted(() => {
  doTableQuery();
});
</script>
<style lang="scss">
.v-tooltip__content {
  pointer-events: initial;
}
</style>
