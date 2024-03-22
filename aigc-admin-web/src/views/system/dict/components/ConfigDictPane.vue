<template>
  <Pane ref="refPane" @close="paneClose">
    <v-row>
      <v-col cols="5">
        <UiParentCard>
          <el-tree
            ref="refTree"
            :data="treeData"
            node-key="id"
            default-expand-all
            :highlight-current="true"
            :expand-on-click-node="false"
            :props="{
              label: 'dictLabel'
            }"
            @node-click="nodeClick"
          >
            <template #default="{ node, data }">
              <span class="d-flex justify-space-between align-center w-100 pr-2 overflow-hidden">
                <span class="text-truncate">{{ node.label }}</span>
                <span class="ml-2" @click.stop>
                  <a class="link text-info" @click="append(data)">添加</a>
                  <a v-if="data.parentId !== 0" class="link text-error ml-2" @click="remove(data)">删除</a>
                </span>
              </span>
            </template>
          </el-tree>
        </UiParentCard>
      </v-col>
      <v-col cols="7">
        <UiParentCard class="dict-card" ref="refDictFormCard" :title="dictFormConfig.title">
          <template #action>
            <AiBtn class="ml-2" size="small" color="primary" variant="flat" @click="onSubmit">提交</AiBtn>
          </template>
          <DictForm ref="refDictForm" :type="dictFormConfig.type" :parentId="dictFormConfig.parentId" />
        </UiParentCard>
      </v-col>
    </v-row>
    <ConfirmByClick ref="refConfirmDelete" @submit="doRemove">
      <template #text>
        这是进行一项操作时必须了解的重要信息<br />
        您将要删除 <span class="text-primary font-weight-black">{{ confirmDelete.currentLabel }}</span> ，确定要继续吗？
      </template>
    </ConfirmByClick>
  </Pane>
</template>
<script setup lang="ts">
import { reactive, ref, nextTick } from "vue";
import { http } from "@/utils";
import UiParentCard from "@/components/shared/UiParentCard.vue";
import ConfirmByClick from "@/components/business/ConfirmByClick.vue";
import { toast } from "vue3-toastify";
import { animate } from "@/utils/animation";
import DictForm from "./DictForm.vue";

const emits = defineEmits(["refresh"]);

const treeData = ref([]);
const refPane = ref();
const refTree = ref();
const refDictFormCard = ref();
const refDictForm = ref();
const dictFormConfig = reactive({
  type: "edit",
  parentId: 0,
  title: "编辑"
});
const currentParentId = ref(0);
const currentNodeKey = ref(null);
const refConfirmDelete = ref();
const confirmDelete = reactive({
  id: null,
  currentLabel: ""
});
const topLevelNodeModify = ref(false);

const getDictTree = async () => {
  const [err, res] = await http.get({
    url: "/sys/dict",
    data: {
      parentId: currentParentId.value
    }
  });
  if (res) {
    treeData.value = res.list || [];
    if (!currentNodeKey.value) {
      // 默认选中第一项
      nodeClick(res.list[0]);
    }
    nextTick(() => {
      // 每次刷新页面后选中上一次
      setTreeHighlight();
    });
  }
};

const setTreeHighlight = () => {
  refTree.value.setCurrentKey(currentNodeKey.value);
};

const nodeClick = data => {
  dictFormConfig.title = `编辑（${data.dictLabel}）`;
  dictFormConfig.type = "edit";
  currentNodeKey.value = data.id;
  const infos = {
    id: data.id,
    parentId: data.parentId,
    code: data.code,
    dictValue: data.dictValue,
    dictLabel: data.dictLabel,
    dictType: data.dictType,
    sort: data.sort,
    remark: data.remark
  };
  dictFormConfig.parentId = data.parentId;
  refDictForm.value.setFormData(infos);
};

const onSubmit = async () => {
  const { valid } = await refDictForm.value.getRef().validate();
  if (valid) {
    const formData = refDictForm.value.getFormData();
    const requestConfig = {
      url: "",
      method: ""
    };
    if (dictFormConfig.type == "add") {
      requestConfig.url = "/sys/dict";
      requestConfig.method = "post";
    } else {
      requestConfig.url = `/sys/dict/${formData.id}`;
      requestConfig.method = "put";
    }
    const [err, res] = await http[requestConfig.method]({
      showLoading: refDictFormCard.value.$el,
      showSuccess: true,
      url: requestConfig.url,
      data: formData
    });

    if (res) {
      await getDictTree();
      if (dictFormConfig.type === "add") {
        refDictForm.value.reset({ code: formData.code });
      }
      if (formData.parentId === 0 && dictFormConfig.type === "edit") {
        topLevelNodeModify.value = true;
      }
    }
  } else {
    toast.warning("请处理页面标错的地方后，再尝试提交");
  }
};

const append = data => {
  dictFormConfig.title = `新增子项(${data.dictLabel})`;
  dictFormConfig.type = "add";
  dictFormConfig.parentId = data.id;
  currentNodeKey.value = data.id;
  setTreeHighlight();
  refDictForm.value.reset({ code: data.code });
  dictFormCardShake();
};

const remove = data => {
  confirmDelete.currentLabel = data.dictLabel;
  confirmDelete.id = data.id;
  refConfirmDelete.value.show({
    width: "450px",
    confirmText: confirmDelete.currentLabel
  });
};

const doRemove = async (options = {}) => {
  const [err, res] = await http.delete({
    ...options,
    showSuccess: true,
    url: `/sys/dict/${confirmDelete.id}`
  });
  if (res) {
    refConfirmDelete.value.hide();
    if (confirmDelete.id === currentNodeKey.value) {
      currentNodeKey.value = null;
    }
    getDictTree();
  }
};

const dictFormCardShake = () => {
  // 动画效果
  refDictFormCard.value &&
    animate(
      refDictFormCard.value.$el,
      [{ transformOrigin: "center" }, { transform: "scale(1.03)" }, { transformOrigin: "center" }],
      {
        duration: 150,
        easing: "cubic-bezier(0.4, 0, 0.2, 1)"
      }
    );
};

function paneClose() {
  currentNodeKey.value = null;
  if (topLevelNodeModify.value) {
    emits("refresh");
  }
}

defineExpose({
  show({ title, id }) {
    refPane.value.show({
      width: 900,
      showActions: false,
      title
    });
    currentParentId.value = id;
    topLevelNodeModify.value = false;
    getDictTree();
  }
});
</script>
<style lang="scss" scoped>
.dict-card :deep(.v-card-title) {
  flex: 1;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
