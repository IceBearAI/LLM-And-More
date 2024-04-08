<template>
  <NavBack backUrl="/voice-print/denoise-list/list">
    {{ page.title }}
  </NavBack>
  <v-row>
    <v-col cols="12">
      <UiParentCard title="降噪">
        <v-row justify="center">
          <v-col cols="12" md="6">
            <v-form ref="formRef" class="my-form">
              <v-row>
                <v-col cols="12">
                  <div class="d-flex flex-column flex-sm-row gap-3 justify-center">
                    <AiAudio :src="audioUrl" />
                  </div>
                  <v-label class="mb-2 font-weight-medium required">音频文件</v-label>
                  <v-file-input
                    v-model="formData.file"
                    accept="audio/*, .mp3, .wma, .amr, .wav, .m4a"
                    prepend-icon="mdi-volume-high"
                    :rules="formRules.file"
                    hide-details="auto"
                    variant="outlined"
                    required
                    multiple
                  ></v-file-input>
                </v-col>
                <v-col cols="12">
                  <v-label class="mb-2 font-weight-medium required"
                    >采样率
                    <Explain class="ml-1"
                      >16K 推荐使用，降噪模型默认使用16K采样率<br />
                      原始采样率 降噪时，保持跟原始音频同样的采样率</Explain
                    ></v-label
                  >
                  <Select
                    v-model="formData.sampleRate"
                    :mapDictionary="{ code: 'denoise_sample_rate' }"
                    placeholder="请选择采样率"
                    :rules="formRules.sampleRate"
                    :clearable="false"
                  ></Select>
                </v-col>
              </v-row>

              <div class="d-flex flex-column">
                <v-btn :loading="submitLoading" color="primary" class="mt-4" block @click="doSubmit"> 提交 </v-btn>
              </div>
            </v-form>
          </v-col>
        </v-row>
      </UiParentCard>
    </v-col>

    <v-col cols="12">
      <UiParentCard title="处理结果">
        <v-row class="text-body-1 pa-5" :ref="resultRef">
          <v-col cols="12">
            <div v-if="resultUrl" class="d-flex flex-column flex-sm-row gap-3 justify-center">
              <AiAudio :src="resultUrl" />
            </div>
            <p v-else class="flex-fill text-center">暂无结果</p>
          </v-col>
        </v-row>
      </UiParentCard>
    </v-col>
  </v-row>
</template>
<script setup lang="ts">
import { ref, reactive, computed, inject } from "vue";
import { useRoute, useRouter } from "vue-router";
import NavBack from "@/components/business/NavBack.vue";
import BaseBreadcrumb from "@/components/shared/BaseBreadcrumb.vue";
import UiParentCard from "@/components/shared/UiParentCard.vue";
import AiAudio from "@/components/business/AiAudio.vue";
import { http } from "@/utils";
import Explain from "@/components/ui/Explain.vue";

const route = useRoute();
const router = useRouter();
const provideAspectPage = inject("provideAspectPage") as any;
const page = ref({ title: "创建降噪" });

const hasSubmitted = ref(false); // 是否已提交

const formRef = ref();
const resultRef = ref();
const submitLoading = ref(false);
const formData = reactive({
  file: [],
  sampleRate: 16 // 采样率
});
const formRules = reactive({
  file: [(v: string) => (v && v.length !== 0) || "请上传音频文件"],
  sampleRate: [v => !!(v || v === 0) || "请选择采样率"]
});

const audioUrl = computed(() => {
  if (formData.file && formData.file.length > 0) {
    return URL.createObjectURL(formData.file[0]);
  }
  return "";
});

// 处理后音频地址
const resultUrl = ref();

const doSubmit = async () => {
  const { valid } = await formRef.value.validate();
  if (valid) {
    // resultUrl.value = ""; // 清空处理结果
    submitLoading.value = true;
    const [err, res] = await http.upload({
      url: "/voice/denoise",
      showSuccess: "处理成功",
      data: formData
      // showLoading: resultRef.value.$el
    });
    submitLoading.value = false;

    if (res) {
      resultUrl.value = res.s3Url;
      provideAspectPage.methods.refreshListPage();
    }
  }
};
</script>
