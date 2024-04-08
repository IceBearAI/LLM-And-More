import { resolve } from "path";
import { visualizer } from "rollup-plugin-visualizer";
import { createHtmlPlugin } from "vite-plugin-html";
import vue from "@vitejs/plugin-vue";
import vueJsx from "@vitejs/plugin-vue-jsx";
import eslintPlugin from "vite-plugin-eslint";
import viteCompression from "vite-plugin-compression";
import vuetify from "vite-plugin-vuetify";

/**
 * 创建 vite 插件
 * @param viteEnv
 */
export const createVitePlugins = viteEnv => {
  return [
    vue(),
    //全局默认引入scss文件
    vuetify({
      autoImport: true,
      styles: { configFile: "src/scss/variables.scss" }
    }),
    // vue 可以使用 jsx/tsx 语法
    vueJsx(),
    // esLint 报错信息显示在浏览器界面上
    eslintPlugin(),
    // 创建打包压缩配置
    createCompression(viteEnv),
    // 注入变量到 html 文件
    createHtmlPlugin({
      inject: {
        data: { title: viteEnv.VITE_GLOB_APP_TITLE }
      }
    }),
    // 是否生成包预览，分析依赖包大小做优化处理
    viteEnv.VITE_REPORT &&
      visualizer({
        emitFile: true, //是否被触摸
        filename: "stats.html",
        open: true, //在默认用户代理中打开生成的文件
        gzipSize: true, //从源代码中收集 gzip 大小并将其显示在图表中
        brotliSize: true //从源代码中收集 brotli 大小并将其显示在图表中
      })
  ];
};

/**
 * 根据 compress 配置，生成不同的压缩规则
 * @param viteEnv
 */
const createCompression = viteEnv => {
  const { VITE_BUILD_COMPRESS = "none", VITE_BUILD_COMPRESS_DELETE_ORIGIN_FILE } = viteEnv;
  const compressList = VITE_BUILD_COMPRESS.split(",");
  const plugins = [];

  if (compressList.includes("gzip")) {
    plugins.push(
      viteCompression({
        ext: ".gz",
        algorithm: "gzip",
        deleteOriginFile: VITE_BUILD_COMPRESS_DELETE_ORIGIN_FILE
      })
    );
  }
  if (compressList.includes("brotli")) {
    plugins.push(
      viteCompression({
        ext: ".br",
        algorithm: "brotliCompress",
        deleteOriginFile: VITE_BUILD_COMPRESS_DELETE_ORIGIN_FILE
      })
    );
  }

  return plugins;
};
