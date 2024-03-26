<script setup lang="ts">
import Logo from "@/layouts/full/logo/Logo.vue";
/* Login form */
import LoginForm from "@/components/auth/LoginForm.vue";
import { useAppStore } from "@/stores";
import Rainbow from "@/components/ui/Rainbow.vue";

const appStore = useAppStore();
</script>

<template>
  <div class="pa-3">
    <v-row class="h-100vh mh-100 auth">
      <v-col cols="12" lg="7" xl="8" class="d-lg-flex align-center justify-center authentication position-relative">
        <div class="auth-header pt-lg-6 pt-2 px-sm-6 px-3 pb-lg-6 pb-0">
          <div class="position-relative z-10">
            <Logo />
          </div>
        </div>
        <div class="position-relative z-10">
          <img
            src="@/assets/images/backgrounds/login2.png"
            height="450"
            class="position-relative d-none d-lg-flex"
            alt="login-background"
          />
        </div>
        <Rainbow class="d-none d-lg-block" />
      </v-col>
      <v-col cols="12" lg="5" xl="4" class="d-flex align-center justify-center bg-surface">
        <div class="mt-xl-0 mt-5 mw-100" style="min-width: 375px">
          <h2 class="text-h3 font-weight-semibold mb-8">欢迎使用{{ appStore.appName }}</h2>
          <LoginForm />
        </div>
      </v-col>
    </v-row>
  </div>
</template>
