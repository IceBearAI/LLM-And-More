<template>
  <Pane ref="refPane">
    <v-row>
      <v-col cols="12" md="6">
        <UiParentCard>
          <div style="aspect-ratio: 1">
            <video
              class="w-100 h-100"
              :src="showData.video ? showData.video : showData.video"
              controls
              :poster="showData.digitalHumanPerson.cover"
            ></video>
          </div>
        </UiParentCard>
      </v-col>
      <v-col cols="12" md="6">
        <UiParentCard>
          <div class="my-form">
            <h3 class="text-h3 my-2">{{ showData.cname }}</h3>

            <v-input hide-details>
              <template #prepend> <label>头像：</label></template>
              <div class="d-flex align-center">
                <v-avatar size="80" class="mr-2" rounded="0">
                  <img :src="showData.cover" class="w-100" />
                </v-avatar>
              </div>
            </v-input>

            <v-input hide-details>
              <template #prepend> <label>时长：</label></template>
              15s
            </v-input>
            <v-input hide-details>
              <template #prepend> <label>大小：</label></template>
              1.0
            </v-input>

            <v-input hide-details>
              <template #prepend> <label>创建时间：</label></template>
              {{ format.dateFormat(showData.createdAt, "YYYY-MM-DD HH:mm:ss") }}
            </v-input>
            <v-input hide-details>
              <template #prepend> <label>更新时间：</label></template>
              {{ format.dateFormat(showData.updatedAt, "YYYY-MM-DD HH:mm:ss") }}
            </v-input>
            <v-input hide-details>
              <template #prepend> <label>备注：</label></template>
              {{ showData.remark }}
            </v-input>
            <v-btn flat color="secondary" @click="onSynthesis">去合成视频&nbsp;></v-btn>
          </div>
        </UiParentCard>
      </v-col>
    </v-row>
  </Pane>
</template>
<script setup lang="ts">
import { ref } from "vue";
import { http, format } from "@/utils";
import UiParentCard from "@/components/shared/UiParentCard.vue";
import { useRouter } from "vue-router";

const router = useRouter();
const refPane = ref();
const currentUuid = ref("");
const showData = ref<Record<string, any>>({
  digitalHumanPerson: {},
  voiceSpeak: {}
});

const getData = async () => {
  let [err, res] = await http.get({
    url: `/api/digitalhuman/synthesis/${currentUuid.value}/view`
  });
  if (res) {
    showData.value = res || {};
  }
};
const onSynthesis = () => {
  router.push("/digital-human/edit");
};
defineExpose({
  show({ title, uuid }) {
    currentUuid.value = uuid;
    refPane.value.show({
      width: 1000,
      showActions: false,
      title
    });
    getData();
  }
});
</script>
