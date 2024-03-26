<template>
  <Pane ref="refPane" @submit="onSubmit">
    <div class="mx-auto mt-10" style="width: 500px">
      <v-form ref="refForm" class="my-form">
        <v-row>
          <v-col cols="12">
            <v-label class="mb-2">{{ $t("audioMark.audio") }}</v-label>
            <AiAudio type="complex" :src="rawData.audioUrl" />
          </v-col>
          <v-col cols="12" class="mt-5">
            <v-label class="mb-2"
              >{{ $t("audioMark.pane.originalContent") }}
              <span class="ml-1">({{ $t("options.markLanguage." + rawData.originalLanguage) }})</span>
            </v-label>
            <v-textarea v-model="rawData.originalContent" readonly> </v-textarea>
          </v-col>

          <!-- <v-col cols="12">
            <v-label class="mb-2">{{ $t("audioMark.pane.originalLanguage") }}</v-label>
            <Select
              :mapDictionary="{ code: 'audio_tagged_lang', i18nKey: 'markLanguage' }"
              v-model="rawData.originalLanguage"
              disabled
            />
          </v-col> -->

          <v-col cols="12">
            <v-label class="mb-2 required">{{ $t("audioMark.pane.taggedContent") }}</v-label>
            <v-textarea v-model.trim="formData.taggedContent" placeholder="请输入" :rules="rules.taggedContent"> </v-textarea>
          </v-col>
          <v-col cols="12">
            <v-label class="mb-2 required">{{ $t("audioMark.pane.taggedLanguage") }}</v-label>
            <Select
              :mapDictionary="{ code: 'audio_tagged_lang', i18nKey: 'markLanguage' }"
              v-model="formData.taggedLanguage"
              :rules="rules.taggedLanguage"
            />
          </v-col>
        </v-row>
      </v-form>
    </div>
  </Pane>
</template>
<script setup>
import { reactive, toRefs, ref, nextTick } from "vue";
import _ from "lodash";
import AiAudio from "@/components/business/AiAudio.vue";
import { http, validator, format } from "@/utils";
import { useI18n } from "vue-i18n";

const { t } = useI18n(); // 解构出t方法

const state = reactive({
  rawData: {},
  formData: {
    taggedContent: "",
    taggedLanguage: ""
  }
});
const { formData, rawData } = toRefs(state);

const emits = defineEmits(["submit"]);

const refPane = ref();
const refForm = ref();
const rules = reactive({
  taggedContent: [v => !!v || t("audioMark.pane.ruleTaggedContent")],
  taggedLanguage: [v => !!v || t("audioMark.pane.ruleTaggedLanguage")]
});

const onSubmit = async ({ valid, showLoading }) => {
  if (valid) {
    const [err, res] = await http.put({
      url: `/api/voice/audio/${state.rawData.id}`,
      showLoading,
      showSuccess: true,
      data: {
        ...state.formData
      }
    });
    if (res) {
      refPane.value.hide();
      emits("submit");
    }
  }
};

defineExpose({
  show({ infos }) {
    refPane.value.show({
      title: t("audioAnnotation"),
      refForm
    });
    let rawData = _.cloneDeep(infos);
    state.rawData = rawData;
    let { taggedContent, taggedLanguage } = rawData;
    if (!taggedContent) {
      //未设置时，使用原内容
      taggedContent = rawData.originalContent;
    }
    if (!taggedLanguage) {
      //未设置时，使用原语言
      taggedLanguage = rawData.originalLanguage;
    }
    _.assign(state.formData, {
      taggedContent,
      taggedLanguage
    });
  }
});
</script>
<style lang="scss" scoped>
label {
  width: 100%;
  text-align: right;
  .compo-explain {
    margin-top: 1px;
  }
}
</style>
