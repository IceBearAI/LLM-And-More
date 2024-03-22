<template>
  <div :class="{ 'flex-row-reverse': msgType === 'answer' }" class="pa-5 d-flex w-100 position-relative">
    <div
      class="d-flex align-center justify-center flex-shrink-0 align-self-start"
      :class="msgType === 'answer' ? 'ml-2' : 'mr-2'"
    >
      <template v-if="msgType === 'answer'">
        <v-avatar>
          <img class="rounded-circle" :src="defaultAvatar" alt="avatar" width="40" />
        </v-avatar>
      </template>
      <template v-else>
        <div class="message-avatar__text">{{ headName }}</div>
      </template>
    </div>
    <div class="overflow-hidden">
      <slot v-if="$slots.text" name="text"></slot>
      <el-popover
        v-else
        :placement="msgType === 'question' ? 'bottom-start' : 'bottom-end'"
        :show-arrow="false"
        width="auto"
        popper-class="text-popover"
        :offset="5"
        :show-after="500"
        :teleported="false"
      >
        <template #reference>
          <div class="d-flex" :class="{ 'flex-row-reverse': msgType === 'question' }">
            <TextComp :text="text" msg-type="answer" />
          </div>
        </template>
        <div class="h-box">
          <div class="operate-btn hv-center" @click="emit('click:edit')">
            <IconEdit :size="20" />
          </div>
          <div class="operate-btn hv-center" @click="emit('click:delete')">
            <IconTrash :size="20" />
          </div>
        </div>
      </el-popover>
    </div>
  </div>
</template>
<script setup lang="ts">
import { computed } from "vue";
import { useUserStore } from "@/stores";
import defaultAvatar from "@/assets/images/favicon.png";
import { IconEdit, IconTrash } from "@tabler/icons-vue";
import TextComp from "@/views/model/chat-playground/components/Message/Text.vue";

interface IProps {
  msgType?: string;
  text?: string;
}

interface IEmits {
  (e: "click:edit"): void;
  (e: "click:delete"): void;
}

const props = withDefaults(defineProps<IProps>(), {
  msgType: "question",
  text: ""
});

const emit = defineEmits<IEmits>();

const userStore = useUserStore();

const headName = computed(() => {
  const name = userStore.userInfo.username;
  return name ? name[0].toUpperCase() : "Y";
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
.operate-btn {
  color: #666;
  background: #fff;
  margin-right: 1px;
  padding: 8px;
  border-radius: 12px;
  cursor: pointer;
  &:hover {
    background: #efefef;
  }
}
</style>
