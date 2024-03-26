<template>
  <v-file-input
    ref="refFileInput"
    v-bind="$attrs"
    density="compact"
    variant="outlined"
    hide-details="auto"
    @update:modelValue="change"
    :disabled="showLoading && uploading"
  >
    <template v-slot:append-inner>
      <v-progress-circular
        v-if="showLoading && uploading"
        indeterminate
        color="primary"
        :size="20"
        :width="2"
      ></v-progress-circular>
    </template>
  </v-file-input>
</template>
<script setup lang="ts">
import { ref } from "vue";
import { http } from "@/utils";
interface IProps {
  modelValue?: string;
  infos?: Record<string, any>;
  purpose?: string;
  loading?: boolean;
  showLoading?: boolean;
}

interface IEmits {
  (e: "update:modelValue", val: string): void;
  (e: "update:infos", val: Record<string, any>): void;
  (e: "upload:success", val: string): void;
  (e: "loading", val: boolean): void;
}

const props = withDefaults(defineProps<IProps>(), {
  modelValue: "",
  infos: () => ({}),
  purpose: "",
  loading: false,
  showLoading: true
});
const emits = defineEmits<IEmits>();

const refFileInput = ref();
const uploading = ref(false);

const change = async val => {
  if (val.length === 0) return;
  const data = {
    file: val
  };
  if (props.purpose) {
    data["purpose"] = props.purpose;
  }
  uploading.value = true;
  emits("loading", true);
  const [err, res] = await http.upload({
    url: "/files",
    data
  });
  if (res) {
    emits("update:modelValue", res.fileId);
    emits("update:infos", res);
    emits("upload:success", res);
  }
  uploading.value = false;
  emits("loading", false);
};
</script>
