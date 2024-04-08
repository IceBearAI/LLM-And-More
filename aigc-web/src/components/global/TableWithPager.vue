<template>
  <div class="compo-tableWithPager" ref="el">
    <Table
      ref="refTable"
      :list="infos.list"
      v-bind="_.omit(attrs, ['class', 'style'])"
      :isRowDragable="isRowDragable"
      :showPager="showPager"
      class="flex-full"
    >
      <slot />
    </Table>
    <Pager
      ref="refPager"
      v-model="state.pageInfo"
      :total="infos.total"
      @query="doQuery"
      v-show="showPager && (infos.list || infos.list.length > 0)"
    />
  </div>
</template>

<script setup lang="ts">
import Table from "./Table.vue";
import Pager from "./Pager.vue";
import _ from "lodash";
import { reactive, ref, watch, provide, useAttrs } from "vue";

const el = ref();
const attrs = useAttrs();
const refTable = ref();
const state = reactive({
  child: "",
  isRequested: false, //是否请求过数据
  page: 1,
  pageSize: 10,
  compoElTable: "", //表格element组件
  compoTable: refTable, //表格组件
  pageInfo: {
    //页码信息
    pageSize: "",
    page: ""
  }
});

provide("provideTableWithPager", state);
const props = defineProps({
  isRowDragable: {
    //行是否可拖拽
    type: Boolean,
    default: false
  },
  infos: {
    type: Object,
    default() {
      return {
        list: [], //表格渲染数据
        total: 0 //数据总量
      };
    }
  },
  showPager: {
    type: Boolean,
    default: true
  }
});
watch(
  () => [props.infos.total, props.infos.list],
  newValue => {
    state.isRequested = true;
  }
);

const emit = defineEmits(["query"]);

const refPager = ref();

const doQuery = options => {
  emit("query", {
    ...options
  });
};
defineExpose({
  query: options => {
    refPager.value.query(options);
  },
  state,
  el
});
</script>

<style lang="scss">
.compo-tableWithPager {
  .el-popper {
    max-width: 400px;
  }
}
</style>
