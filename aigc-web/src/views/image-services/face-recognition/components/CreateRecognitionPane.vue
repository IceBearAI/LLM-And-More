<template>
  <Pane ref="refPane">
    <v-row>
      <v-col v-if="paneConfig.operateType === 'add'" cols="12" md="7">
        <UiParentCard title="输入">
          <v-form ref="refForm" class="my-form">
            <Select
              placeholder="请选择检测类型"
              :rules="[v => !!v || '请选择检测类型']"
              :mapDictionary="{ code: 'face_check_type' }"
              v-model="formData.checkType"
            >
              <template #prepend>
                <label class="required">检测类型</label>
              </template>
            </Select>
            <v-file-input
              v-model="formData.checkImage"
              :prepend-icon="null"
              accept="image/*"
              label="请上传检测图片"
              hide-details="auto"
              variant="outlined"
              :rules="[v => v.length > 0 || '请上传检测图片']"
            >
              <template #prepend>
                <label class="required">检测图片</label>
              </template>
              <template #append>
                <v-img
                  :src="previewImageUrl.checkImage"
                  width="80px"
                  alt="previewImageUrl"
                  cover
                  class="rounded-md align-end text-right"
                ></v-img>
              </template>
            </v-file-input>
            <template v-if="formData.checkType === 2">
              <v-file-input
                v-model="formData.baseImage"
                :prepend-icon="null"
                accept="image/*"
                label="请上传比对图片"
                hide-details="auto"
                variant="outlined"
                :rules="[v => v.length > 0 || '请上传比对图片']"
              >
                <template #prepend>
                  <label class="required">比对图片</label>
                </template>
                <template #append>
                  <v-img
                    :src="previewImageUrl.baseImage"
                    width="80px"
                    alt="previewImageUrl"
                    cover
                    class="rounded-md align-end text-right"
                  ></v-img>
                </template>
              </v-file-input>
              <v-slider
                class="mx-0"
                v-model="formData.tolerance"
                color="primary"
                :max="1"
                :min="0"
                :step="0.1"
                hide-details="auto"
                thumb-label
              >
                <template v-slot:append>
                  <div class="text-center" style="width: 28px">{{ formData.tolerance }}</div>
                </template>
                <template #prepend>
                  <label class="required">比对阈值</label>
                </template>
              </v-slider>
            </template>
          </v-form>
          <v-btn color="primary" block size="large" flat :loading="submitLoading" @click="onSubmit">开始检测</v-btn>
        </UiParentCard>
      </v-col>
      <v-col
        :class="{ 'result-right': paneConfig.operateType === 'edit' }"
        cols="12"
        :md="paneConfig.operateType === 'add' ? 5 : 0"
      >
        <UiParentCard title="输出">
          <v-input v-if="result.inputS3Url">
            <template #prepend> <label>检测图片：</label></template>
            <img :src="result.inputS3Url" width="200" alt="检测图片" class="rounded-md align-end text-right" />
          </v-input>
          <v-input v-if="result.outputS3Url">
            <template #prepend> <label>比对图片：</label></template>
            <img :src="result.outputS3Url" width="200" alt="比对图片" class="rounded-md align-end text-right" />
          </v-input>
          <v-input>
            <template #prepend> <label>人脸个数：</label></template>
            {{ result.faceNum }}
          </v-input>
          <v-input v-if="result.outputS3Url">
            <template #prepend> <label>比对阈值：</label></template>
            {{ result.denoiseStrength }}
          </v-input>
          <v-input v-if="formData.checkType === 2 || result.outputS3Url">
            <template #prepend> <label>同一个人：</label></template>
            {{ result.isSame === undefined ? "" : result.isSame ? "是" : "否" }}
          </v-input>
        </UiParentCard>
      </v-col>
    </v-row>
  </Pane>
</template>
<script setup lang="ts">
import { reactive, ref, computed } from "vue";
import UiParentCard from "@/components/shared/UiParentCard.vue";
import { http } from "@/utils";

interface IFormData {
  checkType: number | null;
  checkImage: any[];
  baseImage: any[];
  tolerance: number;
}

const emits = defineEmits(["submit"]);

const paneConfig = reactive({
  operateType: "add"
});
const refPane = ref();
const refForm = ref();
const formData = reactive<IFormData>({
  checkType: null,
  checkImage: [],
  baseImage: [],
  tolerance: 0.6
});
const result = ref<Record<string, any>>({});
const submitLoading = ref(false);

const previewImageUrl = computed(() => {
  const result = {
    checkImage: "",
    baseImage: ""
  };
  ["checkImage", "baseImage"].forEach(key => {
    if (formData[key] && formData[key].length > 0) {
      result[key] = URL.createObjectURL(formData[key][0]);
    }
  });
  return result;
});

const onSubmit = async () => {
  let { valid } = await refForm.value.validate();
  if (valid) {
    submitLoading.value = true;
    const [err, res] = await http.upload({
      url: "/esrgan/face/recognition",
      data: formData
    });

    if (res) {
      result.value = res;
      emits("submit");
    }
    submitLoading.value = false;
  }
};

defineExpose({
  show({ title, operateType, infos }) {
    refPane.value.show({
      title,
      width: operateType === "add" ? "1000px" : "600px",
      hasSubmitBtn: false
    });
    paneConfig.operateType = operateType;
    if (paneConfig.operateType === "add") {
      formData.checkType = null;
      formData.checkImage = [];
      formData.baseImage = [];
      formData.tolerance = 0.6;
      result.value = {};
    } else {
      result.value = infos;
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
.result-right label {
  width: 150px;
}
</style>
