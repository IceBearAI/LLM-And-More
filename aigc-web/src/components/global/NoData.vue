<template>
  <div class="compo-noData-5 hv-center">
    <div>
      <img src="@/assets/images/blank.svg" alt="" :width="props.size" />
      <div class="text">
        <slot v-if="$slots.default" />
        <span v-else>{{ t("noData") }}</span>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { reactive, toRefs, useSlots } from "vue";
import { useI18n } from "vue-i18n";
const { t } = useI18n(); // 解构出t方法

const $slots = useSlots();
const props = defineProps({
  size: {
    type: String,
    default: "110"
  }
});
</script>
<style lang="scss" scoped>
.compo-noData-5 {
  padding: 50px 0;
  line-height: 1;
  .text {
    margin-top: 10px;
    font-size: 12px;
    line-height: 2;
    color: #666;
  }
}
</style>
