<template>
  <BaseBreadcrumb :title="page.title"></BaseBreadcrumb>
  <v-card elevation="10">
    <div class="chat-playground d-flex">
      <div class="left-part pa-4 d-none d-md-block">
        <h4 class="text-h6 mb-6">助理设置</h4>
        <v-alert class="mb-3 text-body-2" color="textPrimary" variant="tonal">
          系统消息包含在提示的开头，用于为模型提供上下文、说明或与用例相关的其他信息。
          可以使用系统消息来描述助手的个性，定义模型应回答和不应回答的内容，以及定义模型响应的格式。
        </v-alert>
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
      <div class="middle-part">
        <h4 class="text-h6 pa-4">聊天会话</h4>
        <v-divider />
        <perfect-scrollbar ref="scrollRef" class="middle-part__list">
          <template v-if="chatList.length > 0">
            <Message v-for="(item, index) in chatList" :key="index" :chat-item="item" />
          </template>
        </perfect-scrollbar>
        <v-divider />
        <SendMsg v-model="question" @submit="handleTextSend" @clear="handleChatClear" :send-loading="sendLoading" />
      </div>
      <div class="right-part pa-4">
        <h4 class="text-h6 mb-6">配置</h4>
        <Config ref="configRef" @change:model="modelChange" />
      </div>
    </div>
  </v-card>
  <ConfirmByClick ref="refConfirmDelete" @submit="doChatClear">
    <template #text>确定清空会话？</template>
  </ConfirmByClick>
</template>
<script setup lang="ts">
import { ref } from "vue";
import { useScroll } from "@/hooks/useScroll";
import BaseBreadcrumb from "@/components/shared/BaseBreadcrumb.vue";
import { http } from "@/utils";
import Config from "./components/Config.vue";
import SendMsg from "./components/SendMsg.vue";
import Message from "./components/Message/index.vue";
import ConfirmByClick from "@/components/business/ConfirmByClick.vue";

const { scrollRef, scrollToBottom, scrollToBottomIfAtBottom } = useScroll();

// theme breadcrumb
const page = ref({ title: "聊天操场" });

const systemMessage = ref("");
const chatList = ref([]);
const question = ref("");
const sendLoading = ref(false);
const configRef = ref();
const refConfirmDelete = ref();
const systemPrompt = ref("");
const systemPrompts = ref([]);

const handleTextSend = () => {
  const message = question.value;
  question.value = "";
  generateAnswer("text", message);
};

const generateAnswer = async (sendType, sendContent, ext = {}) => {
  let content = sendContent;
  const data = getChatSendParams(content); // 获取参数应该放在addChat方法前
  addChat({
    content,
    contentType: sendType,
    createdAt: new Date(),
    ext,
    msgType: "question"
  });
  scrollToBottom();
  sendLoading.value = true;

  addChat({
    loading: true,
    content: "",
    contentType: "text",
    createdAt: new Date(),
    msgType: "answer"
  });
  scrollToBottom();
  try {
    await http.post({
      url: "/channels/chat/completions",
      timeout: 300 * 1000, // 请求超时时间设置为5分钟
      data,
      onDownloadProgress: event => {
        const xhr = event.target;
        const { responseText } = xhr;
        const lastIndex = responseText.lastIndexOf("\n", responseText.length - 2);
        let chunk = responseText;
        if (lastIndex !== -1) {
          chunk = responseText.substring(lastIndex);
        }
        try {
          const responseJson = JSON.parse(chunk);
          const listLastIndex = chatList.value.length - 1;
          if (responseJson.success) {
            const data = responseJson.data;
            // currentMessageId.value = data.messageId;
            updateChat(listLastIndex, {
              loading: true,
              content: data.fullContent,
              contentType: data.contentType ?? "text",
              createdAt: data.createdAt,
              msgType: "answer"
            });
          } else {
            const currentChat = chatList.value[listLastIndex];
            if (currentChat.content && currentChat.content !== "") {
              updateChatSome(listLastIndex, {
                content: `${currentChat.content}\n[${responseJson.message}]`,
                loading: true
              });
              return;
            }

            updateChat(listLastIndex, {
              loading: true,
              content: responseJson.message,
              contentType: "text",
              createdAt: new Date(),
              msgType: "answer"
            });
          }
          scrollToBottomIfAtBottom();
        } catch (error) {
          //
        }
      }
    });
  } catch (error) {
    console.log(error);
  } finally {
    updateChatSome(chatList.value.length - 1, { loading: false });
    sendLoading.value = false;
    // currentMessageId.value = null;
  }
};

const getChatSendParams = sendMessage => {
  const configData = configRef.value.getData();
  const result = configData.data;
  const sendHistoryCount = configData.sendHistoryCount * 2; // *2 的目的是历史消息要成对发送
  let history = [];

  // 默认系统消息
  const messages = [
    {
      role: "system",
      content: systemMessage.value
    }
  ];
  // 添加历史消息
  if (sendHistoryCount !== 0) {
    history = chatList.value.slice(-sendHistoryCount);
  }
  history.forEach(item => {
    const role = item.msgType === "question" ? "user" : "assistant";
    messages.push({
      role,
      content: item.content
    });
  });
  // 添加当前发送的消息
  messages.push({
    role: "user",
    content: sendMessage
  });
  result["messages"] = messages;
  return result;
};

const addChat = chat => {
  chatList.value.push(chat);
};
const updateChat = (index, chat) => {
  chatList.value[index] = chat;
};

const updateChatSome = (index, chat) => {
  chatList.value[index] = { ...chatList.value[index], ...chat };
};

const handleChatClear = () => {
  refConfirmDelete.value.show({
    width: "400px"
  });
};
const doChatClear = () => {
  chatList.value = [];
  refConfirmDelete.value.hide();
};
const modelChange = async data => {
  const [err, res] = await http.get({
    url: `/api/models/${data.modelName}/info`
  });
  if (res) {
    if (res.fineTuned) {
      const prompts = res.fineTuned.systemPrompts || [];
      systemPrompts.value = prompts;
      systemPrompt.value = prompts[0];
      systemMessage.value = prompts[0];
    } else {
      systemPrompts.value = [];
      systemPrompt.value = "";
      systemMessage.value = "";
    }
  }
};

const systemPromptChange = value => {
  systemMessage.value = value;
};
</script>
<style lang="scss" scoped>
.left-part {
  width: 300px;
  border-right: 1px solid rgb(var(--v-theme-borderColor));
}
.middle-part {
  flex: 1;
  overflow: hidden; // 防止不缩放
}
.right-part {
  width: 300px;
  border-left: 1px solid rgb(var(--v-theme-borderColor));
}
.middle-part__list {
  height: 530px;
}
</style>
