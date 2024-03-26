import { createPinia } from "pinia";
import piniaPluginPersistedstate from "pinia-plugin-persistedstate";

import { useAppStore } from "./modules/app";
import { useUserStore } from "./modules/user";
import { useMapRemoteStore } from "./modules/mapRemote";

const store = createPinia();
store.use(piniaPluginPersistedstate);

export function setupStore(app) {
  app.use(store);
}

export { useAppStore, useUserStore, useMapRemoteStore };
