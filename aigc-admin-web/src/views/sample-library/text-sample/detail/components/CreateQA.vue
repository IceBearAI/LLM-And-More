<template>
  <Pane ref="refPane" @submit="onSubmit">
    <div class="d-flex flex-column h-100">
      <div class="flex-1 position-relative">
        <template v-if="qaList.length > 0">
          <perfect-scrollbar ref="scrollRef" class="h-100">
            <div class="qa-item" v-for="(item, index) in qaList">
              <Message
                :msg-type="item.role === 'user' ? 'question' : 'answer'"
                :text="item.content"
                @click:edit="handleClickEdit(item)"
                @click:delete="handleClickDelete(index, item)"
              >
                <template v-if="item.isEdit" v-slot:text>
                  <v-textarea
                    autofocus
                    style="width: 350px"
                    hide-details
                    no-resize
                    auto-grow
                    :rows="1"
                    :max-rows="4"
                    v-model="item.content"
                    @update:focused="handleUpdateFocused($event, item)"
                  ></v-textarea>
                </template>
              </Message>
            </div>
          </perfect-scrollbar>
        </template>
        <template v-else>
          <div class="no-data text-medium-emphasis">暂无问答数据</div>
        </template>
      </div>
      <div class="footer">
        <UiParentCard>
          <v-form ref="refForm" class="my-form">
            <v-textarea
              v-model="question"
              hide-details="auto"
              density="default"
              placeholder="用户问"
              no-resize
              auto-grow
              :rows="1"
              :max-rows="4"
              :rules="rules.question"
              validate-on="submit"
              @keypress="handleEnter"
            >
              <template #prepend> <label class="required">问</label></template>
            </v-textarea>
            <v-textarea
              class="mb-0"
              v-model="answer"
              hide-details="auto"
              density="default"
              placeholder="客服回答"
              no-resize
              auto-grow
              :rows="1"
              :max-rows="4"
              :rules="rules.answer"
              validate-on="submit"
              @keypress="handleEnter"
            >
              <template #prepend> <label class="required">答</label></template>
              <template #append-inner> <v-btn color="primary" size="small" flat @click="addQA">追加</v-btn></template>
            </v-textarea>
          </v-form>
        </UiParentCard>
      </div>
    </div>
  </Pane>
</template>
<script setup lang="ts">
import { reactive, ref } from "vue";
import { http } from "@/utils";
import { useScroll } from "@/hooks/useScroll";
import Message from "./Message.vue";
import { toast } from "vue3-toastify";
import UiParentCard from "@/components/shared/UiParentCard.vue";

const props = defineProps({
  uuid: String
});
const emits = defineEmits(["submit"]);

const { scrollRef, scrollToBottom } = useScroll();

const qaList = ref([]);
const paneConfig = reactive({
  operateType: "add",
  datasetSampleId: ""
});
const question = ref("");
const answer = ref("");
const refPane = ref();
const refForm = ref();
const rules = reactive({
  question: [v => !!v || "请输入用户问"],
  answer: [v => !!v || "请输入客服回答"]
});
const currentEditContent = ref("");

const handleEnter = event => {
  if (event.key === "Enter" && !event.shiftKey) {
    event.preventDefault();
    addQA();
  }
};

const addQA = async () => {
  const { valid } = await refForm.value.validate();
  if (valid) {
    qaList.value.push(
      {
        isEdit: false,
        role: "user",
        content: question.value
      },
      {
        isEdit: false,
        role: "assistant",
        content: answer.value
      }
    );
    scrollToBottom();
    question.value = "";
    answer.value = "";
  }
};

const handleClickEdit = item => {
  currentEditContent.value = item.content;
  item.isEdit = true;
};

const handleUpdateFocused = (isFocused, item) => {
  if (!isFocused) {
    if (!item.content) {
      item.content = currentEditContent.value;
    }
    item.isEdit = false;
  }
};

const handleClickDelete = (index, item) => {
  // 删除问答一队
  if (item.role === "user") {
    qaList.value.splice(index, 2);
  } else {
    qaList.value.splice(index - 1, 2);
  }
};

const onSubmit = async ({ showLoading }) => {
  if (qaList.value.length === 0) {
    toast.warning("问答数据不能为空");
    return;
  }
  const requestConfig = {
    url: "",
    method: ""
  };
  if (paneConfig.operateType == "add") {
    requestConfig.url = `/datasets/${props.uuid}/samples`;
    requestConfig.method = "post";
  } else {
    requestConfig.url = `/api/datasets/${props.uuid}/samples/${paneConfig.datasetSampleId}`;
    requestConfig.method = "put";
  }
  const [err, res] = await http[requestConfig.method]({
    showLoading,
    showSuccess: true,
    url: requestConfig.url,
    data: {
      messages: qaList.value
    }
  });

  if (res) {
    refPane.value.hide();
    emits("submit");
  }
};

defineExpose({
  show({ title, operateType, infos }) {
    refPane.value.show({
      title,
      confirmText: "save"
    });
    paneConfig.operateType = operateType;
    if (paneConfig.operateType === "add") {
      question.value = "";
      answer.value = "";
      qaList.value = [];
    } else {
      paneConfig.datasetSampleId = infos.uuid;
      qaList.value = infos.messages;
    }
  }
});
</script>
<style lang="scss" scoped>
label {
  width: 50px;
  text-align: right;
  .compo-explain {
    margin-top: 1px;
  }
}
.no-data {
  position: absolute;
  left: 50%;
  top: 50%;
  transform: translate(-50%, -50%);
}
:deep(.v-field--no-label) {
  --v-field-padding-top: 2px;
}
</style>
