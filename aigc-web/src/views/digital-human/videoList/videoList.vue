<template>
  <BaseBreadcrumb :title="page.title" :breadcrumbs="breadcrumbs"></BaseBreadcrumb>
  <UiParentCard class="mb-3">
    <TaskOverview :config="taskDetailConfig" request-url="/digitalhuman/synthesis/count" />
  </UiParentCard>
  <UiParentCard>
    <v-row>
      <v-col cols="12" lg="3" md="4" sm="6">
        <v-text-field
          density="compact"
          v-model="searchData.title"
          label="请输入标题"
          hide-details
          clearable
          variant="outlined"
          @keyup.enter="doQueryFirstPage"
          @click:clear="doQueryFirstPage"
        ></v-text-field>
      </v-col>
      <v-col cols="12" lg="3" md="4" sm="6">
        <Select
          @change="doQueryFirstPage"
          label="请选择状态"
          :mapDictionary="{ code: 'digitalhuman_synthesis_status' }"
          v-model="searchData.status"
        ></Select>
      </v-col>
      <v-col cols="12" lg="3" md="4" sm="6">
        <ButtonsInForm>
          <v-btn color="primary" @click="onAdd">创建视频</v-btn>
          <refresh-button ref="refreshButtonRef" @refresh="doQueryCurrentPage" :disabled="loading" />
        </ButtonsInForm>
      </v-col>
      <v-col cols="12">
        <AlertBlock> 生成视频的时间可能会比较长，请耐心等待！ </AlertBlock>
      </v-col>

      <v-col cols="12">
        <v-row ref="listContentRef" :class="{ 'justify-center': list.length === 0 }">
          <v-col v-if="list.length > 0" cols="12" lg="3" md="4" sm="6" v-for="item in list">
            <v-card elevation="10" rounded="md">
              <a
                class="card-list-item text-textPrimary text-decoration-none"
                href="javascript: void(0)"
                @click="openDetail(item)"
              >
                <TagCorner :class="`bg-${digitalhumanStatusMap[item.status].color}`">
                  {{ digitalhumanStatusMap[item.status].text }}
                </TagCorner>
                <v-img :src="item.digitalHumanPerson.cover" height="180px" cover class="rounded-t-md align-end text-right">
                  <div class="pa-3">
                    <v-chip
                      @click.stop
                      class="bg-surface text-body-2 font-weight-medium"
                      variant="flat"
                      size="small"
                      :text="`${item.videoDuration || '0s'}/${item.videoSize}`"
                    ></v-chip>
                  </div>
                </v-img>
                <v-card-item class="pa-5">
                  <h5 class="text-h5">{{ item.title }}</h5>
                  <p class="text-subtitle-1 mt-1 text-medium-emphasis text-truncate">{{ item.ttsText }}</p>
                  <p class="text-subtitle-1 mt-1 text-medium-emphasis text-truncate">{{ item.synthesisModel }}</p>
                  <p v-copy="item.uuid" class="text-subtitle-2 mt-1 text-medium-emphasis text-truncate">{{ item.uuid }}</p>
                  <div class="d-flex align-center justify-space-between mt-2" style="height: 32px">
                    <div class="flex-1-1 d-flex justify-space-between text-medium-emphasis">
                      <span>{{ format.dateFromNow(item.createdAt) }}</span>
                      <span>
                        {{ item.digitalHumanPerson.cname }} ({{ getLabels([["speak_gender", item.digitalHumanPerson.gender]]) }})
                      </span>
                    </div>
                    <v-btn class="ml-6" size="x-small" color="inherit" icon variant="text">
                      <IconDotsVertical width="14" stroke-width="1.5" />
                      <v-menu activator="parent">
                        <v-list density="compact" @click:select="operateMenuClick($event, item)">
                          <v-list-item
                            v-for="operate in getOperateConfig(item.status)"
                            :key="operate.key"
                            :value="operate.key"
                            hide-details
                            min-height="38"
                          >
                            <template v-slot:prepend>
                              <component
                                :is="operate.icon"
                                :size="16"
                                class="mr-2"
                                :class="[`text-${operate.color}`]"
                              ></component>
                            </template>
                            <v-list-item-title :class="[`text-${operate.color}`]">{{ operate.text }}</v-list-item-title>
                          </v-list-item>
                        </v-list>
                      </v-menu>
                    </v-btn>
                  </div>
                </v-card-item>
              </a>
            </v-card>
          </v-col>
          <NoData v-else></NoData>
        </v-row>
        <Pager
          class="mt-5"
          ref="refPager"
          :total="total"
          :page-sizes="[12, 20, 40, 60, 120]"
          @query="doQuery"
          v-show="list.length > 0"
        />
      </v-col>
    </v-row>
  </UiParentCard>
  <DialogLog ref="refDialogLog" />

  <ConfirmByClick ref="refConfirmDelete" @submit="doDelete">
    <template #text>
      <span class="font-weight-bold">这是进行一项操作时必须了解的重要信息</span><br />
      此操作将会<span class="text-primary">删除</span>该视频<br />
      标题：<span class="text-primary">{{ currentOperateData.title }}</span
      ><br />
      你还要继续吗？
    </template>
  </ConfirmByClick>
  <ConfirmByClick ref="refConfirmCancel" @submit="doCancel">
    <template #text>
      <span class="font-weight-bold">这是进行一项操作时必须了解的重要信息</span><br />
      此操作将会<span class="text-primary">取消</span>该视频<br />
      标题：<span class="text-primary">{{ currentOperateData.title }}</span
      ><br />
      你还要继续吗？
    </template>
  </ConfirmByClick>
  <ConfirmByClick ref="refConfirmFirst" @submit="doFirst">
    <template #text>
      <span class="font-weight-bold">这是进行一项操作时必须了解的重要信息</span><br />
      此操作将会<span class="text-primary">优先合成</span>该视频<br />
      标题：<span class="text-primary">{{ currentOperateData.title }}</span
      ><br />
      你还要继续吗？
    </template>
  </ConfirmByClick>
</template>
<script setup lang="ts">
import { reactive, ref, onMounted } from "vue";
import BaseBreadcrumb from "@/components/shared/BaseBreadcrumb.vue";
import UiParentCard from "@/components/shared/UiParentCard.vue";
import AlertBlock from "@/components/ui/AlertBlock.vue";
import ConfirmByClick from "@/components/business/ConfirmByClick.vue";
import TaskOverview from "@/components/business/TaskOverview.vue";
import Pager from "@/components/global/Pager.vue";
import DialogLog from "@/components/ui/log/DialogLog.vue";
import { useRouter } from "vue-router";
import { http, format } from "@/utils";
import { IconDotsVertical, IconCircleMinus, IconEye, IconBolt } from "@tabler/icons-vue";
import { useMapRemoteStore } from "@/stores";
import { taskDetailConfig, digitalhumanStatusMap } from "./map";
import TagCorner from "@/components/business/TagCorner.vue";

defineOptions({
  name: "DigitalVideoList"
});

const router = useRouter();

const { getLabels, loadDictTree } = useMapRemoteStore();

const page = ref({ title: "视频合成列表" });
const breadcrumbs = ref([]);
const searchData = reactive({
  title: "",
  status: null
});
const list = ref([]);
const total = ref(0);
const refPager = ref();

const refDialogLog = ref();
const refConfirmCancel = ref();
const refConfirmDelete = ref();
const refConfirmFirst = ref();
const currentOperateData = reactive({
  uuid: "",
  title: ""
});
const listContentRef = ref();
const refreshButtonRef = ref();
const loading = ref(false);

loadDictTree(["speak_gender"]);
const doQuery = async (options = {}) => {
  loading.value = true;
  const [err, res] = await http.get({
    url: "/digitalhuman/synthesis/list",
    showLoading: listContentRef.value.$el,
    data: {
      ...searchData,
      ...options
    }
  });
  if (res) {
    list.value = res.list || [];
    total.value = res.total;
  } else {
    list.value = [];
    total.value = 0;
  }
  loading.value = false;
  refreshButtonRef.value.start();
};

const doQueryFirstPage = () => {
  refPager.value.query({ page: 1 });
};

const doQueryCurrentPage = () => {
  refPager.value.query();
};

const getOperateConfig = status => {
  let ret = [];
  if (status !== "waiting") {
    ret.push({
      text: "日志",
      color: "info",
      icon: IconEye,
      key: "log"
    });
  }
  if (status === "waiting" || status === "running") {
    ret.push({
      text: "取消",
      color: "error",
      icon: IconCircleMinus,
      key: "cancel"
    });
  }
  if (status === "failed" || status === "cancel" || status === "success") {
    ret.push({
      text: "删除",
      color: "error",
      icon: IconCircleMinus,
      key: "delete"
    });
  }
  if (status === "waiting") {
    ret.push({
      text: "优速通",
      color: "info",
      icon: IconBolt,
      key: "first"
    });
  }
  return ret;
};

const operateMenuClick = ({ id }, item) => {
  currentOperateData.uuid = item.uuid;
  currentOperateData.title = item.title;
  if (id === "log") {
    onLog();
  } else if (id === "cancel") {
    refConfirmCancel.value.show({
      width: "400px"
    });
  } else if (id === "delete") {
    refConfirmDelete.value.show({
      width: "400px"
    });
  } else if (id === "first") {
    refConfirmFirst.value.show({
      width: "400px"
    });
  }
};

const onLog = async () => {
  refDialogLog.value.show();
  let [err, res] = await http.get({
    url: `/api/digitalhuman/synthesis/${currentOperateData.uuid}/view`
  });
  if (res) {
    refDialogLog.value.setContent(res.synthesisLog);
  }
};

const doCancel = async (options = {}) => {
  const [err, res] = await http.put({
    ...options,
    showSuccess: true,
    url: `/api/digitalhuman/synthesis/${currentOperateData.uuid}/cancel`
  });
  if (res) {
    refConfirmCancel.value.hide();
    doQueryCurrentPage();
  }
};

const doDelete = async (options = {}) => {
  const [err, res] = await http.delete({
    ...options,
    showSuccess: true,
    url: `/api/digitalhuman/synthesis/${currentOperateData.uuid}/delete`
  });
  if (res) {
    refConfirmDelete.value.hide();
    doQueryCurrentPage();
  }
};

const doFirst = async (options = {}) => {
  const [err, res] = await http.put({
    ...options,
    showSuccess: true,
    url: `/digitalhuman/synthesis/${currentOperateData.uuid}/first`
  });
  if (res) {
    refConfirmFirst.value.hide();
    doQueryCurrentPage();
  }
};

const onAdd = () => {
  router.push("/digital-human/video-list/edit");
};

const openDetail = ({ status, uuid }) => {
  // if (status !== "success") return;
  router.push(`/digital-human/video-list/detail?uuid=${uuid}`);
};

onMounted(() => {
  doQueryCurrentPage();
});
</script>
<style lang="scss" scoped></style>
