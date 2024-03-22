<template>
  <Pane ref="refPane" class="" @submit="onSubmit">
    <div class="mx-auto mt-10" style="width: 500px">
      <v-form ref="refForm" class="my-form">
        <v-input
          :rules="rules.sampleFileId"
          v-model="formData.sampleFileId"
          hide-details="auto"
          style="position: relative; width: 100%"
          v-if="state.operateType == 'upload'"
        >
          <template v-if="headSampleInfos">
            <v-chip closable color="info" @click:close="headSampleClose">{{ headSampleInfos.filename }}</v-chip>
          </template>
          <template v-else>
            <v-btn color="info" variant="outlined" prepend-icon="mdi-tray-arrow-up " :disabled="state.upLoading"
              >{{ state.upLoading ? "上传中..." : "上传样本" }}
              <UploadFile
                show-loading
                v-model="formData.sampleFileId"
                v-model:infos="headSampleInfos"
                purpose="fine-tune"
                @upload:success="doQueryFirstPage"
                @loading="val => (state.upLoading = val)"
                style="width: 146px; position: absolute; top: 0; left: -31%; opacity: 0"
              />
            </v-btn>
          </template>
          <template #prepend> <label class="required">样本文件</label></template>
        </v-input>
        <!-- 支持扩展名：.jsonl -->
        <v-text-field
          density="compact"
          variant="outlined"
          type="text"
          placeholder="请输入中文、数字、字母、-、_ "
          hide-details="auto"
          clearable
          :rules="rules.name"
          v-model="formData.name"
        >
          <!-- :disabled="state.disabledField" -->
          <template #prepend> <label class="required">别名</label></template>
        </v-text-field>
        <v-textarea v-model.trim="formData.remark" placeholder="请输入">
          <template #prepend> <label>备注</label></template>
        </v-textarea>
        <div class="down-sample" v-if="state.operateType == 'upload'">
          点击
          <div class="down" @click="onDownload">下载</div>
          数据集模版
        </div>
      </v-form>
    </div>
  </Pane>
</template>
<script setup lang="ts">
import { reactive, toRefs, ref, nextTick } from "vue";
import _ from "lodash";
import UploadFile from "@/components/business/UploadFile.vue";
import { http, validator, format } from "@/utils";
import { ItfModel } from "../textSampleList.ts";
import { toast } from "vue3-toastify";
import { useRouter, useRoute } from "vue-router";
const state = reactive<{
  // formData: ItfModel;
  /** 操作类型 add 添加  、 edit 编辑 ，默认add */
  operateType: "add" | "edit" | "upload";
  [x: string]: any;
}>({
  operateType: "add",
  disabledField: false,
  maxTokens: 500000,
  formData: {
    name: "",
    sampleFileId: "",
    remark: ""
  },
  upLoading: false,
  uuid: ""
});
const { formData } = toRefs(state);
const router = useRouter();
const emits = defineEmits(["submit"]);

const refPane = ref();
const refForm = ref();
const headSampleInfos = ref(null);
const rules = reactive({
  name: [
    // value => {
    //   if (value && value.length > 0) {
    //     return true;
    //   } else {
    //     return "请输入中文、数字、字母、-、_ ";
    //   }
    // }
    value => {
      return validator.isName({ value, required: true, errorValid: "请输入中文、数字、字母、-、_" });
    }
  ],
  sampleFileId: [
    value => {
      if (value && value.length > 0) {
        return true;
      } else {
        return "请上传样本文件";
      }
    }
  ]
});

const doAdd = async (options = {}) => {
  const [err, res] = await http.post({
    ...options,
    showSuccess: true,
    url: `/api/datasets/create`,
    data: {
      ...state.formData
    }
  });
  if (res) {
    refPane.value.hide();
    emits("submit");
  }
};
const onView = ({ id }) => {
  router.push(`/sample-library/text-sample/detail?jobId=${id}`);
};
const doEdit = async (options = {}) => {
  const [err, res] = await http.put({
    ...options,
    showSuccess: true,
    url: `/api/datasets/${state.uuid}`,
    data: {
      ...state.formData
    }
  });
  if (res) {
    refPane.value.hide();
    emits("submit");
    if (res && res.traceId.length > 0) {
      onView(res.traceId);
    }
  }
};

const onSubmit = ({ valid, showLoading }) => {
  if (valid) {
    if (state.operateType == "add") {
      //创建
      doAdd({ showLoading });
    } else if (state.operateType == "upload") {
      doAdd({ showLoading });
    } else {
      doEdit({ showLoading });
    }
  }
};
const onDownload = row => {
  if (headSampleInfos.value != null) {
    http.downloadByUrl({
      fileUrl: headSampleInfos.s3Url,
      suffixName: "jsonl"
    });
  } else {
    toast.warning("请先上传样本文件");
  }
};

defineExpose({
  show({
    title,
    infos = {
      uuid: "",
      name: "",
      remark: ""
    },
    operateType
  }) {
    refPane.value.show({
      title,
      refForm
    });
    state.formData = _.pick(_.cloneDeep(infos), ["name", "remark"]);
    state.operateType = operateType;
    state.uuid = infos.uuid;
    if (operateType == "add" || state.operateType == "upload") {
      //添加
      headSampleInfos.value = null;
      state.disabledField = false;
    } else {
      //编辑
      headSampleInfos.value = null;
      state.disabledField = true;
    }
  }
});
const headSampleClose = () => {
  formData.value.sampleFileId = "";
  headSampleInfos.value = null;
};
const doQueryFirstPage = () => {};
</script>
<style lang="scss" scoped>
label {
  width: 100px;
  text-align: right;
  .compo-explain {
    margin-top: 1px;
  }
}
.down-sample {
  display: flex;
  justify-content: center;
  margin-top: -20px;
  .down {
    color: #539bff;

    &:hover {
      cursor: pointer;
      text-decoration: underline;
    }
  }
}
</style>
