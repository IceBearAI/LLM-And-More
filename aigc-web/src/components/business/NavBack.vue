<!--导航后退条-->
<template>
  <div class="d-flex v-center">
    <el-tooltip :content="tooltip" :disabled="!tooltip">
      <v-btn icon variant="text" @click="onBack" :class="state.btnBackClassName" class="[animation-duration:1s]">
        <IconArrowLeft stroke-width="2.5" :size="24" />
      </v-btn>
    </el-tooltip>
    <h3 ref="refTitle" class="text-h3 ml-2 flex-1"><slot /></h3>
  </div>
</template>
<script setup lang="ts">
import { IconArrowLeft } from "@tabler/icons-vue";
import { useRouter } from "vue-router";
import { animate } from "@/utils";
import { onMounted, ref } from "vue";
import { reactive } from "vue";

/** 外部传入属性 */
type TypePropsOutter = {
  /** 悬浮提示文案 ，未传入时，不展示悬浮文案 */
  tooltip?: string;
  /**
   * 后退页面路径
   * 1. 未传入时，调用 router.back() 方法
   * 2. 传入 'backCallback'，抛出 back 事件
   * 3. 其他字符串，调用 router.push 方法
   */
  backUrl?: string;
};

const props = withDefaults(defineProps<TypePropsOutter>(), {
  tooltip: "",
  backUrl: ""
});

const emits = defineEmits<{
  /** 抛出后退事件 */
  (e: "back"): void;
}>();

const router = useRouter();
const refTitle = ref<HTMLElement>();

const state = reactive({
  btnBackClassName: "opacity-0"
});

const onBack = () => {
  let { backUrl } = props;
  if (backUrl) {
    if (backUrl == "backCallback") {
      emits("back");
    } else {
      router.replace(backUrl);
    }
  } else {
    //backUrl未传入
    router.back();
  }
};

onMounted(() => {
  let { finished } = animate(
    refTitle.value,
    [{ transform: "translateX(-50px)", opacity: 0.1 }, { opacity: 0.2 }, { transform: "translateX(0)" }],
    {
      duration: 300,
      easing: "ease-in"
    }
  );
  finished.then(() => {
    state.btnBackClassName = "opacity-1 animate__bounceIn";
  });
});
</script>
<style lang="scss"></style>
