import { IconPhoto, IconFileTypePdf } from "@tabler/icons-vue";

export const statusMap = {
  running: {
    color: "info"
  },
  success: {
    color: "success"
  },
  failed: {
    color: "error"
  },
  waiting: {
    color: "warning"
  },
  cancel: {
    color: "default"
  }
};

export const mediaType = {
  image: {
    text: "图片",
    className: "!bg-green-500"
  },
  video: {
    text: "视频",
    className: "!bg-cyan-500"
  }
};

export const textMarkStatus = {
  pending: "待标注",
  processing: "标注中",
  completed: "已完成",
  abandoned: "已废弃",
  cleaned: "已取消"
};

export const fileTypes = {
  image: {
    text: "图片",
    icon: IconPhoto
  },
  pdf: {
    text: "PDF",
    icon: IconFileTypePdf
  }
};

export const fileLanguageMap = {
  py: "python"
};
