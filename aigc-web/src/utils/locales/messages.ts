import { zhHans as vuetifyZh, en as vuetifyEn, es as vueifyEs } from "vuetify/locale";
import en from "./en.json";
import zh from "./zh.json";
import es from "./es.json";
import ph from "./ph.json";

zh["$vuetify"] = vuetifyZh;
en["$vuetify"] = vuetifyEn;
es["$vuetify"] = vueifyEs;
ph["$vuetify"] = vuetifyEn; // 菲律宾语言组件库暂时使用英文

const messages = {
  en: en,
  zh: zh,
  es: es,
  ph: ph
};

export default messages;
