<template>
  <BaseBreadcrumb :title="page.title" :breadcrumbs="breadcrumbs"></BaseBreadcrumb>
  <UiParentCard>
    <v-row>
      <v-col cols="12" lg="3" md="4" sm="6">
        <v-text-field
          v-model="searchData.uuid"
          label="请输入uuid"
          hide-details
          clearable
          @keyup.enter="doQueryFirstPage"
          @click:clear="doQueryFirstPage"
        ></v-text-field>
      </v-col>
      <v-col cols="12" lg="3" md="4" sm="6">
        <v-text-field
          v-model="searchData.source"
          label="请输入来源渠道"
          hide-details
          clearable
          @keyup.enter="doQueryFirstPage"
          @click:clear="doQueryFirstPage"
        ></v-text-field>
      </v-col>
      <v-col cols="12" lg="3" md="4" sm="6">
        <ButtonsInForm>
          <v-btn color="primary" @click="onAdd">创建检测</v-btn>
        </ButtonsInForm>
      </v-col>

      <v-col cols="12">
        <TableWithPager @query="doTableQuery" ref="tableWithPagerRef" :infos="tableInfos">
          <el-table-column label="检测视频" min-width="150px">
            <template #default="{ row }">
              <a class="link" href="javascript: void(0)" @click="viewVideo(row)">
                <IconVideo class="align-top" :size="20" />
              </a>
            </template>
          </el-table-column>
          <el-table-column label="返回图片" min-width="150px">
            <template #default="{ row }">
              <template v-if="row.outputS3Url">
                <v-avatar size="80" rounded="md">
                  <img :src="row.outputS3Url" alt="返回图片" height="80" />
                </v-avatar>
              </template>
              <template v-else> -- </template>
            </template>
          </el-table-column>
          <el-table-column label="uuid" min-width="200px">
            <template #default="{ row }">
              <span v-copy="row.uuid">{{ row.uuid }}</span>
            </template>
          </el-table-column>
          <el-table-column label="是否活体" min-width="120px">
            <template #default="{ row }">
              <ChipBoolean v-model="row.isSame"></ChipBoolean>
            </template>
          </el-table-column>
          <el-table-column label="描述" prop="remark" min-width="200px"></el-table-column>
          <el-table-column label="py耗时" prop="durationPy" min-width="90px">
            <template #default="{ row }">
              {{ row.durationPy.toFixed(4) }}
            </template>
          </el-table-column>
          <el-table-column label="接口耗时" prop="duration" min-width="90px">
            <template #default="{ row }">
              {{ row.duration.toFixed(4) }}
            </template>
          </el-table-column>
          <el-table-column label="来源渠道" prop="source" min-width="100px" show-overflow-tooltip></el-table-column>
          <el-table-column label="应用场景" prop="sceneType" min-width="100px" show-overflow-tooltip></el-table-column>
          <el-table-column label="创建时间" min-width="165px">
            <template #default="{ row }">
              {{ format.dateFormat(row.createdAt, "YYYY-MM-DD HH:mm:ss") }}
            </template>
          </el-table-column>
        </TableWithPager>
      </v-col>
    </v-row>
  </UiParentCard>
  <Dialog ref="refDialogVideo" max-width="600px" persistent>
    <template #title> 视频 </template>
    <div class="video-content flex justify-center">
      <video :src="dialogVideoData.url" controls class="rounded-md max-w-full max-h-full"></video>
    </div>
  </Dialog>
  <CreateLivePane ref="createLivePaneRef" @submit="doQueryFirstPage" />
</template>
<script setup lang="ts">
import { reactive, ref, onMounted } from "vue";
import BaseBreadcrumb from "@/components/shared/BaseBreadcrumb.vue";
import UiParentCard from "@/components/shared/UiParentCard.vue";
import ChipBoolean from "@/components/ui/ChipBoolean.vue";
import CreateLivePane from "./components/CreateLivePane.vue";
import { IconVideo } from "@tabler/icons-vue";

import { http, format } from "@/utils";

const page = ref({ title: "活体检测" });
const breadcrumbs = ref([
  {
    text: "人脸服务",
    disabled: false,
    href: "#"
  },
  {
    text: "活体检测",
    disabled: true,
    href: "#"
  }
]);
const searchData = reactive({
  uuid: "",
  source: ""
});
const tableInfos = reactive({
  list: [],
  total: 0
});
const createLivePaneRef = ref();
const tableWithPagerRef = ref();
const refDialogVideo = ref();
const dialogVideoData = reactive({
  url: ""
});

const viewVideo = row => {
  dialogVideoData.url = row.inputS3Url;
  refDialogVideo.value.show({
    showActions: false
  });
};

const doTableQuery = async (options = {}) => {
  const [err, res] = await http.get({
    url: "/esrgan/list",
    showLoading: tableWithPagerRef.value.el,
    data: {
      modelType: "faceLive",
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

const onAdd = () => {
  createLivePaneRef.value.show({
    title: "创建检测",
    operateType: "add"
  });
};

onMounted(() => {
  doTableQuery();
});
</script>
