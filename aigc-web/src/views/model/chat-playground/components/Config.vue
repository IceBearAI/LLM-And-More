<template>
  <v-sheet>
    <v-label class="mb-2 font-weight-medium">模型</v-label>
    <ModelSelect
      class="mb-2"
      v-model="configModel"
      default-first-value
      return-object
      @update:modelValue="modelUpdate"
      hide-details
    />
    <v-label class="mb-2 font-weight-medium">包含过去的消息</v-label>
    <v-slider v-model="sendHistoryCount" color="primary" :max="50" :min="0" step="1" hide-details thumb-label>
      <template v-slot:append>
        <v-text-field
          v-model.number="sendHistoryCount"
          hide-details
          single-line
          density="compact"
          type="number"
          :max="50"
          :min="0"
          style="width: 80px"
        ></v-text-field>
      </template>
    </v-slider>
    <template v-for="item in paramsConfig">
      <v-label class="mb-2 font-weight-medium">{{ item.title }}</v-label>
      <v-slider
        v-model="config[item.key]"
        color="primary"
        :max="item.max"
        :min="item.min"
        :step="item.step"
        hide-details
        thumb-label
      >
        <template v-slot:append>
          <v-text-field
            v-model.number="config[item.key]"
            hide-details
            single-line
            density="compact"
            type="number"
            :max="item.max"
            :min="item.min"
            :step="item.step"
            style="width: 80px"
          ></v-text-field>
        </template>
      </v-slider>
    </template>
  </v-sheet>
</template>
<script setup lang="ts">
import { computed, reactive, ref } from "vue";
import ModelSelect from "@/components/business/ModelSelect.vue";
import { useRoute } from "vue-router";
import { http } from "@/utils";

interface IEmits {
  (e: "update:model:info", val: string): void;
}

const emit = defineEmits<IEmits>();

const route = useRoute();

const { modelName } = route.query;

const sendHistoryCount = ref(20);
const configModel = ref<Record<string, any> | null>(modelName ? { modelName } : null);
const config = reactive({
  maxTokens: 512,
  temperature: 0,
  topP: 0,
  frequencyPenalty: 0,
  presencePenalty: 0
});

const paramsConfig = computed(() => {
  const maxTokens = configModel.value?.maxTokens ?? 0;
  return [
    {
      key: "maxTokens",
      title: "最大响应数",
      max: maxTokens,
      min: 0,
      step: 1
    },
    {
      key: "temperature",
      title: "温度",
      max: 1,
      min: 0,
      step: 0.1
    },
    {
      key: "topP",
      title: "TopP",
      max: 1,
      min: 0,
      step: 0.1
    },
    {
      key: "frequencyPenalty",
      title: "频率惩罚",
      max: 1,
      min: 0,
      step: 0.1
    },
    {
      key: "presencePenalty",
      title: "重复惩罚",
      max: 1,
      min: 0,
      step: 0.1
    }
  ];
});

const modelUpdate = async val => {
  const [err, res] = await http.get({
    url: `/api/models/${val.modelName}/info`
  });
  if (res) {
    emit("update:model:info", res);
    const generateConfig = res.generateConfig;
    config.maxTokens = config.maxTokens > res.maxTokens ? res.maxTokens : 512;
    paramsConfig.value.forEach(item => {
      if (item.key !== "maxTokens") {
        config[item.key] = generateConfig?.[item.key] || 0;
      }
    });
  }
};

const getData = () => {
  return {
    sendHistoryCount: sendHistoryCount.value,
    data: {
      model: configModel.value?.modelName,
      ...config
    }
  };
};

defineExpose({
  getData
});
</script>
