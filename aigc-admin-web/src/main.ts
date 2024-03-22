import { createApp } from "vue";
import { setupStore } from "@/stores";
import App from "./App.vue";
import { router } from "./router";
import { setupPlugins } from "@/plugins";
import "@/scss/style.scss";
import PerfectScrollbar from "vue3-perfect-scrollbar";
import VueApexCharts from "vue3-apexcharts";
import Vue3Toasity, { toast } from "vue3-toastify";

// element plus
import ElementPlus from "element-plus";

// 一定要在main.ts中导入tailwind.css，防止vite每次hmr都会请求src/scss/style.scss整体css文件导致热更新慢的问题
import "@/scss/tailwind.scss";

// Table
import Vue3EasyDataTable from "vue3-easy-data-table";
import "vue3-easy-data-table/dist/style.css";
import "vue3-toastify/dist/index.css";

//ScrollTop
import VueScrollTo from "vue-scrollto";

const app = createApp(App);
app.use(router);
app.component("EasyDataTable", Vue3EasyDataTable);
app.use(PerfectScrollbar);
// app.use(createPinia());

// 挂载pina状态管理
setupStore(app);

app.use(VueApexCharts);
app.use(ElementPlus);
app.use(Vue3Toasity, {
  autoClose: 3000,
  transition: "flip",
  theme: "colored",
  position: "top-center",
  hideProgressBar: true, //隐藏关闭进度条
  pauseOnFocusLoss: false, //页面失焦后，是否暂停
  clearOnUrlChange: false //url变化时，是否立马关闭
});

//扩展vue功能：添加全局组件、指令等
setupPlugins(app);

app.mount("#app");
//ScrollTop Use
// app.use(VueScrollTo);

app.use(VueScrollTo, {
  duration: 300,
  easing: "ease"
});

window.appVersion = appVersion;
window.env = import.meta.env.MODE;

window.errorMsg = function (message: string) {
  console.error(message);
  if (window.env == "development") {
    toast.error(message, { position: "top-right", theme: "light", transition: "bounce" });
  }
};
