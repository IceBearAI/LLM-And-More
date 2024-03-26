<template>
  <v-btn v-bind="compuStyle" :id="props.id"><slot /></v-btn>
</template>
<script setup lang="ts">
import { reactive, computed } from "vue";
import { useAppStore } from "@/stores";

interface Props {
  loadingStatus?: "loading" | "disabled";
  id?: string;
}

const props = withDefaults(defineProps<Props>(), {
  loadingStatus: "loading",
  id: ""
});

const appStore = useAppStore();
const state = reactive({
  loading: false
});

const isMacted = () => {
  let { btnId } = appStore;
  if (btnId) {
    if (btnId == props.id) {
      return true;
    } else {
      return false;
    }
  }
  return true;
};

const compuStyle = computed(() => {
  let ret = {};
  let { loadingStatus } = props;
  if (appStore.isBtnLoading) {
    if (isMacted()) {
      if (loadingStatus == "loading") {
        return { loading: true };
      } else if (loadingStatus == "disabled") {
        return { disabled: true };
      }
    }
  }
  return ret;
});
</script>
<style lang="scss"></style>
