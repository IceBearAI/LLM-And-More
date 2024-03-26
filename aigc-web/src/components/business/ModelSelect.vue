<template>
  <v-select
    v-model="model"
    v-bind="$attrs"
    placeholder="请选择模型"
    :items="items"
    :item-title="MODEL_KEY"
    :item-value="MODEL_KEY"
    variant="outlined"
    :return-object="props.returnObject"
  >
    <template #prepend v-if="$slots.prepend">
      <slot name="prepend"></slot>
    </template>
  </v-select>
</template>
<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useVModel } from "@vueuse/core";
import { http } from "@/utils";

interface IProps {
  modelValue?: any;
  defaultFirstValue?: boolean;
  returnObject?: boolean;
}

interface IEmits {
  (e: "update:modelValue", val: string): void;
}
const MODEL_KEY = "modelName";

const props = withDefaults(defineProps<IProps>(), {
  modelValue: null,
  defaultFirstValue: false,
  returnObject: false
});

const emit = defineEmits<IEmits>();

const items = ref([]);

const model = useVModel(props, "modelValue", emit);

const getItems = async () => {
  const [err, res] = await http.get({
    url: "/channels/models",
    data: {
      modelType: "text-generation"
    }
  });
  if (res) {
    items.value = res.list;
    setDefaultValue();
  }
};

const setDefaultValue = () => {
  const { modelValue, defaultFirstValue, returnObject } = props;

  if (modelValue) {
    if (returnObject) {
      model.value = items.value.find(_ => _[MODEL_KEY] === modelValue.modelName) || modelValue;
    }
  } else {
    if (defaultFirstValue) {
      if (returnObject) {
        model.value = items.value[0];
      } else {
        model.value = items.value[0][MODEL_KEY];
      }
    }
  }
};

onMounted(() => {
  getItems();
});
</script>
