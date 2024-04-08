<template>
  <Pane ref="refPane" @submit="onSubmit">
    <div class="mx-auto mt-3" style="width: 540px">
      <v-form ref="refForm" class="my-form">
        <v-textarea v-model="formData.question" :rules="rules.question" placeholder="用户问" hide-details="auto" :rows="3">
          <template #prepend> <label class="required">标准问</label></template>
        </v-textarea>
        <v-textarea v-model="formData.output" :rules="rules.output" placeholder="客服回答" hide-details="auto" :rows="3">
          <template #prepend> <label class="required">回答</label></template>
        </v-textarea>
        <v-textarea v-model="formData.intent" :rules="rules.intent" placeholder="意图" hide-details="auto" :rows="3">
          <template #prepend> <label class="required">意图</label></template>
        </v-textarea>
        <!-- <template v-for="item in formData.SimilarityQuestion">
          <v-textarea v-model="item.content" placeholder="相似问法" hide-details="auto" :rows="3">
            <template #prepend> <label>相似问法</label></template>
          </v-textarea>
        </template>
        <v-input>
          <v-btn class="w-100 border-dashed" @click="addSimilarityQuestion" flat color="info" variant="outlined"
            ><IconPlus stroke-width="1.5" :size="18" class="mr-1" />添加</v-btn
          >
          <template #prepend> <label></label></template>
        </v-input> -->
      </v-form>
    </div>
  </Pane>
</template>
<script setup lang="ts">
import { reactive, ref } from "vue";
import { http } from "@/utils";
// import { IconPlus } from "@tabler/icons-vue";

interface IProps {
  intentId: string;
}
interface IEmits {
  (e: "submit"): void;
}

const props = withDefaults(defineProps<IProps>(), {
  intentId: ""
});
const emits = defineEmits<IEmits>();

const paneConfig = reactive({
  operateType: "add"
});
const formData = reactive({
  documentId: "",
  question: "",
  output: "",
  intent: ""
  // SimilarityQuestion: [
  //   {
  //     content: ""
  //   }
  // ]
});

const refPane = ref();
const refForm = ref();
const rules = reactive({
  question: [v => !!v || "请输入标准问"],
  output: [v => !!v || "请输入回答"],
  intent: [v => !!v || "请输入意图"]
});

// const addSimilarityQuestion = () => {
//   formData.SimilarityQuestion.push({
//     content: ""
//   });
// };

// const removeSimilarityQuestion = index => {
//   formData.SimilarityQuestion.splice(index, 1);
// };

const onSubmit = async ({ valid, showLoading }) => {
  if (valid) {
    const requestConfig = {
      url: "",
      method: ""
    };
    if (paneConfig.operateType == "add") {
      requestConfig.url = `/intent/${props.intentId}/document/create`;
      requestConfig.method = "post";
    } else {
      requestConfig.url = `/api/intent/${props.intentId}/document/${formData.documentId}/update`;
      requestConfig.method = "put";
    }
    const [err, res] = await http[requestConfig.method]({
      showLoading,
      showSuccess: true,
      url: requestConfig.url,
      data: formData
    });

    if (res) {
      refPane.value.hide();
      emits("submit");
    }
  }
};

defineExpose({
  show({ title, operateType, infos }) {
    refPane.value.show({
      title,
      refForm
    });
    paneConfig.operateType = operateType;
    if (paneConfig.operateType === "add") {
      formData.question = "";
      formData.output = "";
      formData.intent = "";
    } else {
      formData.question = infos.question;
      formData.output = infos.output;
      formData.intent = infos.intent;
      formData.documentId = infos.documentId;
    }
  }
});
</script>
<style lang="scss" scoped>
label {
  width: 130px;
  text-align: right;
  .compo-explain {
    margin-top: 1px;
  }
}
</style>
