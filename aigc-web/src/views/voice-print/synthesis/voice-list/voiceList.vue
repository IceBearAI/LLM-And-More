<template>
  <BaseBreadcrumb :title="page.title" :breadcrumbs="breadcrumbs"></BaseBreadcrumb>
  <UiParentCard>
    <v-row>
      <v-col cols="12" lg="3" md="4" sm="6">
        <v-text-field
          v-model="searchData.speakName"
          label="请输入标识"
          hide-details
          clearable
          @keyup.enter="doQueryFirstPage"
          @click:clear="doQueryFirstPage"
        ></v-text-field>
      </v-col>
      <v-col cols="12" lg="3" md="4" sm="6">
        <Select
          v-model="searchData.provider"
          :mapDictionary="{ code: 'speak_provider' }"
          label="请选择供应商"
          hide-details
          @change="doQueryFirstPage"
        ></Select>
      </v-col>
      <v-col cols="12" lg="3" md="4" sm="6">
        <Select
          v-model="searchData.lang"
          :mapDictionary="{ code: 'speak_lang' }"
          label="请选择语言"
          hide-details
          @change="doQueryFirstPage"
        ></Select>
      </v-col>
      <v-col cols="12" lg="3" md="4" sm="6">
        <ButtonsInForm>
          <v-btn color="primary" @click="onSynthesisVoice">创建TTS</v-btn>
        </ButtonsInForm>
      </v-col>
      <v-col cols="12">
        <AlertBlock>修改之后将实时生效，请谨慎操作！</AlertBlock>
      </v-col>

      <v-col cols="12">
        <TableWithPager @query="doTableQuery" ref="tableWithPagerRef" :infos="tableInfos">
          <el-table-column label="试听" min-width="330px">
            <template #default="{ row }">
              <AiAudio :src="row?.s3Url" />
            </template>
          </el-table-column>
          <el-table-column label="内容" prop="text" min-width="300px" show-overflow-tooltip></el-table-column>
          <el-table-column label="标识" width="150px" show-overflow-tooltip>
            <template #default="{ row }">
              <span v-copy="row.speakName">{{ row.speakName }}</span>
            </template>
          </el-table-column>
          <el-table-column label="姓名" prop="speakCname" width="100px"></el-table-column>
          <el-table-column label="供应" width="100px">
            <template #default="{ row }">
              <span>{{ getLabels([["speak_provider", row.provider]]) }}</span>
            </template>
          </el-table-column>
          <el-table-column label="语言" width="160px">
            <template #default="{ row }">
              {{ row.lang }}
              <span>{{ getLabels([["speak_lang", row.lang]]) }}</span>
            </template>
          </el-table-column>
          <el-table-column label="音色" width="120px">
            <template #default="{ row }">
              <span>
                {{
                  getLabels(
                    [
                      ["speak_age_group", row.ageGroup],
                      ["speak_gender", row.gender]
                    ],
                    ret => {
                      if (ret.length) {
                        return ret.join("") + "声";
                      } else {
                        return "未知";
                      }
                    }
                  )
                }}
              </span>
            </template>
          </el-table-column>
          <el-table-column label="用时" prop="duration" width="120px"></el-table-column>
          <el-table-column label="语速" prop="speed" width="120px"></el-table-column>
          <el-table-column label="大小" width="120px">
            <template #default="{ row }">
              {{ format.getFileSize(row.size) }}
            </template>
          </el-table-column>
          <el-table-column label="更新时间" min-width="165px">
            <template #default="{ row }">
              {{ format.dateFormat(row.updatedAt, "YYYY-MM-DD HH:mm:ss") }}
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
  <ConfirmByClick ref="refConfirmDelete" @submit="doDelete">
    <template #text>
      这是进行一项操作时必须了解的重要信息<br />
      您将要删除 <span class="text-primary font-weight-black">{{ confirmDelete.name }}</span> 人音频，确定要继续吗？
    </template>
  </ConfirmByClick>
</template>
<script setup lang="ts">
import { reactive, ref, onMounted } from "vue";
import BaseBreadcrumb from "@/components/shared/BaseBreadcrumb.vue";
import UiParentCard from "@/components/shared/UiParentCard.vue";
import AlertBlock from "@/components/ui/AlertBlock.vue";
import ConfirmByClick from "@/components/business/ConfirmByClick.vue";
import AiAudio from "@/components/business/AiAudio.vue";

import { http, format } from "@/utils";
import { useMapRemoteStore } from "@/stores";
import { useRouter } from "vue-router";

const { loadDictTree, getLabels } = useMapRemoteStore();
const router = useRouter();
loadDictTree(["speak_age_group", "speak_gender", "speak_provider", "speak_lang"]);

const page = ref({ title: "TTS合成列表" });
const breadcrumbs = ref([
  {
    text: "声音合成",
    disabled: false,
    href: "#"
  },
  {
    text: "TTS合成列表",
    disabled: true,
    href: "#"
  }
]);
const searchData = reactive({
  speakName: "",
  provider: null,
  lang: null
});
const tableInfos = reactive({
  list: [],
  total: 0
});
const tableWithPagerRef = ref();
const refConfirmDelete = ref();
const confirmDelete = reactive({
  id: "",
  name: ""
});

const onSynthesisVoice = () => {
  router.push("/voice-print/synthesis/synthesis-voice");
};

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
    text: "下载",
    color: "info",
    click() {
      onDownload(row);
    }
  });
  return ret;
};

const doTableQuery = async (options = {}) => {
  const [err, res] = await http.get({
    url: "/voice/tts",
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

const onDelete = row => {
  confirmDelete.name = row.speakCname;
  confirmDelete.id = row.id;
  refConfirmDelete.value.show({
    width: "400px",
    confirmText: confirmDelete.name
  });
};

const doDelete = async (options = {}) => {
  const [err, res] = await http.delete({
    ...options,
    showSuccess: true,
    url: `/voice/tts/${confirmDelete.id}`
  });
  if (res) {
    refConfirmDelete.value.hide();
    doTableQuery();
  }
};

const onDownload = row => {
  window.open(row.s3Url);
};

onMounted(() => {
  doTableQuery();
});
</script>
