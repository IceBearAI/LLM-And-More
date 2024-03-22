<script setup lang="ts">
import { ref, shallowRef } from "vue";
import { useCustomizerStore } from "@/stores/customizer";
import sidebarItems from "./sidebarItem";

import NavGroup from "./NavGroup/index.vue";
import NavItem from "./NavItem/index.vue";
import NavCollapse from "./NavCollapse/NavCollapse.vue";
import Logo from "../logo/Logo.vue";
import Darklogo from "@/assets/images/favicon.png";
const customizer = useCustomizerStore();
const sidebarMenu = shallowRef(sidebarItems);
</script>

<template>
  <v-navigation-drawer
    left
    v-model="customizer.Sidebar_drawer"
    elevation="0"
    rail-width="75"
    :mobile-breakpoint="1200"
    app
    class="leftSidebar"
    :rail="customizer.mini_sidebar"
    expand-on-hover
    width="256"
  >
    <div class="pa-5 pl-4">
      <div class="logo" style="display: flex; flex-direction: row; align-items: center">
        <RouterLink to="/" style="text-decoration: none; display: flex; flex-direction: row; align-items: center">
          <Logo />
        </RouterLink>
      </div>
    </div>
    <!-- ---------------------------------------------- -->
    <!---Navigation -->
    <!-- ---------------------------------------------- -->
    <perfect-scrollbar class="scrollnavbar">
      <v-list class="py-6 px-4">
        <!---Menu Loop -->
        <template v-for="(item, i) in sidebarMenu">
          <!---Item Sub Header -->
          <NavGroup :item="item" v-if="item.header" :key="item.title" />
          <!---If Has Child -->
          <NavCollapse class="leftPadding" :item="item" :level="0" v-else-if="item.children" />
          <!---Single Item-->
          <NavItem :item="item" v-else class="leftPadding" />
          <!---End Single Item-->
        </template>
      </v-list>
    </perfect-scrollbar>
  </v-navigation-drawer>
</template>
