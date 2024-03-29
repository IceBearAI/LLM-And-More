import { App } from "vue";
import { setupGlobalComponents } from "./globalComponents";
import { setVuetify } from "./vuetify";
import directives from "@/directives/index";

export function setupPlugins(app: App) {
  // 注册全局组件,如：<svg-icon />
  setupGlobalComponents(app);
  app.use(directives);
  setVuetify(app);
}
