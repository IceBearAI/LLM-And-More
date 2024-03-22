<template>
  <v-card>
    <!-- <div class="d-flex"><IconZoomIn /></div> -->
    <div ref="refBox" class="chart-item w-100"></div>
  </v-card>
</template>
<script setup>
import { reactive, toRefs, ref, watchEffect, inject } from "vue";
import * as echarts from "echarts";
// 导入所需的语言包
import ZH from "echarts/lib/i18n/langZH.js";

import { useResizeObserver } from "@vueuse/core";
const state = reactive({
  style: {},
  formData: {}
});
const refBox = ref();

let myChart;

const provideFineTuningDetail = inject("provideFineTuningDetail");

const onChart = () => {
  let { loss } = provideFineTuningDetail.rawData?.trainAnalysis || {};
  let xData = [];
  let seriesData = [];
  loss.list.forEach((item, index) => {
    xData.push(index);
    seriesData.push(item.value);
  });

  // 注册语言包
  echarts.registerLocale("ZH", ZH);

  // 基于准备好的dom，初始化echarts实例
  myChart = echarts.init(refBox.value, null, {
    locale: "ZH"
  });
  // 绘制图表
  myChart.setOption({
    title: {
      text: "train/loss",
      x: "center"
    },
    toolbox: {
      feature: {
        saveAsImage: {}, // 保存为图片按钮
        dataView: {}, // 数据视图按钮
        dataZoom: {} // 数据缩放按钮
        // restore: {} // 还原按钮
      }
    },
    tooltip: {
      trigger: "axis",
      axisPointer: {
        type: "cross"
      }
    },
    xAxis: {
      type: "category",
      data: xData,
      axisLine: {
        onZero: false
      },
      axisLabel: {
        //X轴 不显示 0
        formatter(value) {
          if (value == 0) {
            return "";
          } else {
            return value;
          }
        },
        //显示10个
        interval: Math.floor((xData.length - 1) / 10)
      }
    },
    yAxis: {
      type: "value",
      axisLine: { onZero: false }
    },
    series: [
      {
        data: seriesData,
        type: "line",
        smooth: true
      }
    ]
  });
};

watchEffect(() => {
  let { loss } = provideFineTuningDetail.rawData?.trainAnalysis || {};
  if (loss) {
    onChart();
  }
});

useResizeObserver(refBox, entries => {
  myChart?.resize();
});
</script>
<style lang="scss" scoped>
.chart-item {
  aspect-ratio: 3 /1;
}
</style>
