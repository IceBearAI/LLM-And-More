<template>
  <v-row class="my-0" ref="contentRef">
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
    <v-col cols="12">
      <TableWithPager
        row-class-name="cursor-pointer"
        @row-click="onSelect"
        @query="doTableQuery"
        ref="tableWithPagerRef"
        :infos="tableInfos"
        height="400px"
      >
        <el-table-column label="标准问法" prop="question" show-overflow-tooltip></el-table-column>
        <el-table-column label="答案" prop="output" show-overflow-tooltip></el-table-column>
        <el-table-column label="意图" prop="intent" show-overflow-tooltip></el-table-column>
      </TableWithPager>
    </v-col>
  </v-row>
</template>
<script setup lang="ts">
import { reactive, ref, onMounted } from "vue";
import { http } from "@/utils";

interface IProps {
  intentId: string;
}
interface IEmits {
  (e: "onSelect", val: any): void;
}

const props = withDefaults(defineProps<IProps>(), {
  intentId: ""
});
const emits = defineEmits<IEmits>();

const searchData = reactive({
  keyword: ""
});
const tableInfos = reactive({
  list: [],
  total: 0
});
const tableWithPagerRef = ref();
const contentRef = ref();

const doTableQuery = async (options = {}) => {
  const [err, res] = await http.get({
    url: `/intent/${props.intentId}/document/list`,
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
  emits("onSelect", row);
};

onMounted(() => {
  doQueryFirstPage();
});
</script>
