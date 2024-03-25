import { createRouter, createWebHistory, createWebHashHistory } from "vue-router";
import MainRoutes from "./MainRoutes";
import AuthRoutes from "./AuthRoutes";

import { useUserStore } from "@/stores";
import { toast } from "vue3-toastify";
import { useMapRemoteStore } from "@/stores";
import NProgress from "./nprogress";

export const router = createRouter({
  // history: createWebHistory(import.meta.env.BASE_URL),
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/:pathMatch(.*)*",
      component: () => import("@/views/authentication/Error.vue")
    },
    MainRoutes,
    AuthRoutes
  ]
});

router.beforeEach(async (to, from, next) => {
  const mapRemoteStore = useMapRemoteStore();

  // 1.NProgress 开始
  NProgress.start();

  //2. 访问页面，无需登录验证
  if (!to.matched.some(record => record.meta.requiresAuth)) {
    return next();
  }

  const user = useUserStore();
  let token = user.userInfo.token;
  //3. 未登录，跳转登录
  if (!token) {
    toast.error("会话超时，请重新登录");
    return next("/auth/login?next=" + encodeURIComponent(to.fullPath));
  }

  //3. mapping 本地字典
  mapRemoteStore.mergeLocalData();

  //4. 正常访问页面
  next();
});

/**
 * @description 路由跳转错误
 * */
router.onError(error => {
  NProgress.done();
  console.warn("路由错误", error.message);
});

/**
 * @description 路由跳转结束
 * */
router.afterEach(() => {
  NProgress.done();
});
