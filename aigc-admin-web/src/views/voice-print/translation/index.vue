<template>
  <BaseBreadcrumb :title="page.title" :breadcrumbs="breadcrumbs"></BaseBreadcrumb>

  <UiParentCard>
    <v-tabs v-model="state.tabIndex" color="primary" @update:modelValue="onTab">
      <v-tab value="tabVoice" class="text-h5">实时语音转换</v-tab>
      <v-tab value="tabFile" class="text-h5">上传音频文件</v-tab>
    </v-tabs>

    <v-window v-model="state.tabIndex" class="py-6">
      <v-window-item value="tabVoice">
        <div style="min-height: 300px">
          <TabVoice ref="refTabVoice" @voiceStart="onVoiceStart" @voiceEnd="onVoiceEnd" @translate="onAddResult" />
        </div>
      </v-window-item>
      <v-window-item value="tabFile">
        <div style="min-height: 300px">
          <TabFile ref="refTabFile" @translate="onReplaceResult" />
        </div>
      </v-window-item>
    </v-window>

    <v-row>
      <v-col cols="12">
        <UiParentCard title="转换结果">
          <v-row ref="refResult" class="text-body-1 pa-5 scrollbar-auto shower-result">
            <p v-if="showResult" v-text="state.result" style="white-space: pre-wrap"></p>
            <p v-else class="flex-fill text-center">暂无结果</p>
          </v-row>
        </UiParentCard>
      </v-col>
    </v-row>
  </UiParentCard>
</template>
<script setup lang="ts">
import { ref, computed, reactive, nextTick, onUnmounted, onBeforeUnmount } from "vue";
import BaseBreadcrumb from "@/components/shared/BaseBreadcrumb.vue";
import UiParentCard from "@/components/shared/UiParentCard.vue";
import TabFile from "./components/TabFile.vue";
import TabVoice from "./components/TabVoice.vue";
import { animate } from "@/utils";

const refResult = ref();
const refTabVoice = ref();
const refTabFile = ref();

const state = reactive({
  tabIndex: "tabVoice",
  // tabIndex: "tabFile",
  result: "",
  statusVoice: ""
});

const page = ref({ title: "声音转文本" });
const breadcrumbs = ref([
  {
    text: "语音服务",
    disabled: false,
    href: "#"
  },
  {
    text: "声音转文本",
    disabled: true,
    href: "#"
  }
]);

const onAddResult = text => {
  console.log("onVoice", text);
  state.result += text;
};

const onReplaceResult = text => {
  if (text == state.result) {
    //文本无变化
    return;
  }
  let { finished } = animate(
    refResult.value.$el,
    [
      { opacity: 0.7, transform: "scale(1.02)" },
      { opacity: 0, transform: "scale(1.05)" }
    ],
    {
      duration: 300,
      fill: "forwards"
    }
  );
  finished.then(() => {
    state.result = text || "";
    nextTick(() => {
      animate(
        refResult.value.$el,
        [
          {
            opacity: 0.2,
            transform: "none"
            // , transform: "scale(0.97)"
          },
          {
            opacity: 1,
            transform: "none"
            // , transform: "none"
          }
        ],
        {
          duration: 600,
          fill: "forwards"
        }
      );
    });
  });
};

const showResult = computed(() => {
  let { statusVoice, result } = state;
  if (statusVoice == "working") {
    //处于收音时
    return true;
  } else {
    return !!result;
  }
});

const onVoiceStart = () => {
  state.statusVoice = "working";
  onReplaceResult("");
};

const onVoiceEnd = () => {
  state.statusVoice = "";
};

const onTab = () => {
  refTabVoice.value?.reset();
  refTabFile.value?.reset();
  onReplaceResult("");
};

onBeforeUnmount(() => {
  refTabVoice.value?.reset();
  refTabFile.value?.reset();
});
</script>
<style lang="scss" scoped>
.shower-result {
  max-height: 400px;
  overscroll-behavior: none;
}
:deep() {
  .compo-aiAudio {
    width: 80%;
  }
}
</style>
