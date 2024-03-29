<script setup lang="ts">
import { RouterView } from "vue-router";
import VerticalSidebarVue from "./vertical-sidebar/VerticalSidebar.vue";
import VerticalHeaderVue from "./vertical-header/VerticalHeader.vue";
import HorizontalHeader from "./horizontal-header/HorizontalHeader.vue";
import HorizontalSidebar from "./horizontal-sidebar/HorizontalSidebar.vue";
import Customizer from "./customizer/Customizer.vue";
import { useCustomizerStore } from "../../stores/customizer";
import { pl, zhHans } from "vuetify/locale";
import { IconSettings } from "@tabler/icons-vue";
const customizer = useCustomizerStore();
</script>

<template>
  <!-----RTL LAYOUT------->
  <v-locale-provider v-if="customizer.setRTLLayout" rtl>
    <v-app
      :theme="customizer.actTheme"
      :class="[
        customizer.actTheme,
        customizer.mini_sidebar ? 'mini-sidebar' : '',
        customizer.setHorizontalLayout ? 'horizontalLayout' : 'verticalLayout',
        customizer.setBorderCard ? 'cardBordered' : '',
        customizer.inputBg ? 'inputWithbg' : ''
      ]"
    >
      <!---Customizer location left side--->
      <v-navigation-drawer
        app
        temporary
        elevation="10"
        location="left"
        v-model="customizer.Customizer_drawer"
        width="320"
        class="left-customizer"
      >
        <Customizer />
      </v-navigation-drawer>
      <VerticalSidebarVue v-if="!customizer.setHorizontalLayout" />
      <VerticalHeaderVue v-if="!customizer.setHorizontalLayout" />
      <HorizontalHeader v-if="customizer.setHorizontalLayout" />
      <HorizontalSidebar v-if="customizer.setHorizontalLayout" />

      <v-main>
        <v-container fluid class="page-wrapper pb-sm-15 pb-10">
          <div :class="customizer.boxed ? 'maxWidth' : ''">
            <RouterView />
            <!-- <v-btn
              class="customizer-btn"
              size="large"
              icon
              variant="flat"
              color="primary"
              @click.stop="customizer.SET_CUSTOMIZER_DRAWER(!customizer.Customizer_drawer)"
            >
              <IconSettings />
            </v-btn> -->
          </div>
        </v-container>
      </v-main>
    </v-app>
  </v-locale-provider>
  <!-----LTR LAYOUT------->
  <v-locale-provider v-else>
    <v-app
      :theme="customizer.actTheme"
      :class="[
        customizer.actTheme,
        customizer.mini_sidebar ? 'mini-sidebar' : '',
        customizer.setHorizontalLayout ? 'horizontalLayout' : 'verticalLayout',
        customizer.setBorderCard ? 'cardBordered' : '',
        customizer.inputBg ? 'inputWithbg' : ''
      ]"
    >
      <!---Customizer location left side--->
      <v-navigation-drawer app temporary elevation="10" location="right" v-model="customizer.Customizer_drawer" width="320">
        <Customizer />
      </v-navigation-drawer>
      <VerticalSidebarVue v-if="!customizer.setHorizontalLayout" />
      <VerticalHeaderVue v-if="!customizer.setHorizontalLayout" />
      <HorizontalHeader v-if="customizer.setHorizontalLayout" />
      <HorizontalSidebar v-if="customizer.setHorizontalLayout" />

      <v-main>
        <v-container fluid class="page-wrapper pb-sm-15 pb-10">
          <div :class="customizer.boxed ? 'maxWidth' : ''">
            <RouterView />
            <!-- <v-btn
              class="customizer-btn"
              size="large"
              icon
              variant="flat"
              color="primary"
              @click.stop="customizer.SET_CUSTOMIZER_DRAWER(!customizer.Customizer_drawer)"
            >
              <IconSettings />
            </v-btn> -->
          </div>
        </v-container>
      </v-main>
    </v-app>
  </v-locale-provider>
</template>
