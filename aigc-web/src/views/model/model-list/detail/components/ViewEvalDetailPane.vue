<template>
  <Pane ref="refPane">
    <v-row>
      <v-col :cols="isShowRemind ? 7 : 12">
        <div style="height: 500px">
          <GradeRadarChart ref="gradeRadarChartRef" />
        </div>
      </v-col>
      <v-col v-if="isShowRemind" cols="5">
        <v-alert v-if="data['current'].riskOver" :type="data['current'].riskOver ? 'error' : 'success'" variant="tonal"
          >过拟合风险</v-alert
        >
        <v-alert
          v-if="data['current'].riskUnder"
          :type="data['current'].riskUnder ? 'error' : 'success'"
          variant="tonal"
          class="mt-4"
          >欠拟合风险</v-alert
        >
        <v-alert type="info" variant="tonal" class="mt-4">
          <h5 class="text-h6 text-capitalize">建议进一步操作</h5>
          <div>{{ data["current"].remind }}</div>
        </v-alert>
      </v-col>
      <v-col cols="12">
        <v-row>
          <v-col v-for="item in modelMap" cols="4">
            <Select
              placeholder="请选择模型"
              v-model="searchData[item.value]"
              :disabled="item.disabled"
              :clearable="false"
              :mapAPI="{
                url: '/channels/models',
                data: { pageSize: -1, providerName: 'LocalAI', modelType: 'text-generation', evalTag: 'five' },
                labelField: 'modelName',
                valueField: 'id'
              }"
              @change="getData(item.disabled)"
            ></Select>
            <div class="text-center mt-1">{{ data[item.key]?.score }}</div>
          </v-col>
        </v-row>
      </v-col>
    </v-row>
  </Pane>
</template>
<script setup lang="ts">
import { ref, reactive, computed } from "vue";
import GradeRadarChart from "@/components/business/charts/GradeRadarChart.vue";
import { http } from "@/utils";

const refPane = ref();
const gradeRadarChartRef = ref();
const searchData = reactive({
  currentModelId: null,
  compare1ModelId: null,
  compare2ModelId: null
});
const data = ref({});
const modelMap = [
  {
    key: "current",
    value: "currentModelId",
    disabled: true
  },
  {
    key: "compare1",
    value: "compare1ModelId",
    disabled: false
  },
  {
    key: "compare2",
    value: "compare2ModelId",
    disabled: false
  }
];

const paneConfig = reactive({
  id: ""
});

const isShowRemind = computed(() => {
  return data.value["current"] && (data.value["current"].riskOver || data.value["current"].riskUnder);
});

const renderRadarChart = () => {
  const radar = {
    indicator: [
      { name: "中文能力", max: 10, axisLabel: { show: true } },
      { name: "推理能力", max: 10 },
      { name: "指令遵从能力", max: 10 },
      { name: "创新能力", max: 10 },
      { name: "阅读理解", max: 10 }
    ]
  };

  const seriesData = [];
  if (data.value) {
    modelMap.forEach(item => {
      if (data.value[item.key].value) {
        seriesData.push({
          name: data.value[item.key].name,
          value: data.value[item.key].value
        });
      }
    });
  }
  gradeRadarChartRef.value.initChart({
    title: "",
    radar,
    seriesData
  });
};

const getData = async (isDisabled = false) => {
  if (isDisabled) return;
  const [err, res] = await http.post({
    url: `/evaluate/fivegraph`,
    showLoading: refPane.value.el,
    data: {
      currentModelId: searchData.currentModelId,
      compare1ModelId: searchData.compare1ModelId,
      compare2ModelId: searchData.compare2ModelId,
      currentModelEvaluateId: paneConfig.id
    }
  });
  if (res) {
    data.value = res;
    renderRadarChart();
  }
};

defineExpose({
  show({ title, modelId, id }) {
    refPane.value.show({
      width: 900,
      title,
      hasSubmitBtn: false
    });
    searchData.currentModelId = parseInt(modelId);
    searchData.compare1ModelId = null;
    searchData.compare2ModelId = null;
    paneConfig.id = id;
    getData();
  }
});
</script>
<style lang="scss" scoped>
label {
  width: 130px;
  text-align: right;
  .compo-explain {
    margin-top: 1px;
  }
}
</style>
