<template>
  <Pane ref="refPane" @submit="onSubmit">
    <div class="mx-auto mt-3" style="width: 540px">
      <DictForm ref="refDictForm" />
    </div>
  </Pane>
</template>
<script setup lang="ts">
import { reactive, ref, nextTick } from "vue";
import { http } from "@/utils";
import { toast } from "vue3-toastify";
import DictForm from "./DictForm.vue";

const emits = defineEmits(["submit"]);

const paneConfig = reactive({
  operateType: "add"
});
const refPane = ref();
const refDictForm = ref();

const onSubmit = async ({ showLoading }) => {
  const { valid } = await refDictForm.value.getRef().validate();
  if (valid) {
    const formData = refDictForm.value.getFormData();
    const requestConfig = {
      url: "",
      method: ""
    };
    if (paneConfig.operateType == "add") {
      requestConfig.url = "/sys/dict";
      requestConfig.method = "post";
    } else {
      requestConfig.url = `/sys/dict/${formData.id}`;
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
  } else {
    toast.warning("请处理页面标错的地方后，再尝试提交");
  }
};

defineExpose({
  show({ title, operateType, infos }) {
    refPane.value.show({
      title
    });
    paneConfig.operateType = operateType;
    nextTick(() => {
      if (paneConfig.operateType === "add") {
        refDictForm.value.reset();
      } else {
        refDictForm.value.setFormData(infos);
      }
    });
  }
});
</script>
