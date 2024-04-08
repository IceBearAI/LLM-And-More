<template>
  <div :class="`text--${msgType}`" class="text bg-grey100 rounded-md px-3 py-2 mb-1">
    <div ref="textRef" class="text-content">
      <div v-if="msgType === 'answer'" class="markdown-body" v-html="renderText" />
      <div v-else class="whitespace-pre-wrap" v-text="renderText" />
      <div v-if="loading" class="cursor animation-blink" />
    </div>
  </div>
</template>
<script setup lang="ts">
import { computed, ref, onMounted, onUnmounted, watch } from "vue";
import MarkdownIt from "markdown-it";
import hljs from "highlight.js";
import mdImsize from "@steelydylan/markdown-it-imsize";
import { copyToClip } from "@/utils/copy";

interface IProps {
  text?: string;
  msgType?: string;
  loading?: boolean;
}

const props = withDefaults(defineProps<IProps>(), {
  text: "",
  msgType: "question",
  loading: false
});
const textRef = ref<HTMLElement>();

const mdi = new MarkdownIt({
  html: false,
  linkify: true,
  highlight(code, language) {
    const validLang = !!(language && hljs.getLanguage(language));
    if (validLang) {
      const lang = language ?? "";
      return highlightBlock(hljs.highlight(code, { language: lang }).value, lang);
    }
    return highlightBlock(hljs.highlightAuto(code).value, "");
  }
});
mdi.use(mdImsize, { autofill: true });

const renderText = computed(() => {
  const value = props.text || "";
  if (props.msgType === "answer") {
    return mdi.render(value);
  } else {
    return value;
  }
});

function highlightBlock(str: string, lang?: string) {
  return `<pre class="code-block-wrapper"><div class="code-block-header"><span class="code-block-header__lang">${lang}</span><span class="code-block-header__copy">复制代码</span></div><code class="hljs code-block-body ${lang}">${str}</code></pre>`;
}

if (textRef.value) {
  const copyBtn = textRef.value.querySelectorAll(".code-block-header__copy");
  copyBtn.forEach(btn => {
    btn.addEventListener("click", () => {
      const code = btn.parentElement?.nextElementSibling?.textContent;
      if (code) {
        copyToClip(code).then(() => {
          btn.textContent = "复制成功";
          setTimeout(() => {
            btn.textContent = "复制代码";
          }, 1000);
        });
      }
    });
  });
}

function addCopyEvents() {
  if (textRef.value) {
    const copyBtn = textRef.value.querySelectorAll(".code-block-header__copy");
    copyBtn.forEach(btn => {
      btn.addEventListener("click", () => {
        let code = null;
        if (btn.parentElement && btn.parentElement.nextElementSibling) {
          code = btn.parentElement.nextElementSibling.textContent;
        }
        if (code) {
          copyToClip(code).then(() => {
            btn.textContent = "复制成功";
            setTimeout(() => {
              btn.textContent = "复制代码";
            }, 1000);
          });
        }
      });
    });
  }
}

function removeCopyEvents() {
  if (textRef.value) {
    const copyBtn = textRef.value.querySelectorAll(".code-block-header__copy");
    copyBtn.forEach(btn => {
      btn.removeEventListener("click", () => {});
    });
  }
}

watch(
  () => props.loading,
  val => {
    if (!val) {
      addCopyEvents();
    }
  }
);

onMounted(() => {
  addCopyEvents();
});

onUnmounted(() => {
  removeCopyEvents();
});
</script>
<style lang="scss">
@import "./style.scss";
</style>
<style lang="scss" scoped>
.text {
  min-width: 20px;
}
.text-content {
  font-size: 14px;
  line-height: 1.625;
  overflow-wrap: break-word;
}
.cursor {
  width: 4px;
  height: 20px;
}
</style>
