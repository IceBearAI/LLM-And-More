<template>
  <BaseBreadcrumb :title="$t(page.title)" :breadcrumbs="breadcrumbs"></BaseBreadcrumb>
  <v-row>
    <v-col>
      <UiParentCard>
        <v-row>
          <v-col cols="12" lg="3" md="4" sm="6">
            <Select
              v-model="formData.taggedStatus"
              :label="$t('audioMark.taggedStatus')"
              :mapDictionary="{ code: 'boolean' }"
              @change="doQueryFirstPage"
            >
            </Select>
          </v-col>

          <v-col cols="12" lg="3" md="4" sm="6">
            <Select
              v-model="formData.originalLanguage"
              :label="$t('language')"
              :mapDictionary="{ code: 'audio_tagged_lang', i18nKey: 'markLanguage' }"
              @change="doQueryFirstPage"
            >
            </Select>
          </v-col>
          <v-col cols="12" lg="3" md="4" sm="6">
            <v-menu>
              <template v-slot:activator="{ props }">
                <v-btn color="primary" v-bind="props">{{ $t("audioMark.tableHeader.exportBtn") }}</v-btn>
              </template>
              <v-list>
                <v-list-item v-for="(item, index) in onButtons()" :key="index" :value="index">
                  <v-list-item-title @click="onExport(item)">{{ item.text }}</v-list-item-title>
                </v-list-item>
              </v-list>
            </v-menu>
          </v-col>
          <v-col cols="12">
            <AlertBlock> {{ $t("audioMark.alertBlock") }} </AlertBlock>
          </v-col>
          <v-col cols="12">
            <TableWithPager @query="doTableQuery" ref="refTableWithPager" :infos="state.tableInfos">
              <el-table-column :label="$t('fileName')" width="200px" fixed="left">
                <template #default="{ row }">
                  <el-tooltip effect="dark" :content="$t('audioMark.tableHeader.audioNameTip')" placement="top">
                    <span class="link" @click="onEdit(row)">{{ row.audioName }}</span>
                  </el-tooltip>
                </template>
              </el-table-column>
              <el-table-column :label="$t('language')" width="110px" show-overflow-tooltip>
                <template #default="{ row }">
                  <span>
                    {{ $t("options.markLanguage." + row.originalLanguage) }}
                  </span>
                </template>
              </el-table-column>
              <!-- <el-table-column :label="$t('audioMark.tableHeader.taggedStatus')" min-width="100px">
                <template #default="{ row }">
                  <ChipBoolean v-model="row.taggedStatus"></ChipBoolean>
                </template>
              </el-table-column>
              <el-table-column :label="$t('duration')" min-width="100px">
                <template #default="{ row }">
                  {{ row.duration }}
                </template>
              </el-table-column>
              <el-table-column :label="$t('size')" min-width="100px">
                <template #default="{ row }">
                  {{ format.getFileSize(row.audioSize) }}
                </template>
              </el-table-column> -->
              <el-table-column :label="$t('tryListening')" min-width="330px">
                <template #default="{ row }">
                  <AiAudio type="simple" :src="row.audioUrl" />
                </template>
              </el-table-column>
              <el-table-column :label="$t('audioMark.tableHeader.originalContent')" min-width="200px" show-overflow-tooltip>
                <template #default="{ row }">{{ row.originalContent }} </template>
              </el-table-column>
              <el-table-column :label="$t('audioMark.tableHeader.taggedContent')" min-width="200px" show-overflow-tooltip>
                <template #default="{ row }">{{ row.taggedContent }} </template>
              </el-table-column>
              <el-table-column :label="$t('audioMark.tableHeader.taggedLanguage')" min-width="150px">
                <template #default="{ row }">
                  {{ row.taggedLanguage ? $t("options.markLanguage." + row.taggedLanguage, "") : "" }}
                </template>
              </el-table-column>
              <el-table-column :label="$t('audioMark.tableHeader.updatedAt')" min-width="200px">
                <template #default="{ row }">
                  {{ row.taggedStatus ? format.dateFormat(row.updatedAt, "YYYY-MM-DD HH:mm:ss") : "" }}
                </template>
              </el-table-column>
              <el-table-column :label="$t('audioMark.tableHeader.lastOperator')" min-width="150px">
                <template #default="{ row }">{{ row.lastOperator }} </template>
              </el-table-column>
              <el-table-column :label="$t('operation')" min-width="140px" fixed="right">
                <template #default="{ row }">
                  <ButtonsInTable :onlyOne="true" :buttons="getButtons(row)" />
                </template>
              </el-table-column>
            </TableWithPager>
          </v-col>
        </v-row>
      </UiParentCard>
    </v-col>
  </v-row>
  <PaneAudioMark ref="refPaneAudioMark" @submit="doQueryFirstPage" />
  <ConfirmByClick ref="refConfirmDelete" @submit="doDelete">
    <template #text> {{ $t("audioMark['confirm delete']") }}？</template>
  </ConfirmByClick>
</template>
<script setup>
import { reactive, toRefs, ref, onMounted, computed } from "vue";
import BaseBreadcrumb from "@/components/shared/BaseBreadcrumb.vue";
import UiParentCard from "@/components/shared/UiParentCard.vue";
import AlertBlock from "@/components/ui/AlertBlock.vue";
import { useMapRemoteStore } from "@/stores";
import ChipBoolean from "@/components/ui/ChipBoolean.vue";
import { http, format } from "@/utils";
import AiAudio from "@/components/business/AiAudio.vue";
import PaneAudioMark from "./components/PaneAudioMark.vue";
import { useI18n } from "vue-i18n";
import ConfirmByClick from "@/components/business/ConfirmByClick.vue";

const mapRemoteStore = useMapRemoteStore();
const { t } = useI18n(); // 解构出t方法
const refTableWithPager = ref();
const refPaneAudioMark = ref();
const refConfirmDelete = ref();
const state = reactive({
  style: {},
  formData: {
    taggedStatus: false,
    originalLanguage: null
  },
  selectedInfo: {},
  tableInfos: {
    list: [],
    total: ""
  }
});

const { formData, tableInfos } = toRefs(state);

const doQueryFirstPage = () => {
  refTableWithPager.value.query({ page: 1 });
};

const doTableQuery = async (options = {}) => {
  const [err, res] = await http.get({
    url: "/api/voice/audio",
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

const onEdit = row => {
  refPaneAudioMark.value.show({
    infos: row
  });
};

const onDelete = info => {
  state.selectedInfo = info;
  refConfirmDelete.value.show({
    width: "400px"
  });
};

const doDelete = async () => {
  const [err, res] = await http.delete({
    showSuccess: true,
    url: `/api/voice/audio/${state.selectedInfo.id}`
  });
  if (res) {
    refConfirmDelete.value.hide();
    doQueryFirstPage();
  }
};

const onConfirm = async row => {
  const [err, res] = await http.put({
    url: `/api/voice/audio/${row.id}`,
    showLoading: true,
    showSuccess: true,
    data: {
      taggedContent: row.taggedContent || row.originalContent,
      taggedLanguage: row.taggedLanguage || row.originalLanguage
    }
  });
};

const getButtons = row => {
  let ret = [];
  //标注
  ret.push({
    text: t("audioMark.tableHeader.mark"),
    color: "info",
    click() {
      onEdit(row);
    }
  });
  //确认
  if (row.taggedStatus == false) {
    ret.push({
      text: t("audioMark.tableHeader.confirm"),
      color: "info",
      click() {
        onConfirm(row);
      }
    });
  }
  //废弃
  ret.push({
    text: t("audioMark.tableHeader.discard"),
    color: "info",
    click() {
      onDelete(row);
    }
  });

  return ret;
};
const onButtons = () => {
  let ret = [];
  //标注
  ret.push({
    text: "CSV",
    // t("audioMark.tableHeader.mark"),
    minTitle: "csv"
  });
  ret.push({
    text: "JSONL",
    // t("audioMark.tableHeader.mark"),
    minTitle: "jsonl"
  });
  return ret;
};

const page = ref({ title: "audioAnnotation" });
const breadcrumbs = ref([
  {
    text: "audioManagement",
    disabled: false,
    href: "#"
  },
  {
    text: "audioAnnotation",
    disabled: true,
    href: "#"
  }
]);
const onExport = async row => {
  const [err, res] = await http.post({
    showSuccess: true,
    url: `/api/voice/audio/export`,
    data: {
      format: row.minTitle
    }
  });
  if (res) {
    http.downloadByUrl({
      fileUrl: res.url,
      suffixName: res.format
    });
  }
};
onMounted(() => {
  doQueryFirstPage();
});
</script>
<style lang="scss"></style>
