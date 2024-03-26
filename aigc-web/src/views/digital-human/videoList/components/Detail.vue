<template>
  <NavBack backUrl="/digital-human/video-list/list">视频详情</NavBack>
  <v-row class="mt-4">
    <v-col cols="12" md="7">
      <UiParentCard>
        <div style="aspect-ratio: 1">
          <video
            class="w-100 h-100"
            :src="showData.isGfpgan ? showData.gfpganS3Url : showData.wav2LipS3Url"
            controls
            :poster="showData.digitalHumanPerson.cover"
          ></video>
        </div>
      </UiParentCard>
    </v-col>
    <v-col cols="12" md="5">
      <UiParentCard>
        <div class="my-form">
          <h3 class="text-h3 my-2">{{ showData.title }}</h3>
          <p>{{ showData.ttsText }}</p>
          <v-input hide-details>
            <template #prepend> <label>ID</label></template>
            <span v-copy="showData.uuid">{{ showData.uuid }}</span>
          </v-input>
          <v-input hide-details>
            <template #prepend> <label>状态</label></template>
            <v-chip :color="digitalhumanStatusMap[showData.status]?.color" variant="outlined" size="small">{{
              digitalhumanStatusMap[showData.status]?.text
            }}</v-chip>
          </v-input>
          <v-input hide-details>
            <template #prepend> <label>形象</label></template>
            <div class="d-flex align-center">
              <v-avatar size="40" class="mr-2">
                <img :src="showData.digitalHumanPerson.cover" class="w-100" />
              </v-avatar>
              {{ showData.digitalHumanPerson.cname }}
              ({{ getLabels([["speak_gender", showData.digitalHumanPerson.gender]]) }})
            </div>
          </v-input>
          <v-input hide-details>
            <template #prepend> <label>声音</label></template>
            <div class="d-flex align-center">
              <v-avatar size="40" class="mr-2">
                <img :src="showData.voiceSpeak.headImg" class="w-100" />
              </v-avatar>
              {{ showData.voiceSpeak.speakCname }}({{
                getLabels(
                  [
                    ["speak_age_group", showData.voiceSpeak.ageGroup],
                    ["speak_gender", showData.voiceSpeak.gender]
                  ],
                  ret => {
                    if (ret.length) {
                      return ret.join("") + "声";
                    } else {
                      return "未知";
                    }
                  }
                )
              }})
            </div>
          </v-input>
          <v-input hide-details>
            <template #prepend> <label>是否超分</label></template>
            {{ showData.isGfpgan ? "是" : "否" }}
          </v-input>
          <v-input hide-details>
            <template #prepend> <label>语速</label></template>
            1.0
          </v-input>
          <v-input hide-details>
            <template #prepend> <label>时长</label></template>
            {{ showData.videoDuration || "0s" }}
          </v-input>
          <v-input hide-details>
            <template #prepend> <label>大小</label></template>
            {{ showData.videoSize }}
          </v-input>
          <v-input hide-details>
            <template #prepend> <label>创建时间</label></template>
            {{ format.dateFormat(showData.createdAt, "YYYY-MM-DD HH:mm:ss") }}
          </v-input>
          <v-input hide-details>
            <template #prepend> <label>开始时间</label></template>
            {{ format.dateFormat(showData.startTime, "YYYY-MM-DD HH:mm:ss") }}
          </v-input>
          <v-input hide-details>
            <template #prepend> <label>结束时间</label></template>
            {{ format.dateFormat(showData.endTime, "YYYY-MM-DD HH:mm:ss") }}
          </v-input>
          <v-input hide-details>
            <template #prepend> <label>更新时间</label></template>
            {{ format.dateFormat(showData.updatedAt, "YYYY-MM-DD HH:mm:ss") }}
          </v-input>
          <v-input hide-details :center-affix="false">
            <template #prepend> <label>备注</label></template>
            {{ showData.remark }}
          </v-input>
        </div>
      </UiParentCard>
    </v-col>
  </v-row>
</template>
<script setup lang="ts">
import { onMounted, ref } from "vue";
import { http, format } from "@/utils";
import NavBack from "@/components/business/NavBack.vue";
import UiParentCard from "@/components/shared/UiParentCard.vue";
import { useMapRemoteStore } from "@/stores";
import { useRoute } from "vue-router";
import { digitalhumanStatusMap } from "../map";

const route = useRoute();

const { loadDictTree, getLabels } = useMapRemoteStore();

loadDictTree(["speak_gender", "speak_age_group"]);

const { uuid } = route.query;
const showData = ref<Record<string, any>>({
  digitalHumanPerson: {},
  voiceSpeak: {}
});

const getData = async () => {
  let [err, res] = await http.get({
    url: `/api/digitalhuman/synthesis/${uuid}/view`
  });
  if (res) {
    showData.value = res || {};
  }
};

onMounted(() => {
  getData();
});
</script>
