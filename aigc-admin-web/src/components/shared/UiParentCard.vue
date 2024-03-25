<script setup lang="ts">
const props = defineProps({
  title: String
});
</script>

// ===============================|| Ui Parent Card||=============================== //
<template>
  <v-card elevation="10">
    <template v-if="title || $slots.header">
      <v-card-item class="py-4 px-6">
        <div class="d-sm-flex align-center justify-space-between">
          <v-card-title class="text-h5">
            <template v-if="$slots.header">
              <slot name="header" />
            </template>
            <template v-else>{{ title }}</template>
          </v-card-title>
          <!-- <template v-slot:append> -->
          <slot name="action"></slot>
          <!-- </template> -->
        </div>
      </v-card-item>
      <v-divider></v-divider>
    </template>
    <v-card-text>
      <slot />
    </v-card-text>
  </v-card>
</template>
