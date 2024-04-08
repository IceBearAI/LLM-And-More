import { IconClock, IconCircleCheckFilled, IconLoader, IconAlarm } from "@tabler/icons-vue";

export const digitalhumanStatusMap = {
  running: {
    text: "合成中",
    bgColor: "bg-info"
  },
  success: {
    text: "已完成",
    bgColor: "bg-success"
  },
  failed: {
    text: "失败",
    bgColor: "bg-error"
  },
  waiting: {
    text: "等待中",
    bgColor: "bg-warning"
  },
  cancel: {
    text: "已取消",
    bgColor: "bg-error"
  }
};

export const taskDetailConfig = [
  {
    statusText: "等待中",
    valueText: "个任务",
    key: "waitingJobCount",
    color: "warning",
    bgColor: "lightwarning",
    icon: IconLoader
  },
  {
    statusText: "已完成",
    valueText: "个合成任务",
    key: "successJobCount",
    color: "success",
    bgColor: "lightsuccess",
    icon: IconCircleCheckFilled
  },
  {
    statusText: "总合成时间",
    valueText: "",
    key: "totalDurationCount",
    color: "secondary",
    bgColor: "lightprimary",
    icon: IconAlarm
  }
];
