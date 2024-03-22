<template>
  <Pane ref="refPane">
    <v-row>
      <v-col cols="12" md="6">
        <UiParentCard title="输入">
          <v-img
            v-if="previewImageUrl"
            :src="previewImageUrl"
            height="350px"
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
            </v-form>
            <v-btn color="primary" block size="large" flat :loading="submitLoading" @click="onSubmit">开始生成</v-btn>
          </template>
        </UiParentCard>
      </v-col>
      <v-col cols="12" md="6">
        <UiParentCard title="输出">
          <PreviewImage :loading="submitLoading" :url="outputS3Url" />
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
}

const emits = defineEmits(["submit"]);

const paneConfig = reactive({
  operateType: "add"
});
const refPane = ref();
const refForm = ref();
const formData = reactive<IFormData>({
  file: []
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
    submitLoading.value = true;
    const [err, res] = await http.upload({
      url: "/esrgan/rmbg",
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
