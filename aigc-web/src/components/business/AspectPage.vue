<template>
  <router-view v-slot="{ Component }">
    <keep-alive :include="(route.meta.aspectPageInclude || []) as string[]">
      <component :is="Component" />
    </keep-alive>
  </router-view>
</template>
<script setup lang="ts">
import { provide, reactive, watch } from "vue";
import { useRoute, onBeforeRouteLeave, onBeforeRouteUpdate } from "vue-router";
import $ from "jquery";
import { ItfAspectPageState } from "@/types/AspectPageTypes.ts";

const route = useRoute();

const state = reactive<ItfAspectPageState>({
  scrollTop: {},
  methods: {
    refreshListPage() {}
  }
});
provide("provideAspectPage", state);

const getComponentName = routeItem => {
  return routeItem.matched[routeItem.matched.length - 1].components.default.__name;
};

const isKeepAliveComponent = componentName => {
  return ((route.meta.aspectPageInclude || []) as string[]).includes(componentName);
};

onBeforeRouteUpdate((to, from, next) => {
  let fromComponent = getComponentName(from);
  let toComponent = getComponentName(to);
  if (isKeepAliveComponent(fromComponent)) {
    //keep-alive组件，记录离开组件的scrollTop位置
    state.scrollTop[fromComponent] = $(document).scrollTop();
  }
  if (isKeepAliveComponent(toComponent)) {
    //keep-alive组件，滚到离开时的位置
    $(document).scrollTop(state.scrollTop[toComponent] || "0");
  } else {
    //非keep-alive组件,滚到页面顶部
    $(document).scrollTop(0);
  }
  next();
});
</script>
<style lang="scss"></style>
