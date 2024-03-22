<template>
  <div v-show="state.isReady !== ''" style="min-height: 390px">
    <template v-if="state.isReady">
      <div class="w-100">
        <div class="d-flex tab-voice">
          <div class="flex-1 text-h2 overflow-hidden d-flex justify-end align-center mr-10" style="margin-left: -40px">
            <div ref="refRecordShowText" class="text-medium-emphasis">{{ state.statusText }}</div>
          </div>
          <div class="btn-record cursor-pointer" :class="getClassName4BtnRecord()" @click="onToggleRecord">
            <IconMicrophone class="icon-record d-block" :size="80" strokeWidth="2" />
          </div>
          <div class="flex-1"></div>
        </div>
      </div>
      <div class="shower-momentWords text-center" v-show="state.momentWords">
        <p class="line1">
          {{ state.momentWords }}
        </p>
      </div>
      <div v-show="state.videoUrl" class="h-center mt-10">
        <AiAudio ref="refAiAudio" :src="state.videoUrl" type="complex" />
      </div>
    </template>
    <div v-else></div>
  </div>
</template>
<script setup lang="ts">
import { reactive, toRefs, ref, nextTick, onMounted } from "vue";
import { animate, file } from "@/utils";
import { IconMicrophone } from "@tabler/icons-vue";
import { toast } from "vue3-toastify";
import AiAudio from "@/components/business/AiAudio.vue";

const refRecordShowText = ref<HTMLElement>();
const refAiAudio = ref();

enum EnumStatusRecord {
  /** 空闲 */
  Idle = "idle",
  /** 连接中 */
  Connecting = "connecting",
  /** 录音中 */
  Recording = "recording",
  /** 异常状态（网络等） */
  Error = "error"
}

const state = reactive<{
  /** 录音状态 */
  statusRecord: EnumStatusRecord;
  [x: string]: any;
  /** 音频地址 */
  videoUrl: string;
  /** 状态文本 */
  statusText: string;
  /** 当前只言片语 */
  momentWords: string;
  /** 是否ready */
  isReady: boolean | "";
}>({
  statusRecord: EnumStatusRecord.Idle,
  statusText: "",
  momentWords: "",
  videoUrl: "",
  isReady: ""
});

interface IEmits {
  (e: "translate", val: string): void;
  (e: "voiceStart"): void;
  (e: "voiceEnd"): void;
}
const emits = defineEmits<IEmits>();

const onToggleRecord = () => {
  let { statusRecord } = state;
  if (statusRecord == EnumStatusRecord.Idle) {
    onConnect();
  } else if (statusRecord == EnumStatusRecord.Recording) {
    onIdel();
  } else if (statusRecord == EnumStatusRecord.Error) {
    onConnect();
  }
};

const getJsonMessage = jsonMsg => {
  try {
    let text = JSON.parse(jsonMsg.data)["text"];
    emits("translate", text);
    state.momentWords = text;
    console.log("getJsonMessage ============== " + JSON.parse(jsonMsg.data)["text"]);
  } catch (e) {
    toast.error("实时语音异常");
    console.log("getJsonMessage error", jsonMsg);
  }
};

// asr 工具方法、对象
let Recorder;
let WebSocketConnectMethod;
// asr 实例变量
var sampleBuf = new Int16Array();
let wsconnecter;
let rec;

const stopAll = () => {
  var chunk_size = new Array(5, 10, 5);
  var request = {
    chunk_size: chunk_size,
    wav_name: "h5",
    is_speaking: false,
    chunk_interval: 10,
    mode: "online"
  };
  if (sampleBuf.length > 0) {
    wsconnecter.wsSend(sampleBuf);
    console.log("sampleBuf.length" + sampleBuf.length);
    sampleBuf = new Int16Array();
  }
  //1. 关闭socket连接
  wsconnecter.wsSend(JSON.stringify(request));

  setTimeout(function () {
    console.log("call stop ws!");
    wsconnecter.wsStop();
    // }, 3000);
  }, 0);

  //2. 关闭录音
  rec.stop(
    function (blob, duration) {
      console.log("stop record", blob);
      var audioBlob = Recorder.pcm2wav(
        { sampleRate: 16000, bitRate: 16, blob: blob },
        function (theblob, duration) {
          console.log(theblob);
          state.videoUrl = (window.URL || webkitURL).createObjectURL(theblob);
          state.momentWords = "";
        },
        function (msg) {
          console.log(msg);
        }
      );
    },
    function (errMsg) {
      console.log("stop record , errMsg: " + errMsg);
    }
  );
};

const recProcess = (buffer, powerLevel, bufferDuration, bufferSampleRate, newBufferIdx, asyncEnd) => {
  console.log("rec process ... ");
  var data_48k = buffer[buffer.length - 1];

  var array_48k = new Array(data_48k);
  var data_16k = Recorder.SampleData(array_48k, bufferSampleRate, 16000).data;

  sampleBuf = Int16Array.from([...sampleBuf, ...data_16k]);
  var chunk_size = 960; // for asr chunk_size [5, 10, 5]
  while (sampleBuf.length >= chunk_size) {
    let sendBuf = sampleBuf.slice(0, chunk_size);
    sampleBuf = sampleBuf.slice(chunk_size, sampleBuf.length);
    // console.log("sendBuf", sendBuf);
    wsconnecter.wsSend(sendBuf);
  }
};

const onConnect = () => {
  state.statusRecord = EnumStatusRecord.Connecting;
  changeRecordShowText("正在等待...");
  wsconnecter = new WebSocketConnectMethod({
    msgHandle: getJsonMessage,
    stateHandle: getConnState,
    toast: toast
  });
  var ret = wsconnecter.wsStart();
  if (ret == "broswerNotSupport") {
    onIdel();
    toast.error("当前浏览器不支持 WebSocket");
  }
};

// 连接状态响应
const getConnState = connState => {
  console.log("getConnState  ", connState);
  if (connState === "ws-success") {
    //连接成功
    onRecord();
  } else if (connState === "ws-close") {
    //连接关闭
    // stop();
  } else if (connState === "ws-error") {
    //异常
    onError();
  }
};

const onRecord = () => {
  state.videoUrl = "";
  rec.open(
    function () {
      rec.start();
      emits("voiceStart");
      state.statusRecord = EnumStatusRecord.Recording;
      changeRecordShowText("正在听取...");
    },
    errMsg => {
      toast.error(errMsg);
      onError();
    }
  );
};

const onError = () => {
  state.statusRecord = EnumStatusRecord.Error;
  changeRecordShowText("点击重试");
  emits("voiceEnd");
  stopAll();
};

const onIdel = () => {
  state.statusRecord = EnumStatusRecord.Idle;
  changeRecordShowText("");
  emits("voiceEnd");
  state.momentWords = "";
  stopAll();
};

const changeRecordShowText = showText => {
  if (showText == state.statusText) {
    return;
  }
  state.statusText = showText;
  nextTick(() => {
    if (!refRecordShowText.value) {
      //页面销毁时，refRecordShowText不存在
      return;
    }
    animate(
      refRecordShowText.value,
      [
        { transform: "translateX(60px)", opacity: 0 },
        { transform: "translateX(15px)", opacity: 0.2 },
        { transform: "translateX(0px)", opacity: 1 }
      ],
      {
        duration: 500,
        fill: "forwards"
      }
    );
  });
};

const getClassName4BtnRecord = () => {
  let { statusRecord } = state;
  return statusRecord;
};

const init = () => {
  Recorder = window.Recorder;
  WebSocketConnectMethod = window.WebSocketConnectMethod;
  // 录音; 定义录音对象,wav格式
  rec = Recorder({
    type: "pcm",
    bitRate: 16,
    sampleRate: 16000,
    onProcess: recProcess
  });
  state.isReady = true;
};

onMounted(async () => {
  let { origin } = location;
  file.asyncImportJavaScript(`${origin}/assets/js/asr/recorder-core.js`).then(() => {
    Promise.all([
      file.asyncImportJavaScript(`${origin}/assets/js/asr/wav.js`),
      file.asyncImportJavaScript(`${origin}/assets/js/asr/pcm.js`),
      file.asyncImportJavaScript(`${origin}/assets/js/asr/wsconnecter.js`)
    ])
      .then(() => {
        init();
      })
      .catch(e => {
        //文件加载失败、init方法抛错等
        console.error(e);
        toast.error("初始化失败，请检查");
        state.isReady = false;
      });
  });
});

defineExpose({
  reset() {
    //1. 停止音频播放
    refAiAudio.value.pause();
    state.videoUrl = "";
    if ([EnumStatusRecord.Connecting, EnumStatusRecord.Recording].includes(state.statusRecord)) {
      onIdel();
    } else if (EnumStatusRecord.Error == state.statusRecord) {
      onError();
    }
  }
});
</script>
<style lang="scss" scoped>
.tab-voice {
  margin-top: 20px;
  .btn-record {
    padding: 30px;
    border: solid 1px #999;
    border-radius: 100%;
    color: #999;
    position: relative;
    .icon-record {
      position: relative;
      z-index: 10;
    }
    &.idle {
      &:hover {
        // color: rgb(var(--v-theme-info));
        transition: all 0.2s linear;
        color: #333;
      }
    }
    &.connecting {
      color: #333;
    }
    &.recording {
      color: #16a34a;
      background: #f0fdf4;
      border: none;
      &:after {
        position: absolute;
        z-index: 1;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        content: "";
        border-radius: 100%;
        background: #bbf7d0;
        animation: flicker 1.2s ease-out infinite;
      }
    }
    &.error {
      color: #dc2626;
      border-color: #fecdd3;
    }
  }
}

@keyframes flicker {
  0% {
    opacity: 1;
    transform: scale(0.5);
  }
  30% {
    opacity: 1;
  }
  100% {
    opacity: 0;
    transform: scale(var(--point-scale));
  }
}
.shower-momentWords {
  position: absolute;

  left: 20px;
  right: 20px;
  bottom: 100px;
  display: flex;
  justify-content: center;
  p {
    font-size: 50px;
    color: #fff;
    background: #404040;
    border-radius: 2px;
    padding: 0px 12px;
    line-height: 1.5;
  }
}
</style>
