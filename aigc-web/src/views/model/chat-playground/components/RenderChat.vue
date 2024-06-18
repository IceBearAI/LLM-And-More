<template>
  <div class="render-chat border rounded-md h-100 d-flex flex-column">
    <div class="d-flex align-center justify-space-between pa-3">
      <h4 class="text-h6">
        <slot v-if="$slots['title-left']" name="title-left" />
        <template v-else> 聊天 </template>
      </h4>
      <v-menu eager :close-on-content-click="false">
        <template v-slot:activator="{ props }">
          <v-btn class="menu-model-btn" height="32" v-bind="props" flat variant="text">
            <span class="truncate">{{ currentModel.modelName }}</span
            ><IconChevronDown stroke-width="1.5" :size="18" class="ml-1" />
          </v-btn>
        </template>
        <v-sheet class="pa-4" rounded="md" width="350" elevation="10">
          <Config ref="configRef" @update:model:info="modelChange" />
        </v-sheet>
      </v-menu>
      <slot v-if="$slots['title-right']" name="title-right"></slot>
    </div>
    <v-divider />
    <perfect-scrollbar ref="scrollRef" class="flex-1">
      <template v-if="chatList.length > 0">
        <Message v-for="(item, index) in chatList" :key="index" :chat-item="item" />
      </template>
    </perfect-scrollbar>
  </div>
</template>
<script setup lang="ts">
import { ref } from "vue";
import { useScroll } from "@/hooks/useScroll";
import { http } from "@/utils";
import Config from "./Config.vue";
import Message from "./Message/index.vue";
import { IconChevronDown } from "@tabler/icons-vue";
import awaitTo from "await-to-js";
import { toast } from "vue3-toastify";

interface IProps {
  systemMessage?: string;
}
const props = withDefaults(defineProps<IProps>(), {
  systemMessage: ""
});

interface IEmits {
  (e: "change:model", val: string): void;
}
const emit = defineEmits<IEmits>();

const { scrollRef, scrollToBottom, scrollToBottomIfAtBottom } = useScroll();

const chatList = ref([]);
const sendLoading = ref(false);
const configRef = ref();
const currentModel = ref<Record<string, any>>({});

const modelChange = val => {
  currentModel.value = val;
  emit("change:model", val);
};

const generateAnswer = async ({ data, resolve, reject }) => {
  if (sendLoading.value) {
    toast.error("请等上个会话生成结束后再试");
    reject("请等上个会话生成结束后再试");
    return;
  }
  let content = data.prompt;
  const sendData = getChatSendParams(content); // 获取参数应该放在addChat方法前
  addChat({
    content,
    contentType: data.sendType || "text",
    createdAt: new Date(),
    ext: data.ext,
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
    const [err, res] = await http.post({
      url: "/channels/chat/completions",
      timeout: 300 * 1000, // 请求超时时间设置为5分钟
      data: sendData,
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
    if (res) {
      resolve(res);
    } else {
      reject(err);
    }
  } catch (error) {
    reject(error);
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
      content: props.systemMessage
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

const clearChat = () => {
  chatList.value = [];
};

const send = data => {
  return awaitTo(
    new Promise((resolve, reject) => {
      generateAnswer({
        data,
        resolve,
        reject
      });
    })
  );
};

defineExpose({
  send,
  clearChat
});
</script>
<style lang="scss" scoped>
:deep() {
  .menu-model-btn .v-btn__content {
    max-width: 305px;
  }
}
</style>
