<template>
  <div :class="{ 'flex-row-reverse': props.chatItem.msgType === 'question' }" class="pa-5 d-flex w-100 overflow-hidden">
    <div
      class="d-flex align-center justify-center flex-shrink-0 align-self-start"
      :class="props.chatItem.msgType === 'question' ? 'ml-2' : 'mr-2'"
    >
      <template v-if="props.chatItem.msgType === 'answer'">
        <v-avatar>
          <img class="rounded-circle" :src="defaultAvatar" alt="avatar" width="40" />
        </v-avatar>
      </template>
      <template v-else>
        <div class="message-avatar__text">{{ headName }}</div>
      </template>
    </div>
    <div class="overflow-hidden">
      <p class="message-content__time" :class="{ 'text-right': props.chatItem.msgType === 'question' }">{{ showTime }}</p>
      <div class="d-flex mt-2" :class="{ 'flex-row-reverse': props.chatItem.msgType === 'question' }">
        <template v-if="props.chatItem.contentType === 'text'">
          <TextComp :text="props.chatItem.content" :msg-type="props.chatItem.msgType" :loading="props.chatItem.loading" />
        </template>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { computed } from "vue";
import { useUserStore } from "@/stores";
import TextComp from "./Text.vue";
import defaultAvatar from "@/assets/images/favicon.png";
import { format } from "@/utils";

interface IProps {
  chatItem?: Record<string, any>;
}

const props = withDefaults(defineProps<IProps>(), {
  chatItem: () => ({})
});

const userStore = useUserStore();

const headName = computed(() => {
  const name = userStore.userInfo.username;
  return name ? name[0].toUpperCase() : "Y";
});

const showTime = computed(() => {
  return format.dateFormat(props.chatItem.createdAt, "YYYY-MM-DD HH:mm:ss");
});
</script>
<style lang="scss" scoped>
.message-avatar__text {
  width: 36px;
  height: 36px;
  line-height: 36px;
  background-color: #0488d2;
  color: #fff;
  border-radius: 50%;
  text-align: center;
}
.message-content__time {
  color: #b4bbc4;
  font-size: 12px;
  line-height: 1;
  margin: 0;
}
</style>
