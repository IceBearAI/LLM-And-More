// import { useClipboard } from "@vueuse/core";
// 如果剪贴板API不可用，则设置legacy:true以保留复制功能。它将使用execCommand作为回退来处理复制。(https://vueuse.org/core/useClipboard/#useclipboard)
// const { isSupported, copy } = useClipboard({ legacy: true });

import useClipboard from "vue-clipboard3";
import { toast } from "vue3-toastify";

const { toClipboard } = useClipboard();

export const clipboard = {
  mounted(el, binding, vnode, prevVnode) {
    vnode.key = binding.value + parseInt(Math.random() * 1000);
    el.style.cursor = "copy";
    el.addEventListener("click", function (e) {
      let arr = (binding.value + "").split("の");
      // copy(arr[0]);
      toClipboard(arr[0]);
      if (arr[1] != "") {
        toast.success(arr[1] ? arr[1] : `已复制：${arr[0]}`, {
          icon: false,
          autoClose: 2000
        });
      }
      e.stopPropagation(); //阻止事件冒泡
    });
  }
};
