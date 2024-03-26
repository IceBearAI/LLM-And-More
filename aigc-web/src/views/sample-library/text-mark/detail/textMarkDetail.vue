<template>
  <NavBack backUrl="backCallback" @back="onBack">文本标注</NavBack>
  <UiParentCard class="mt-4">
    <div>
      <v-row class="my-form waterfall">
        <v-col cols="12" md="6">
          <v-col>
            <v-input hide-details>
              <template #prepend> <label>任务ID</label></template>
              <span v-copy="state.annotationId" class="link">{{ state.annotationId }}</span>
            </v-input>
          </v-col>
          <v-col>
            <v-input hide-details>
              <template #prepend> <label>任务名称</label></template>
              {{ displayData.name }}
            </v-input>
          </v-col>
          <!-- <v-col>
            <v-input hide-details>
              <template #prepend> <label>任务类型</label></template>
              {{ getLabels([["textannotation_type", displayData.annotationType]]) }}
            </v-input>
          </v-col> -->
        </v-col>
        <v-col cols="12" md="6" class="text-xl text-slate-500 flex justify-end">
          <div class="v-box item-end">
            <v-col class="flex">
              <div class="pl-5">
                标注进度：<span v-show="state.displayData.total">
                  <span class="link font-bold">{{ displayData.completed + displayData.abandoned + 1 }}</span
                  >/{{ displayData.total }}
                </span>
              </div>

              <div class="pl-8">
                任务类型：<span class="link font-bold">{{
                  getLabels([["textannotation_type", displayData.annotationType]])
                }}</span>
              </div>
            </v-col>
            <v-col class="ml-auto flex flex-wrap divide-x">
              <div class="px-5">
                已完成：<span class="link font-bold">{{ displayData.completed }}</span>
              </div>
              <div class="px-5">
                未完成：<span class="link font-bold">{{ displayData.uncompleted }}</span>
              </div>
              <div class="px-5">
                丢弃：<span class="link font-bold">{{ displayData.abandoned }}</span>
              </div>
            </v-col>
          </div>
        </v-col>
      </v-row>
    </div>
    <v-divider class="mb-5"></v-divider>
    <v-row>
      <v-col cols="12" md="5">
        <v-card class="h-full flex p-4 overflow-hidden min-h-[300px] max-h-[665px] text-[15px]" elevation="3">
          <div class="flex-1 scrollbar-auto space-y-3 [line-height:1.4]" v-html="displayData.segmentContent"></div>
        </v-card>
      </v-col>
      <v-col cols="12" md="7">
        <v-form ref="refForm" class="my-form space-y-4">
          <v-textarea
            type="textarea"
            placeholder="请输入信息"
            hide-details="auto"
            :rules="rules.instruction"
            validate-on="submit"
            v-model="formData.instruction"
            class="block"
          >
            <template #prepend><label class="required">Instruction</label></template>
          </v-textarea>
          <v-textarea
            v-if="['rag'].includes(state.displayData.annotationType)"
            type="textarea"
            placeholder="请输入信息"
            hide-details="auto"
            :rules="rules.document"
            validate-on="submit"
            v-model="formData.document"
            class="block"
          >
            <template #prepend><label class="required">Document</label></template>
          </v-textarea>
          <v-textarea
            v-if="['faq', 'rag'].includes(state.displayData.annotationType)"
            type="text"
            placeholder="请输入信息"
            hide-details="auto"
            :rules="rules.question"
            validate-on="submit"
            v-model="formData.question"
            class="block"
          >
            <template #prepend><label class="required">Question</label></template>
          </v-textarea>
          <v-textarea
            v-if="['faq'].includes(state.displayData.annotationType)"
            type="text"
            placeholder="请输入信息"
            hide-details="auto"
            :rules="rules.intent"
            validate-on="submit"
            v-model="formData.intent"
            class="block"
          >
            <template #prepend><label class="required">Intent</label></template>
          </v-textarea>
          <v-textarea
            v-if="['general'].includes(state.displayData.annotationType)"
            type="text"
            placeholder="请输入信息"
            hide-details="auto"
            clearable
            :rules="rules.input"
            validate-on="submit"
            v-model="formData.input"
            class="block"
          >
            <template #prepend><label class="required">Input</label></template>
          </v-textarea>
          <v-textarea
            type="text"
            placeholder="请输入信息"
            hide-details="auto"
            clearable
            :rules="rules.output"
            validate-on="submit"
            v-model="formData.output"
            class="block"
          >
            <template #prepend><label class="required">Output</label></template>
          </v-textarea>
        </v-form>
      </v-col>
    </v-row>

    <div class="flex justify-between mt-5">
      <!-- <v-btn color="primary">返回未完成标注</v-btn> -->
      <div></div>
      <ButtonsInForm>
        <v-btn class="w-[120px] btn-danger" plain @click="onCancel" :disabled="!formData.uuid">丢弃</v-btn>
        <AiBtn color="secondary" class="w-[120px]" @click="onNext" :disabled="!formData.uuid">提交</AiBtn>
      </ButtonsInForm>
    </div>
  </UiParentCard>
</template>
<script setup lang="ts">
import { reactive, toRefs, ref, onMounted, inject } from "vue";
import NavBack from "@/components/business/NavBack.vue";
import UiParentCard from "@/components/shared/UiParentCard.vue";
import { http, format } from "@/utils";
import { useRouter, useRoute } from "vue-router";
import { useMapRemoteStore } from "@/stores";
import _ from "lodash";
import { toast } from "vue3-toastify";
import { useI18n } from "vue-i18n";
import $ from "jquery";
import { ItfAspectPageState } from "@/types/AspectPageTypes.ts";

const { t } = useI18n(); // 解构出t方法
const provideAspectPage = inject("provideAspectPage") as ItfAspectPageState;

const { getLabels, loadDictTree } = useMapRemoteStore();
loadDictTree(["textannotation_type"]);

const route = useRoute();
const router = useRouter();
const refForm = ref();

/** 是否修改过数据 */
let hasModifyed = false;

const state = reactive({
  style: {},
  annotationId: "",
  displayData: {
    name: "",
    annotationType: "",
    segmentContent: "",
    total: 0, //总量
    uncompleted: 0, //未完成
    completed: 0, //已完成
    abandoned: 0 //已废弃
  },
  formData: {
    uuid: "",
    instruction: "",
    document: "",
    question: "",
    intent: "",
    input: "",
    output: ""
  }
});
const { style, formData, displayData } = toRefs(state);

const rules = reactive({
  instruction: [v => !!v || "请输入信息"],
  question: [v => !!v || "请输入信息"],
  intent: [v => !!v || "请输入信息"],
  output: [v => !!v || "请输入信息"],
  document: [v => !!v || "请输入信息"],
  input: [v => !!v || "请输入信息"]
});

const getDetail = async () => {
  const [err, res] = await http.get({
    showLoading: true,
    url: `/api/mgr/annotation/task/${state.annotationId}/info`
  });
  if (res) {
    state.displayData = {
      ...state.displayData,
      ...res
    };
    state.displayData.uncompleted = res.total - res.completed - res.abandoned;
    $(document).scrollTop(0);
  }
};

const onBack = () => {
  if (hasModifyed) {
    provideAspectPage.methods.refreshListPage();
  }
  router.push("/sample-library/text-mark/list");
};

const getNext = async () => {
  refForm.value.resetValidation();
  state.formData = {
    uuid: "",
    instruction: "",
    document: "",
    question: "",
    intent: "",
    input: "",
    output: ""
  };
  const [err, res] = await http.get({
    showLoading: true,
    url: `/api/mgr/annotation/task/${state.annotationId}/segment/next`
  });
  if (res) {
    _.assign(state.formData, _.pick(res, ["uuid", "instruction", "question", "intent", "output"]));

    if (res.segmentContent) {
      state.displayData.segmentContent = "<div>" + res.segmentContent.split("\n").join("</div><div>") + "</div>";
    } else {
      state.displayData.segmentContent = "";
    }
  }
};

const afterSubmit = () => {
  hasModifyed = true;
  if (state.displayData.uncompleted == 1) {
    //已处理到最后一条数据，返回列表页
    onBack();
  } else {
    getDetail();
    getNext();
  }
};

const onCancel = async () => {
  const [err, res] = await http.put({
    showLoading: true,
    showSuccess: true,
    url: `/api/mgr/annotation/task/${state.annotationId}/segment/${state.formData.uuid}/abandoned`
  });
  if (res) {
    afterSubmit();
  }
};

const onNext = async () => {
  let { valid, errors } = await refForm.value.validate();
  if (valid) {
    const [err, res] = await http.post({
      showLoading: true,
      showSuccess: true,
      url: `/api/mgr/annotation/task/${state.annotationId}/segment/${state.formData.uuid}/mark`,
      data: state.formData
    });
    if (res) {
      afterSubmit();
    }
  } else {
    let errorMsg = t("pane.errorMsg");
    toast.warning(errorMsg);
  }
};

onMounted(() => {
  let { annotationId } = route.query;
  state.annotationId = annotationId as string;
  getDetail();
  getNext();
});
</script>
<style lang="scss"></style>
