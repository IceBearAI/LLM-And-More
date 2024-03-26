<script setup lang="ts">
const props = defineProps({
  title: String
});
</script>

<template>
  <v-card variant="outlined">
    <v-card-item class="py-4 px-6" v-if="$slots.header || title">
      <v-card-title class="text-h5">
        <template v-if="$slots.header">
          <slot name="header" />
        </template>
        <template v-else>{{ title }}</template>
      </v-card-title>
    </v-card-item>
    <v-divider />
    <v-card-text class="px-6">
      <slot />
    </v-card-text>
  </v-card>
</template>
