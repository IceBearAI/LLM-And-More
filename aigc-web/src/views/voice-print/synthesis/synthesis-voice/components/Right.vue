<template>
  <v-form ref="refForm" class="my-form">
    <v-row>
      <v-col cols="12">
        <v-label class="mb-2 required">标题<Explain class="ml-2">用于列表展示和搜索，能够快速了解基本信息</Explain></v-label>
        <v-text-field
          density="compact"
          variant="outlined"
          placeholder="请输入标题"
          hide-details="auto"
          :rules="rules.title"
          v-model.trim="formData.title"
          clearable
        >
        </v-text-field>
      </v-col>
      <v-col cols="12">
        <v-label class="mb-2 required" style="white-space: inherit">
          <div>
            请输入<span class="text-primary font-weight-bold" v-if="selectedLanguageCN">「 {{ selectedLanguageCN }} 」</span
            >语音播放文本，文本内容小于200个字(包括标点符号)。
          </div>
        </v-label>
        <v-textarea v-model.trim="formData.text" :rules="rules.text" placeholder="语音播放文本" counter rows="5" maxlength="200">
        </v-textarea>
      </v-col>
      <v-col cols="12">
        <v-label class="mb-2">语气</v-label>
        <Select :mapDictionary="{ code: 'speak_tone' }" placeholder="请选择语气" v-model="formData.tone"> </Select>
      </v-col>
      <v-col cols="12">
        <v-label class="mb-2 required">语速</v-label>
        <div style="width: 300px">
          <v-slider
            density="compact"
            v-model="formData.speed"
            color="primary"
            :max="speedConfig.max"
            :min="speedConfig.min"
            :step="speedConfig.step"
            thumb-label
            :rules="rules.speed"
            hide-details="auto"
          >
            <template v-slot:append>
              <v-text-field
                v-model.number="formData.speed"
                type="number"
                density="compact"
                :max="speedConfig.max"
                :min="speedConfig.min"
                :step="speedConfig.step"
                style="width: 80px"
              ></v-text-field>
            </template>
          </v-slider>
        </div>
      </v-col>
      <v-col v-if="tenantId === '5f9b3b3d-9b9c-4e1a-8e1a-5a4b4b4b4b4b'" cols="12">
        <v-label class="mb-2">设置为样音</v-label>
        <v-switch v-model="formData.setDemo" color="primary" hide-details="auto"></v-switch>
      </v-col>
      <slot></slot>
    </v-row>
  </v-form>
</template>
<script setup>
import { reactive, toRefs, ref, onMounted, inject, computed } from "vue";
import { http, validator } from "@/utils";
import Explain from "@/components/ui/Explain.vue";
import { useUserStore } from "@/stores";

const { tenantId } = useUserStore().userInfo;
const refForm = ref();
const provideSynthesisVoice = inject("provideSynthesisVoice");
const { formData } = toRefs(provideSynthesisVoice);

const state = reactive({
  speedConfig: {
    min: 0.5,
    max: 2,
    step: 0.1
  }
});
const { speedConfig } = toRefs(state);

const rules = reactive({
  title: [value => !!value || "请输入标题"],
  text: [value => !!value || "请输入语音播放文本"],
  speed: [
    value => {
      if (value) {
        let { min, max } = state.speedConfig;
        if (value < min) {
          return "语速不能低于" + min + "倍";
        } else if (value > max) {
          return "语速不能高于" + max + "倍";
        } else {
          return true;
        }
      } else {
        return "请设置语速";
      }
    }
  ]
});

const onSubmit = () => {};

const selectedLanguageCN = computed(() => {
  let { label } = provideSynthesisVoice.selectedLanguage;
  if (label) {
    /*
      中文（粤语，繁体） -> 中文
      中文（简体） -> 中文
      中文 -> 中文
    */
    return label.trim().replace(/([^(（]+)(.*)$/, "$1");
  }
  return "";
});

defineExpose({
  validate() {
    return refForm.value.validate();
  }
});
</script>
<style lang="scss"></style>
