<template>
  <div :class="textClass">
    <span>{{ text }}</span>
    <div v-if="['running', 'success'].includes(props.item.trainStatus)" class="link ml-1" @click="openTrainLog">(日志)</div>
    <Explain v-if="props.item.trainStatus === 'failed' && props.item.errorMessage" class="ml-1">
      <div class="tooltip-error-message">
        {{ props.item.errorMessage }}
      </div>
    </Explain>
  </div>
</template>
<script setup lang="ts">
import { computed, ref } from "vue";
import Explain from "@/components/ui/Explain.vue";
import { trainStatusMap } from "../map";

interface IProps {
  item?: Record<string, any>;
}
interface IEmits {
  (e: "open:log"): void;
}
const props = withDefaults(defineProps<IProps>(), {
  item: () => ({})
});
const emit = defineEmits<IEmits>();

const text = computed(() => {
  return trainStatusMap[props.item.trainStatus]?.text;
});
const textClass = computed(() => {
  return [trainStatusMap[props.item.trainStatus]?.textColor, "d-flex", "justify-center", "align-center"];
});

const openTrainLog = () => {
  emit("open:log");
};
</script>
<style lang="scss" scoped>
.tooltip-error-message {
  max-width: 1200px;
}
</style>
