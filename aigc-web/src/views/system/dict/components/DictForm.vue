<template>
  <v-form ref="refForm" class="my-form">
    <v-text-field
      type="text"
      placeholder="请输入字典编号"
      hide-details="auto"
      clearable
      :rules="rules.code"
      v-model="formData.code"
      :disabled="disabledCode"
    >
      <template #prepend><label class="required">字典编号</label></template>
    </v-text-field>
    <v-text-field
      type="text"
      placeholder="请输入字典名称"
      hide-details="auto"
      clearable
      :rules="rules.dictLabel"
      v-model="formData.dictLabel"
    >
      <template #prepend><label class="required">字典名称</label></template>
    </v-text-field>
    <template v-if="parentId !== 0">
      <template v-if="formData.dictType === 'json'">
        <v-input v-model="formData.dictValue" :rules="rules.dictValue" :center-affix="false" validate-on="submit">
          <CodeMirror v-model="formData.dictValue" placeholder="请输入" />
          <template #prepend><label class="required">字典键值</label></template>
        </v-input>
      </template>
      <template v-else>
        <v-text-field
          type="text"
          placeholder="请输入字典键值"
          hide-details="auto"
          clearable
          :rules="rules.dictValue"
          v-model="formData.dictValue"
        >
          <template #prepend><label class="required">字典键值</label></template>
        </v-text-field>
      </template>
    </template>
    <v-text-field
      type="number"
      placeholder="请输入字典排序"
      hide-details="auto"
      :rules="rules.sort"
      v-model.number="formData.sort"
    >
      <template #prepend> <label class="required">字典排序</label></template>
    </v-text-field>
    <Select
      placeholder="请选择值类型"
      :rules="rules.dictType"
      :mapDictionary="{ code: 'sys_dict_type' }"
      v-model="formData.dictType"
    >
      <template #prepend>
        <label class="required">值类型</label>
      </template>
    </Select>
    <v-textarea v-model.trim="formData.remark" placeholder="请输入字典备注" clearable>
      <template #prepend> <label>字典备注</label></template>
    </v-textarea>
  </v-form>
</template>
<script setup lang="ts">
import { ref, reactive, computed } from "vue";
import { validator } from "@/utils";

interface IProps {
  type?: string;
  parentId?: number;
}

interface IFormData {
  id?: number;
  parentId: number;
  code: string;
  dictLabel: string;
  dictType: string | null;
  dictValue?: string;
  sort: number;
  remark: string;
}
const initFormData = {
  parentId: 0,
  code: "",
  dictLabel: "",
  dictValue: "",
  dictType: null,
  sort: 0,
  remark: ""
};

const props = withDefaults(defineProps<IProps>(), {
  type: "add",
  parentId: 0
});

const formData = ref<IFormData>({ ...initFormData });
const refForm = ref();
const rules = reactive({
  code: [v => !!v || "请输入字典编号"],
  dictLabel: [v => !!v || "请输入字典名称"],
  dictValue: [v => validateDictValue(v)],
  sort: [v => v !== "" || "请输入字典排序"],
  dictType: [v => !!v || "请选择值类型"]
});

const validateDictValue = val => {
  if (!val) {
    return "请输入字典键值";
  } else {
    if (formData.value.dictType === "json") {
      if (!validator.isJson(val)) {
        return "JSON格式不正确";
      }
    }
    return true;
  }
};

const isEdit = computed(() => {
  return props.type === "edit";
});
const disabledCode = computed(() => {
  return isEdit.value || (props.type === "add" && props.parentId !== 0);
});

const reset = (ext = {}) => {
  formData.value = { ...initFormData, ...ext };
};

const setFormData = infos => {
  formData.value = { ...infos };
};

defineExpose({
  reset,
  setFormData,
  getFormData() {
    return { ...formData.value, parentId: props.parentId };
  },
  getRef() {
    return refForm.value;
  }
});
</script>
<style lang="scss" scoped>
label {
  width: 80px;
  text-align: right;
  .compo-explain {
    margin-top: 1px;
  }
}
</style>
