import { defineStore } from "pinia";
import { http, dataDictionary } from "@/utils";
import _ from "lodash";
import { toast } from "vue3-toastify";

type TypeOption = {
  label: string;
  value: string | number | boolean;
};

type TypeGetLabels = [any, any];

type TypeState = {
  /**
   * options 存放的是数组格式，用于构建下拉选项
   * 如：options.gender = [ {label:'男',value:1} , {label:'女',value:2} ]
   */
  options: {
    [key: string]: TypeOption[];
  };
  /**
   * mappings 存放的是对象格式，用于字段转义
   * 如：mapping.gender = {1:'男',2:'女'}
   */
  mappings: any;
  /**
   * 正在拉取的字典code
   *   多次调用loadDictTree方法，重复拉取同个code，只拉取一次，可减少 /api/sys/dict/tree 接口调用次数
   */
  loadingCodes: {
    [key: string]: 1;
  };
};

export const useMapRemoteStore = defineStore({
  id: "mapRemote",
  getters: {},
  state: (): TypeState => ({
    mappings: {},
    options: {},
    loadingCodes: {}
  }),
  actions: {
    /**
     * 合成本地数据字典数据
     */
    mergeLocalData() {
      let { mappings, options } = this;
      let { localData } = dataDictionary;

      // itemField 取值如 gender
      for (let itemField of Object.keys(localData)) {
        let itemType = localData[itemField];
        mappings[itemField] = itemType;

        options[itemField] = [];
        for (let itemValue of Object.keys(itemType)) {
          let value = itemValue;
          try {
            value = JSON.parse(itemValue);
          } catch (e) {}

          options[itemField].push({
            label: itemType[itemValue],
            value
          });
        }
      }
    },

    /**
     * 获取多个字段的label
     * @param inValue
     * @param cb  回调函数，传入时覆盖默认回调函数
     * @returns
     * 比如 getLabels([ ['gender','1'] ] ) =>  ret =['男']
     *     getLabels([ ['gender','1'],['gender','2'] ] ) =>  ret = ['男','女']
     */
    getLabels(inValue: Array<TypeGetLabels>, cb?: (a: string[]) => string) {
      let { mappings } = this;
      let ret = [];

      inValue.forEach(async item => {
        let fieldName = item[0];
        let fieldValue = item[1];
        if (["", null, undefined].includes(fieldValue)) {
          //无效值
          return;
        }
        let cn = mappings[fieldName]?.[fieldValue];
        if (cn) {
          ret.push(cn);
        } else {
          // console.error(`数据字典里，无法转义 ${fieldName}、${fieldValue}`);
        }
      });
      if (cb) {
        return cb(ret);
      } else {
        //默认回调
        if (ret.length) {
          return ret.join("");
        } else {
          return "未知";
        }
      }
    },
    /**
     * 查询数据字典
     *   页面使用 mappings 转义code，但页面里没用到对应code 的<Select  控件时，需要调用此方法
     */
    async loadDictTree(code: string | string[]) {
      let { options, mappings, loadingCodes } = this;
      let arrCode = [];
      if (typeof code == "string") {
        arrCode.push(code);
      } else {
        arrCode = code;
      }
      let queryCode = arrCode.filter(item => {
        if (!options[item] && !loadingCodes[item]) {
          loadingCodes[item] = 1;
          //过滤出没有缓存的code
          return true;
        }
      });
      if (queryCode.length == 0) {
        //没有要查询的code
        if (Object.keys(loadingCodes).length == 0) {
          return;
        } else {
          //仍有在拉取的code
          return new Promise(resolve => {
            let startTime = +new Date();
            let timer = setInterval(() => {
              if (Object.keys(loadingCodes).length == 0) {
                //数据已取回
                clearInterval(timer);
                resolve("");
              } else if (+new Date() - startTime > 5000) {
                //5秒后，不再等待了
                clearInterval(timer);
                resolve("");
              }
            }, 100);
          });
        }
      }
      const [err, res] = await http.get({
        url: `/api/sys/dict/tree`,
        data: {
          code: queryCode
        }
      });
      if (res && res.length) {
        res.forEach(itemType => {
          let code = itemType.code; //如 gender
          delete loadingCodes[code];
          if (!mappings[code]) {
            mappings[code] = {};
          }
          if (!options[code]) {
            options[code] = [];
          }
          (itemType.children || []).forEach(item => {
            let { dictLabel: label, dictValue: value } = item;
            mappings[code][value] = label;
            options[code].push({
              label,
              value
            });
          });
        });
      } else {
        //接口失败了，清除对应code的拉取中的标识
        for (let item of queryCode) {
          delete loadingCodes[item];
        }
      }
    },
    async getOptionsByCode(code: string): Promise<TypeOption[]> {
      await this.loadDictTree(code);
      return this.options[code];
    }
  }
});
