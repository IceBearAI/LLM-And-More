export type TypeButtonInTable = {
  text: string;
  color: "error" | "success" | "warning" | "info" | "primary" | "default";
  size: "x-small" | "small" | "medium" | "large" | "x-large";
  click: () => void;
};
