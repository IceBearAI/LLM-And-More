<script setup lang="ts">
import { ref } from "vue";
import { useUserStore } from "@/stores";
import { TypeTenant } from "@/stores/types/stores.type.ts";
import { useRouter } from "vue-router";
import { IconWorld } from "@tabler/icons-vue";
import { nextTick } from "process";

const router = useRouter();
const userStore = useUserStore();

const onSelect = (item: TypeTenant) => {
  //变更租户id
  userStore.userInfo.tenantId = item.id;
  //跳转首页
  // router.replace({ path: "/" }).then(() => {
  //   window.location.reload();
  // });
  window.location.replace("/");
};
</script>
<template>
  <v-menu location="bottom">
    <template v-slot:activator="{ props }">
      <v-btn icon variant="text" color="primary" v-bind="props">
        <IconWorld stroke-width="1.2" :size="24" />
      </v-btn>
    </template>
    <v-sheet rounded="md" width="120" elevation="10">
      <v-list class="theme-list">
        <v-list-item
          v-for="(item, index) in userStore.userInfo.tenants"
          :key="index"
          :value="index"
          active-color="primary"
          :active="userStore.userInfo.tenantId == item.id"
          class="d-flex align-center"
          @click="onSelect(item)"
        >
          <v-list-item-title class="text-subtitle-1 font-weight-regular">
            {{ item.name }}
          </v-list-item-title>
        </v-list-item>
      </v-list>
    </v-sheet>
  </v-menu>
</template>
