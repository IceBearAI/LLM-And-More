<template>
  <!--发声人列表-->
  <div ref="refContainer" class="w-100" :class="{ 'opacity-0': style.isReady == false }">
    <template v-if="state.listSpeaker?.length > 0">
      <div class="d-flex">
        <div class="d-flex flex-wrap voice-list flex-1 overflow-auto pt-4 scrollbar-auto">
          <v-card
            :id="item.speakName"
            variant="outlined"
            elevation="0"
            class="voice-item my-1 mr-5 bg-hover-secondary d-flex align-items"
            rounded="sm"
            pointer
            v-for="item in state.listSpeaker"
            :class="{ active: item.speakName === data }"
          >
            <v-card-text @click="onSelectSpeaker(item)" class="d-flex align-center py-0 px-2">
              <v-avatar size="40 ">
                <img :src="item.headImg" :alt="item.speakCname" class="w-100" />
              </v-avatar>
              <div class="ml-3 text-body-2 text-black">
                {{ item.speakCname }}<span class="text-light" v-if="item.subTitle">（{{ item.subTitle }}）</span>
              </div>

              <IconChecked v-if="item.speakName === data" />
            </v-card-text>
          </v-card>
        </div>
      </div>
      <!--选择发声人，音频播放-->
      <v-card variant="outlined" class="mt-8">
        <v-card-title class="d-flex align-center">
          <v-avatar size="40 ">
            <img :src="optionInfo.headImg" :alt="optionInfo.speakCname" class="w-100" />
          </v-avatar>
          <div class="ml-3 text-body-2">
            {{ optionInfo.speakCname }}<span class="text-light" v-if="optionInfo.subTitle">（{{ optionInfo.subTitle }}）</span>
          </div>
        </v-card-title>
        <v-card-text class="mt-4">
          <div class="hv-center">
            <AiAudio :src="optionInfo.speakDemo" type="simple" />
          </div>
        </v-card-text>
      </v-card>
    </template>
    <NoData v-else></NoData>
  </div>
</template>
<script setup>
import { reactive, toRefs, ref, onMounted, nextTick } from "vue";
import { useMapRemoteStore } from "@/stores";
import IconChecked from "@/components/ui/IconChecked.vue";
import { useVModel } from "@vueuse/core";
import { http } from "@/utils";
import $ from "jquery";
import AiAudio from "@/components/business/AiAudio.vue";

const { loadDictTree, getLabels } = useMapRemoteStore();
const refContainer = ref();

const state = reactive({
  listSpeaker: [],
  queryParams: {},
  optionInfo: {
    speakName: "",
    speakCname: "",
    headImg: "",
    speakDemo: "",
    gender: "",
    ageGroup: "",
    subTitle: ""
  },
  style: {
    isReady: false
  }
});
const { optionInfo, style } = toRefs(state);

const props = defineProps({
  modelValue: String,
  infos: Object
});
const emits = defineEmits(["update:modelValue", "update:infos"]);
const data = useVModel(props, "modelValue", emits);

const getSubTitle = item => {
  return getLabels(
    [
      ["speak_age_group", item.ageGroup],
      ["speak_gender", item.gender]
    ],
    ret => {
      if (ret.length) {
        return ret.join("") + "声";
      } else {
        return "未知";
      }
    }
  );
};

// 获取数字人列表
const getList = async () => {
  await loadDictTree(["speak_age_group", "speak_gender"]);
  const [err, res] = await http.get({
    showLoading: refContainer.value,
    url: "/api/voice/speak",
    data: {
      pageSize: -1,
      ...state.queryParams
    }
  });
  if (res) {
    state.listSpeaker = res.list.map(item => {
      return {
        ...item,
        subTitle: getSubTitle(item)
      };
    });
    if (state.listSpeaker.length) {
      let initValue = data.value;
      if (initValue) {
        //存在初始值
        let matchedIndex = state.listSpeaker.findIndex(itemOption => {
          return itemOption.speakName == initValue;
        });
        let matchedItem = state.listSpeaker[matchedIndex];
        if (matchedItem) {
          onSelectSpeaker(matchedItem);
          if (matchedIndex > 6) {
            //超出首屏，滚到对应位置
            nextTick(() => {
              $("#" + initValue)[0]?.scrollIntoView();
            });
          }
        } else {
          window.errorMsg(`未找到初始值对应的数字人 ${initValue}`);
          //选中第一项
          let firstItem = state.listSpeaker[0];
          onSelectSpeaker(firstItem);
        }
      } else {
        //没有初始值时，默认选中第一个
        let firstItem = state.listSpeaker[0];
        onSelectSpeaker(firstItem);
      }
    } else {
      onSelectSpeaker({});
    }
    state.style.isReady = true;
  }
};

const onSelectSpeaker = item => {
  data.value = item.speakName;
  state.optionInfo = item;
  emits("update:infos", item);
};

onMounted(() => {
  getList();
});

defineExpose({
  reload(params) {
    state.queryParams = params;
    getList();
  }
});
</script>
<style lang="scss" scoped>
.voice-list {
  max-height: 200px;
  .voice-item {
    overflow: visible;
    width: 180px;
    height: 60px;
  }
}
</style>
