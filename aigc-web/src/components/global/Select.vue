<template>
  <v-autocomplete
    v-if="state.isReady"
    class="compo-select"
    density="compact"
    v-model="data"
    :items="state.cellOptions"
    item-title="label"
    item-value="value"
    @update:modelValue="onChange"
    @update:search="onSearch"
    :item-props="true"
    variant="outlined"
    hide-details="auto"
    clearable
    :return-object="false"
    :placeholder="placeholder"
  >
    <template #prepend v-if="$slots.prepend">
      <slot name="prepend" />
    </template>
  </v-autocomplete>

  <v-autocomplete v-else class="compo-select" density="compact" variant="outlined" hide-details="auto" :return-object="false">
    <template #prepend v-if="$slots.prepend">
      <slot name="prepend" />
    </template>
    <template #loader>
      <v-progress-circular indeterminate color="primary" :size="20" :width="2"></v-progress-circular>
    </template>
  </v-autocomplete>
</template>
<script setup lang1="ts">
import { reactive, toRefs, ref, computed, watch, onMounted } from "vue";
import { useVModel, useDebounceFn } from "@vueuse/core";
import { dataDictionary } from "@/utils";
import _ from "lodash";
import { useMapRemoteStore } from "@/stores";
import { useI18n } from "vue-i18n";
import { useAppStore } from "@/stores";

const appStore = useAppStore();

const { t } = useI18n(); // 解构出t方法

const props = defineProps({
  modelValue: {
    //所选项 value 的值
    type: [String, Boolean, Number]
  },
  infos: {
    //所选项 整个的值 ，附加 index 字段
    type: Object
  },
  defaultFirst: {
    //默认选中第一项
    type: Boolean,
    default: false
  },
  options: {
    //下拉选项，手动传入
    type: Array,
    default() {
      return [];
    }
  },
  placeholder: String,
  mapAPI: {
    //下拉选项，通过接口拉取
    type: Object,
    default() {
      return {
        url: "", //api 接口url地址
        data: {}, //api 接口入参
        labelField: "", //给人看的字段名
        valueField: "", //给程序传值的字段名
        search_keywordField: "", //使用关键词搜索，关键词对应的字段名。未指定时静态搜索，指定了动态搜索。
        responsePath: "" // 取值路径
      };
    }
  },
  mapDictionary: {
    /**
     * 下拉选项，通过数据字典获得
     */
    type: Object,
    default() {
      return {
        code: "", //code 传入 'gender' 之类的
        i18nKey: "" //多语言匹配 options 下的key
      };
    }
  }
  //获取本地数据字典（dataDictionary）中已配置的数据
  // mapLocal: String
});

const emits = defineEmits(["update:modelValue", "update:infos", "change"]);
let data = useVModel(props, "modelValue", emits);

const state = reactive({
  isReady: false,
  cellOptions: []
});

const getOptionItemLabel = ({ code, itemValue, defaultLabel }) => {
  let i18nKey = props.mapDictionary.i18nKey || code;
  let key = `options.${i18nKey}.${itemValue}`;
  let ret = t(key);
  if (ret == key) {
    return defaultLabel;
  } else {
    return ret;
  }
};

const onSearch = useDebounceFn(keywordValue => {
  if (!lockSearch) {
    let { mapAPI } = props;
    if (mapAPI.url && mapAPI.search_keywordField) {
      //使用关键词，重新搜索下拉数据
      props.mapAPI.data[mapAPI.search_keywordField] = keywordValue;
      initOptions();
    }
  }
}, 500);

watch(
  () => appStore.localLanguage,
  () => {
    initOptions();
  }
);

const initOptions = async () => {
  if (props.options.length) {
    state.cellOptions = _.cloneDeep(props.options);
  } else if (props.mapAPI?.url) {
    state.cellOptions = await dataDictionary.getOptionsByAPI(props.mapAPI);
  } else if (props.mapDictionary?.code) {
    const mapRemoteStore = useMapRemoteStore();
    let options = await mapRemoteStore.getOptionsByCode(props.mapDictionary.code);
    state.cellOptions = options.map(item => {
      return {
        value: item.value,
        label: getOptionItemLabel({
          code: props.mapDictionary.code,
          itemValue: item.value,
          defaultLabel: item.label
        }),
        rawData: item
      };
      // return item;
    });
  }
  //空处理
  state.cellOptions = state.cellOptions || [];

  if (props.defaultFirst && !props.modelValue && props.modelValue !== 0) {
    //无值时候，才默认选中第一项
    let defaultValue;
    if (props.multiple) {
      //多选
      defaultValue = [state.cellOptions[0].value];
    } else {
      //单选
      defaultValue = state.cellOptions[0].value;
    }
    data.value = defaultValue;
    onChange(defaultValue);
  }
};

/**
 * 获得选项信息
 * @param {*} optionValue
 */
const getOptionInfo = optionValue => {
  let { cellOptions } = state;
  let index = cellOptions.findIndex(itemOption => {
    return itemOption.value == optionValue;
  });
  let matched = cellOptions[index];
  if (matched) {
    //有匹配项
    return {
      ...matched,
      index
    };
  } else {
    //未匹配
    return {
      label: "",
      value: null,
      index
    };
  }
};

//引入 lockSearch相关逻辑，选中某项以后，依然可保留搜索结果
let lockSearch = false;
let timerLockSearch;
const onChange = selectedValue => {
  lockSearch = true;
  clearTimeout(timerLockSearch);
  timerLockSearch = setTimeout(() => {
    lockSearch = false;
  }, 800);
  let infos;
  if (selectedValue instanceof Array) {
    //多选
    infos = [];
    selectedValue.forEach(itemValue => {
      let optionInfo = getOptionInfo(itemValue);
      infos.push(optionInfo);
    });
  } else {
    infos = getOptionInfo(selectedValue);
  }
  emits("update:infos", infos);
  emits("change", selectedValue, infos);
};

onMounted(async () => {
  await initOptions();
  let initValue = props.modelValue;
  state.isReady = true;
  if (initValue) {
    onChange(initValue);
  }
});
</script>
<style lang="scss">
.compo-select {
  .v-field__overlay {
    background: #fff;
  }
  .v-field__input {
    min-height: 40px;
    padding: 7px 16px;
  }
  .v-field__loader {
    top: 10px;
    left: 10px;
    bottom: 0;
  }
}
</style>
