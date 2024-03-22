import { IconCircleCheckFilled, IconClockPause, IconClock, IconClockCancel } from "@tabler/icons-vue";

export const trainStatusMap = {
  running: {
    text: "训练中",
    textColor: "text-info",
    icon: IconClock,
    iconColor: "rgb(var(--v-theme-info))"
  },
  success: {
    text: "已完成",
    textColor: "text-success",
    icon: IconCircleCheckFilled,
    iconColor: "#67C23A"
  },
  failed: {
    text: "失败",
    textColor: "text-error",
    icon: IconClockCancel,
    iconColor: "rgb(var(--v-theme-error))"
  },
  waiting: {
    text: "等待中",
    textColor: "text-warning",
    icon: IconClockPause,
    iconColor: "rgb(var(--v-theme-warning))"
  },
  cancel: {
    text: "已取消",
    textColor: "text-error",
    icon: IconClockPause,
    iconColor: "#ccc"
  }
};
