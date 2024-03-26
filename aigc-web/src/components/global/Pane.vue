<template>
  <el-drawer
    v-model="state.visible"
    :show-close="false"
    :size="state.width"
    :destroy-on-close="true"
    :before-close="onClose"
    :close-on-click-modal="false"
    class="compo-pane"
    :class="'v-theme--' + customizer.actTheme"
    append-to-body
    ref="refPane"
  >
    <template #header>
      <template v-if="$slots.header">
        <slot name="header" />
      </template>
      <div v-else class="py-2 px-4 d-flex align-center justify-space-between">
        <div class="text-h6">{{ state.title }}</div>
        <v-btn @click="onClose" :icon="IconX" flat class="mr-n3" :title="$t('pane.closeTip')" width="34px" height="34px"> </v-btn>
      </div>
    </template>
    <slot />
    <template #footer v-if="state.showActions">
      <div class="py-2 px-4 d-flex align-center justify-end">
        <template v-if="$slots.txtLeft">
          <slot name="txtLeft" />
        </template>
        <template v-if="$slots.buttons">
          <slot name="buttons" />
        </template>
        <template v-else>
          <template v-if="state.hasSubmitBtn">
            <v-btn size="small" color="secondary" class="mr-3" variant="outlined" @click="onClose">{{ $t("cancel") }}</v-btn>
            <AiBtn id="btnPaneSubmit" size="small" color="primary" variant="flat" @click="onSubmit">{{
              $t(state.confirmText)
            }}</AiBtn>
          </template>
          <template v-else>
            <v-btn size="small" color="secondary" class="mr-3" variant="outlined" @click="onClose">{{ $t("close") }}</v-btn>
          </template>
        </template>
      </div>
    </template>
  </el-drawer>
</template>
<script setup>
import { reactive, toRefs, ref, nextTick } from "vue";
import { toast } from "vue3-toastify";
import { IconX } from "@tabler/icons-vue";
import { useCustomizerStore } from "@/stores/customizer";
import $ from "jquery";
import { useI18n } from "vue-i18n";
const { t } = useI18n(); // 解构出t方法
const customizer = useCustomizerStore();
const refPane = ref();

const state = reactive({
  visible: false,
  hasSubmitBtn: true,
  showActions: true,
  title: "",
  width: "",
  showError: true,
  refForm: "",
  confirmText: ""
});
const { style, formData } = toRefs(state);
const emits = defineEmits(["close", "submit"]);
const closePane = () => {
  state.visible = false;
  document.getElementsByTagName("html")[0].classList.remove("el-popup-parent--hidden");
};

const onClose = () => {
  closePane();
  emits("close");
};

const onSubmit = async () => {
  let showLoading = $(`[aria-labelledby=${refPane.value.titleId}]`)[0];
  if (state.refForm) {
    let { valid, errors } = await state.refForm.validate();
    emits("submit", { valid, errors, showLoading });
    if (!valid) {
      let { showError } = state;
      if (showError) {
        let errorMsg = t("pane.errorMsg");
        if (typeof showError == "string") {
          errorMsg = showError;
        }
        toast.warning(errorMsg);
      }
    }
  } else {
    emits("submit", { showLoading });
  }
};

defineExpose({
  /**
   *
   * @param {*} param0
   *   refForm : pane中内嵌的form
   */
  show({
    title = "",
    width = "800",
    showError = true,
    hasSubmitBtn = true,
    refForm,
    showActions = true,
    confirmText = "confirm"
  }) {
    state.title = title;
    state.visible = true;
    state.width = width;
    state.refForm = refForm;
    state.hasSubmitBtn = hasSubmitBtn;
    state.showActions = showActions;
    state.showError = showError;
    state.confirmText = confirmText;
    //防止body滚动
    document.getElementsByTagName("html")[0].classList.add("el-popup-parent--hidden");

    nextTick(() => {
      state.refForm?.resetValidation();
    });
  },
  hide() {
    closePane();
  }
});
</script>
<style lang="scss">
.compo-pane {
  .el-drawer__header,
  .el-drawer__footer {
    padding: 0;
    margin-bottom: 0;
  }

  .el-drawer__header {
    height: 44px;
    border-bottom: 1px solid #e8e8e8;
  }

  .el-drawer__footer {
    border-top: 1px solid #e8e8e8;
  }

  > .el-loading-mask {
    //锁body+footer
    top: 44px;
  }
}
</style>
