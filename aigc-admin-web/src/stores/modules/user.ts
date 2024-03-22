import { defineStore } from "pinia";
import { router } from "@/router";
import { http } from "@/utils";
import { unref } from "vue";
import { TypeTenant, TypeUserStore, TypeUserInfo } from "../types/stores.type.ts";

type loginResponse = {
  token: string;
  username: string;
};

export const useUserStore = defineStore({
  id: "user",
  state: (): TypeUserStore => ({
    userInfo: {
      roleId: "",
      token: "",
      tenantId: "", //商户Id
      username: "",
      tenants: []
    }
  }),
  actions: {
    async login(username: string, password: string) {
      const [err, res] = await http.post<loginResponse>({
        url: "/auth/login",
        data: { username, password },
        headers: {
          "X-Token": "",
          "X-Tenant-Id": ""
        }
      });
      if (res) {
        this.userInfo = {
          token: res.token,
          username: res.username
        };
        await this.getAccount();
        let next = unref(router.currentRoute).query.next || "/dashboards/index";
        if (Array.isArray(next)) {
          next = next[0];
        }
        router.push(next);
      }
    },
    async getAccount() {
      const [err, res] = await http.get<{ tenants: TypeTenant[] }>({
        url: "/auth/account",
        headers: {
          "X-Tenant-Id": ""
        }
      });
      this.userInfo.tenants = res.tenants;
      if (res.tenants?.length) {
        this.userInfo.tenantId = res.tenants[0].id;
      }
    },
    /** 退出登录 */
    logout() {
      this.userInfo = {
        token: "",
        tenantId: "", //商户Id
        username: "",
        tenants: []
      };
      location.replace("#auth/login");
      // router.push("/auth/login");
    }
  },
  persist: {
    storage: localStorage,
    paths: ["userInfo"]
  }
});
