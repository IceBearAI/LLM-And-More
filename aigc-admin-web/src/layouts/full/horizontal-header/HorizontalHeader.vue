<script setup lang="ts">
import { ref, watch, computed } from "vue";
import { useCustomizerStore } from "../../../stores/customizer";
// Icon Imports
import { IconMenu2, IconShoppingCart } from "@tabler/icons-vue";

import Logo from "../logo/Logo.vue";
// dropdown imports
import LanguageDD from "../vertical-header/LanguageDD.vue";
import ProfileDD from "../vertical-header/ProfileDD.vue";

const customizer = useCustomizerStore();
const showSearch = ref(false);
const drawer = ref(false);
const appsdrawer = ref(false);
const priority = ref(customizer.setHorizontalLayout ? 0 : 0);
function searchbox() {
  showSearch.value = !showSearch.value;
}
watch(priority, newPriority => {
  // yes, console.log() is a side effect
  priority.value = newPriority;
});
</script>

<template>
  <v-app-bar elevation="0" :priority="priority" height="64" class="horizontal-header" color="background">
    <div :class="customizer.boxed ? 'maxWidth v-toolbar__content px-lg-0 px-4' : 'v-toolbar__content px-6'">
      <div class="hidden-md-and-down">
        <Logo />
      </div>
      <v-btn class="hidden-md-and-up" icon variant="text" @click.stop="customizer.SET_SIDEBAR_DRAWER" size="small">
        <IconMenu2 :size="25" />
      </v-btn>

      <v-spacer />
      <LanguageDD />

      <div class="ml-3">
        <ProfileDD />
      </div>
    </div>
  </v-app-bar>
</template>
