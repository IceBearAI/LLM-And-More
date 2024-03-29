<template>
  <v-row class="my-form waterfall">
    <v-col xs="12" sm="6" md="4" lg="3">
      <v-input hide-details>
        <template #prepend> <label>创建人</label></template>
        {{ props.info.operator }}
      </v-input>
    </v-col>
    <v-col xs="12" sm="6" md="4" lg="3">
      <v-input hide-details>
        <template #prepend> <label>模型</label></template>
        {{ props.info.modelName }}
      </v-input>
    </v-col>
    <v-col xs="12" sm="6" md="4" lg="3">
      <v-input hide-details>
        <template #prepend> <label>创建时间</label></template>
        {{ format.dateFormat(props.info.createdAt, "YYYY-MM-DD HH:mm:ss") }}
      </v-input>
    </v-col>
    <v-col xs="12" sm="6" md="4" lg="3">
      <v-input hide-details>
        <template #prepend> <label>更新时间</label></template>
        {{ format.dateFormat(props.info.updatedAt, "YYYY-MM-DD HH:mm:ss") }}
      </v-input>
    </v-col>
    <v-col xs="12" sm="6" md="4" lg="3">
      <v-input hide-details>
        <template #prepend> <label>备注</label></template>
        {{ props.info.remark }}
      </v-input>
    </v-col>
  </v-row>
</template>
<script setup lang="ts">
import { format } from "@/utils";

interface IProps {
  info: Record<string, any>;
}

const props = withDefaults(defineProps<IProps>(), {
  info: () => ({})
});
</script>
<style lang="scss" scoped>
.my-form > * {
  margin-bottom: 0;
}
</style>
