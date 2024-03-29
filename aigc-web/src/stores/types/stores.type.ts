/**
 * 租户
 */
export type TypeTenant = {
  /** 租户id */
  id: string;
  /** 租户名称 */
  name: string;
};

export type TypeUserInfo = {
  /** token，请求headers 需携带 */
  token: string;
  /** 租户id，请求headers 需携带 */
  tenantId: string;
  /** 用户名 */
  username: string;
  /* 租户options列表 */
  tenants: TypeTenant[];
  /** 角色id */
  roleId: string;
};

export type TypeUserStore = {
  userInfo: TypeUserInfo;
};
