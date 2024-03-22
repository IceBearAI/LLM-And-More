<template>
  <div class="compo-aiAudio">
    <audio v-if="type == 'simple'" style="vertical-align: top" :src="src" controls>
      Your browser does not support the audio element.
    </audio>
    <div
      class="aiAudio-complex"
      v-loading="!style.isReadyComplex"
      :element-loading-text="$t('aiAudio.loadingText')"
      v-else-if="type == 'complex'"
    >
      <div v-if="!src" class="box-error">
        <el-divider class="tips">
          <span class="text-info">{{ $t("aiAudio.noUrl") }}</span></el-divider
        >
      </div>
      <div ref="refBoxComplex" :class="style.classNameComplex"></div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { reactive, toRefs, ref, onMounted, watchEffect } from "vue";
import WaveSurfer from "wavesurfer.js";
import $ from "jquery";

enum EnumGender {
  Male = 1, //  男
  Female = 2, //女
  Unknown = 3 //未知
}

interface Props {
  /** 音频地址 */
  src: string;
  /** 性别 1:男， 2:女  ，3:未知*/
  gender?: EnumGender;
  /** 模式 simple简易、complex复杂 */
  type?: "simple" | "complex";
  time?: number;
}

const refBoxComplex = ref<HTMLElement>();
const state = reactive({
  style: {
    isReadyComplex: false,
    classNameComplex: "opacity-0"
  },
  formData: {}
});
const { style, formData } = toRefs(state);

const props = withDefaults(defineProps<Props>(), {
  src: "",
  gender: 3,
  type: "simple"
});

let wavesurfer;
const renderWave = () => {
  let { type, gender, src } = props;
  if (type == "complex") {
    //清空内容，避免生成多个
    $(refBoxComplex.value).html("");
    let colors = {
      [EnumGender.Male]: {
        waveColor: "#38bdf8",
        progressColor: "#38bdf850",
        cursorColor: "#CCC"
      },
      [EnumGender.Female]: {
        waveColor: "#f980e9",
        progressColor: "#f980e950",
        cursorColor: "#ccc"
      },

      [EnumGender.Unknown]: {
        waveColor: "#475569",
        progressColor: "#47556950",
        cursorColor: "#ccc"
      }
    };
    let { waveColor, progressColor, cursorColor } = colors[gender];
    wavesurfer = WaveSurfer.create({
      container: refBoxComplex.value,
      waveColor,
      progressColor,
      cursorColor,
      cursorWidth: 3,
      // media: audio,
      mediaControls: true,
      // url: "https://wavesurfer.xyz/wavesurfer-code/examples/audio/audio.wav",
      url: src,
      autoplay: false,
      /** Pass false to disable clicks on the waveform */
      interact: true
    });

    if (src) {
      wavesurfer.on("ready", () => {
        renderWaveWithNoProblem();
      });
    } else {
      state.style.isReadyComplex = true;
      renderWaveWithNoProblem();
    }
  }
};

/**
 * 绘制UI，修复播放条宽度问题
 */
const renderWaveWithNoProblem = () => {
  state.style.classNameComplex = "w-75 opacity-0";

  setTimeout(() => {
    state.style.classNameComplex = "w-100 opacity-0";
  }, 100);

  setTimeout(() => {
    state.style.isReadyComplex = true;
    state.style.classNameComplex = "";
  }, 200);
};

watchEffect(() => {
  let { type, src } = props;
  if (type || src) {
    state.style.classNameComplex = "opacity-0";
    setTimeout(() => {
      renderWave();
    }, 100);
  }
});

defineExpose({
  pause() {
    wavesurfer.pause();
  }
});
</script>
<style lang="scss">
.compo-aiAudio {
  position: relative;
  .aiAudio-complex {
    min-height: 180px;
  }
  .box-error {
    position: absolute;
    left: 0;
    right: 0;
    top: 0;
    border-left: solid 4px #47556950;
    bottom: 54px;
    .tips {
      position: absolute;
      left: 0;
      right: 0;
      top: 50%;
      transform: translateY(-50%);
      text-align: center;
      font-size: 14px;
      margin: 0;
    }
  }
}
</style>
