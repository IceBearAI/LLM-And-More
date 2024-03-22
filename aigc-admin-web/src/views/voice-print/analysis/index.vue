<script setup>
import { ref, reactive, inject } from "vue";
import { useRoute } from "vue-router";
import $ from "jquery";
import NavBack from "@/components/business/NavBack.vue";
import BaseBreadcrumb from "@/components/shared/BaseBreadcrumb.vue";
import UiParentCard from "@/components/shared/UiParentCard.vue";
import { http } from "@/utils";
import Explain from "@/components/ui/Explain.vue";

const provideAspectPage = inject("provideAspectPage");

const route = useRoute();

const modelType = ref(); // 模型字段值

const dialog = ref(false);
const loading = ref(false);
const rating = ref(0);
const voice1 = ref("");
const voice2 = ref("");
const canUpload = ref(false);
const page = ref({ title: "声纹识别" });
const detail = reactive({ data: {} });
const resultRef = ref(); // 分析结果组件ref
const breadcrumbs = ref([
  {
    text: "智能声纹",
    disabled: false,
    href: "#"
  },
  {
    text: "声纹识别",
    disabled: true,
    href: "#"
  }
]);

function uploadFile1(event) {
  const selectedAudio = $("#file1")[0].files[0];
  const audioURL = URL.createObjectURL(selectedAudio);
  const audio1 = $("#audio1")[0];
  audio1.src = audioURL;
  if ($("#file1")[0].files[0] && $("#file2")[0].files[0]) {
    canUpload.value = true;
  } else {
    canUpload.value = false;
  }
}
function uploadFile2(event) {
  const selectedAudio = $("#file2")[0].files[0];
  const audioURL = URL.createObjectURL(selectedAudio);
  const audio2 = $("#audio2")[0];
  audio2.src = audioURL;
  if ($("#file1")[0].files[0] && $("#file2")[0].files[0]) {
    canUpload.value = true;
  } else {
    canUpload.value = false;
  }
}

async function voiceAnalysis() {
  console.log("---voiceAnalysis---");
  if (loading.value) {
    return;
  }
  let file1 = $("#file1")[0].files[0];
  let file2 = $("#file2")[0].files[0];
  if (!file1) {
    return;
  }
  if (!file2) {
    return;
  }
  console.log(file1);

  // 清空分析结果
  // detail.data = {};
  // rating.value = 0;

  loading.value = true;
  const [err, res] = await http.upload({
    url: "/voice/compare",
    data: {
      file1,
      file2,
      modelType: modelType.value,
      threshold: 0.8
    }
    // showLoading: resultRef.value.$el
  });
  loading.value = false;

  if (res) {
    detail.data = res.data;
    rating.value = detail.data.dist * 10;
    provideAspectPage.methods.refreshListPage();
  }
}
</script>
<template>
  <NavBack v-if="route.path === '/voice-print/compare-list/compare'" backUrl="/voice-print/compare-list/list"> 声纹比对 </NavBack>
  <BaseBreadcrumb v-else :title="page.title" :breadcrumbs="breadcrumbs"></BaseBreadcrumb>
  <v-row>
    <v-col cols="12" md="6">
      <UiParentCard title="第一段声音">
        <div class="d-flex flex-column flex-sm-row gap-3 justify-center">
          <audio id="audio1" controls>Your browser does not support the audio element.</audio>
        </div>
        <div class="d-flex flex-column flex-sm-row gap-3 justify-center" style="margin-top: 20px">
          <!-- accept="audio/*, .mp3, .wma, .amr, .wav, .m4a" -->
          <v-file-input
            id="file1"
            ref="file1"
            accept="audio/*, .mp3, .wma, .amr, .wav, .m4a"
            label="选择第一段音频"
            @change="uploadFile1"
            prepend-icon="mdi-volume-high"
            hide-details
            variant="outlined"
          ></v-file-input>
        </div>
      </UiParentCard>
    </v-col>
    <v-col cols="12" md="6">
      <UiParentCard title="第二段声音">
        <div class="d-flex flex-column flex-sm-row gap-3 justify-center">
          <audio id="audio2" controls>Your browser does not support the audio element.</audio>
        </div>
        <div class="d-flex flex-column flex-sm-row gap-3 justify-center" style="margin-top: 20px">
          <!-- accept="audio/*, .mp3, .wma, .amr, .wav, .m4a" -->
          <v-file-input
            id="file2"
            ref="file2"
            accept="audio/*, .mp3, .wma, .amr, .wav, .m4a"
            label="选择第二段音频"
            @change="uploadFile2"
            prepend-icon="mdi-volume-high"
            hide-details
            variant="outlined"
          ></v-file-input>
        </div>
      </UiParentCard>
    </v-col>
    <v-col cols="12">
      <v-card elevation="10">
        <v-card-text>
          <v-row>
            <v-col cols="12" sm="6">
              <Select
                v-model="modelType"
                :mapDictionary="{ code: 'vrp_model_type' }"
                placeholder="请选择模型"
                hide-details
                defaultFirst
                :clearable="false"
              >
                <template #prepend>
                  <Explain>
                    CAMPPlus 训练CN-Celeb数据集<br />
                    CAMPPlus-big 训练更大数据集<br />
                    CAMPPlus-bbig 训练其他超大数据集<br />
                    ERes2Net 训练CN-Celeb数据集<br />
                    ERes2Net-big 训练其他超大数据集<br />
                    EcapaTdnn训练CN-Celeb数据集<br />
                    Res2Net训练CN-Celeb数据集<br />
                    ResNetSE训练CN-Celeb数据集<br />
                    TDNN训练CN-Celeb数据集
                  </Explain>
                </template>
              </Select>
            </v-col>
            <v-col cols="12" sm="6">
              <v-btn :loading="loading" :disabled="!canUpload" color="secondary" width="200" height="40" @click="voiceAnalysis">
                开始分析
              </v-btn>
            </v-col>
          </v-row>
        </v-card-text>
      </v-card>
    </v-col>
    <v-col cols="12">
      <UiParentCard title="分析结果" ref="resultRef">
        <v-card-actions class="flex-column">
          <v-card-text style="font-size: 20px"> {{ detail.data.message }} </v-card-text>
          <v-rating v-model="rating" length="10" readonly density="compact" size="large" color="warning"></v-rating>
        </v-card-actions>
      </UiParentCard>
    </v-col>
  </v-row>
  <v-dialog v-model="dialog" max-width="800" title="分析结果">
    <v-card>
      <v-card-actions class="flex-column">
        <v-card-text style="font-size: 20px">
          {{ detail.data.message }}
        </v-card-text>
        <v-rating v-model="rating" length="10" readonly density="compact" size="large" color="warning"></v-rating>
        <v-btn
          color="secondary"
          class="px-4 mx-auto"
          @click="dialog = false"
          variant="tonal"
          width="200"
          height="48"
          style="margin-top: 20px"
          >关闭</v-btn
        >
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
<style scoped>
.voice_1 {
  background-color: #ebf9fa;
}
.voice_2 {
  background-color: #fff5f2;
}
</style>
