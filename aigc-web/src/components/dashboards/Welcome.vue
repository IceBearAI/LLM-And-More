<script setup lang="ts">
import { useI18n } from "vue-i18n";
import { useAppStore } from "@/stores";
const appStore = useAppStore();
const { t } = useI18n(); // 解构出t方法
</script>
<template>
  <v-card elevation="10" class="compo-welcome">
    <v-card-text class="text-center">
      <h1 class="text-h1">{{ t("welcome") + " , " + t("appName") }}</h1>
    </v-card-text>
    <img class="bg-img-1" src="@/assets/images/backgrounds/welcome-bg2-2x-svg.svg" />
  </v-card>
</template>

<style scoped lang="scss">
.compo-welcome {
  height: calc(100vh - 90px);
  margin-bottom: -100px;
  min-height: 600px;
  h1 {
    margin-top: 160px;
  }
}
.bg-img-1 {
  position: absolute;
  right: 0px;
  width: 377px;
  bottom: 30px;
}
</style>
