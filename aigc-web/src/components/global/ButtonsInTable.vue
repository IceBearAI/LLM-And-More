<!---->
<template>
  <div v-if="buttons.length" class="compo-buttonsInList d-flex">
    <!--第一个按钮-->

    <template v-if="buttons.length == 1">
      <!--只有一个按钮-->
      <div class="flex-1 hv-center" v-if="onlyOne">
        <div :class="getClassName(buttons[0])" @click="buttons[0].click">
          {{ buttons[0].text }}
        </div>
      </div>
      <template v-else>
        <div class="flex-1 hv-center">
          <div :class="getClassName(buttons[0])" @click="buttons[0].click">
            {{ buttons[0].text }}
          </div>
        </div>
        <v-divider vertical class="mx-2 position-relative opacity-0"></v-divider>
        <div class="flex-1 opacity-0"></div>
      </template>
    </template>
    <template v-else-if="buttons.length == 2">
      <!--有两个按钮-->
      <div class="flex-1 hv-center">
        <div :class="getClassName(buttons[0])" @click="buttons[0].click">
          {{ buttons[0].text }}
        </div>
      </div>
      <v-divider vertical class="mx-2 position-relative"></v-divider>
      <div class="flex-1 hv-center">
        <div :class="getClassName(buttons[1])" @click="buttons[1].click">
          {{ buttons[1].text }}
        </div>
      </div>
    </template>
    <template v-else-if="buttons.length > 2">
      <!-- 有多个按钮-->
      <div class="flex-1 hv-center">
        <div :class="getClassName(buttons[0])" @click="buttons[0].click">
          {{ buttons[0].text }}
        </div>
      </div>
      <v-divider vertical class="mx-2 position-relative"></v-divider>
      <div class="flex-1 hv-center">
        <div class="btn-more" :class="getClassName()">
          {{ $t("more") }}...
          <v-menu activator="parent" width="100px" open-on-hover>
            <v-list density="compact">
              <template v-for="(itemBtn, indexBtn) of buttons">
                <v-list-item
                  v-if="indexBtn > 0"
                  :value="indexBtn"
                  class="compo-buttonsInList-list-item"
                  :disabled="itemBtn.color == 'disabled'"
                >
                  <v-list-item-title @click="itemBtn.click">
                    <div :class="getClassName(itemBtn)">
                      {{ itemBtn.text }}
                    </div>
                  </v-list-item-title>
                </v-list-item>
              </template>
            </v-list>
          </v-menu>
        </div>
      </div>
    </template>
  </div>
</template>
<script setup lang="ts">
import { reactive, toRefs, ref } from "vue";
import { TypeButtonsInTable } from "@/components/types/components.ts";
import { useI18n } from "vue-i18n";

const { t } = useI18n(); // 解构出t方法

const props = withDefaults(
  defineProps<{
    buttons: Array<TypeButtonsInTable>;
    /** 只有一个按钮 */
    onlyOne?: boolean;
  }>(),
  {
    buttons: null,
    onlyOne: false
  }
);

const getClassName = (info: TypeButtonsInTable = {} as TypeButtonsInTable) => {
  let { color } = info;
  let textClassName = "";
  if (color == "disabled") {
    //不可用状态
    textClassName = `text-gray-400 !cursor-not-allowed`;
  } else if (color) {
    textClassName = `text-${color}`;
  } else if (!color) {
    textClassName = `text-info`;
  }

  return `btn link-hover text-13 font-weight-medium ${textClassName}`;
};
</script>
<style lang="scss">
.compo-buttonsInList,
.compo-buttonsInList-list-item {
  font-size: 13px;
  .text-error {
    color: rgb(249, 89, 80) !important;
  }
}

.compo-buttonsInList {
  .btn {
    cursor: pointer;
    :not(.btn-more) {
      &:hover {
        opacity: 0.6;
      }
      &:active {
        opacity: 1;
        filter: brightness(0.75);
      }
    }
  }
}
</style>
