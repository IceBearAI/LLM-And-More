<template>
  <v-row class="my-form waterfall">
    <v-col xs="12" sm="6" md="4" lg="3">
      <v-input hide-details>
        <template #prepend> <label>供应</label></template>
        {{ getLabels([["speak_provider", props.info.provider]]) }}
      </v-input>
    </v-col>
    <v-col xs="12" sm="6" md="4" lg="3">
      <v-input hide-details>
        <template #prepend> <label>语言</label></template>
        {{ getLabels([["speak_lang", props.info.lang]]) }}
      </v-input>
    </v-col>
    <v-col xs="12" sm="6" md="4" lg="3">
      <v-input hide-details>
        <template #prepend> <label>音色</label></template>
        {{
          getLabels(
            [
              ["speak_age_group", props.info.ageGroup],
              ["speak_gender", props.info.gender]
            ],
            ret => {
              if (ret.length) {
                return ret.join("") + "声";
              } else {
                return "未知";
              }
            }
          )
        }}
      </v-input>
    </v-col>
    <v-col xs="12" sm="6" md="4" lg="3">
      <v-input hide-details>
        <template #prepend> <label>适应范围</label></template>
        {{ getLabels([["speak_area", props.info.area]]) }}
      </v-input>
    </v-col>
    <v-col xs="12" sm="6" md="4" lg="3">
      <v-input hide-details>
        <template #prepend> <label>年龄段</label></template>
        {{ getLabels([["speak_age_group", props.info.ageGroup]]) }}
      </v-input>
    </v-col>
    <v-col xs="12" sm="6" md="4" lg="3">
      <v-input hide-details>
        <template #prepend> <label>性别</label></template>
        {{ getLabels([["speak_gender", props.info.gender]]) }}
      </v-input>
    </v-col>
    <v-col xs="12" sm="6" md="4" lg="3">
      <v-input hide-details>
        <template #prepend> <label>备注</label></template>
        {{ props.info.remark }}
      </v-input>
    </v-col>
    <v-col xs="12" sm="6" md="4" lg="3">
      <v-input hide-details>
        <template #prepend> <label>创建时间</label></template>
        {{ format.dateFormat(props.info.createdAt, "YYYY-MM-DD HH:mm:ss") }}
      </v-input>
    </v-col>
  </v-row>
</template>
<script setup lang="ts">
import { format } from "@/utils";
import { useMapRemoteStore } from "@/stores";
const { loadDictTree, getLabels } = useMapRemoteStore();

loadDictTree(["speak_provider", "speak_lang", "speak_age_group", "speak_gender", "speak_area"]);

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
