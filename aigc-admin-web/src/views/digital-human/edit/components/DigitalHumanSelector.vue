<template>
  <div>
    <!--预览区-->
    <div class="aspect-square">
      <template v-if="selectedDigitalHuman.video">
        <video
          v-if="selectedDigitalHuman.mediaType == 'video'"
          style="width: 100%; height: 100%"
          :src="selectedDigitalHuman.video"
          controls
          :poster="selectedDigitalHuman.cover"
          class="block"
        ></video>
        <v-img v-else-if="selectedDigitalHuman.mediaType == 'image'" :src="selectedDigitalHuman.video">
          <template #placeholder>
            <div class="d-flex align-center justify-center fill-height">
              <v-progress-circular color="primary" indeterminate :width="2" :size="42"></v-progress-circular>
            </div>
          </template>
        </v-img>
      </template>
      <div v-else class="h-full w-full bg-gray-100"></div>
    </div>

    <v-row class="mt-5">
      <v-col>
        <v-text-field
          density="compact"
          v-model="formData.name"
          label="搜索角色名称"
          hide-details
          clearable
          variant="outlined"
          color="red"
          @keyup.enter="getList"
          @click:clear="getList"
        ></v-text-field>
      </v-col>
    </v-row>

    <div class="box-list scrollbar-auto flex-1 pt-4 overflow-auto">
      <template v-if="state.list.length">
        <v-card
          variant="outlined"
          rounded="sm"
          elevation="0"
          class="pa-2 bg-hover-secondary d-inline-block list-item ma-2 overflow-visible"
          pointer
          v-for="(item, index) in state.list"
        >
          <div @click="onSelect(item)">
            <div class="relative overflow-hidden">
              <TagCorner :class="map.mediaType[item.mediaType].className">
                {{ map.mediaType[item.mediaType].text }}
              </TagCorner>
              <div class="image-human overflow-hidden">
                <img :src="item.cover" :alt="item.cname" />
              </div>
              <div class="hv-center mt-2">
                <div class="line1">{{ item.cname }}</div>
              </div>
            </div>
            <IconChecked v-if="item.name === selectedDigitalHuman?.name" />
          </div>
        </v-card>
      </template>
      <NoData v-else>暂无数据，请更换搜索条件</NoData>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref, reactive, computed, onMounted, toRefs, inject } from "vue";
import { useRoute } from "vue-router";
import { http, map } from "@/utils";
import IconChecked from "@/components/ui/IconChecked.vue";
import TagCorner from "@/components/business/TagCorner.vue";

import { ItfProvideState, ItfDigitalHuman } from "../types";

const route = useRoute();

const provieDigitalHumanEdit = <ItfProvideState>inject("provieDigitalHumanEdit");

const state = reactive<{
  list: Array<ItfDigitalHuman>;
  [x: string]: any;
}>({
  list: [],
  formData: {
    name: ""
  }
});

const { formData } = toRefs(state);

const { selectedDigitalHuman } = toRefs(provieDigitalHumanEdit);

function onSelect(item) {
  provieDigitalHumanEdit.selectedDigitalHuman = item;
}

const getList = async () => {
  await http
    .get({
      showLoading: true,
      url: "/api/digitalhuman/person/list",
      data: {
        pageSize: 20,
        page: 1,
        ...state.formData
      }
    })
    .then(([err, res]) => {
      if (res) {
        state.list = res.list;

        if (state.list.length > 0) {
          onSelect(state.list[0]);
        } else {
          provieDigitalHumanEdit.selectedDigitalHuman = {
            name: "",
            cname: "",
            cover: "",
            video: "",
            mediaType: ""
          };
        }
      }
    });
};

onMounted(() => {
  getList();
});
</script>
<style lang="scss" scoped>
.box-list {
  overflow: scroll;
  white-space: nowrap;
  width: 100%;
}
.list-item {
  display: inline-block !important;
  width: 180px;

  &:hover {
    img {
      transform: scale(1.1);
    }
  }
  .image-human {
    aspect-ratio: 4/3;

    img {
      width: 100%;
      height: 100%;
      object-fit: cover;
      transition: all 0.2s linear;
    }
  }
}
</style>
