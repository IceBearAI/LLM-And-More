<template>
  <v-avatar :size="props.size" class="compo-iconChecked position-absolute flag bg-surface animate__bounceIn">
    <IconCircleCheckFilled color="#67C23A" :size="props.size"></IconCircleCheckFilled>
  </v-avatar>
</template>
<script setup>
import { reactive, toRefs, ref, computed } from "vue";
import { IconCircleCheckFilled } from "@tabler/icons-vue";
const props = defineProps({
  size: {
    type: Number,
    default: 32
  }
});

const width = computed(() => {
  return -(props.size / 2) + "px";
});
</script>
<style lang="scss">
.compo-iconChecked {
  z-index: 100;
  right: v-bind(width);
  top: v-bind(width);
}
</style>
