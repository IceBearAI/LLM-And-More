<template>
  <v-row class="text-center" v-if="data">
    <template v-for="(item, index) in config">
      <v-col :cols="12 / config.length" :class="{ 'border-right': index !== config.length - 1 }">
        <v-btn color="inherit" icon :class="'pa-0 text-' + item.color + ' bg-' + item.bgColor">
          <component :is="item.icon" stroke-width="1.5" :size="24" />
        </v-btn>
        <h4 class="text-h4 mt-3">{{ item.statusText }}</h4>
        <p class="text-subtitle-1 font-weight-medium text-medium-emphasis mt-1">{{ `${data[item.key]}${item.valueText}` }}</p>
      </v-col>
    </template>
  </v-row>
</template>
<script setup lang="ts">
import { onMounted, ref } from "vue";
import { http } from "@/utils";
import type { FunctionalComponent } from "vue";

interface IConfigItem {
  statusText: string;
  valueText: string;
  key: string;
  color: string;
  bgColor: string;
  icon: FunctionalComponent;
}

interface IProps {
  config?: IConfigItem[];
  requestUrl?: string;
}

const props = withDefaults(defineProps<IProps>(), {
  config: () => [],
  requestUrl: ""
});

let data = ref(null);

const getData = async () => {
  const [err, res] = await http.get({
    url: props.requestUrl
  });
  if (res) {
    data.value = res;
  }
};

onMounted(() => {
  getData();
});
</script>
