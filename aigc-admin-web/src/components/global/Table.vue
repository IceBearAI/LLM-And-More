<template>
  <div ref="refBox" class="table-main">
    <el-table
      border
      ref="refTable"
      size="default"
      class="compo-table"
      :key="state.key"
      :header-row-style="{
        fontWeight: 'bolder',
        fontSize: '12px',
        color: '#333'
      }"
      header-cell-class-name="is-center"
      cell-class-name="is-center"
      :data="tableList"
      row-key="id"
      stripe
      highlight-current-row
      v-bind="{ ...attrs, ...compuHeight }"
    >
      <!-- <el-table-column align="center" width="35" v-if="isRowDragable" class-name="td-dragable">
        <template #default>
          <div class="cell-dragable">
            <IFont size="12px" color="gray" title="拖拽行调整顺序">dragable</IFont>
          </div>
        </template>
      </el-table-column> -->

      <el-table-column type="selection" v-if="compuShowSelection" fixed="left" />

      <el-table-column
        v-if="compuShowIndex"
        :label="compuShowIndex"
        align="center"
        fixed="left"
        type="index"
        width="60"
        :index="indexMethod"
      >
      </el-table-column>
      <slot />
      <template #empty>
        <!--已拉取过数据，仍为空，显示对应组件-->
        <NoData v-if="provideTableWithPager.isRequested && showEmpty" />
        <!--初始化的时候，不显示空数据提示-->
        <div v-else></div>
      </template>
    </el-table>
  </div>
</template>
<script setup>
import { ref, onMounted, nextTick, reactive, inject, watch, computed, useAttrs } from "vue";
import Sortable from "sortablejs";
import { useWindowSize, useResizeObserver, useDebounceFn, useVModel } from "@vueuse/core";
import jquery from "jquery";
import _ from "lodash";

const provideTableWithPager = inject("provideTableWithPager") || {
  isRequested: true
};

const refTable = ref();
provideTableWithPager.compoTable = refTable;

/**
 * 属性
 *  1. showIndex ，是否显示序号，默认显示，不想显示的话，传 :showIndex="false"
 *  2. @selection-change="" 传入函数，显示多选框，默认不显示
 */
const attrs = useAttrs();

const props = defineProps({
  isRowDragable: {
    //行是否可拖拽
    type: Boolean,
    default: false
  },
  minHeight: {
    //elTable 原本没有此属性，扩展
    type: [Number, String],
    default() {
      return "";
    }
  },
  list: {
    type: Array,
    default: () => []
  },
  showEmpty: {
    type: Boolean,
    default: true
  },
  showIndex: {
    //是否显示序号，默认显示，不想显示的话，传 :showIndex="false"，默认此列列标题为"序号"，要改为其他的标题，showIndex直接指定为标题名
    type: [Boolean, String],
    default: false
  },
  showSelection: {
    // 是否显示多选框，默认不显示，需显示的话，传 showSelection 或 :showSelection="true"
    type: [Boolean, String],
    default: false
  }
});

const refBox = ref();

var state = reactive({
  key: +new Date(),
  boxHeight: "auto",
  height: 300
});

const emit = defineEmits(["update:list"]);

const tableList = useVModel(props, "list", emit);

watch(
  () => {
    return props.isRowDragable;
  },
  () => {
    nextTick(() => {
      rowDrop();
    });
  },
  {
    immediate: true
  }
);

watch(
  () => {
    return props.list;
  },
  () => {
    refTable.value?.setScrollTop(0);
  }
);

const compuHeight = computed(() => {
  if (props.list.length == 0) {
    //没有数据时
    if (props.minHeight) {
      //指定了最小高度
      return { height: (props.minHeight + "px").replace("pxpx", "px") };
    } else if (attrs.maxHeight) {
      // 配置了 maxHeight ，就不是自适应高度表格，默认 300px
      return { height: "300px" };
    } else if (attrs.class?.includes("flex-full")) {
      //设置了flex-full，通过样式撑开高度，不进行干预
      return {};
    } else if (attrs.height) {
      //设置了height，不进行干预
      return {};
    }
  }
  return {};
});

//kings todo 后续调整，请求还未发送，表格的index就开始跳的问题
const compuShowIndex = computed(() => {
  let { showIndex } = props;
  if (typeof showIndex == "boolean") {
    return showIndex ? "序号" : "";
  } else {
    return showIndex;
  }
});

const compuShowSelection = computed(() => {
  let { onSelectionChange } = attrs;
  if (typeof onSelectionChange == "function") {
    return true;
  } else {
    return false;
  }
});

const indexMethod = index => {
  let { page = 1, pageSize = 10 } = provideTableWithPager;
  return index + 1 + (page - 1) * pageSize;
};

const rowDrop = async () => {
  const tbody = document.querySelector(".el-table__body tbody", refBox.value);
  Sortable.create(tbody, {
    handle: ".cell-dragable",
    animation: 150,
    async onEnd({ newIndex, oldIndex }) {
      const currRow = tableList.value.splice(oldIndex, 1)[0];

      //防止页面抖动
      state.boxHeight = document.querySelector(".compo-table", refBox.value).getBoundingClientRect().height + "px";

      tableList.value.splice(newIndex, 0, currRow);
      {
        //更改了表格list排序，但界面没有跟着变化，调用 $forceUpdate也不好使，感觉框架的Bug
        state.key = +new Date();
        await nextTick();
        rowDrop();
      }
      state.boxHeight = "auto";
    }
  });
};

defineExpose({
  compoTable: refTable
});
</script>
<style lang="scss">
.el-table {
  &.compo-table {
    font-size: 13px;
    th {
      font-size: 15px;
    }
    .el-table__inner-wrapper {
      &::before {
        // 表格底部边框线
        // display: none;
      }
    }
  }
  .el-table__cell {
    padding: 16px 0 !important;
  }
  .cell-dragable {
    cursor: move;
  }
  .el-table__body tr.current-row > td.el-table__cell {
    background-color: var(--el-table-current-row-bg-color);
  }
}
</style>
