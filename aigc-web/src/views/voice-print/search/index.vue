<template>
  <NavBack v-if="route.path === '/voice-print/library-list/search'" backUrl="/voice-print/library-list/list">
    {{ page.title }}
  </NavBack>
  <BaseBreadcrumb v-else :title="page.title" :breadcrumbs="breadcrumbs"></BaseBreadcrumb>
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
            >开始查询</v-btn
          >
        </v-card-text>
      </v-card>
    </v-col>
    <v-col cols="12">
      <UiParentCard title="查询结果">
        <v-row v-if="result" class="text-body-1 justify-center pa-5">
          <template v-if="result.isSame">
            <v-col cols="12" class="d-flex align-center">
              <v-label class="font-weight-medium">用户姓名：</v-label>
              <span>{{ result.userName }}</span>
            </v-col>
            <v-col cols="12" class="d-flex align-center">
              <v-label class="font-weight-medium">用户标识：</v-label>
              <span>{{ result.userKey }}</span>
            </v-col>
            <v-col cols="12" class="d-flex align-center">
              <v-label class="font-weight-medium">相似度：</v-label>
              <v-rating
                :model-value="result.dist * 10"
                length="10"
                readonly
                density="compact"
                size="large"
                color="warning"
              ></v-rating>
              <span>({{ result.dist }})</span>
            </v-col>
          </template>
          <template v-else>未匹配到相似的声纹</template>
        </v-row>
        <v-row class="text-body-1 justify-center pa-5" v-else> 暂无结果 </v-row>
      </UiParentCard>
    </v-col>
  </v-row>
</template>
<script setup lang="ts">
import { ref, computed } from "vue";
import { useRoute } from "vue-router";
import NavBack from "@/components/business/NavBack.vue";
import BaseBreadcrumb from "@/components/shared/BaseBreadcrumb.vue";
import UiParentCard from "@/components/shared/UiParentCard.vue";
import AiAudio from "@/components/business/AiAudio.vue";
import { http } from "@/utils";

const route = useRoute();
const page = ref({ title: "声纹查询" });
const breadcrumbs = ref([
  {
    text: "智能声纹",
    disabled: false,
    href: "#"
  },
  {
    text: "声纹查询",
    disabled: true,
    href: "#"
  }
]);

const searchBtnLoading = ref(false);
const audioFile = ref([]);
const result = ref(null);

const audioUrl = computed(() => {
  if (audioFile.value && audioFile.value.length > 0) {
    return URL.createObjectURL(audioFile.value[0]);
  }
  return "";
});

const handleSearch = async () => {
  if (!audioUrl) return;

  searchBtnLoading.value = true;
  const [err, res] = await http.upload({
    url: "/voice/search",
    data: {
      file: audioFile.value
    }
  });
  if (res) {
    result.value = res.data;
  }
  searchBtnLoading.value = false;
};
</script>
