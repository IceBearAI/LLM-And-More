<template>
  <BaseBreadcrumb :title="page.title" :breadcrumbs="breadcrumbs"></BaseBreadcrumb>
  <v-row>
    <v-col>
      <UiParentCard>
        <v-row>
          <v-col cols="12" class="d-flex justify-space-between align-center">
            <div style="width: 300px">
              <v-text-field
                density="compact"
                v-model="state.selectedInfo.name"
                label="请输入样本名称"
                hide-details
                clearable
                variant="outlined"
                color="red"
                @keyup.enter="doQueryFirstPage"
                @click:clear="doQueryFirstPage"
              ></v-text-field>
            </div>
            <ButtonsInForm>
              <v-btn color="primary" @click="onAdd('add')">创建样本</v-btn>
              <v-btn color="primary" @click="onAdd('upload')">上传样本</v-btn>
            </ButtonsInForm>
          </v-col>

          <v-col cols="12">
            <AlertBlock> 修改之后将实时生效，请谨慎操作！ </AlertBlock>
          </v-col>
          <v-col cols="12">
            <TableWithPager @query="doTableQuery" ref="refTableWithPager" :infos="state.tableInfos">
              <el-table-column label="样本ID" width="200px" fixed="left">
                <template #default="{ row }">
                  <el-tooltip content="查看详情" placement="top">
                    <span class="link" @click="onView(row.uuid)">{{ row.uuid }}</span>
                  </el-tooltip>
                </template>
              </el-table-column>
              <el-table-column label="样本名称" width="200px" fixed="left">
                <template #default="{ row }">
                  <el-tooltip content="查看详情" placement="top">
                    <span class="link" @click="onView(row.uuid)">{{ row.name }}</span>
                  </el-tooltip>
                </template>
              </el-table-column>
              <el-table-column label="样本数量" width="110px">
                <template #default="{ row }">
                  <span>{{ row.samples }}</span>
                </template>
              </el-table-column>

              <el-table-column label="备注" min-width="200px" show-overflow-tooltip>
                <template #default="{ row }"> {{ row.remark }} </template>
              </el-table-column>
              <el-table-column label="修改时间" min-width="200px">
                <template #default="{ row }">
                  {{ format.dateFormat(row.updatedAt, "YYYY-MM-DD HH:mm:ss") }}
                </template>
              </el-table-column>
              <el-table-column label="创建时间" min-width="200px">
                <template #default="{ row }">
                  {{ format.dateFormat(row.createdAt, "YYYY-MM-DD HH:mm:ss") }}
                </template>
              </el-table-column>
              <el-table-column label="操作人" min-width="200px">
                <template #default="{ row }">
                  {{ row.creator }}
                </template>
              </el-table-column>
              <el-table-column label="操作" min-width="120px" fixed="right">
                <template #default="{ row }">
                  <ButtonsInTable :buttons="getButtons(row)" />
                </template>
              </el-table-column>
            </TableWithPager>
          </v-col>
        </v-row>
      </UiParentCard>
    </v-col>
  </v-row>

  <PaneModel ref="refPaneModel" @submit="doQueryFirstPage" />

  <ConfirmByInput ref="refConfirmDelete" @submit="doDelete">
    <template #text>
      此操作将会<span class="text-primary font-weight-black">删除</span>该样本数据集<br />
      数据集ID：<span class="text-primary font-weight-black">{{ state.formData.uuid }}</span
      ><br />你还要继续吗？
    </template>
  </ConfirmByInput>

  <ConfirmByClick ref="refConfirmByClick" @submit="onConfirmByClick">
    <template #text>
      <div>
        <AlertBlock class="mb-6"
          >该操作{{ state.subTitle }}
          <!-- <span>，导出成功之后，在创建微调时可以直接选择该数据集</span> -->
        </AlertBlock>
        <v-form ref="refForm" class="my-form">
          <v-text-field type="text" placeholder="eg: 普泽AI客服" hide-details="auto" clearable>
            <template #prepend>
              <label style="display: flex"
                >角色
                <div class="toolTip">
                  <IconInfoCircle
                    :size="16"
                    color="#bbb"
                    pointer
                    class="iconTool"
                    @mouseenter="state.showIcon = true"
                    @mouseleave="state.showIcon = false"
                  >
                  </IconInfoCircle>
                  <div v-show="state.showIcon" class="iconToolTip">AI角色</div>
                </div>
              </label>
            </template>
          </v-text-field>
          <v-text-field type="text" placeholder="eg: 普泽基金的研究人员" hide-details="auto" clearable>
            <template #prepend>
              <label style="display: flex"
                >作者
                <div class="toolTip">
                  <IconInfoCircle
                    :size="16"
                    color="#bbb"
                    pointer
                    class="iconTool"
                    @mouseenter="state.showIconRole = true"
                    @mouseleave="state.showIconRole = false"
                  >
                  </IconInfoCircle>
                  <div v-show="state.showIconRole" class="iconToolTipRole">AI的制作人及使用的数据等</div>
                </div>
              </label>
            </template>
          </v-text-field>
        </v-form>
        <v-chip-group
          mandatory
          selected-class="text-primary"
          style="margin-left: 20%"
          v-model="state.radioCheck"
          @click="changeChip"
        >
          <v-chip v-for="(tag, index) in state.exportList" :key="index" filter variant="outlined">
            {{ tag.text }}
          </v-chip>
        </v-chip-group>
      </div>
    </template>
  </ConfirmByClick>
</template>
<script setup lang="ts">
import { reactive, toRefs, ref, onMounted } from "vue";
import BaseBreadcrumb from "@/components/shared/BaseBreadcrumb.vue";
import UiParentCard from "@/components/shared/UiParentCard.vue";
import AlertBlock from "@/components/ui/AlertBlock.vue";
import PaneModel from "./components/PaneModel.vue";

import ConfirmByInput from "@/components/business/ConfirmByInput.vue";
import ConfirmByClick from "@/components/business/ConfirmByClick.vue";
import { IconInfoCircle } from "@tabler/icons-vue";
import { http, format } from "@/utils";
import { useRouter, useRoute } from "vue-router";
import _ from "lodash";
import { type ItfModel } from "./textSampleList.ts";
import { TypeButtonsInTable } from "@/components/types/components.ts";

const router = useRouter();
const route = useRoute();

const state = reactive<{
  /** 所选行，编辑操作 */
  // selectedInfo: Partial<ItfModel>;
  [x: string]: any;
}>({
  style: {},
  formData: {
    name: ""
  },
  selectedInfo: {
    name: ""
  },
  tableInfos: {
    list: [],
    total: ""
  },
  confirmByClickInfo: {
    html: "",
    action: "",
    row: {}
  },
  radioCheck: 0,
  showIcon: false,
  showIconRole: false,
  subTitle: "",
  exportList: [
    {
      text: "导出jsonl",
      minTitle: "jsonl",
      subTitle: "将导出成alpaca的训练格式文件，可在创建微调时上传该文件进行训练。"
    },
    // {
    //   text: "导出csv",
    //   minTitle: "csv",
    //   subTitle: "将导出成csv格式的文件"
    // },
    {
      text: "生成微调数据集",
      minTitle: "data",
      subTitle: "将直接生成微调数据集，在创建微调时可以直接选择该数据集"
    }
  ]
});
const { style, formData } = toRefs(state);

const refPaneModel = ref();

const refConfirmDelete = ref();
const refTableWithPager = ref();
const refConfirmByClick = ref();

const page = ref({ title: "微调样本列表" });
const breadcrumbs = ref([
  {
    text: "样本库",
    disabled: false,
    href: "#"
  },
  {
    text: "微调样本列表",
    disabled: true,
    href: "#"
  }
]);

const onView = id => {
  router.push(`/sample-library/text-sample/detail?uuid=${id}`);
};

const getButtons = (row): Array<TypeButtonsInTable> => {
  // let ret: Array<TypeButtonsInTable> = [];
  let ret = [
    {
      text: "编辑",
      color: "info",
      click() {
        onEdit(row);
      }
    },
    {
      text: "删除",
      color: "error",
      click() {
        onDelete(row);
      }
    },
    {
      text: "导出",
      color: "info",
      click() {
        refConfirmByClick.value.show({ width: "550px" });
      }
    }
  ];

  return ret;
};

const onConfirmByClick = (options = {}) => {
  state.subTitle = state.exportList[state.radioCheck].subTitle;
  // let { action, row } = state.confirmByClickInfo;
  // if (action == "deploy") {
  //   onDeploy(row, options);
  // } else if (action == "undeploy") {
  //   onUndeploy(row, options);
  // }
};
const changeChip = (options = {}) => {
  state.subTitle = state.exportList[state.radioCheck].subTitle;
  // let { action, row } = state.confirmByClickInfo;
  // if (action == "deploy") {
  //   onDeploy(row, options);
  // } else if (action == "undeploy") {
  //   onUndeploy(row, options);
  // }
};

const doTableQuery = async (options = {}) => {
  const [err, res] = await http.get({
    url: "/api/datasets/list",
    showLoading: refTableWithPager.value.el,
    data: {
      ...state.selectedInfo,
      ...options
    }
  });
  if (res) {
    state.tableInfos.list = res.list || [];
    state.tableInfos.total = res.total;
  } else {
    state.tableInfos.list = [];
    state.tableInfos.total = 0;
  }
};

const doQueryFirstPage = () => {
  refTableWithPager.value.query({ page: 1 });
};

const onDelete = info => {
  state.formData = info;
  refConfirmDelete.value.show({
    width: "400px",
    confirmText: state.formData.uuid
  });
};

const doDelete = async (options = {}) => {
  const [err, res] = await http.delete({
    ...options,
    showSuccess: true,
    url: `/api/datasets/${state.formData.uuid}`
  });
  if (res) {
    refConfirmDelete.value.hide();
    doTableQuery();
  }
};

const onAdd = type => {
  refPaneModel.value.show({
    title: type == "add" ? "创建数据集" : "上传数据集",
    operateType: type
  });
};
const onEdit = info => {
  refPaneModel.value.show({
    title: "编辑数据集",
    infos: info,
    operateType: "edit"
  });
};

onMounted(() => {
  state.subTitle = state.exportList[state.radioCheck].subTitle;
  doTableQuery();
});
</script>
<style lang="scss">
.export-title {
  color: #333;
  font-weight: 700 !important;
}
.export-sub {
  color: #fb9678;
  font-weight: 600 !important;
}

.v-sheet {
  color: #888;
}
.compo-explain {
  // display: inline-flex;
  // vertical-align: text-top;
  // svg {
  //   outline: none;
  // }
}
.toolTip {
  position: relative;
  width: 40px;
  .iconTool {
    position: absolute;
    top: 3px;
    left: 2px;
  }
  .iconToolTip,
  .iconToolTipRole {
    background: #333;
    color: #fff;
    border-radius: 4px;
    font-size: 11px;
    padding: 4px 6px;
    z-index: 12;
  }
  .iconToolTip {
    position: absolute;
    top: -26px;
    left: -14px;
  }
  .iconToolTipRole {
    position: absolute;
    top: -25px;
    left: -52px;
    width: 150px;
  }
}
</style>
./types/modelList.type.ts ./types/modelList.ts
