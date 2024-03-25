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
                v-model="formData.name"
                label="搜索场景名称"
                hide-details
                clearable
                variant="outlined"
                @keyup.enter="doQueryFirstPage"
                @click:clear="doQueryFirstPage"
              ></v-text-field>
            </div>
            <ButtonsInForm>
              <v-btn color="primary" @click="onAdd">添加场景</v-btn>
            </ButtonsInForm>
          </v-col>

          <v-col cols="12">
            <AlertBlock> 修改之后将实时生效，请谨慎操作！ </AlertBlock>
          </v-col>
          <v-col cols="12">
            <TableWithPager @query="doTableQuery" ref="refTableWithPager" :infos="state.tableInfos">
              <!-- <el-table-column label="唯一标识" width="120px" fixed="left">
                <template #default="{ row }">
                  <span class1="text-primary">{{ row.id }}</span>
                </template>
              </el-table-column> -->

              <el-table-column label="模型数量" min-width="100px">
                <template #default="{ row }">
                  <span pointer>
                    {{ row.model.num }}
                    <v-menu open-on-hover location="end" activator="parent" v-if="row.model.num > 0">
                      <v-sheet rounded="md" width="200" elevation="10">
                        <v-list class="px-4">
                          <v-list-item v-for="(item, index) in row.model.list" :key="index">
                            <v-list-item-title class="text-subtitle-1 font-weight-regular">
                              {{ index + 1 }}. {{ item.modelName }}
                            </v-list-item-title>
                          </v-list-item>
                        </v-list>
                      </v-sheet>
                    </v-menu>
                  </span>
                </template>
              </el-table-column>

              <!-- <el-table-column label="名称" width="150px">
                <template #default="{ row }">
                  <span>{{ row.name }}</span>
                </template>
              </el-table-column> -->

              <el-table-column label="别名" width="150px">
                <template #default="{ row }">
                  <span>{{ row.alias }}</span>
                </template>
              </el-table-column>
              <!-- <el-table-column label="项目" width="120px">
                <template #default="{ row }">
                  <span>{{ row.projectName }}</span>
                </template>
              </el-table-column>
              <el-table-column label="服务" width="120px">
                <template #default="{ row }">
                  <span>{{ row.serviceName }}</span>
                </template>
              </el-table-column> -->

              <el-table-column label="配额" min-width="100px">
                <template #default="{ row }">
                  {{ row.quota }}
                </template>
              </el-table-column>
              <!-- <el-table-column label="模型数量" min-width="100px">
                <template #default="{ row }">
                  <v-menu :close-on-content-click="false" location="bottom">
                    <template v-slot:activator="{ props }"> {{ row.model.num }} </template>
                    <v-sheet rounded="md" width="200" elevation="10">
                      <v-list class="theme-list">
                        <v-list-item-title class="text-subtitle-1 font-weight-regular">
                          aaaaaa
                          <span class="text-disabled text-subtitle-1 pl-2">(aaaaaaaaa)</span>
                        </v-list-item-title>

                        <v-list-item-title class="text-subtitle-1 font-weight-regular">
                          bbbbbb
                          <span class="text-disabled text-subtitle-1 pl-2">(bbbbbb)</span>
                        </v-list-item-title>
                      </v-list>
                    </v-sheet>
                  </v-menu>
                </template>
              </el-table-column> -->
              <el-table-column label="ApiKey" min-width="120px" show-overflow-tooltip>
                <template #default="{ row }">
                  <span v-copy="row.apiKey"> {{ row.apiKey }}</span>
                </template>
              </el-table-column>
              <!-- <el-table-column label="标签" min-width="200px">
                <template #default="{ row }">
                  {{ row.updateAt }}
                </template>
              </el-table-column> -->
              <el-table-column label="备注" min-width="200px">
                <template #default="{ row }">
                  {{ row.remark }}
                </template>
              </el-table-column>
              <el-table-column label="更新时间" min-width="200px">
                <template #default="{ row }">
                  {{ format.dateFormat(row.updatedAt, "YYYY-MM-DD HH:mm:ss") }}
                </template>
              </el-table-column>
              <el-table-column label="最后操作人" min-width="200px" show-overflow-tooltip>
                <template #default="{ row }">
                  {{ row.lastOperator }}
                </template>
              </el-table-column>

              <el-table-column label="负责人" min-width="120px" show-overflow-tooltip>
                <template #default="{ row }">
                  <span v-copy="row.email"> {{ row.email }}</span>
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

  <PaneScene ref="refPaneScene" @submit="doQueryFirstPage" />

  <ConfirmByInput ref="refConfirmDelete" @submit="doDelete">
    <template #text>
      您将要删除 <span class="text-primary font-weight-black">{{ state.selectedInfo.alias }}</span> 场景，删除之后，<br />
      若有使用该场景的应用场景都将取消授权。<br />
      确定要继续吗？
    </template>
  </ConfirmByInput>
</template>
<script setup>
import { reactive, toRefs, ref, onMounted } from "vue";
import BaseBreadcrumb from "@/components/shared/BaseBreadcrumb.vue";
import UiParentCard from "@/components/shared/UiParentCard.vue";
import AlertBlock from "@/components/ui/AlertBlock.vue";
import ChipBoolean from "@/components/ui/ChipBoolean.vue";

import PaneScene from "./components/PaneScene.vue";
import ConfirmByInput from "@/components/business/ConfirmByInput.vue";
import { http, format } from "@/utils";

const state = reactive({
  style: {},
  formData: {
    name: ""
  },
  selectedInfo: {},
  tableInfos: {
    list: [],
    total: ""
  }
});
const { style, formData } = toRefs(state);

const refPaneScene = ref();
const refConfirmDelete = ref();
const refTableWithPager = ref();

const page = ref({ title: "场景列表" });
const breadcrumbs = ref([
  {
    text: "场景管理",
    disabled: false,
    href: "#"
  },
  {
    text: "场景列表",
    disabled: true,
    href: "#"
  }
]);

const getButtons = row => {
  let ret = [];
  ret.push({
    text: "删除",
    color: "error",
    click() {
      onDelete(row);
    }
  });
  ret.push({
    text: "编辑",
    color: "info",
    click() {
      onEdit(row);
    }
  });

  return ret;
};

const doTableQuery = async (options = {}) => {
  const [err, res] = await http.get({
    url: "/api/channels",
    showLoading: refTableWithPager.value.el,
    data: {
      ...state.formData,
      ...options
    }
  });
  if (res) {
    state.tableInfos.list = res.list;
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
  state.selectedInfo = info;
  refConfirmDelete.value.show({
    width: "400px",
    confirmText: state.selectedInfo.alias
  });
};

const doDelete = async (options = {}) => {
  const [err, res] = await http.delete({
    ...options,
    showSuccess: true,
    url: `/api/channels/${state.selectedInfo.id}`
  });
  if (res) {
    refConfirmDelete.value.hide();
    doQueryFirstPage();
  }
};

const onAdd = () => {
  refPaneScene.value.show({
    title: "添加场景",
    operateType: "add"
  });
};

const onEdit = info => {
  refPaneScene.value.show({
    title: "编辑场景",
    infos: info,
    operateType: "edit"
  });
};

onMounted(() => {
  doQueryFirstPage();
});
</script>
<style lang="scss"></style>
