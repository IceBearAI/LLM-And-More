<template>
  <Pane ref="refPane">
    <v-row>
      <v-col cols="12" md="6">
        <UiParentCard title="输入">
          <v-img
            v-if="previewImageUrl"
            :src="previewImageUrl"
            height="250px"
            alt="previewImageUrl"
            cover
            class="rounded-md align-end text-right mb-5"
          ></v-img>
          <template v-if="paneConfig.operateType === 'add'">
            <v-form ref="refForm" class="my-form">
              <v-file-input
                v-model="formData.file"
                :prepend-icon="null"
                accept="image/*"
                label="请上传图片"
                hide-details="auto"
                variant="outlined"
                :rules="[v => v.length > 0 || '请上传图片']"
              >
                <template #prepend>
                  <label class="required">上传图像</label>
                </template>
              </v-file-input>
              <Select
                placeholder="请选择模型"
                :mapDictionary="{ code: 'esrgan_model_type' }"
                v-model="formData.modelName"
                :rules="[v => !!v || '请选择模型']"
                defaultFirst
              >
                <template #prepend>
                  <label class="required">模型</label>
                </template>
              </Select>
              <v-slider
                v-if="formData.modelName === 'realesr-general-x4v3'"
                v-model="formData.denoiseStrength"
                color="primary"
                :max="1"
                :min="0"
                :step="0.1"
                hide-details="auto"
                thumb-label
              >
                <template v-slot:append>
                  <div class="text-center" style="width: 28px">{{ formData.denoiseStrength }}</div>
                </template>
                <template #prepend>
                  <label>降噪强度</label>
                </template>
              </v-slider>
              <v-checkbox v-model="formData.faceEnhance" color="primary">
                <template #prepend>
                  <label>面部增强</label>
                </template>
              </v-checkbox>
            </v-form>
            <v-btn color="primary" block size="large" flat :loading="submitLoading" @click="onSubmit">开始生成</v-btn>
          </template>
        </UiParentCard>
      </v-col>
      <v-col cols="12" md="6">
        <UiParentCard title="输出">
          <PreviewImage :loading="submitLoading" :url="outputS3Url" height="250px" />
        </UiParentCard>
      </v-col>
    </v-row>
  </Pane>
</template>
<script setup lang="ts">
import { reactive, ref, computed } from "vue";
import UiParentCard from "@/components/shared/UiParentCard.vue";
import { http } from "@/utils";
import PreviewImage from "@/components/business/PreviewImage.vue";

interface IFormData {
  file: any[];
  modelName: string | null;
  denoiseStrength: number;
  faceEnhance: boolean;
}

const emits = defineEmits(["submit"]);

const paneConfig = reactive({
  operateType: "add"
});
const refPane = ref();
const refForm = ref();
const formData = reactive<IFormData>({
  file: [],
  modelName: null,
  denoiseStrength: 0,
  faceEnhance: false
});
const submitLoading = ref(false);
const outputS3Url = ref("");
const inputS3Url = ref("");

const previewImageUrl = computed({
  get() {
    if (paneConfig.operateType === "edit") return inputS3Url.value;
    if (formData.file && formData.file.length > 0) {
      return URL.createObjectURL(formData.file[0]);
    }
    return "";
  },
  set(newValue) {
    inputS3Url.value = newValue;
  }
});

const onSubmit = async () => {
  let { valid } = await refForm.value.validate();
  if (valid) {
    if (formData.modelName !== "realesr-general-x4v3") {
      formData.denoiseStrength = 0;
    }
    submitLoading.value = true;
    const [err, res] = await http.upload({
      url: "/esrgan/image",
      data: formData
    });

    if (res) {
      outputS3Url.value = res.data.s3Url;
      emits("submit");
    }
    submitLoading.value = false;
  }
};

defineExpose({
  show({ title, operateType, infos }) {
    refPane.value.show({
      title,
      width: "1000px",
      hasSubmitBtn: false
    });
    paneConfig.operateType = operateType;
    if (paneConfig.operateType === "add") {
      formData.file = [];
      formData.modelName = "";
      formData.denoiseStrength = 0;
      formData.faceEnhance = false;
      outputS3Url.value = "";
    } else {
      inputS3Url.value = infos.inputS3Url;
      outputS3Url.value = infos.outputS3Url;
    }
  }
});
</script>
<style lang="scss" scoped>
label {
  width: 80px;
  text-align: right;
  .compo-explain {
    margin-top: 1px;
  }
}
</style>
