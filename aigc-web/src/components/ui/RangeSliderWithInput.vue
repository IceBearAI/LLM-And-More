<template>
  <v-input>
    <v-range-slider
      :disabled="min == max"
      color="primary"
      v-model="data"
      :min="min"
      :max="max"
      :step="1"
      hide-details
      class="align-center"
    >
      <template #prepend>
        <v-text-field
          v-model="data[0]"
          hide-details
          single-line
          type="number"
          variant="outlined"
          density="compact"
          class="input-text"
          @change="onCheckValid"
          @keyup="onCheckValid"
        ></v-text-field>
      </template>
      <template #append>
        <v-text-field
          v-model="data[1]"
          hide-details
          single-line
          type="number"
          variant="outlined"
          class="input-text"
          density="compact"
          @change="onCheckValid"
          @keyup="onCheckValid"
        ></v-text-field>
      </template>
    </v-range-slider>
    <template #prepend>
      <slot name="prepend" />
    </template>
  </v-input>
</template>
<script setup lang="ts">
import { reactive, toRefs, ref, watch } from "vue";
import { useVModel } from "@vueuse/core";

interface Props {
  /** 所选项 value 的值 */
  modelValue: Array<number>;
  /** 最小值 */
  min: number;
  /** 最大值*/
  max: number;
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: () => [0, 0],
  min: 0,
  max: 0
});

interface IEmits {
  (e: "update:modelValue", value: any): void;
}

const emits = defineEmits<IEmits>();

let data = useVModel(props, "modelValue", emits);

const onCheckValid = () => {
  let [currentMin, currentMax] = data.value;
  let { min, max } = props;
  let validMin = Math.min(currentMin, currentMax);
  let validMax = Math.max(currentMin, currentMax);

  if (validMin < min) {
    validMin = min;
  } else if (validMin > max) {
    validMin = max;
  }

  if (validMax < min) {
    validMax = min;
  } else if (validMax > max) {
    validMax = max;
  }
  data.value = [validMin, validMax];
};
</script>
<style lang="scss" scoped>
:deep() {
  // .v-input__details {
  //   display: none;
  // }
  .v-input {
    margin: 0 !important;
  }
  input {
    padding: 10px 6px;
    text-align: center;
  }
}
.input-text {
  width: 75px;
}
</style>
