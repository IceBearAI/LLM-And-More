<template>
  <NavBack :backUrl="`/ai-assistant/assistants/detail?assistantId=${assistantId}`">助手操场</NavBack>
  <v-card elevation="10" class="mt-4" style="height: calc(100vh - 188px)">
    <div class="chat-playground d-flex h-100">
      <div class="left-part pa-4 d-none d-md-block">
        <h4 class="text-h6 mb-6">助理设置</h4>
        <v-sheet>
          <v-label class="mb-2 font-weight-medium">系统提示词</v-label>
          <v-textarea
            v-model="systemMessage"
            variant="outlined"
            placeholder="可以给定助理角色及相关要求，不要超过2000个字"
            rows="3"
            no-resize
            color="primary"
            row-height="25"
            shaped
          ></v-textarea>
        </v-sheet>
        <v-sheet>
          <v-label class="mb-2 font-weight-medium">指令</v-label>
          <v-textarea
            v-model="sendParams.instructions"
            variant="outlined"
            placeholder="请输入指令"
            rows="6"
            no-resize
            color="primary"
            row-height="25"
            shaped
          ></v-textarea>
        </v-sheet>
        <v-sheet>
          <v-label class="mb-2 font-weight-medium">模型</v-label>
          <ModelSelect v-model="sendParams.modelName" />
        </v-sheet>
        <v-sheet>
          <v-label class="mb-2 font-weight-medium">工具</v-label>
          <template v-for="item in assistantData.tools">
            <v-checkbox v-model="sendParams.toolIds" color="primary" hide-details :value="item.toolId">
              <template #label>
                <el-tooltip :content="item.description" placement="top">
                  {{ item.name }}
                </el-tooltip>
              </template>
            </v-checkbox>
          </template>
        </v-sheet>
      </div>
      <div class="middle-part d-flex flex-column">
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
    </div>
  </v-card>
  <ConfirmByClick ref="refConfirmDelete" @submit="doChatClear">
    <template #text>确定清空会话？</template>
  </ConfirmByClick>
</template>
<script setup lang="ts">
import { ref, onMounted, reactive } from "vue";
import { useScroll } from "@/hooks/useScroll";
import NavBack from "@/components/business/NavBack.vue";
import { http } from "@/utils";
import SendMsg from "@/views/model/chat-playground/components/SendMsg.vue";
import Message from "@/views/model/chat-playground/components/Message/index.vue";
import ConfirmByClick from "@/components/business/ConfirmByClick.vue";
import ModelSelect from "@/components/business/ModelSelect.vue";
import { useRoute } from "vue-router";

const { scrollRef, scrollToBottom, scrollToBottomIfAtBottom } = useScroll();
const route = useRoute();

const { assistantId } = route.query;

// theme breadcrumb
const page = ref({ title: "助手操场" });

const sendParams = reactive({
  toolIds: [],
  instructions: "",
  modelName: "",
  stream: true
});
const chatList = ref([]);
const question = ref("");
const sendLoading = ref(false);
const refConfirmDelete = ref();
const assistantData = ref<Record<string, any>>({});
const systemMessage = ref("");

const handleTextSend = () => {
  const message = question.value;
  question.value = "";
  generateAnswer("text", message);
};

const getChatSendParams = sendMessage => {
  const result = { ...sendParams };
  const sendHistoryCount = 6 * 2; // *2 的目的是历史消息要成对发送
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
      url: `/assistants/${assistantId}/playground`,
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
              contentType: data.contentType || "text",
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

const getAssistantData = async () => {
  let [err, res] = await http.get({
    showLoading: true,
    url: `/assistants/${assistantId}`
  });
  if (res) {
    assistantData.value = res;
    sendParams.instructions = res.instructions;
    sendParams.modelName = res.modelName;
    sendParams.toolIds = assistantData.value.tools.map(_ => _.toolId);
  }
};

onMounted(() => {
  getAssistantData();
});
</script>
<style lang="scss" scoped>
.left-part {
  width: 350px;
  overflow-y: auto;
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
  flex: 1;
}
</style>
