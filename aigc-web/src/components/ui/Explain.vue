<template>
  <span class="compo-explain">
    <el-tooltip effect="dark" placement="top">
      <template #content> <slot /></template>
      <IconInfoCircle :size="16" color="#bbb" pointer tabindex="-1"> </IconInfoCircle>
    </el-tooltip>
  </span>
</template>
<script setup lang="ts">
import { IconInfoCircle } from "@tabler/icons-vue";
</script>
<style lang="scss">
.compo-explain {
  display: inline-flex;
  vertical-align: text-top;
  svg {
    outline: none;
  }
}
</style>
