<template>
  <Pane ref="refPane">
    <v-row>
      <v-col cols="6">
        <UiParentCard title="输入">
          <v-form ref="refForm" class="my-form">
            <v-file-input
              v-model="formData.file"
              :prepend-icon="null"
              accept="video/*"
              label="请上传检测视频"
              hide-details="auto"
              variant="outlined"
              :rules="[v => v.length > 0 || '请上传检测视频']"
            >
              <template #prepend>
                <label class="required">检测视频</label>
              </template>
            </v-file-input>
            <v-slider
              v-model="formData.liveThreshold"
              color="primary"
              :max="10"
              :min="1"
              :step="1"
              hide-details="auto"
              thumb-label
            >
              <template v-slot:append>
                <div class="text-center" style="width: 28px">{{ formData.liveThreshold }}</div>
              </template>
              <template #prepend>
                <label class="required">要求活体动作</label>
              </template>
            </v-slider>
            <v-switch v-model="formData.returnImg" color="primary" hide-details="auto">
              <template #prepend><label>是否返回图片</label></template>
            </v-switch>
          </v-form>
          <v-btn color="primary" block size="large" flat :loading="submitLoading" @click="onSubmit">开始检测</v-btn>
        </UiParentCard>
      </v-col>
      <v-col cols="6">
        <UiParentCard title="输出">
          <v-input>
            <template #prepend> <label>是否为活体：</label></template>
            {{ result.isLive === undefined ? "" : result.isLive ? "是" : "否" }}
          </v-input>
          <v-input>
            <template #prepend> <label>描述：</label></template>
            {{ result.desc }}
          </v-input>
          <v-input>
            <template #prepend> <label>返回图片：</label></template>
            <img
              v-if="result.imgS3Url"
              :src="result.imgS3Url"
              width="200"
              alt="返回图片"
              class="rounded-md align-end text-right"
            />
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
  file: any[];
  liveThreshold: number;
  returnImg: boolean;
}

const emits = defineEmits(["submit"]);

const paneConfig = reactive({
  operateType: "add"
});
const refPane = ref();
const refForm = ref();
const formData = reactive<IFormData>({
  file: [],
  liveThreshold: 1,
  returnImg: false
});
const result = ref<Record<string, any>>({});
const submitLoading = ref(false);

const onSubmit = async () => {
  let { valid } = await refForm.value.validate();
  if (valid) {
    submitLoading.value = true;
    const [err, res] = await http.upload({
      url: "/esrgan/face/live",
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
      width: "1000px",
      hasSubmitBtn: false
    });
    paneConfig.operateType = operateType;
    if (paneConfig.operateType === "add") {
      formData.file = [];
      formData.liveThreshold = 1;
      formData.returnImg = false;
      result.value = {};
    } else {
      //
    }
  }
});
</script>
<style lang="scss" scoped>
label {
  width: 100px;
  text-align: right;
  .compo-explain {
    margin-top: 1px;
  }
}
.result-right label {
  width: 150px;
}
</style>
