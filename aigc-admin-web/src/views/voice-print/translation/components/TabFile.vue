<template>
  <v-row ref="refBox" class="justify-center py-4">
    <v-col cols="12" class="h-center" :class="{ 'opacity-0': style.isLoading }">
      <AiAudio ref="refAiAudio" :src="audioUrl" type="complex" />
    </v-col>
    <v-col cols="12" md="6" class="mt-5">
      <v-file-input
        v-model="state.audioFile"
        accept="audio/*, .mp3, .wma, .amr, .wav, .m4a"
        label="请选择音频"
        prepend-icon="mdi-volume-high"
        hide-details
        variant="outlined"
        @change="onTranslate"
      ></v-file-input>
    </v-col>
  </v-row>
</template>
<script setup lang="ts">
import { reactive, toRefs, ref, computed } from "vue";
import AiAudio from "@/components/business/AiAudio.vue";
import { http } from "@/utils";
const refBox = ref();
const refAiAudio = ref();
const state = reactive({
  style: {
    isLoading: false
  },
  formData: {},
  audioFile: []
});
const { style, formData } = toRefs(state);

interface IEmits {
  (e: "translate", val: string): void;
}

const emits = defineEmits<IEmits>();

const audioUrl = computed(() => {
  let { audioFile } = state;
  if (audioFile && audioFile.length > 0) {
    return URL.createObjectURL(audioFile[0]);
  }
  return "";
});

const onTranslate = async () => {
  if (!audioUrl) return;

  state.style.isLoading = true;
  const [err, res] = await http.upload({
    url: "/voice/translation",
    showLoading: refBox.value.$el,
    data: {
      file: state.audioFile[0]
    }
  });
  state.style.isLoading = false;
  if (res) {
    emits("translate", res.data.text);
  }
};

defineExpose({
  reset() {
    state.style.isLoading = false;
    refAiAudio.value.pause();
    state.audioFile = [];
  }
});
</script>
<style lang="scss"></style>
