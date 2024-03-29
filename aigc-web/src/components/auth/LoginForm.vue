<script setup lang="ts">
import { ref } from "vue";
import { Form } from "vee-validate";
import { useUserStore } from "@/stores";

const userStore = useUserStore();

const valid = ref(false);
const password = ref("admin");
const username = ref("admin");
const passwordRules = ref([(v: string) => !!v || "请输入密码"]);
const emailRules = ref([(v: string) => !!v.trim() || "请输入账号"]);
const showPwd = ref(false);

function validate(values: any, { setErrors }: any) {
  return userStore.login(username.value, password.value);
  // .catch(error => setErrors({ apiError: error.message }));
}
</script>

<template>
  <div class="d-flex align-center text-center mb-6">
    <div class="text-h6 w-100 px-5 font-weight-regular auth-divider position-relative">
      <span class="bg-surface px-5 py-3 position-relative">登录</span>
    </div>
  </div>
  <Form @submit="validate" v-slot="{ errors, isSubmitting }" class="mt-5">
    <v-label class="text-subtitle-1 font-weight-medium pb-2 text-lightText">用户名</v-label>
    <VTextField
      v-model="username"
      :rules="emailRules"
      class="mb-8"
      required
      hide-details="auto"
      placeholder="账号:admin"
    ></VTextField>
    <v-label class="text-subtitle-1 font-weight-medium pb-2 text-lightText">密码</v-label>
    <VTextField
      v-model="password"
      :rules="passwordRules"
      required
      autocomplete="current-password"
      hide-details="auto"
      :type="showPwd ? 'text' : 'password'"
      class="pwdInput"
      placeholder="密码:admin"
      :append-icon="showPwd ? 'mdi-eye' : 'mdi-eye-off'"
      @click:append="showPwd = !showPwd"
    ></VTextField>
    <v-btn class="mt-6" size="large" :loading="isSubmitting" color="primary" :disabled="valid" block type="submit" flat
      >登录</v-btn
    >
    <div v-if="errors.apiError" class="mt-2">
      <v-alert color="error">{{ errors.apiError }}</v-alert>
    </div>
  </Form>
</template>
<style lang="scss">
.pwdInput {
  position: relative;
  .v-input__append {
    position: absolute;
    right: 12px;
    top: 9px;
  }
}
</style>
