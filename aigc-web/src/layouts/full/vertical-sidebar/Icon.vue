<script setup>
const props = defineProps({ item: Object, level: Number });
</script>

<template>
  <template v-if="level > 0">
    <component :is="item" size="14" stroke-width="1.5" class="iconClass"></component>
  </template>
  <template v-else>
    <component :is="item" size="20" stroke-width="1.5" class="iconClass"></component>
  </template>
</template>
