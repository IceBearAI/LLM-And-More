<template>
  <NavBack backUrl="/voice-print/synthesis/voice-list">创建TTS</NavBack>
  <v-row class="mt-1">
    <v-col cols="12" md="6">
      <UiParentCard>
        <Left ref="refFormLeft" />
      </UiParentCard>
    </v-col>
    <v-col cols="12" md="6">
      <UiParentCard>
        <Right ref="refFormRight">
          <v-col cols="12" class="mt-4">
            <div class="hv-center py-3">
              <AiBtn id="btnSubmit" color="secondary" width="200" height="48" size="large" @click="onSubmit">开始合成</AiBtn>
            </div>
          </v-col>
        </Right>
      </UiParentCard>
    </v-col>

    <v-expand-transition>
      <v-col cols="12" v-show="style.showPreview">
        <UiParentCard>
          <v-col cols="12" class="pt-10">
            <UiChildCard title="合成音频" v-loading="style.loadingPreview">
              <AiAudio v-bind="{ ...state.audioInfo }" ref="refAiAudio" />
            </UiChildCard>
          </v-col>
        </UiParentCard>
      </v-col>
    </v-expand-transition>
  </v-row>
</template>
<script setup>
import { reactive, toRefs, ref, onMounted, provide } from "vue";
import UiParentCard from "@/components/shared/UiParentCard.vue";
import UiChildCard from "@/components/shared/UiChildCard.vue";

import { useMapRemoteStore } from "@/stores";
import { http } from "@/utils";
import { toast } from "vue3-toastify";
import NavBack from "@/components/business/NavBack.vue";

import Left from "./components/Left.vue";
import Right from "./components/Right.vue";
import VueScrollTo from "vue-scrollto";
import { useRoute } from "vue-router";
import _ from "lodash";
import $ from "jquery";
import AiAudio from "@/components/business/AiAudio.vue";

const { mappings, loadDictTree } = useMapRemoteStore();
const route = useRoute();
const refFormLeft = ref();
const refFormRight = ref();
const refAiAudio = ref();

const state = reactive({
  style: {
    showPreview: false,
    loadingPreview: false
  },
  formData: {
    provider: "",
    lang: "",
    gender: "",
    ageGroup: "",
    speakStyle: "",
    area: "",
    speakName: "",
    text: "",
    title: "",
    speed: 1,
    tone: null,
    setDemo: false
  },
  /** 所选语言 */
  selectedLanguage: {
    label: "",
    value: ""
  },
  selectedSpeaker: {
    speakName: "",
    speakCname: "",
    headImg: "",
    speakDemo: "",
    gender: "",
    ageGroup: "",
    subTitle: ""
  },
  selectedDigitalHuman: {
    name: "",
    cname: "",
    cover: "",
    video: ""
  },
  audioInfo: {
    src: "",
    gender: "",
    type: "complex"
  }
});
const { style, formData, selectedSpeaker } = toRefs(state);
provide("provideSynthesisVoice", state);

const onSelectSpeaker = item => {
  state.selectedSpeaker = item;
};

const initAudio = info => {
  state.style.loadingPreview = true;
  state.style.showPreview = true;
  setTimeout(() => {
    VueScrollTo.scrollTo(refAiAudio.value.$el, 500);
    _.assign(state.audioInfo, {
      src: info.s3Url,
      gender: info.gender
    });
  }, 100);
  setTimeout(() => {
    state.style.loadingPreview = false;
  }, 500);
};

const onSubmit = async () => {
  const { valid: validLeft } = await refFormLeft.value.validate();
  const { valid: validRight } = await refFormRight.value.validate();
  if (validLeft && validRight) {
    //验证通过
    const [err, res] = await http.post({
      showLoading: "btn#btnSubmit",
      showSuccess: true,
      url: "/api/voice/tts",
      data: _.omit(state.formData, ["provider", "lang", "gender", "ageGroup", "speakStyle", "area"])
    });
    if (res) {
      initAudio(res);
    } else {
      state.style.showPreview = false;
    }
  } else {
    let errorMsg = "请处理页面标错的地方后，再尝试提交";
    toast.warning(errorMsg);
  }
};
onMounted(async () => {
  await loadDictTree(["speak_provider", "speak_lang"]);
  let { provider, lang, speakName } = route.query;

  if (provider && mappings.speak_provider?.[provider]) {
    state.formData.provider = provider;
  }
  if (lang && mappings.speak_lang?.[lang]) {
    state.formData.lang = lang;
  }
  if (speakName) {
    state.formData.speakName = speakName;
  }
});
</script>

<style lang="scss" scoped>
.preview-voice {
  width: 200px;
}
</style>
