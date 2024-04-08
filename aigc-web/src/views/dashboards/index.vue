<script setup lang="ts">
import { ref } from "vue";

import Welcome from "@/components/dashboards/Welcome.vue";

const page = ref({ title: "Analytical" });
const breadcrumbs = ref([
  {
    text: "Dashboard",
    disabled: false,
    href: "#"
  },
  {
    text: "Analytical",
    disabled: true,
    href: "#"
  }
]);
</script>

<template>
  <v-row>
    <Welcome />
  </v-row>
</template>
