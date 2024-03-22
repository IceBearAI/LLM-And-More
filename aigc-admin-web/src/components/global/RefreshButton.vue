<template>
  <v-btn :color="color" @click="handleClick" v-bind="$attrs"
    >{{ text }}<v-progress-circular :model-value="percent" size="14" width="2" class="ml-1"></v-progress-circular
  ></v-btn>
</template>
<script setup lang="ts">
import { ref, onUnmounted } from "vue";
import { useIntervalFn } from "@vueuse/core";

interface IProps {
  text?: string;
  color?: string;
  interval?: number;
}

interface IEmits {
  (e: "refresh"): void;
}

const props = withDefaults(defineProps<IProps>(), {
  text: "刷新",
  color: "primary",
  interval: 600 // 100/10s, 300/30s, 600/1min
});

const emits = defineEmits<IEmits>();

const percent = ref(0);
const intervalFn = ref(null);

const handleClick = () => {
  if (intervalFn.value) {
    intervalFn.value.pause();
  }
  emits("refresh");
};

const start = () => {
  percent.value = 0;
  if (!intervalFn.value) {
    // 初始化 intervalFn
    intervalFn.value = useIntervalFn(() => {
      percent.value += 1;
      if (percent.value === 100) {
        intervalFn.value.pause();
        emits("refresh");
      }
    }, props.interval);
  } else {
    if (!intervalFn.value.isActive) {
      intervalFn.value.resume();
    }
  }
};

onUnmounted(() => {
  if (intervalFn.value) {
    intervalFn.value.pause(); // 没有立即执行 useIntervalFn 需要手动销毁定时器
  }
});

defineExpose({
  start
});
</script>
