<template>
  <Pane ref="refPane" class="" @submit="onSubmit">
    <div class="mx-auto mt-10" style="width: 500px">
      <!-- formData :{{ formData }} -->
      <v-form ref="refForm" class="my-form">
        <!-- <v-text-field
          density="compact"
          variant="outlined"
          type="text"
          placeholder="请输入名称"
          hide-details="auto"
          class="required"
          :rules="rules.name"
          v-model="formData.name"
          clearable
          :disabled="state.disabledField"
        >
          <template #prepend> <label>名称</label></template>
        </v-text-field> -->
        <v-text-field
          density="compact"
          variant="outlined"
          type="text"
          placeholder="请输入别名"
          hide-details="auto"
          :rules="rules.alias"
          v-model="formData.alias"
          clearable
        >
          <template #prepend> <label class="required">别名</label></template>
        </v-text-field>

        <v-text-field
          density="compact"
          variant="outlined"
          type="number"
          placeholder="请输入配额"
          hide-details="auto"
          :rules="rules.quota"
          v-model.number="formData.quota"
          clearable
        >
          <template #prepend>
            <label class="required">配额 <Explain>每天最多请求的次数</Explain></label>
          </template>
        </v-text-field>

        <Select
          :rules="rules.modelId"
          v-model="formData.modelId"
          :mapAPI="{ url: '/api/channels/models', data: { pageSize: -1 }, labelField: 'modelName', valueField: 'id' }"
          hide-details="auto"
          multiple
        >
          <template #prepend>
            <label class="required">支持模型 <Explain>指该应用场景可以使用的模型</Explain></label>
          </template>
        </Select>

        <v-text-field
          density="compact"
          variant="outlined"
          type="text"
          placeholder="请输入负责人邮箱"
          hide-details="auto"
          clearable
          v-model="formData.email"
          :rules="rules.email"
        >
          <template #prepend>
            <label class="required">邮箱 <Explain>该应用场景的负责人邮箱</Explain></label>
          </template>
        </v-text-field>

        <!--

          <Select v-model="formData.projectName" :mapAPI="{ url: '/api/channels', data: { pageSize: -1 } }">
            <template #prepend>
              <label>项目 <Explain>提示AAA</Explain></label>
            </template>
          </Select>

          <Select v-model="formData.serviceName" :mapAPI="{ url: '/api/channels', data: { pageSize: -1 } }">
          <template #prepend>
            <label>服务 <Explain>提示AAA</Explain></label>
          </template>
        </Select>
      -->

        <v-textarea v-model.trim="formData.remark" placeholder="请输入">
          <template #prepend> <label>备注</label></template>
        </v-textarea>
      </v-form>
    </div>
  </Pane>
</template>
<script setup>
import { reactive, toRefs, ref, nextTick } from "vue";
import _ from "lodash";
import Explain from "@/components/ui/Explain.vue";
import { http, validator } from "@/utils";
const state = reactive({
  operateType: "", //add 添加  、 edit 编辑
  disabledField: false,
  maxTokens: 4096,
  formData: {
    name: "",
    alias: "",
    quota: "",
    modelId: [],
    email: "",
    projectName: "",
    serviceName: "",
    modelId: "",
    remark: ""
  }
});
const { formData } = toRefs(state);

const emits = defineEmits(["submit"]);

const refPane = ref();
const refForm = ref();
const rules = reactive({
  name: [v => !!v || "请输入名称"],
  alias: [v => !!v || "请输入别名"],
  modelId: [
    value => {
      if (value && value.length > 0) {
        return true;
      } else {
        return "请选择模型";
      }
    }
  ],
  quota: [v => !!v || "请输入配额"],
  email: [
    value => {
      return validator.isEmail({ value, required: true });
    }
  ]
});

const doAdd = async (options = {}) => {
  const [err, res] = await http.post({
    ...options,
    showSuccess: true,
    url: `/api/channels`,
    data: {
      ...state.formData
    }
  });
  if (res) {
    refPane.value.hide();
    emits("submit");
  }
};

const doEdit = async (options = {}) => {
  const [err, res] = await http.put({
    ...options,
    showSuccess: true,
    url: `/api/channels/${state.formData.id}`,
    data: {
      ...state.formData
    }
  });
  if (res) {
    refPane.value.hide();
    emits("submit");
  }
};

const onSubmit = ({ valid, showLoading }) => {
  if (valid) {
    if (state.operateType == "add") {
      doAdd({ showLoading });
    } else {
      doEdit({ showLoading });
    }
  }
};

defineExpose({
  show({
    title,
    infos = {
      id: "",
      name: "",
      alias: "",
      quota: "",
      email: "",
      projectName: "",
      serviceName: "",
      remark: "",
      modelId: []
    },
    operateType
  }) {
    refPane.value.show({
      title,
      refForm
    });
    state.formData = _.pick(_.cloneDeep(infos), [
      "id",
      "name",
      "alias",
      "quota",
      "email",
      "projectName",
      "serviceName",
      "remark"
    ]);

    state.operateType = operateType;
    if (operateType == "add") {
      //添加
      state.disabledField = false;
    } else {
      //编辑
      state.disabledField = true;
      state.formData.modelId = infos.model.list.map(item => {
        return item.id;
      });
    }
  }
});
</script>
<style lang="scss" scoped>
label {
  width: 100px;
  text-align: right;
  .compo-explain {
    margin-top: 1px;
  }
}
</style>
