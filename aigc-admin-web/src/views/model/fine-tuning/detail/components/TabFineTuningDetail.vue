<template>
  <div>
    <UiChildCard title="训练进度">
      <div class="pb-4">
        <div style="padding: 20px 60px 10px">
          <v-progress-linear v-model="process.value" :color="process.color" height="25" :striped="process.striped">
            <template v-slot:default="{ value }">
              <strong>{{ process.valueCN }}</strong>
            </template>
          </v-progress-linear>
        </div>
        <div class="d-flex points justify-space-between">
          <div class="item text-center">
            <div class="text-h6">创建训练任务</div>
            <div class="text-subtitle-1 text-medium-emphasis mt-1">
              {{ format.dateFormat(rawData.createdAt, "YYYY-MM-DD HH:mm:ss") }}
            </div>
          </div>
          <div class="item text-center">
            <div class="text-h6 hv-center">
              <component
                :is="trainStatusMap[rawData.trainStatus]?.icon"
                :size="24"
                :color="trainStatusMap[rawData.trainStatus]?.iconColor"
                class="mr-1"
              ></component>
              {{ dataDictionary.localData.local_trainStatus[rawData.trainStatus] }}
            </div>
            <div class="text-subtitle-1 text-medium-emphasis mt-1">
              {{ format.dateFormat(rawData.finishedAt, "YYYY-MM-DD HH:mm:ss") }}
            </div>
          </div>
        </div>
      </div>
    </UiChildCard>
    <UiChildCard title="训练统计" class="mt-4">
      <v-row>
        <v-col md="12" sm="12">
          <Chart1 />
        </v-col>
        <v-col md="12" sm="12">
          <Chart2 />
        </v-col>
        <v-col md="12" sm="12">
          <Chart3 />
        </v-col>
      </v-row>
    </UiChildCard>
  </div>
</template>

<script setup>
import { reactive, toRefs, ref, computed, inject, onMounted } from "vue";

import UiChildCard from "@/components/shared/UiChildCard.vue";
import { format, dataDictionary } from "@/utils";
import { trainStatusMap } from "@/views/model/fine-tuning/list/map";
import Chart1 from "./Chart1.vue";
import Chart2 from "./Chart2.vue";
import Chart3 from "./Chart3.vue";

import * as echarts from "echarts";

const provideFineTuningDetail = inject("provideFineTuningDetail");
const state = reactive({
  process: {
    value: 0,
    icon: "",
    striped: false,
    color: "",
    iconColor: ""
  }
});

const { process } = toRefs(state);

const rawData = computed(() => {
  let ret = provideFineTuningDetail.rawData;

  let { process } = state;
  process.value = ret.process * 100;
  process.valueCN = format.toPercent(ret.process, 2);
  let { trainStatus } = ret;
  if (["running"].includes(trainStatus)) {
    process.striped = true;
  }
  if (trainStatus == "cancel") {
    process.color = "#ccc";
  } else if (trainStatus == "failed") {
    process.color = "rgb(var(--v-theme-error))";
  } else {
    process.color = "rgb(var(--v-theme-info))";
  }

  return ret;
});
</script>
<style lang="scss" scoped>
.points {
  .item {
    min-width: 128px;
  }
}
</style>
