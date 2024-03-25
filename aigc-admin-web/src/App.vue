<template>
  <el-config-provider :locale="appStore.localeElement" size="small">
    <RouterView></RouterView>
  </el-config-provider>
</template>

<script setup lang="ts">
import { watch } from "vue";
import { ElConfigProvider } from "element-plus";
import { RouterView } from "vue-router";
import { useI18n } from "vue-i18n";
import { useAppStore } from "@/stores";
import moment from "moment";

const appStore = useAppStore();
const i18n = useI18n();

appStore.methods.t = i18n.t;

watch(
  () => appStore.localLanguage,
  val => {
    i18n.locale.value = val;
    moment.locale(appStore.localeMoment);
  },
  {
    immediate: true
  }
);
</script>
