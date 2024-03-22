import type { Ref, ComponentPublicInstance } from "vue";
import { nextTick, ref } from "vue";

type ScrollElement = HTMLDivElement | ComponentPublicInstance | null;

interface ScrollReturn {
  scrollRef: Ref<ScrollElement>;
  scrollToBottom: () => Promise<void>;
  scrollToTop: () => Promise<void>;
  scrollToBottomIfAtBottom: () => Promise<void>;
}

export function useScroll(): ScrollReturn {
  const scrollRef = ref<ScrollElement>();

  const scrollToBottom = async () => {
    await nextTick();
    if (scrollRef.value) {
      const refEl = "$el" in scrollRef.value ? scrollRef.value.$el : scrollRef.value;
      refEl.scrollTop = refEl.scrollHeight;
    }
  };

  const scrollToTop = async () => {
    await nextTick();
    if (scrollRef.value) {
      const refEl = "$el" in scrollRef.value ? scrollRef.value.$el : scrollRef.value;
      refEl.scrollTop = 0;
    }
  };

  const scrollToBottomIfAtBottom = async () => {
    await nextTick();
    if (scrollRef.value) {
      const refEl = "$el" in scrollRef.value ? scrollRef.value.$el : scrollRef.value;
      const threshold = 100; // 阈值，表示滚动条到底部的距离阈值
      const distanceToBottom = refEl.scrollHeight - refEl.scrollTop - refEl.clientHeight;
      if (distanceToBottom <= threshold) refEl.scrollTop = refEl.scrollHeight;
    }
  };

  return {
    scrollRef,
    scrollToBottom,
    scrollToTop,
    scrollToBottomIfAtBottom
  };
}
