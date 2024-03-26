<template>
  <div class="compo-pager flex hv-center">
    <el-pagination
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
      :current-page="currentPage"
      :page-sizes="pageSizes"
      :page-size="pageSize"
      layout="total, sizes, prev, pager, next, jumper"
      :total="Number(total)"
      class="pager"
    >
    </el-pagination>
  </div>
</template>

<script setup>
import { ref, reactive, toRefs, inject, computed } from "vue";

const props = defineProps({
  total: {
    //数据总量
    type: [Number, String],
    default: 0
  },
  pageSizes: {
    type: Array,
    default: () => [10, 20, 40, 60, 120]
  }
});

const emit = defineEmits(["query", "update:modelValue"]);
const provideTableWithPager = inject("provideTableWithPager", null);

const state = reactive({
  currentPage: 1,
  pageSize: props.pageSizes[0]
});

let { currentPage, pageSize } = toRefs(state);

const handleCurrentChange = val => {
  state.currentPage = val;
  query();
};

const handleSizeChange = val => {
  state.currentPage = 1;
  state.pageSize = val;
  query();
};

const query = (options = {}) => {
  if (typeof options.page == "number") {
    let { page } = options;
    state.currentPage = Math.max(page, 1);
  }
  if (provideTableWithPager) {
    provideTableWithPager.page = state.currentPage;
    provideTableWithPager.pageSize = state.pageSize;
  }

  emit("query", {
    pageSize: state.pageSize, //数据抓取量
    page: state.currentPage,
    ...options
  });

  emit("update:modelValue", {
    pageSize: state.pageSize, //数据抓取量
    page: state.currentPage
  });
};

defineExpose({
  query
});
</script>

<style lang="scss">
.compo-pager {
  --height-pager: 50px;
  &.large {
    --height-pager: 50px;
    .el-pagination {
      font-size: var(--el-font-size-large);
      .el-input__wrapper {
        font-size: var(--el-font-size-large);
      }
    }
    .el-input {
      --el-input-height: 36px;
    }
  }
  &.default {
    --height-pager: 40px;
    .el-pagination {
      font-size: var(--el-font-size-default);
      .el-input__wrapper {
        font-size: var(--el-font-size-default);
      }
    }
    .el-input {
      --el-input-height: 30px;
    }
  }
  &.small {
    --height-pager: 32px;
    .el-pagination {
      font-size: var(--el-font-size-small);
      .el-input__wrapper {
        font-size: var(--el-font-size-small);
      }
    }
    .el-input {
      --el-input-height: 26px;
    }
  }

  height: var(--height-pager);
  .pager {
    // margin: 0 auto;
    margin-bottom: -5px;
  }
}
</style>
