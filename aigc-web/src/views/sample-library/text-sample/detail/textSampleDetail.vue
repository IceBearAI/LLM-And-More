<template>
  <NavBack backUrl="/sample-library/text-sample/list">样本详情</NavBack>
  <UiParentCard class="mt-4">
    <template #header> 名称：{{ data.name }}</template>
    <template #action>
      <v-btn color="primary" @click="onExport">导出数据集</v-btn>
    </template>
    <DetailBaseInfo :info="data" />
  </UiParentCard>
  <UiParentCard class="mt-5">
    <DetailTableInfo />
  </UiParentCard>
  <ConfirmByClick ref="refConfirmByClick" @submit="onConfirmByClick">
    <template #text>
      <div>
        <AlertBlock class="mb-6"
          >该操作{{ state.subTitle }}
          <!-- <span>，导出成功之后，在创建微调时可以直接选择该数据集</span> -->
        </AlertBlock>
        <v-form ref="refForm" class="my-form">
          <v-text-field type="text" placeholder="eg: 普泽AI客服" hide-details="auto" clearable>
            <template #prepend>
              <label style="display: flex"
                >角色
                <div class="toolTip">
                  <IconInfoCircle
                    :size="16"
                    color="#bbb"
                    pointer
                    class="iconTool"
                    @mouseenter="state.showIcon = true"
                    @mouseleave="state.showIcon = false"
                  >
                  </IconInfoCircle>
                  <div v-show="state.showIcon" class="iconToolTip">AI角色</div>
                </div>
              </label>
            </template>
          </v-text-field>
          <v-text-field type="text" placeholder="eg: 普泽基金的研究人员" hide-details="auto" clearable>
            <template #prepend>
              <label style="display: flex"
                >作者
                <div class="toolTip">
                  <IconInfoCircle
                    :size="16"
                    color="#bbb"
                    pointer
                    class="iconTool"
                    @mouseenter="state.showIconRole = true"
                    @mouseleave="state.showIconRole = false"
                  >
                  </IconInfoCircle>
                  <div v-show="state.showIconRole" class="iconToolTipRole">AI的制作人及使用的数据等</div>
                </div>
              </label>
            </template>
          </v-text-field>
        </v-form>
        <v-chip-group
          mandatory
          selected-class="text-primary"
          style="margin-left: 20%"
          v-model="state.radioCheck"
          @click="changeChip"
        >
          <v-chip v-for="(tag, index) in state.exportList" :key="index" filter variant="outlined">
            {{ tag.text }}
          </v-chip>
        </v-chip-group>
      </div>
    </template>
  </ConfirmByClick>
</template>
<script setup lang="ts">
import { reactive, ref, onMounted } from "vue";
import NavBack from "@/components/business/NavBack.vue";
import { http } from "@/utils";
import { useRoute } from "vue-router";
import UiParentCard from "@/components/shared/UiParentCard.vue";
import DetailBaseInfo from "./components/DetailBaseInfo.vue";
import DetailTableInfo from "./components/DetailTableInfo.vue";
import ConfirmByClick from "@/components/business/ConfirmByClick.vue";
import AlertBlock from "@/components/ui/AlertBlock.vue";
import { IconInfoCircle } from "@tabler/icons-vue";

const state = reactive<{
  /** 所选行，编辑操作 */
  // selectedInfo: Partial<ItfModel>;
  [x: string]: any;
}>({
  formData: {
    name: ""
  },
  radioCheck: 0,
  showIcon: false,
  showIconRole: false,
  subTitle: "",
  exportList: [
    {
      text: "导出jsonl",
      minTitle: "jsonl",
      subTitle: "将导出成alpaca的训练格式文件，可在创建微调时上传该文件进行训练。"
    },
    // {
    //   text: "导出csv",
    //   minTitle: "csv",
    //   subTitle: "将导出成csv格式的文件"
    // },
    {
      text: "生成微调数据集",
      minTitle: "data",
      subTitle: "将直接生成微调数据集，在创建微调时可以直接选择该数据集"
    }
  ]
});
const route = useRoute();
const refConfirmByClick = ref();
const { uuid } = route.query;
const data = ref<Record<string, any>>({});

const getData = async () => {
  let [err, res] = await http.get({
    showLoading: true,
    url: `/api/datasets/${uuid}`
  });
  if (res) {
    data.value = res;
  }
};
const onExport = () => {
  refConfirmByClick.value.show({ width: "550px" });
};
const onConfirmByClick = (options = {}) => {
  state.subTitle = state.exportList[state.radioCheck].subTitle;
  // let { action, row } = state.confirmByClickInfo;
  // if (action == "deploy") {
  //   onDeploy(row, options);
  // } else if (action == "undeploy") {
  //   onUndeploy(row, options);
  // }
};

const changeChip = (options = {}) => {
  state.subTitle = state.exportList[state.radioCheck].subTitle;
  // let { action, row } = state.confirmByClickInfo;
  // if (action == "deploy") {
  //   onDeploy(row, options);
  // } else if (action == "undeploy") {
  //   onUndeploy(row, options);
  // }
};
onMounted(() => {
  state.subTitle = state.exportList[state.radioCheck].subTitle;
  getData();
});
</script>
<style lang="scss">
.export-title {
  color: #333;
  font-weight: 700 !important;
}
.export-sub {
  color: #fb9678;
  font-weight: 600 !important;
}

.v-sheet {
  color: #888;
}
.compo-explain {
  // display: inline-flex;
  // vertical-align: text-top;
  // svg {
  //   outline: none;
  // }
}
.toolTip {
  position: relative;
  width: 40px;
  .iconTool {
    position: absolute;
    top: 3px;
    left: 2px;
  }
  .iconToolTip,
  .iconToolTipRole {
    background: #333;
    color: #fff;
    border-radius: 4px;
    font-size: 11px;
    padding: 4px 6px;
    z-index: 12;
  }
  .iconToolTip {
    position: absolute;
    top: -26px;
    left: -14px;
  }
  .iconToolTipRole {
    position: absolute;
    top: -25px;
    left: -52px;
    width: 150px;
  }
}
</style>
