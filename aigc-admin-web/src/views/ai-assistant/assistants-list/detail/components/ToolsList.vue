<template>
  <v-row class="my-0">
    <v-col cols="12" lg="3" md="4" sm="6">
      <v-text-field
        v-model="searchData.name"
        label="请输入工具名称"
        hide-details
        clearable
        @keyup.enter="doQueryFirstPage"
        @click:clear="doQueryFirstPage"
      ></v-text-field>
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
        <el-table-column label="工具ID" prop="toolId" min-width="220px"></el-table-column>
        <el-table-column label="工具名称" prop="name" min-width="200px" show-overflow-tooltip></el-table-column>
        <el-table-column label="工具描述" prop="description" min-width="300px" show-overflow-tooltip></el-table-column>
        <el-table-column label="工具类型" width="100px">
          <template #default="{ row }">
            {{ getLabels([["assistant_tool_type", row.toolType]]) }}
          </template>
        </el-table-column>
      </TableWithPager>
    </v-col>
  </v-row>
</template>
<script setup lang="ts">
import { reactive, ref, onMounted } from "vue";
import { http } from "@/utils";
import { useMapRemoteStore } from "@/stores";

interface IEmits {
  (e: "selected", val: any): void;
}

const emits = defineEmits<IEmits>();

const { getLabels, loadDictTree } = useMapRemoteStore();
loadDictTree(["assistant_tool_type"]);
const searchData = reactive({
  name: ""
});
const tableInfos = reactive({
  list: [],
  total: 0
});
const tableWithPagerRef = ref();

const doTableQuery = async (options = {}) => {
  const [err, res] = await http.get({
    url: "/tools/list",
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

const onSelect = row => {
  emits("selected", row);
};

onMounted(() => {
  doTableQuery();
});
</script>
<style scoped lang="scss">
// vuetify 中dialog添加了contain 属性导致tooltip定位便宜。
:deep(.el-table) {
  transform: scale(1);
}
</style>
