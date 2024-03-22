<template>
  <v-row class="my-0">
    <v-col cols="12" lg="3" md="4" sm="6">
      <UploadFile
        show-loading
        label="上传微调文件"
        :purpose="purpose"
        @upload:success="doQueryFirstPage"
        @loading="val => (upLoading = val)"
      />
    </v-col>
    <v-col cols="12">
      <TableWithPager
        @query="doTableQuery"
        ref="tableWithPagerRef"
        row-class-name="cursor-pointer"
        :infos="tableInfos"
        height="400px"
        @row-click="onSelect"
      >
        <el-table-column label="文件名" prop="filename" min-width="250px"></el-table-column>
        <el-table-column label="文件大小" width="150px">
          <template #default="{ row }">
            {{ format.getFileSize(row.size) }}
          </template>
        </el-table-column>
        <el-table-column label="文件地址" prop="s3Url" min-width="400px">
          <template #default="{ row }">
            <span v-copy="row.s3Url">{{ row.s3Url }}</span>
          </template>
        </el-table-column>
      </TableWithPager>
    </v-col>
  </v-row>
</template>
<script setup lang="ts">
import { reactive, ref, onMounted } from "vue";
import UploadFile from "@/components/business/UploadFile.vue";
import { http, format } from "@/utils";

interface IProps {
  purpose?: string;
}

interface IEmits {
  (e: "selected", val: any): void;
}

const props = withDefaults(defineProps<IProps>(), {
  purpose: ""
});
const emits = defineEmits<IEmits>();

const tableInfos = reactive({
  list: [],
  total: 0
});
const tableWithPagerRef = ref();
const upLoading = ref(false);

const doTableQuery = async (options = {}) => {
  const [err, res] = await http.get({
    url: "/files",
    showLoading: tableWithPagerRef.value.el,
    data: {
      purpose: props.purpose,
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

const onSelect = row => {
  if (upLoading.value) return;
  emits("selected", row);
};

onMounted(() => {
  doTableQuery();
});
</script>
