<template>
  <Pane ref="refPane" @submit="onSubmit">
    <div class="mx-auto mt-3" style="width: 540px">
      <v-form ref="refForm" class="my-form">
        <Select
          placeholder="请选择授权ID"
          :rules="rules.consentId"
          :mapAPI="{
            url: '/personal/consent',
            data: {
              speakName: props.speakName
            },
            labelField: 'consentId',
            valueField: 'consentId'
          }"
          v-model="formData.consentId"
        >
          <template #prepend>
            <label class="required">授权ID <Explain>如果授权ID选项为空，请先去上传授权语音</Explain></label>
          </template>
        </Select>
        <v-file-input
          v-model="formData.file"
          accept="audio/*, .mp3, .wma, .amr, .wav, .m4a"
          :prepend-icon="null"
          hide-details="auto"
          variant="outlined"
          :rules="rules.file"
        >
          <template #prepend> <label class="required">参考音频</label></template>
        </v-file-input>
        <v-input v-if="audioUrl" hide-details="auto">
          <AiAudio :src="audioUrl" />
          <template #prepend> <label></label></template>
        </v-input>
      </v-form>
    </div>
  </Pane>
</template>
<script setup lang="ts">
import { reactive, ref, computed } from "vue";
import AiAudio from "@/components/business/AiAudio.vue";
import Explain from "@/components/ui/Explain.vue";
import { http } from "@/utils";

interface IProps {
  speakName: string;
}
interface IEmits {
  (e: "submit"): void;
}

const props = withDefaults(defineProps<IProps>(), {
  speakName: ""
});
const emits = defineEmits<IEmits>();

const paneConfig = reactive({
  operateType: "add"
});
const formData = reactive({
  file: [],
  consentId: null
});

const refPane = ref();
const refForm = ref();
const rules = reactive({
  file: [v => v.length > 0 || "请上传参考音频"],
  consentId: [v => !!v || "请选择授权ID"]
});

const audioUrl = computed(() => {
  if (formData.file && formData.file.length > 0) {
    return URL.createObjectURL(formData.file[0]);
  }
  return "";
});

const onSubmit = async ({ valid, showLoading }) => {
  if (valid) {
    const [err, res] = await http.upload({
      timeout: 600 * 1000, // 请求超时时间设置为10分钟
      showLoading,
      showSuccess: true,
      url: "/api/personal/voice",
      data: formData
    });

    if (res) {
      refPane.value.hide();
      emits("submit");
    }
  }
};

defineExpose({
  show({ title, operateType }) {
    refPane.value.show({
      title,
      refForm
    });
    paneConfig.operateType = operateType;
    if (paneConfig.operateType === "add") {
      formData.file = [];
      formData.consentId = null;
    }
  }
});
</script>
<style lang="scss" scoped>
label {
  width: 130px;
  text-align: right;
  .compo-explain {
    margin-top: 1px;
  }
}
</style>
