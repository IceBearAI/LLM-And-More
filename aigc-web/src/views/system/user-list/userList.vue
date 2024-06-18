<template>
  <BaseBreadcrumb :title="page.title" :breadcrumbs="breadcrumbs"></BaseBreadcrumb>
  <UiParentCard>
    <v-row>
      <v-col cols="12" lg="3" md="4" sm="6">
        <v-text-field
          v-model="searchData.email"
          label="请输入邮箱"
          hide-details
          clearable
          @keyup.enter="doQueryFirstPage"
          @click:clear="doQueryFirstPage"
        ></v-text-field>
      </v-col>
      <v-col cols="12" lg="3" md="4" sm="6">
        <Select
          @change="doQueryFirstPage"
          label="请选择是否是ldap用户"
          :mapDictionary="{ code: 'boolean' }"
          v-model="searchData.isLdap"
        ></Select>
      </v-col>
      <v-col cols="12" lg="3" md="4" sm="6">
        <ButtonsInForm>
          <v-btn color="primary" @click="onAdd">添加用户</v-btn>
        </ButtonsInForm>
      </v-col>
      <v-col cols="12">
        <AlertBlock>修改之后将实时生效，请谨慎操作！</AlertBlock>
      </v-col>

      <v-col cols="12">
        <TableWithPager @query="doTableQuery" ref="tableWithPagerRef" :infos="tableInfos">
          <el-table-column label="昵称" prop="nickname" min-width="200px" show-overflow-tooltip></el-table-column>
          <el-table-column label="邮箱" prop="email" min-width="180px" show-overflow-tooltip></el-table-column>
          <el-table-column label="ldap用户" min-width="100px">
            <template #default="{ row }">
              <ChipBoolean v-model="row.isLdap"></ChipBoolean>
            </template>
          </el-table-column>
          <el-table-column label="是否启用" min-width="100px">
            <template #default="{ row }">
              <ChipStatus v-model="row.status"></ChipStatus>
            </template>
          </el-table-column>
          <el-table-column label="默认语言" min-width="110px" show-overflow-tooltip>
            <template #default="{ row }">
              <span>{{ getLabels([["system_language", row.language]]) }}</span>
            </template>
          </el-table-column>
          <el-table-column label="创建时间" min-width="165px">
            <template #default="{ row }">
              {{ format.dateFormat(row.createdAt, "YYYY-MM-DD HH:mm:ss") }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="120px" fixed="right">
            <template #default="{ row }">
              <ButtonsInTable :buttons="getButtons(row)" />
            </template>
          </el-table-column>
        </TableWithPager>
      </v-col>
    </v-row>
  </UiParentCard>
  <ConfirmByInput ref="refConfirmDelete" @submit="doDelete">
    <template #text>
      您将要删除 <span class="text-primary font-weight-black">{{ confirmData.nickname }}</span> 用户<br />
      确定要继续吗？
    </template>
  </ConfirmByInput>
  <CreateUserPane ref="createUserPaneRef" @submit="doQueryFirstPage" />
</template>
<script setup lang="ts">
import { reactive, ref, onMounted } from "vue";
import BaseBreadcrumb from "@/components/shared/BaseBreadcrumb.vue";
import UiParentCard from "@/components/shared/UiParentCard.vue";
import AlertBlock from "@/components/ui/AlertBlock.vue";
import ConfirmByInput from "@/components/business/ConfirmByInput.vue";
import CreateUserPane from "./components/CreateUserPane.vue";
import ChipBoolean from "@/components/ui/ChipBoolean.vue";
import ChipStatus from "@/components/ui/ChipStatus.vue";
import { http, format } from "@/utils";
import { useMapRemoteStore } from "@/stores";

const { loadDictTree, getLabels } = useMapRemoteStore();

const page = ref({ title: "用户列表" });
const breadcrumbs = ref([
  {
    text: "系统管理",
    disabled: false,
    href: "#"
  },
  {
    text: "用户列表",
    disabled: true,
    href: "#"
  }
]);
const searchData = reactive({
  email: "",
  isLdap: null
});
const tableInfos = reactive({
  list: [],
  total: 0
});
const createUserPaneRef = ref();
const tableWithPagerRef = ref();
const refConfirmDelete = ref();
const confirmData = reactive({
  nickname: "",
  id: ""
});

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

const onDelete = row => {
  confirmData.nickname = row.nickname;
  confirmData.id = row.id;
  refConfirmDelete.value.show({
    width: "450px",
    confirmText: confirmData.nickname
  });
};

const doDelete = async (options = {}) => {
  const [err, res] = await http.delete({
    ...options,
    showSuccess: true,
    url: `/accounts/${confirmData.id}`
  });
  if (res) {
    refConfirmDelete.value.hide();
    doQueryCurrentPage();
  }
};

const doTableQuery = async (options = {}) => {
  await loadDictTree(["system_language"]);
  const [err, res] = await http.get({
    url: "/accounts",
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

const onAdd = () => {
  createUserPaneRef.value.show({
    title: "添加用户",
    operateType: "add"
  });
};

const onEdit = info => {
  createUserPaneRef.value.show({
    title: "编辑用户",
    infos: info,
    operateType: "edit"
  });
};

onMounted(() => {
  doTableQuery();
});
</script>
