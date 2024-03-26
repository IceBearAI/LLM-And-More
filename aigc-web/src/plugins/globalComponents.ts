import type { App } from "vue";

import Pane from "@/components/global/Pane.vue";
import ButtonsInForm from "@/components/global/ButtonsInForm.vue";
import ButtonsInTable from "@/components/global/ButtonsInTable.vue";

import Table from "@/components/global/Table.vue";
import Pager from "@/components/global/Pager.vue";
import TableWithPager from "@/components/global/TableWithPager.vue";
import NoData from "@/components/global/NoData.vue";
import Select from "@/components/global/Select.vue";
import Dialog from "@/components/global/Dialog.vue";
import AiBtn from "@/components/global/AiBtn.vue";
import CodeMirror from "@/components/global/CodeMirror.vue";
import RefreshButton from "@/components/global/RefreshButton.vue";

/**
 * 全局注册组件
 * @param app
 */
export function setupGlobalComponents(app: App) {
  app.component("Pane", Pane);
  app.component("ButtonsInForm", ButtonsInForm);
  app.component("ButtonsInTable", ButtonsInTable);
  app.component("Select", Select);
  app.component("NoData", NoData);
  app.component("Table", Table);
  app.component("Pager", Pager);
  app.component("TableWithPager", TableWithPager);
  app.component("Dialog", Dialog);
  app.component("AiBtn", AiBtn);
  app.component("CodeMirror", CodeMirror);
  app.component("RefreshButton", RefreshButton);
}
