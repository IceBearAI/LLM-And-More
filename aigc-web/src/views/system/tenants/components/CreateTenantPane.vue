<template>
  <Pane ref="refPane" @submit="onSubmit">
    <div class="mx-auto mt-3" style="width: 600px">
      <v-form ref="refForm" class="my-form">
        <v-text-field
          type="text"
          placeholder="请输入名称"
          hide-details="auto"
          clearable
          :rules="rules.name"
          v-model="formData.name"
        >
          <template #prepend><label class="required">名称</label></template>
        </v-text-field>
        <v-text-field
          type="text"
          placeholder="请输入联系人邮箱"
          hide-details="auto"
          clearable
          :rules="rules.contactEmail"
          v-model="formData.contactEmail"
        >
          <template #prepend><label class="required">联系人邮箱</label></template>
        </v-text-field>
        <Select
          :rules="rules.modelNames"
          v-model="formData.modelNames"
          :mapAPI="{ url: '/channels/models', data: { pageSize: -1 }, labelField: 'modelName', valueField: 'modelName' }"
          hide-details="auto"
          multiple
          placeholder="请选择模型"
        >
          <template #prepend>
            <label class="required">模型</label>
          </template>
        </Select>
      </v-form>
    </div>
  </Pane>
</template>
<script setup lang="ts">
import { reactive, ref, computed } from "vue";
import { http, validator } from "@/utils";

interface IFormData {
  id?: number;
  name: string;
  contactEmail: string;
  modelNames: string[] | null;
}
const initFormData = {
  name: "",
  contactEmail: "",
  modelNames: null
};

const emits = defineEmits(["submit"]);

const paneConfig = reactive({
  operateType: "add"
});
const formData = ref<IFormData>({ ...initFormData });
const refPane = ref();
const refForm = ref();
const rules = reactive({
  name: [v => !!v || "请输入名称"],
  contactEmail: [v => validator.isEmail({ value: v, required: true })],
  modelNames: [v => (v && v.length > 0) || "请选择模型"]
});

const isEdit = computed(() => {
  return paneConfig.operateType === "edit";
});

const onSubmit = async ({ valid, showLoading }) => {
  if (valid) {
    const requestConfig = {
      url: "",
      method: ""
    };
    if (paneConfig.operateType == "add") {
      requestConfig.url = "/tenants";
      requestConfig.method = "post";
    } else {
      requestConfig.url = `/tenants/${formData.value.id}`;
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
  width: 120px;
  text-align: right;
  .compo-explain {
    margin-top: 1px;
  }
}
</style>
