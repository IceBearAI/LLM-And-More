const AuthRoutes = {
  path: "/auth",
  component: () => import("@/layouts/blank/BlankLayout.vue"),
  meta: {
    requiresAuth: false
  },
  children: [
    {
      name: "Side Login", // 登录
      path: "/auth/login",
      component: () => import("@/views/authentication/SideLogin.vue")
    },
    {
      name: "Error", // 404页面
      path: "/auth/404",
      component: () => import("@/views/authentication/Error.vue")
    },
    {
      name: "Maintenance", // 维护页面
      path: "/auth/maintenance",
      component: () => import("@/views/authentication/Maintenance.vue")
    }
  ]
};

export default AuthRoutes;
