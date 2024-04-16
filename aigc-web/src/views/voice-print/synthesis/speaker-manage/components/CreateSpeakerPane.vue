<template>
  <Pane ref="refPane" @submit="onSubmit">
    <div class="mx-auto mt-3" style="width: 540px">
      <v-form ref="refForm" class="my-form">
        <Select
          placeholder="请选择供应商"
          :rules="rules.provider"
          :mapDictionary="{ code: 'speak_provider' }"
          v-model="formData.provider"
          :disabled="isEdit"
        >
          <template #prepend>
            <label class="required">供应 <Explain>供应商指的是外部服务提供，自己有服务请选择Local</Explain></label>
          </template>
        </Select>
        <v-text-field
          type="text"
          placeholder="请输入标识"
          hide-details="auto"
          clearable
          :rules="rules.speakName"
          v-model="formData.speakName"
          :disabled="isEdit"
        >
          <template #prepend><label class="required">标识</label></template>
        </v-text-field>
        <v-text-field
          type="text"
          placeholder="请输入姓名"
          hide-details="auto"
          clearable
          :rules="rules.speakCname"
          v-model="formData.speakCname"
        >
          <template #prepend><label class="required">姓名</label></template>
        </v-text-field>
        <v-radio-group hide-details="auto" v-model="formData.gender" inline :disabled="isEdit">
          <template v-for="item in genderOptions">
            <v-radio :label="item.label" color="primary" :value="item.value"></v-radio>
          </template>
          <template #prepend><label class="required">性别</label></template>
        </v-radio-group>
        <Select
          placeholder="请选择年龄段"
          :rules="rules.ageGroup"
          :mapDictionary="{ code: 'speak_age_group' }"
          v-model="formData.ageGroup"
        >
          <template #prepend>
            <label class="required">年龄段</label>
          </template>
        </Select>
        <Select
          placeholder="请选择语言"
          :rules="rules.lang"
          :mapDictionary="{ code: 'speak_lang' }"
          v-model="formData.lang"
          :disabled="isEdit"
        >
          <template #prepend>
            <label class="required">语言</label>
          </template>
        </Select>
        <Select
          placeholder="请选择风格"
          :rules="rules.speakStyle"
          :mapDictionary="{ code: 'speak_style' }"
          v-model="formData.speakStyle"
        >
          <template #prepend>
            <label class="required">风格</label>
          </template>
        </Select>
        <Select placeholder="请选择适应范围" :rules="rules.area" :mapDictionary="{ code: 'speak_area' }" v-model="formData.area">
          <template #prepend>
            <label class="required">适应范围</label>
          </template>
        </Select>
        <v-input hide-details="auto">
          <template v-if="headImgInfos && headImgInfos.s3Url">
            <v-alert color="borderColor" variant="outlined" density="compact">
              <v-avatar size="60">
                <v-img :transition="false" :src="headImgInfos.s3Url" alt="上传成功后的头像" cover />
              </v-avatar>
              <template #close>
                <v-icon class="text-24 opacity-50 cursor-pointer" color="textPrimary" @click="headImgClose"
                  >mdi-close-circle</v-icon
                >
              </template>
            </v-alert>
          </template>
          <template v-else>
            <UploadFile
              accept="image/*"
              v-model="formData.headImgFileId"
              v-model:infos="headImgInfos"
              :prepend-icon="null"
              prepend-inner-icon="mdi-camera"
            />
          </template>
          <template #prepend> <label>头像</label></template>
        </v-input>
        <v-switch v-model="formData.enabled" color="primary" hide-details="auto">
          <template #prepend><label class="required">启用</label></template>
        </v-switch>
        <v-textarea v-model.trim="formData.remark" placeholder="请输入备注" clearable>
          <template #prepend> <label>备注</label></template>
        </v-textarea>
      </v-form>
    </div>
  </Pane>
</template>
<script setup lang="ts">
import { computed, reactive, ref } from "vue";
import Explain from "@/components/ui/Explain.vue";
import UploadFile from "@/components/business/UploadFile.vue";
import { http } from "@/utils";
import { useMapRemoteStore } from "@/stores";

interface IFormData {
  id?: number;
  provider: string | null;
  speakName: string;
  speakCname: string;
  gender: number | null;
  ageGroup: number | null;
  lang: string | null;
  speakStyle: number | null;
  area: number | null;
  headImgFileId: string;
  enabled: boolean;
  remark: string;
}
const initFormData = {
  provider: null,
  speakName: "",
  speakCname: "",
  gender: 1,
  ageGroup: null,
  lang: null,
  speakStyle: null,
  area: null,
  headImgFileId: "",
  enabled: false,
  remark: ""
};

const emits = defineEmits(["submit"]);

const { options } = useMapRemoteStore(); // 主页面已经请求过speak_gender

const paneConfig = reactive({
  operateType: "add"
});
const formData = ref<IFormData>({ ...initFormData });
const headImgInfos = ref(null);
const refPane = ref();
const refForm = ref();
const rules = reactive({
  provider: [v => !!v || "请选择供应商"],
  speakName: [v => !!v || "请输入标识"],
  speakCname: [v => !!v || "请输入姓名"],
  ageGroup: [v => !!v || "请选择年龄段"],
  lang: [v => !!v || "请选择语言"],
  speakStyle: [v => !!v || "请选择风格"],
  area: [v => !!v || "请选择适应范围"]
});

const isEdit = computed(() => {
  return paneConfig.operateType === "edit";
});

const genderOptions = computed(() => {
  return options["speak_gender"];
});

const headImgClose = () => {
  formData.value.headImgFileId = "";
  headImgInfos.value = null;
};

const onSubmit = async ({ valid, showLoading }) => {
  if (valid) {
    const requestConfig = {
      url: "",
      method: ""
    };
    if (paneConfig.operateType == "add") {
      requestConfig.url = "/voice/speak";
      requestConfig.method = "post";
    } else {
      requestConfig.url = `/voice/speak/${formData.value.id}`;
      requestConfig.method = "put";
    }
    const [err, res] = await http[requestConfig.method]({
      showLoading,
      showSuccess: true,
      url: requestConfig.url,
      data: formData.value
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
      formData.value = { ...initFormData };
      headImgInfos.value = null;
    } else {
      formData.value = { ...infos };
      headImgInfos.value = {
        s3Url: infos.headImg
      };
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
