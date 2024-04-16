<template>
  <div ref="terminalRef" class="terminal"></div>
</template>
<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from "vue";
import SockJS from "sockjs-client/dist/sockjs.min.js";
import { Terminal } from "@xterm/xterm";
import { FitAddon } from "@xterm/addon-fit";
import { CanvasAddon } from "@xterm/addon-canvas";
import { toast } from "vue3-toastify";
import "@xterm/xterm/css/xterm.css";

interface IProps {
  wsUrl: string;
  startData: Record<string, any>;
}

interface IEmits {
  (e: "save:cmd", val: string): void;
}

const props = withDefaults(defineProps<IProps>(), {
  wsUrl: "",
  startData: () => ({})
});
const emit = defineEmits<IEmits>();

const terminalRef = ref();
const socket = ref(null);
const terminal = ref(null);

const initSocket = () => {
  // 连接路径
  if (props.wsUrl == "") {
    return;
  }
  console.log(props.wsUrl);
  socket.value = new SockJS(props.wsUrl);
  socketOnOpen();
  socketOnMessage();
  socketOnClose();
  socketOnError();
};

const socketOnOpen = () => {
  socket.value.onopen = () => {
    console.log("web链接成功");
    const startData = {
      op: "bind",
      sessionId: props.startData.sessionId,
      data: JSON.stringify({
        container: props.startData.container,
        sessionId: props.startData.sessionId,
        serviceName: props.startData.serviceName
      })
    };
    // 发送格式与后台保持一致要不发送也是失败的
    socket.value.send(JSON.stringify(startData));
    initTerminal();
  };
};

const initTerminal = () => {
  terminal.value = new Terminal({
    disableStdin: false, // 是否应禁用输入
    windowsMode: true, // 根据窗口换行
    cursorBlink: true, // 光标闪烁
    cursorStyle: "underline", // 光标样式
    theme: {
      foreground: "#ececec", // 字体
      background: "#000" // 背景色
    }
  });
  const element = terminalRef.value;
  const fitAddon = new FitAddon(); // 全屏插件
  terminal.value.loadAddon(fitAddon);
  terminal.value.loadAddon(new CanvasAddon());
  terminal.value.open(element);
  fitAddon.fit();
  terminal.value.focus();
  terminal.value.onData(data => {
    socket.value.send(
      JSON.stringify({
        op: "stdin",
        cols: terminal.value.cols,
        rows: terminal.value.rows,
        data: data,
        sessionId: props.startData.sessionId
      })
    );
  });
};

const socketOnMessage = () => {
  socket.value.onmessage = evt => {
    let msg = JSON.parse(evt.data);
    try {
      if (msg["op"] === "stdout") {
        if (msg["data"].toString().indexOf("executable file not found in $PATH: unknown") === -1) {
          if (msg["data"] === "") {
            socket.value.send(
              JSON.stringify({
                op: "resize",
                cols: terminal.value.cols,
                rows: terminal.value.rows
              })
            );
          } else {
            terminal.value.write(msg["data"]);
          }
        }
      } else if (msg["op"] === "toast") {
        terminal.value.write(msg["data"]);
      } else {
        console.error("Unexpected message type:", msg);
      }
    } catch (e) {
      console.log("parse json error.", evt.data);
    }
  };
};

const socketOnClose = () => {
  socket.value.onclose = () => {
    // socket.value.close();
    // toast.error("关闭 socket, 请重新刷新页面连接socket");
    console.log("关闭 socket");
  };
};
const socketOnError = () => {
  socket.value.onerror = () => {
    toast.error("socket 链接失败");
    console.log("socket 链接失败");
  };
};

onMounted(() => {
  initSocket();
});

onBeforeUnmount(() => {
  socket.value && socket.value.close();
});
</script>
<style scoped lang="scss">
.terminal {
  touch-action: none;
}
</style>
