<template>
  <Pane ref="refPane" @submit="onSubmit">
    <div class="mx-auto mt-3" style="width: 540px">
      <v-form ref="refForm" class="my-form">
        <Select
          placeholder="请选择项目ID"
          :rules="rules.projectId"
          :mapDictionary="{ code: 'voice_personal_project' }"
          v-model="formData.projectId"
          defaultFirst
          :clearable="false"
        >
          <template #prepend>
            <label class="required">项目ID</label>
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
          <template #prepend> <label class="required">授权音频</label></template>
          <template #append>
            <a class="link" href="http://docs.paas.creditease.corp/paas-gpt/voice/voice_personal.html" target="_blank"
              >授权音频要求</a
            >
          </template>
        </v-file-input>
        <v-input v-if="audioUrl" hide-details="auto">
          <AiAudio :src="audioUrl" />
          <template #prepend> <label></label></template>
        </v-input>
        <v-textarea v-model.trim="formData.description" placeholder="请输入备注">
          <template #prepend> <label>备注</label></template>
        </v-textarea>
      </v-form>
    </div>
  </Pane>
</template>
<script setup lang="ts">
import { reactive, ref, computed, watch } from "vue";
import AiAudio from "@/components/business/AiAudio.vue";
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
  speakName: props.speakName,
  file: [],
  projectId: "",
  description: ""
});
const projectIdFirst = ref("");

const refPane = ref();
const refForm = ref();
const rules = reactive({
  file: [v => v.length > 0 || "请上传授权音频"],
  projectId: [v => !!v || "请选择项目ID"]
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
      timeout: 6 * 1000, // 请求超时时间设置为10分钟
      showLoading,
      showSuccess: true,
      url: "/personal/consent",
      data: formData
    });

    if (res) {
      refPane.value.hide();
      emits("submit");
    }
  }
};

watch(
  () => formData.projectId,
  (newVal, oldVal) => {
    // 记录第一次默认值
    if (!oldVal) {
      projectIdFirst.value = newVal;
    }
  }
);

defineExpose({
  show({ title, operateType }) {
    refPane.value.show({
      title,
      refForm
    });
    paneConfig.operateType = operateType;
    if (paneConfig.operateType === "add") {
      formData.file = [];
      formData.projectId = projectIdFirst.value;
      formData.description = "";
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
