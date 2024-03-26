<template>
  <div ref="refChart" class="chart-item w-100 h-100"></div>
</template>
<script setup>
import { ref, onBeforeUnmount } from "vue";
import * as echarts from "echarts";
import { useResizeObserver } from "@vueuse/core";

const refChart = ref();
var chart = null;

const initChart = ({ title, radar, seriesData }) => {
  chart = echarts.init(refChart.value);
  const legendData = seriesData.map(item => item.name);
  chart.setOption({
    title: {
      text: title
    },
    legend: {
      icon: "circle",
      data: legendData
    },
    radar,
    series: [
      {
        type: "radar",
        data: seriesData
      }
    ]
  });
};

onBeforeUnmount(() => {
  if (!chart) return;
  chart.dispose();
  chart = null;
});

useResizeObserver(refChart, entries => {
  chart?.resize();
});

defineExpose({
  initChart
});
</script>
<style lang="scss" scoped>
.chart-item {
  aspect-ratio: 3 /1;
}
</style>
