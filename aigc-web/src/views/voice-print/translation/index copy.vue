<template>
  <BaseBreadcrumb :title="page.title" :breadcrumbs="breadcrumbs"></BaseBreadcrumb>
  <v-row>
    <v-col cols="12">
      <UiParentCard title="上传音频">
        <v-row class="justify-center">
          <v-col cols="12" class="text-center">
            <AiAudio :src="audioUrl" />
          </v-col>
          <v-col cols="12" md="6">
            <v-file-input
              v-model="audioFile"
              accept="audio/*, .mp3, .wma, .amr, .wav, .m4a"
              label="请选择音频"
              prepend-icon="mdi-volume-high"
              hide-details
              variant="outlined"
            ></v-file-input>
          </v-col>
        </v-row>
      </UiParentCard>
    </v-col>
    <v-col cols="12">
      <v-card elevation="10">
        <v-card-text class="text-center">
          <v-btn :loading="searchBtnLoading" :disabled="!audioUrl" color="primary" width="200" @click="handleSearch"
            >开始转换</v-btn
          >
        </v-card-text>
      </v-card>
    </v-col>
    <v-col cols="12">
      <UiParentCard title="转换结果">
        <v-row class="text-body-1 pa-5">
          <p v-if="result" v-text="result" v-copy="result" style="white-space: pre-wrap"></p>
          <p v-else class="flex-fill text-center">暂无结果</p>
        </v-row>
      </UiParentCard>
    </v-col>
  </v-row>
</template>
<script setup lang="ts">
import { ref, computed } from "vue";
import BaseBreadcrumb from "@/components/shared/BaseBreadcrumb.vue";
import UiParentCard from "@/components/shared/UiParentCard.vue";
import { http } from "@/utils";
import AiAudio from "@/components/business/AiAudio.vue";

const page = ref({ title: "声音转文本" });
const breadcrumbs = ref([
  {
    text: "语音服务",
    disabled: false,
    href: "#"
  },
  {
    text: "声音转文本",
    disabled: true,
    href: "#"
  }
]);

const searchBtnLoading = ref(false);
const audioFile = ref([]);
const result = ref("");

const audioUrl = computed(() => {
  if (audioFile.value && audioFile.value.length > 0) {
    return URL.createObjectURL(audioFile.value[0]);
  }
  return "";
});

const handleSearch = async () => {
  if (!audioUrl) return;
  const formData = new FormData();
  formData.append("file", audioFile.value[0]);

  searchBtnLoading.value = true;

  const [err, res] = await http.upload({
    url: "/voice/translation",
    data: {
      file: audioFile.value
    }
  });
  if (res) {
    result.value = res.data.text;
  }
  searchBtnLoading.value = false;
};
</script>
