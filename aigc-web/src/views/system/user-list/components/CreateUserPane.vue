<template>
  <Pane ref="refPane" @submit="onSubmit">
    <div class="mx-auto mt-3" style="width: 600px">
      <v-form ref="refForm" class="my-form">
        <v-text-field
          type="text"
          placeholder="请输入昵称"
          hide-details="auto"
          clearable
          :rules="rules.nickname"
          v-model="formData.nickname"
        >
          <template #prepend><label class="required">昵称</label></template>
        </v-text-field>
        <v-text-field
          type="text"
          placeholder="请输入邮箱"
          hide-details="auto"
          clearable
          :rules="rules.email"
          v-model="formData.email"
        >
          <template #prepend><label class="required">邮箱</label></template>
        </v-text-field>
        <v-switch v-model="formData.isLdap" color="primary" hide-details="auto">
          <template #prepend><label>ldap用户</label></template>
        </v-switch>
        <v-text-field
          v-if="!formData.isLdap"
          :append-inner-icon="showPassword ? 'mdi-eye' : 'mdi-eye-off'"
          :type="showPassword ? 'text' : 'password'"
          @click:appendInner="showPassword = !showPassword"
          placeholder="大写英文字母+小写英文字母+数字+特殊字符，8-20位"
          hide-details="auto"
          :rules="rules.password"
          v-model="formData.password"
        >
          <template #prepend><label :class="{ required: !isEdit }">密码</label></template>
        </v-text-field>
        <Select
          :rules="rules.tenantPublicTenantIdItems"
          v-model="formData.tenantPublicTenantIdItems"
          :mapAPI="{ url: '/tenants', data: { pageSize: -1 }, labelField: 'name', valueField: 'publicTenantId' }"
          hide-details="auto"
          multiple
          placeholder="请选择租户"
        >
          <template #prepend>
            <label class="required">绑定租户</label>
          </template>
        </Select>
        <Select
          placeholder="请选择默认语言"
          :rules="rules.language"
          :mapDictionary="{ code: 'system_language' }"
          v-model="formData.language"
        >
          <template #prepend>
            <label class="required">默认语言</label>
          </template>
        </Select>
        <v-switch v-if="isEdit" v-model="formData.status" color="primary" hide-details="auto">
          <template #prepend> <label>启用</label></template>
        </v-switch>
      </v-form>
    </div>
  </Pane>
</template>
<script setup lang="ts">
import { reactive, ref, computed } from "vue";
import { http, validator } from "@/utils";

interface IFormData {
  id?: number;
  nickname: string;
  email: string;
  isLdap: boolean;
  password: string;
  language: string | null;
  tenantPublicTenantIdItems: string[];
  status?: boolean;
}
const initFormData = {
  nickname: "",
  email: "",
  isLdap: false,
  password: "",
  language: null,
  tenantPublicTenantIdItems: null
};

const emits = defineEmits(["submit"]);

const paneConfig = reactive({
  operateType: "add"
});
const formData = ref<IFormData>({ ...initFormData });
const refPane = ref();
const refForm = ref();
const rules = reactive({
  nickname: [v => !!v || "请输入昵称"],
  email: [v => validator.isEmail({ value: v, required: true })],
  password: [v => validator.isPassword({ value: v, required: isEdit.value ? false : true })],
  tenantPublicTenantIdItems: [v => (v && v.length > 0) || "请选择租户"],
  language: [v => !!v || "请选择默认语言"]
});
const showPassword = ref(false);

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
      requestConfig.url = "/accounts";
      requestConfig.method = "post";
    } else {
      requestConfig.url = `/accounts/${formData.value.id}`;
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
      formData.value.id = infos.id;
      formData.value.nickname = infos.nickname;
      formData.value.email = infos.email;
      formData.value.isLdap = infos.isLdap;
      formData.value.language = infos.language;
      formData.value.tenantPublicTenantIdItems = infos.tenants.map(_ => _.publicTenantId);
      formData.value.status = infos.status;
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
