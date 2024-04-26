<template>
  <v-form ref="refForm" class="my-form">
    <v-row>
      <v-col cols="4">
        <v-label class="mb-2">供应</v-label>
        <Select
          :mapDictionary="{ code: 'speak_provider' }"
          placeholder="请选择供应"
          v-model="formData.provider"
          @change="reloadSpeaker"
        >
        </Select>
      </v-col>
      <v-col cols="4">
        <v-label class="mb-2">语言</v-label>
        <Select
          :mapDictionary="{ code: 'speak_lang' }"
          placeholder="请选择语言"
          v-model="formData.lang"
          v-model:infos="provideSynthesisVoice.selectedLanguage"
          @change="reloadSpeaker"
        >
        </Select>
      </v-col>
      <v-col cols="4">
        <v-label class="mb-2">性别</v-label>
        <Select
          :mapDictionary="{ code: 'speak_gender' }"
          placeholder="请选择性别"
          v-model="formData.gender"
          @change="reloadSpeaker"
        >
        </Select>
      </v-col>
      <v-col cols="4">
        <v-label class="mb-2">年龄段</v-label>
        <Select
          :mapDictionary="{ code: 'speak_age_group' }"
          placeholder="请选择年龄段"
          v-model="formData.ageGroup"
          @change="reloadSpeaker"
        >
        </Select>
      </v-col>
      <v-col cols="4">
        <v-label class="mb-2">说话风格</v-label>
        <Select
          :mapDictionary="{ code: 'speak_style' }"
          placeholder="请选择说话风格"
          v-model="formData.speakStyle"
          @change="reloadSpeaker"
        >
        </Select>
      </v-col>
      <v-col cols="4">
        <v-label class="mb-2">适应范围</v-label>
        <Select
          :mapDictionary="{ code: 'speak_area' }"
          placeholder="请选择适应范围"
          v-model="formData.area"
          @change="reloadSpeaker"
        >
        </Select>
      </v-col>
      <v-col cols="12">
        <v-label class="mb-2 required">请选择需要合成的发声人</v-label>
        <v-input :rules="rules.speakName" v-model="formData.speakName">
          <SpeakerSelector ref="refSpeakerSelector" v-model="formData.speakName" v-model:infos="state.selectedSpeaker" />
        </v-input>
      </v-col>
    </v-row>
  </v-form>
</template>
<script setup>
import { reactive, toRefs, ref, onMounted, inject } from "vue";
import { http, validator } from "@/utils";
import SpeakerSelector from "@/components/business/SpeakerSelector.vue";

const refForm = ref();
const refSpeakerSelector = ref();

const provideSynthesisVoice = inject("provideSynthesisVoice");
const { formData } = toRefs(provideSynthesisVoice);

const state = reactive({
  selectedSpeaker: {}
});

const rules = reactive({
  provider: [value => !!value || "请选择供应"],
  lang: [value => !!value || "请选择语言"],
  speakName: [value => !!value || "请选择发声人"]
});
const reloadSpeaker = () => {
  refSpeakerSelector.value.reload({
    lang: formData.value.lang,
    provider: formData.value.provider,
    gender: formData.value.gender,
    ageGroup: formData.value.ageGroup,
    speakStyle: formData.value.speakStyle,
    area: formData.value.area
  });
};

defineExpose({
  validate() {
    return refForm.value.validate();
  }
});
</script>

<style lang="scss"></style>
