<template>
  <Pane ref="refPane" @submit="onSubmit">
    <div class="mx-auto mt-3" style="width: 540px">
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
          <template #append>
            <v-img
              v-if="previewImageUrl"
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
          :rules="rules.userKey"
          v-model="formData.userKey"
        >
          <template #prepend><label class="required">用户标识</label></template>
        </v-text-field>
        <v-text-field type="text" placeholder="请输入用户名称" hide-details="auto" clearable v-model="formData.userName">
          <template #prepend><label>用户名称</label></template>
        </v-text-field>
      </v-form>
    </div>
  </Pane>
</template>
<script setup lang="ts">
import { reactive, ref, computed } from "vue";
import { http } from "@/utils";

const emits = defineEmits(["submit"]);

const paneConfig = reactive({
  operateType: "add"
});
const formData = reactive({
  file: null,
  userKey: "",
  userName: ""
});
const refPane = ref();
const refForm = ref();
const rules = reactive({
  userKey: [v => /^[a-zA-Z0-9-_]+$/.test(v) || "只允许字母、数字、“-” 、“_”"]
});

const previewImageUrl = computed(() => {
  if (formData.file && formData.file.length > 0) {
    return URL.createObjectURL(formData.file[0]);
  }
  return "";
});

const onSubmit = async ({ valid, showLoading }) => {
  if (valid) {
    const [err, res] = await http.upload({
      showLoading,
      showSuccess: true,
      url: "/esrgan/face/reg",
      data: formData
    });

    if (res) {
      refPane.value.hide();
      emits("submit");
    }
  }
};

defineExpose({
  show({ title, operateType, infos }) {
    refPane.value.show({
      title,
      refForm
    });
    paneConfig.operateType = operateType;
    if (paneConfig.operateType === "add") {
      formData.file = null;
      formData.userKey = "";
      formData.userName = "";
    }
  }
});
</script>
<style lang="scss" scoped>
label {
  width: 120px;
  text-align: right;
  .compo-explain {
    margin-top: 1px;
  }
}
</style>
