<script setup lang="ts">
import { useAppStore } from "@/stores";
// import { useI18n } from "vue-i18n";
import type { languageType } from "@/types/HeaderTypes";
import flagZh from "@/assets/images/flag/icon-flag-zh.svg";
import flagEn from "@/assets/images/flag/icon-flag-en.svg";
import flagEs from "@/assets/images/flag/icon-flag-es.svg";
import flagPh from "@/assets/images/flag/icon-flag-ph.svg";
import { EnumLanguage } from "@/types/apps/System.ts";

const languageDD: languageType[] = [
  { title: "中国人", subtext: "Chinese", value: EnumLanguage.Chinese, avatar: flagZh },
  { title: "English", subtext: "UK", value: EnumLanguage.English, avatar: flagEn },
  { title: "Espanyol", subtext: "Espana", value: EnumLanguage.Espanyol, avatar: flagEs },
  { title: "Pilipino", subtext: "Filipino", value: EnumLanguage.Pilipino, avatar: flagPh }
];

const appStore = useAppStore();
// const i18n = useI18n();

const changeLanguage = (lang: string) => {
  // i18n.locale.value = lang;
  appStore.setLocalLanguage(lang);
};
</script>
<template>
  <!-- ---------------------------------------------- -->
  <!-- language DD -->
  <!-- ---------------------------------------------- -->
  <v-menu location="bottom">
    <template v-slot:activator="{ props }">
      <v-btn icon variant="text" color="primary" v-bind="props">
        <v-avatar size="22">
          <img
            v-if="$i18n.locale === EnumLanguage.English"
            :src="flagEn"
            :alt="$i18n.locale"
            width="24"
            height="24"
            class="obj-cover"
          />
          <img
            v-else-if="$i18n.locale === EnumLanguage.Espanyol"
            :src="flagEs"
            :alt="$i18n.locale"
            width="24"
            height="24"
            class="obj-cover"
          />
          <img
            v-else-if="$i18n.locale === EnumLanguage.Pilipino"
            :src="flagPh"
            :alt="$i18n.locale"
            width="24"
            height="24"
            class="obj-cover"
          />
          <!--默认使用中文-->
          <img v-else :src="flagZh" :alt="$i18n.locale" width="24" height="24" class="obj-cover" />
        </v-avatar>
      </v-btn>
    </template>
    <v-sheet rounded="md" width="200" elevation="10">
      <v-list class="theme-list">
        <v-list-item
          v-for="(item, index) in languageDD"
          :key="index"
          color="primary"
          :active="$i18n.locale == item.value"
          class="d-flex align-center"
          @click="changeLanguage(item.value)"
        >
          <template v-slot:prepend>
            <v-avatar size="22">
              <img :src="item.avatar" :alt="item.avatar" width="22" height="22" class="obj-cover" />
            </v-avatar>
          </template>
          <v-list-item-title class="text-subtitle-1 font-weight-regular">
            {{ item.title }}
            <span class="text-disabled text-subtitle-1 pl-2">({{ item.subtext }})</span>
          </v-list-item-title>
        </v-list-item>
      </v-list>
    </v-sheet>
  </v-menu>
</template>
