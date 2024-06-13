<template>
  <Pane ref="refPane">
    <v-row>
      <v-col cols="6">
        <UiParentCard title="输入">
          <v-form ref="refForm" class="my-form">
            <v-file-input
              v-model="formData.file"
              :prepend-icon="null"
              accept="image/*"
              label="请上传人脸照片"
              hide-details="auto"
              variant="outlined"
              :rules="[v => (v && v.length > 0) || '请上传人脸照片']"
            >
              <template #prepend>
                <label class="required">人脸照片</label>
              </template>
              <template v-if="previewImageUrl" #append>
                <v-img
                  :src="previewImageUrl"
                  width="80px"
                  alt="previewImageUrl"
                  cover
                  class="rounded-md align-end text-right"
                ></v-img>
              </template>
            </v-file-input>
            <v-text-field
              type="text"
              placeholder="只允许字母、数字、“-” 、“_”"
              hide-details="auto"
              clearable
              v-model="formData.userKey"
            >
              <template #prepend><label>用户标识</label></template>
            </v-text-field>
          </v-form>
          <v-btn color="primary" block size="large" flat :loading="submitLoading" @click="onSubmit">开始搜索</v-btn>
        </UiParentCard>
      </v-col>
      <v-col cols="6">
        <UiParentCard title="输出">
          <v-row v-if="result.length > 0">
            <v-col :cols="result.length > 1 ? 6 : 12" v-for="item in result">
              <img class="w-full rounded-md align-top" :src="item.imgUrl" alt="人脸图片" />
              <p class="text-center text-body-1 font-weight-medium mt-2">{{ item.dist.toFixed(4) }}</p>
            </v-col>
          </v-row>
          <el-empty v-else :image-size="42" />
        </UiParentCard>
      </v-col>
    </v-row>
  </Pane>
</template>
<script setup lang="ts">
import { reactive, ref, computed } from "vue";
import UiParentCard from "@/components/shared/UiParentCard.vue";
import { http } from "@/utils";

const emits = defineEmits(["search-success"]);

const formData = reactive({
  file: null,
  userKey: ""
});
const refPane = ref();
const refForm = ref();
const submitLoading = ref(false);
const result = ref([]);

const previewImageUrl = computed(() => {
  if (formData.file && formData.file.length > 0) {
    return URL.createObjectURL(formData.file[0]);
  }
  return "";
});

const onSubmit = async () => {
  let { valid } = await refForm.value.validate();
  if (valid) {
    submitLoading.value = true;
    const [err, res] = await http.upload({
      url: "/esrgan/face/search",
      data: formData
    });

    if (res) {
      result.value = res.list || [];
      emits("search-success", result.value);
    }
    submitLoading.value = false;
  }
};

defineExpose({
  show({ title }) {
    refPane.value.show({
      title,
      hasSubmitBtn: false,
      width: "1000px"
    });
    formData.file = null;
    formData.userKey = "";
    result.value = [];
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
