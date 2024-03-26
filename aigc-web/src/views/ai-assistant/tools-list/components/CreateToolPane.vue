<template>
  <Pane ref="refPane" @submit="onSubmit">
    <div class="mx-auto mt-3" style="width: 750px">
      <v-form ref="refForm" class="my-form">
        <v-text-field
          type="text"
          placeholder="只允许字母、数字、“-” 、“_”"
          hide-details="auto"
          clearable
          :rules="rules.name"
          v-model="formData.name"
        >
          <template #prepend
            ><label class="required">工具名称 <Explain>工具名称，通常为函数名</Explain></label></template
          >
        </v-text-field>
        <v-textarea
          v-model.trim="formData.description"
          :rules="rules.description"
          hide-details="auto"
          placeholder="请输入工具描述"
        >
          <template #prepend>
            <label class="required"
              >工具描述
              <Explain>
                <div style="max-width: 365px">
                  工具的描述，给个案例<br />例：<br />对于生成一张图片很有帮助。 输入应该是相要生成的图片的描述。输入参考:
                  {"prompt": "图片描述", "model": "dall-e-3"} 结果里的url是图片地址。
                </div>
              </Explain></label
            ></template
          >
        </v-textarea>
        <Select
          placeholder="请选择工具类型"
          :rules="rules.toolType"
          :mapDictionary="{ code: 'assistant_tool_type' }"
          v-model="formData.toolType"
        >
          <template #prepend>
            <label class="required">工具类型</label>
          </template>
        </Select>
        <template v-if="formData.toolType === 'function'">
          <v-text-field
            type="text"
            placeholder="请输入url"
            hide-details="auto"
            clearable
            :rules="rules.url"
            v-model="fnMetadata.url"
          >
            <template #prepend><label class="required">http地址</label></template>
          </v-text-field>
          <Select
            placeholder="请选择请求方式"
            :rules="rules.method"
            :mapDictionary="{ code: 'http_method' }"
            v-model="fnMetadata.method"
          >
            <template #prepend>
              <label class="required">method</label>
            </template>
          </Select>
          <v-input hide-details="auto" v-model="fnMetadata.body" :center-affix="false">
            <CodeMirror v-model="fnMetadata.body" placeholder="请输入" />
            <template #prepend> <label>body</label></template>
          </v-input>
          <!-- <v-checkbox hide-details="auto" label="Header透传" color="primary" v-model="fnMetadata.headerPass">
            <template #prepend>
              <label>Header</label>
            </template>
          </v-checkbox> -->
          <v-input hide-details class="header-input mb-0" :center-affix="false">
            <v-row v-for="(item, index) in fnMetadata.headerValues" class="form-item-inner w-100 ml-n5 mb-4" no-gutters>
              <v-col cols="5">
                <v-text-field type="text" hide-details="auto" placeholder="请输入Key" clearable v-model="item.key">
                  <template #prepend><label>Key</label></template>
                </v-text-field>
              </v-col>
              <v-col cols="5">
                <v-text-field hide-details="auto" type="text" placeholder="请输入Value" clearable v-model="item.value">
                  <template #prepend><label>Value</label></template>
                </v-text-field>
              </v-col>
              <v-col cols="2" class="mt-1">
                <v-btn v-if="index === 0" class="ml-3" icon flat color="info" size="x-small"
                  ><IconPlus stroke-width="1.5" :size="20" @click="addHeaderItem"
                /></v-btn>
                <v-btn v-else icon flat class="ml-3" color="error" size="x-small" @click="removeHeaderItem(index)"
                  ><IconMinus stroke-width="1.5" :size="20"
                /></v-btn>
              </v-col>
            </v-row>
            <template #prepend> <label class="mt-n1">header</label></template>
          </v-input>
        </template>
        <template v-else-if="formData.toolType === 'code_interpreter'">
          <Select
            placeholder="请选择编程语言"
            :rules="rules.language"
            :mapDictionary="{ code: 'programming_language' }"
            v-model="codeMetadata.language"
          >
            <template #prepend>
              <label class="required">编程语言</label>
            </template>
          </Select>
          <v-input hide-details="auto" :rules="rules.code" v-model="codeMetadata.code" :center-affix="false">
            <CodeMirror v-model="codeMetadata.code" :language="codeMetadata.language" placeholder="请输入" />
            <template #prepend> <label class="required">代码</label></template>
          </v-input>
        </template>
        <v-textarea v-model.trim="formData.remark" placeholder="请输入备注">
          <template #prepend> <label>备注</label></template>
        </v-textarea>
      </v-form>
    </div>
  </Pane>
</template>
<script setup lang="ts">
import { reactive, ref, computed } from "vue";
import Explain from "@/components/ui/Explain.vue";
import { http } from "@/utils";
import { IconPlus, IconMinus } from "@tabler/icons-vue";

interface IFormData {
  toolId?: string;
  name: string;
  description: string;
  toolType: string | null;
  metadata: string;
  remark: string;
}

const initFormData = {
  name: "",
  description: "",
  toolType: null,
  metadata: "",
  remark: ""
};

const emits = defineEmits(["submit"]);

const paneConfig = reactive({
  operateType: "add"
});
const formData = ref<IFormData>({ ...initFormData });
const fnMetadata = reactive({
  url: "",
  method: null,
  body: "",
  headerValues: [
    {
      key: "",
      value: ""
    }
  ]
});
const codeMetadata = reactive({
  language: null,
  code: ""
});
const refPane = ref();
const refForm = ref();
const rules = reactive({
  name: [v => /^[a-zA-Z0-9-_]+$/.test(v) || "只允许字母、数字、“-” 、“_”"],
  description: [v => !!v || "请输入工具描述"],
  toolType: [v => !!v || "请选择工具类型"],
  url: [v => !!v || "请输入url"],
  method: [v => !!v || "请选择请求方式"],
  language: [v => !!v || "请选择编程语言"],
  code: [v => !!v || "请输入代码"],
  content: [v => !!v || "请输入脚本模版"]
});

const isEdit = computed(() => {
  return paneConfig.operateType === "edit";
});

const addHeaderItem = () => {
  fnMetadata.headerValues.push({
    key: "",
    value: ""
  });
};

const removeHeaderItem = index => {
  fnMetadata.headerValues.splice(index, 1);
};

const getMetadata = () => {
  let result = null;
  if (formData.value.toolType === "function") {
    const header = fnMetadata.headerValues.reduce((obj, item) => {
      if (item.key && item.value) {
        obj[item.key] = item.value;
      }
      return obj;
    }, {});
    result = {
      url: fnMetadata.url,
      method: fnMetadata.method,
      body: fnMetadata.body,
      header
    };
  } else if (formData.value.toolType === "code_interpreter") {
    result = codeMetadata;
  } else {
    result = {};
  }
  return JSON.stringify(result);
};

const onSubmit = async ({ valid, showLoading }) => {
  if (valid) {
    const requestConfig = {
      url: "",
      method: ""
    };
    if (paneConfig.operateType == "add") {
      requestConfig.url = "/tools/create";
      requestConfig.method = "post";
    } else {
      requestConfig.url = `/tools/${formData.value.toolId}`;
      requestConfig.method = "put";
    }
    const metadata = getMetadata();
    const [err, res] = await http[requestConfig.method]({
      showLoading,
      showSuccess: true,
      url: requestConfig.url,
      data: { ...formData.value, metadata }
    });

    if (res) {
      refPane.value.hide();
      emits("submit");
    }
  }
};

const setMetadata = data => {
  try {
    const toolType = data.toolType;
    const json = JSON.parse(data.metadata);
    if (toolType === "function") {
      const headerValues = Object.keys(json.header).map(key => {
        return { key: key, value: json.header[key] };
      });
      fnMetadata.url = json.url;
      fnMetadata.method = json.method;
      fnMetadata.body = json.body;
      fnMetadata.headerValues =
        headerValues.length === 0
          ? [
              {
                key: "",
                value: ""
              }
            ]
          : headerValues;
    } else if (toolType === "code_interpreter") {
      codeMetadata.language = json.language;
      codeMetadata.code = json.code;
    }
  } catch (error) {
    console.log(error);
  }
};

const init = () => {
  formData.value = { ...initFormData };
  fnMetadata.url = "";
  fnMetadata.method = null;
  (fnMetadata.body = ""),
    (fnMetadata.headerValues = [
      {
        key: "",
        value: ""
      }
    ]);
  codeMetadata.language = null;
  codeMetadata.code = "";
};

defineExpose({
  show({ title, operateType, infos }) {
    refPane.value.show({
      title,
      refForm
    });
    paneConfig.operateType = operateType;
    if (paneConfig.operateType === "add") {
      init();
    } else {
      setMetadata(infos);
      formData.value = { ...infos };
    }
  }
});
</script>
<style lang="scss" scoped>
label {
  width: 120px;
  text-align: right;
  .compo-explain {
    margin-top: 1px;
  }
}
.form-item-inner {
  label {
    width: 50px;
  }
}
.header-input :deep(.v-input__control) {
  flex-wrap: wrap;
}
</style>
