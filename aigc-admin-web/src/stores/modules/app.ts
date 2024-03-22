import { defineStore } from "pinia";
// element
import eZh from "element-plus/es/locale/lang/zh-cn";
import eEn from "element-plus/es/locale/lang/en";
import eEs from "element-plus/es/locale/lang/es";

// moment
import "moment/dist/locale/zh-cn";
import "moment/dist/locale/es"; // 西班牙语
import "moment/dist/locale/fil"; // 菲律宾语

const languageMap = {
  zh: {
    element: eZh,
    moment: "zh-cn"
  },
  en: {
    element: eEn,
    moment: "en"
  },
  es: {
    element: eEs,
    moment: "es"
  },
  ph: {
    element: eEn,
    moment: "fil"
  }
};

export const useAppStore = defineStore({
  id: "app",
  state: () => ({
    isBtnLoading: false,
    appName: import.meta.env.VITE_GLOB_APP_TITLE,
    btnId: "",
    localLanguage: "zh",
    methods: {
      //存放一些全局方法，（vue中定义，js中用到的），比如 i18n的t方法
      /** 国际化方法 */
      t: (key: string): string => {
        return "";
      }
    }
  }),
  getters: {
    localeElement() {
      return languageMap[this.localLanguage]["element"];
    },
    localeMoment() {
      return languageMap[this.localLanguage]["moment"];
    }
  },
  actions: {
    setBtnLoading(newStatus, btnId = "") {
      this.isBtnLoading = newStatus;
      if (newStatus) {
        this.btnId = btnId;
      } else {
        this.btnId = "";
      }
    },
    setLocalLanguage(payload: any) {
      this.localLanguage = payload;
    }
  },
  persist: {
    storage: localStorage,
    paths: ["localLanguage"]
  }
});
