<script setup lang="ts">
import { computed } from "vue";
import { IconMail } from "@tabler/icons-vue";

import { useUserStore } from "@/stores";
const userStore = useUserStore();

const user = computed(() => userStore.userInfo || { username: "" });
</script>

<template>
  <!-- ---------------------------------------------- -->
  <!-- notifications DD -->
  <!-- ---------------------------------------------- -->
  <v-menu :close-on-content-click="false">
    <template v-slot:activator="{ props }">
      <v-btn class="custom-hover-primary" variant="text" v-bind="props" icon>
        <v-avatar size="35">
          <img src="@/assets/images/profile/user2.jpg" width="35" alt="Julia" />
        </v-avatar>
      </v-btn>
    </template>
    <v-sheet rounded="md" width="360" elevation="10">
      <div class="px-8 pt-6">
        <h6 class="text-h5 font-weight-medium">用户信息</h6>
        <div class="d-flex align-center mt-4 pb-6">
          <v-avatar size="80">
            <img src="@/assets/images/profile/user2.jpg" width="80" />
          </v-avatar>
          <div class="ml-3">
            <h6 class="text-h6 mb-n1">{{ user.username }}</h6>
            <div class="d-flex align-center mt-3">
              <IconMail :size="18" stroke-width="1.5" />
              <span class="text-subtitle-1 font-weight-regular textSecondary ml-2">{{ user.username }}</span>
            </div>
          </div>
        </div>
        <v-divider></v-divider>
      </div>
      <div class="pt-4 pb-6 px-8 text-center">
        <v-btn color="primary" variant="outlined" block @click="userStore.logout()">退出</v-btn>
      </div>
    </v-sheet>
  </v-menu>
</template>
