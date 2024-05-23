<template>
  <BaseBreadcrumb class="chat-breadcrumb" :title="page.title"></BaseBreadcrumb>
  <v-card elevation="10">
    <div class="chat-playground flex" style="height: calc(100vh - 180px)">
      <div class="left-part pa-4 d-none d-md-block">
        <h4 class="text-h6 mb-6">助理设置</h4>
        <!-- <v-alert class="mb-3 text-body-2" color="textPrimary" variant="tonal">
          系统消息包含在提示的开头，用于为模型提供上下文、说明或与用例相关的其他信息。
          可以使用系统消息来描述助手的个性，定义模型应回答和不应回答的内容，以及定义模型响应的格式。
        </v-alert> -->
        <v-sheet>
          <v-select
            v-if="systemPrompts.length > 1"
            class="mb-3 pt-2"
            v-model="systemPrompt"
            @update:modelValue="systemPromptChange"
            :clearable="false"
            hide-details
            :items="systemPrompts"
            label="请选择系统提示词"
          >
          </v-select>
          <v-label class="mb-2 font-weight-medium">系统消息</v-label>
          <v-textarea
            v-model="systemMessage"
            variant="outlined"
            placeholder="可以给定助理角色及相关要求，不要超过2000个字"
            rows="5"
            no-resize
            color="primary"
            row-height="25"
            shaped
            hide-details
          ></v-textarea>
        </v-sheet>
      </div>
      <div class="right-part flex flex-col pa-2">
        <div class="text-right mb-2">
          <v-btn size="small" flat color="error" variant="outlined" @click="addChat"
            ><IconPlus stroke-width="1.5" :size="18" class="mr-2" />{{
              `添加模型(${currentChatCounts}/${totalModelCount})`
            }}</v-btn
          >
        </div>
        <v-row class="overflow-auto" dense>
          <v-col :cols="currentChatCounts === 1 ? 12 : 6" class="h-100" v-for="(item, index) in renderChats" :key="item.key">
            <RenderChat
              :ref="el => (item.comInstance = el)"
              @change:model="modelChange($event, index)"
              :system-message="systemMessage"
            >
              <template #title-left>{{ `#${index + 1}` }}</template>
              <template v-if="currentChatCounts > 1" #title-right>
                <v-btn size="x-small" color="inherit" icon variant="text">
                  <IconDots width="14" stroke-width="1.5" />
                  <v-menu activator="parent">
                    <v-list density="compact" @click:select="operateMenuClick($event, index)">
                      <v-list-item key="remove" value="remove" hide-details min-height="38">
                        <v-list-item-title>移除</v-list-item-title>
                      </v-list-item>
                    </v-list>
                  </v-menu>
                </v-btn>
              </template>
            </RenderChat>
          </v-col>
        </v-row>
        <SendMsg v-model="question" @submit="handleTextSend" @clear="handleChatClear" />
      </div>
    </div>
  </v-card>
  <ConfirmByClick ref="refConfirmDelete" @submit="doChatClear">
    <template #text>确定清空会话？</template>
  </ConfirmByClick>
</template>
<script setup lang="ts">
import { ref, computed } from "vue";
import BaseBreadcrumb from "@/components/shared/BaseBreadcrumb.vue";
import { http } from "@/utils";
import ConfirmByClick from "@/components/business/ConfirmByClick.vue";
import RenderChat from "./components/RenderChat.vue";
import { IconPlus } from "@tabler/icons-vue";
import SendMsg from "./components/SendMsg.vue";
import { IconDots } from "@tabler/icons-vue";

const totalModelCount = 2;

// theme breadcrumb
const page = ref({ title: "聊天操场" });

const renderChats = ref([
  {
    key: Date.now(),
    comInstance: null
  }
]);
const question = ref("");
const systemMessage = ref("");
const refConfirmDelete = ref();
const systemPrompt = ref("");
const systemPrompts = ref([]);

const currentChatCounts = computed(() => {
  return renderChats.value.length;
});

const addChat = () => {
  if (currentChatCounts.value === totalModelCount) return;
  renderChats.value.push({
    key: Date.now(),
    comInstance: null
  });
};

const operateMenuClick = ({ id }, index) => {
  if (id === "remove") {
    renderChats.value.splice(index, 1);
  }
};

const handleTextSend = () => {
  const message = question.value;
  question.value = "";
  renderChats.value.forEach(item => {
    if (item.comInstance) {
      item.comInstance.send({
        sendType: "text",
        prompt: message
      });
    }
  });
};

const handleChatClear = () => {
  refConfirmDelete.value.show({
    width: "400px"
  });
};
const doChatClear = () => {
  renderChats.value.forEach(item => {
    if (item.comInstance) {
      item.comInstance.clearChat();
    }
  });
  refConfirmDelete.value.hide();
};
const modelChange = (data, index) => {
  if (index !== 0) return;
  if (data.fineTuned) {
    const prompts = data.fineTuned.systemPrompts || [];
    systemPrompts.value = prompts;
    systemPrompt.value = prompts[0];
    systemMessage.value = prompts[0];
  } else {
    systemPrompts.value = [];
    systemPrompt.value = "";
    systemMessage.value = "";
  }
};

const systemPromptChange = value => {
  systemMessage.value = value;
};
</script>
<style lang="scss" scoped>
.chat-breadcrumb.mb-6 {
  margin-bottom: 8px !important;
}
.left-part {
  width: 300px;
  border-right: 1px solid rgb(var(--v-theme-borderColor));
}
.right-part {
  flex: 1;
  overflow: hidden; // 防止不缩放
}
</style>
