<template>
  <NavBack backUrl="/digital-human/video-list/list">创建视频</NavBack>
  <v-row>
    <v-col cols="12" md="6">
      <UiParentCard title1="人物视频"> <DigitalHumanSelector /> </UiParentCard>
    </v-col>
    <v-col cols="12" md="6">
      <UiParentCard title="开始数字人合成">
        <v-form ref="refForm" class="my-form">
          <v-row>
            <v-col cols="12">
              <v-label class="mb-2 required"
                >标题<Explain class="ml-2">用于列表展示和搜索，能够快速了解基本信息</Explain></v-label
              >
              <v-text-field
                density="compact"
                variant="outlined"
                placeholder="请输入标题"
                hide-details="auto"
                :rules="rules.title"
                v-model="formData.title"
              >
              </v-text-field>
            </v-col>

            <v-col cols="12">
              <v-label class="mb-2 required" style="white-space: inherit">
                <div>请输入语音播放文本，文本内容小于500个字(包括标点符号)。</div>
              </v-label>
              <div class="generate-text">
                <v-row>
                  <v-col cols="12" class="flex space-x-2">
                    <div class="flex-1">
                      <v-text-field
                        density="compact"
                        variant="outlined"
                        placeholder="请输入关键字"
                        hide-details="auto"
                        v-model="state.textKey"
                        @keyup.enter="onAIText"
                        clearable
                      >
                      </v-text-field>
                    </div>
                    <AiBtn id="btnGenerateText" color="secondary" height="40" @click="onAIText" :disabled="!state.textKey"
                      >AI生成文案</AiBtn
                    >
                  </v-col>
                </v-row>
              </div>
              <v-textarea
                ref="refTextarea"
                v-model.trim="formData.text"
                :rules="rules.text"
                placeholder="请输入语音播放文本"
                counter
                rows="5"
                maxlength="500"
              >
              </v-textarea>
            </v-col>

            <v-col cols="12">
              <v-label class="mb-2 required">模型选择</v-label>
              <Select
                :mapDictionary="{ code: 'digitalhuman_synthesis_model' }"
                v-model="formData.synthesisModel"
                defaultFirst
                :clearable="false"
              >
              </Select>
            </v-col>

            <v-col v-if="formData.synthesisModel === 'wav2lip'" cols="12">
              <v-label class="mb-2 required">是否超分</v-label>
              <v-switch v-model="formData.isGfpgan" color="primary" hide-details="auto"></v-switch>
            </v-col>

            <v-col cols="12">
              <v-label class="mb-2 required">请选择需要合成的发声人</v-label>
              <v-expansion-panels class="mb-3">
                <v-expansion-panel elevation="1">
                  <v-expansion-panel-title>高级搜索</v-expansion-panel-title>
                  <v-expansion-panel-text class="mt-4" eager>
                    <v-row>
                      <v-col cols="4">
                        <Select
                          :mapDictionary="{ code: 'speak_provider' }"
                          label="请选择供应"
                          v-model="speakerSearchData.provider"
                          @change="reloadSpeaker"
                        >
                        </Select>
                      </v-col>
                      <v-col cols="4">
                        <Select
                          :mapDictionary="{ code: 'speak_lang' }"
                          label="请选择语言"
                          v-model="speakerSearchData.lang"
                          @change="reloadSpeaker"
                        >
                        </Select>
                      </v-col>
                      <v-col cols="4">
                        <Select
                          :mapDictionary="{ code: 'speak_gender' }"
                          label="请选择性别"
                          v-model="speakerSearchData.gender"
                          @change="reloadSpeaker"
                        >
                        </Select>
                      </v-col>
                      <v-col cols="4">
                        <Select
                          :mapDictionary="{ code: 'speak_age_group' }"
                          label="请选择年龄段"
                          v-model="speakerSearchData.ageGroup"
                          @change="reloadSpeaker"
                        >
                        </Select>
                      </v-col>
                      <v-col cols="4">
                        <Select
                          :mapDictionary="{ code: 'speak_style' }"
                          label="请选择说话风格"
                          v-model="speakerSearchData.speakStyle"
                          @change="reloadSpeaker"
                        >
                        </Select>
                      </v-col>
                      <v-col cols="4">
                        <Select
                          :mapDictionary="{ code: 'speak_area' }"
                          label="请选择适应范围"
                          v-model="speakerSearchData.area"
                          @change="reloadSpeaker"
                        >
                        </Select>
                      </v-col>
                    </v-row>
                  </v-expansion-panel-text>
                </v-expansion-panel>
              </v-expansion-panels>
              <SpeakerSelector ref="refSpeakerSelector" v-model="formData.speakName" />
            </v-col>
            <v-col cols="12" class="hv-center">
              <AiBtn id="btnSubmit" color="secondary" width="200" height="48" size="large" @click="onSubmit">开始合成</AiBtn>
            </v-col>
          </v-row>
        </v-form>
      </UiParentCard>
    </v-col>
  </v-row>
</template>
<script setup lang="ts">
import { reactive, toRefs, ref, onMounted, provide, nextTick } from "vue";
import BaseBreadcrumb from "@/components/shared/BaseBreadcrumb.vue";
import UiParentCard from "@/components/shared/UiParentCard.vue";
import Explain from "@/components/ui/Explain.vue";
import DigitalHumanSelector from "./components/DigitalHumanSelector.vue";
import { http, validator } from "@/utils";
import { toast } from "vue3-toastify";
import { ItfProvideState, ItfSpeaker } from "./types";
import SpeakerSelector from "@/components/business/SpeakerSelector.vue";
import { useRouter } from "vue-router";
import { useScroll } from "@/hooks/useScroll";
import NavBack from "@/components/business/NavBack.vue";
const { scrollRef, scrollToBottom } = useScroll();

// import type * as Types from "./types";
const router = useRouter();

const refForm = ref();
const refTextarea = ref(); // 生成语音播放文本的输入框

const state = reactive<ItfProvideState>({
  style: {},
  textKey: "", // 生成语音播放文本的关键字
  formData: {
    text: "",
    title: "",
    speakName: "",
    isGfpgan: false,
    synthesisModel: null
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
    video: "",
    mediaType: ""
  },
  listSpeaker: []
});
const { style, formData, selectedSpeaker } = toRefs(state);
provide("provieDigitalHumanEdit", state);

const speakerSearchData = reactive({
  provider: null,
  lang: null,
  gender: null,
  ageGroup: null,
  speakStyle: null,
  area: null
});
const refSpeakerSelector = ref();

const rules = reactive({
  title: [value => !!value || "请输入标题"],
  text: [value => !!value || "请输入语音播放文本", value => value.length <= 500 || "文本内容不能超过500个字"]
});

const reloadSpeaker = () => {
  refSpeakerSelector.value.reload({
    lang: speakerSearchData.lang,
    provider: speakerSearchData.provider,
    gender: speakerSearchData.gender,
    ageGroup: speakerSearchData.ageGroup,
    speakStyle: speakerSearchData.speakStyle,
    area: speakerSearchData.area
  });
};

// 生成语音播放文本
const onAIText = async () => {
  if (state.textKey) {
    const data = {
      model: "gpt-4-turbo-preview",
      maxTokens: 512,
      temperature: 0,
      topP: 0,
      messages: [
        {
          role: "system",
          content: "你是一个文案专家，请根据输入的关键词，生成一段优秀的营销文案"
        },
        {
          role: "user",
          content: state.textKey // 关键字
        }
      ]
    };
    let hasToast = false; // 是否已弹出异常提示信息
    await http.post({
      url: "/channels/chat/completions",
      timeout: 300 * 1000, // 请求超时时间设置为5分钟
      data,
      showLoading: "btn#btnGenerateText",
      onDownloadProgress: event => {
        const xhr = event.target;
        const { responseText } = xhr;
        const lastIndex = responseText.lastIndexOf("\n", responseText.length - 2);
        let chunk = responseText;
        if (lastIndex !== -1) {
          chunk = responseText.substring(lastIndex);
        }
        try {
          const responseJson = JSON.parse(chunk);
          if (responseJson.success) {
            const data = responseJson.data;
            if (data && data.fullContent) {
              state.formData.text = data.fullContent;
              scrollToBottom();
            }
          } else {
            if (!hasToast) {
              toast.error(responseJson.message || "AI生成文案异常");
              hasToast = true;
            }
            console.error("AI生成文案异常：", responseJson);
          }
        } catch (error) {
          // console.error(error, chunk);
        }
      }
    });
  }
};

const onSubmit = async () => {
  const { valid } = await refForm.value.validate();
  if (valid) {
    if (state.formData.synthesisModel !== "wav2lip") {
      state.formData.isGfpgan = false;
    }
    const [err, res] = await http.post({
      url: "/api/digitalhuman/synthesis/create",
      showSuccess: true,
      showLoading: "btn#btnSubmit",
      data: {
        ...state.formData,
        dhpName: state.selectedDigitalHuman.name
      }
    });
    if (res) {
      router.push("/digital-human/video-list/list");
    }
  } else {
    let errorMsg = "请处理页面标错的地方后，再尝试提交";
    toast.warning(errorMsg);
  }
};
onMounted(() => {
  scrollRef.value = refTextarea.value.$el.querySelector("textarea");
});
</script>
<style lang="scss" scoped>
.preview-voice {
  width: 200px;
}
.generate-text {
  margin-bottom: 10px;
}
</style>
