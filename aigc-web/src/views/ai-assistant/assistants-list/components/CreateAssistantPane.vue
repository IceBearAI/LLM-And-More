<template>
  <Pane ref="refPane" @submit="onSubmit">
    <div class="mx-auto mt-3" style="width: 540px">
      <v-form ref="refForm" class="my-form">
        <v-row class="justify-center">
          <div class="avatar position-relative">
            <v-avatar size="100">
              <img :src="showAvatar" width="100" alt="avatar" />
            </v-avatar>
            <custom-upload :file-type="['image/jpeg', 'image/png', 'image/jpg']" :file-size="5" @after-upload="handleAfterUpload">
              <template #trigger>
                <div class="avatar-add-btn">
                  <v-btn icon flat color="info" size="x-small"><IconPlus stroke-width="2.5" :size="20" /></v-btn>
                </div>
              </template>
            </custom-upload>
          </div>
        </v-row>
        <v-text-field
          type="text"
          placeholder="请输入助手名称"
          hide-details="auto"
          clearable
          :rules="rules.name"
          v-model="formData.name"
        >
          <template #prepend><label class="required">助手名称</label></template>
        </v-text-field>
        <v-textarea v-model.trim="formData.instructions" hide-details="auto" placeholder="请输入给AI的指令">
          <template #prepend> <label>指令</label></template>
        </v-textarea>
        <Model-Select placeholder="请选择模型" :rules="rules.modelName" v-model="formData.modelName" hide-details="auto">
          <template #prepend>
            <label class="required">模型</label>
          </template>
        </Model-Select>
        <v-textarea v-model.trim="formData.remark" placeholder="请输入备注">
          <template #prepend> <label>备注</label></template>
        </v-textarea>
      </v-form>
    </div>
  </Pane>
</template>
<script setup lang="ts">
import { computed, reactive, ref } from "vue";
import ModelSelect from "@/components/business/ModelSelect.vue";
import { IconPlus } from "@tabler/icons-vue";
import { http } from "@/utils";
import CustomUpload from "@/components/business/CustomUpload.vue";

interface IFormData {
  assistantId?: string;
  avatar: string;
  name: string;
  instructions: string;
  modelName: string | null;
  remark: string;
}
const initFormData = {
  avatar: "",
  name: "",
  instructions: "",
  modelName: null,
  remark: ""
};

const emits = defineEmits(["submit"]);

const paneConfig = reactive({
  operateType: "add"
});
const formData = ref<IFormData>({ ...initFormData });
const refPane = ref();
const refForm = ref();
const rules = reactive({
  name: [v => !!v || "请输入助手名称"],
  modelName: [v => !!v || "请选择模型"]
});

const showAvatar = computed(() => {
  return formData.value.avatar || new URL("../images/defaultAvatar.png", import.meta.url).href;
});

const handleAfterUpload = ({ res }) => {
  if (res) {
    formData.value.avatar = res.s3Url;
  }
};

const onSubmit = async ({ valid, showLoading }) => {
  if (valid) {
    const requestConfig = {
      url: "",
      method: ""
    };
    if (paneConfig.operateType == "add") {
      requestConfig.url = "/assistants/create";
      requestConfig.method = "post";
    } else {
      requestConfig.url = `/assistants/${formData.value.assistantId}`;
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
    } else {
      formData.value = { ...infos };
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
.avatar-add-btn {
  position: absolute;
  right: 0;
  bottom: 0;
  border: 2px solid #fff;
  border-radius: 50%;
}
</style>
