<script setup lang="ts">
import { ref, watch, computed } from "vue";
import { useCustomizerStore } from "../../../stores/customizer";
import { IconMenu2 } from "@tabler/icons-vue";
import LanguageDD from "./LanguageDD.vue";
import ProfileDD from "./ProfileDD.vue";
import TenantDD from "./TenantDD.vue";

const customizer = useCustomizerStore();
const priority = ref(customizer.setHorizontalLayout ? 0 : 0);

watch(priority, newPriority => {
  priority.value = newPriority;
});
</script>

<template>
  <v-app-bar elevation="0" :priority="priority" height="64" color="background" id="top" class="vertical-header">
    <v-btn
      class="hidden-md-and-down"
      icon
      color="primary"
      variant="text"
      @click.stop="customizer.SET_MINI_SIDEBAR(!customizer.mini_sidebar)"
    >
      <IconMenu2 :size="25" />
    </v-btn>
    <v-btn class="hidden-lg-and-up" icon variant="text" @click.stop="customizer.SET_SIDEBAR_DRAWER" size="small">
      <IconMenu2 :size="25" />
    </v-btn>
    <v-spacer />
    <LanguageDD />
    <TenantDD />
    <div class="ml-2">
      <ProfileDD />
    </div>
  </v-app-bar>
</template>
