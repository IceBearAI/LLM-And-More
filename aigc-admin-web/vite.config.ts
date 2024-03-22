import { fileURLToPath, URL } from "url";
import { defineConfig, loadEnv } from "vite";
import { wrapperEnv } from "./build/getEnv";
import { createVitePlugins } from "./build/plugins";
import vue from "@vitejs/plugin-vue";
import vuetify from "vite-plugin-vuetify";
import moment from "moment";

// 版本号
const appVersion = moment().format("YYYY-MM-DD HH:mm:ss");
/**
 * 接口请求域名，构建时传入。未指定时，接口域名为当前服务域名；指定以后接口域名为指定域名。
 *  构建命令如下
 *   apiOrigin=https://www.test.com  pnpm build
 */
const { apiOrigin } = process.env;

// @see: https://vitejs.dev/config/
export default defineConfig(({ mode }) => {
  const root = process.cwd();
  const env = loadEnv(mode, root);
  const viteEnv = wrapperEnv(env);

  return {
    plugins: createVitePlugins(viteEnv),
    resolve: {
      alias: {
        "@": fileURLToPath(new URL("./src", import.meta.url))
      }
    },
    css: {
      preprocessorOptions: {
        scss: {}
      }
    },
    define: {
      appVersion: JSON.stringify(appVersion),
      apiOrigin: JSON.stringify(apiOrigin)
    },
    server: {
      logLevel: "info",
      port: 5173,
      proxy: {
        "/api": {
          // target: "", //接口代理地址
          changeOrigin: true
        }
      }
    },
    optimizeDeps: {
      exclude: ["vuetify"],
      entries: ["./src/**/*.vue"]
    },
    build: {
      rollupOptions: {
        output: {
          manualChunks(id) {
            if (id.includes("/src/utils/http.ts") || id.includes("/src/utils/index.ts")) {
              // 将 'http.ts' 和 'index.ts' 强制打包到同一个 chunk 中
              return "utils";
            }
          }
        }
      }
    }
  };
});
