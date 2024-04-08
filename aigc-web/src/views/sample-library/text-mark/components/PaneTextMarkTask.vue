<template>
  <Pane ref="refPane" class="" @submit="onSubmit">
    <div class="mx-auto mt-10" style="width: 500px">
      <v-form ref="refForm" class="my-form">
        <v-text-field
          type="text"
          placeholder="请输入负责人"
          hide-details="auto"
          clearable
          :rules="rules.principal"
          v-model="formData.principal"
        >
          <template #prepend><label class="required">负责人</label></template>
        </v-text-field>

        <Select
          v-model="formData.datasetId"
          :mapAPI="{
            url: '/api/mgr/datasets/list',
            data: { page: 1, pageSize: 20 },
            labelField: 'name',
            valueField: 'uuid',
            search_keywordField: 'name'
          }"
          @update:infos="onChangeDataSetId"
          :rules="rules.datasetId"
        >
          <template #prepend>
            <label class="required">样本名称</label>
          </template>
        </Select>
        <div ref="refRange">
          <RangeSliderWithInput
            :rules="rules.dataSequence"
            v-model="formData.dataSequence"
            :min="rangeConfig.min"
            :max="rangeConfig.max"
          >
            <template #prepend><label class="required">数据序列</label></template>
          </RangeSliderWithInput>
          <div class="text-gray-400 text-sm mt-2 ml-[44px]">数据的随机10%会成为测试集，其余部分会作为训练集</div>
        </div>

        <Select
          placeholder="请选择任务类型"
          :rules="rules.annotationType"
          :mapDictionary="{ code: 'textannotation_type' }"
          v-model="formData.annotationType"
        >
          <template #prepend>
            <label class="required">任务类型 <Explain>标注完成后，会有两个数据集：同名测试集和训练集</Explain></label>
          </template>
        </Select>

        <v-text-field
          type="text"
          placeholder="请输入任务名称"
          hide-details="auto"
          clearable
          :rules="rules.name"
          v-model="formData.name"
        >
          <template #prepend><label class="required">任务名称</label></template>
        </v-text-field>
      </v-form>
    </div>
  </Pane>
</template>
<script setup lang="ts">
import { reactive, toRefs, ref, nextTick } from "vue";
import _ from "lodash";
import { http, validator, format, doAnimation } from "@/utils";
import { toast } from "vue3-toastify";
import { useRouter, useRoute } from "vue-router";
import Explain from "@/components/ui/Explain.vue";
import RangeSliderWithInput from "@/components/ui/RangeSliderWithInput.vue";
import { ItfTextMarkTask } from "../types";
const state = reactive<{
  formData: ItfTextMarkTask;
  [x: string]: any;
}>({
  formData: {
    name: "",
    datasetId: "",
    principal: "",
    dataSequence: [],
    annotationType: ""
  },
  rangeConfig: {
    min: 0,
    max: 0
  }
});

const { formData, rangeConfig } = toRefs(state);
const emits = defineEmits(["submit"]);

const refPane = ref();
const refForm = ref();
const refRange = ref<HTMLElement>();
const rules = reactive({
  name: [v => !!v || "请输入任务名称"],
  datasetId: [v => !!v || "请选择样本"],
  principal: [v => !!v || "请输入负责人"],
  annotationType: [v => !!v || "请选择任务类型"],
  dataSequence: [
    () => {
      let [min, max] = state.formData.dataSequence;
      if (max > min) {
        return true;
      } else {
        return "数据序列上限值需大于下限值";
      }
    }
  ]
});

const onChangeDataSetId = ({ rawData }) => {
  let { rangeConfig } = state;
  if (rawData) {
    doAnimation.scaleIn(refRange.value);
    rangeConfig.max = rawData.segmentCount;
  } else {
    doAnimation.scaleIn(refRange.value);
    rangeConfig.max = 0;
  }
  state.formData.dataSequence = [rangeConfig.min, rangeConfig.max];
};

const doAdd = async (options = {}) => {
  const [err, res] = await http.post({
    ...options,
    url: `/api/mgr/annotation/task/create`,
    data: {
      ...state.formData
    }
  });
  if (res) {
    refPane.value.hide();
    emits("submit");
  }
};

const onSubmit = async ({ valid, showLoading }) => {
  if (valid) {
    doAdd({ showLoading });
  }
};

defineExpose({
  show({
    title,
    infos = {
      name: "",
      datasetId: "",
      principal: "",
      dataSequence: [0, 0],
      annotationType: ""
    } as ItfTextMarkTask,
    operateType
  }) {
    refPane.value.show({
      title,
      refForm
    });
    if (operateType == "add") {
      //添加
      state.formData = _.cloneDeep(infos);
    } else {
      // state.formData = _.pick(_.cloneDeep(infos), ["name", "datasetId", "principal", "annotationType"]);
      //编辑
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
