import { StreamLanguage } from "@codemirror/language";
import { shell } from "@codemirror/legacy-modes/mode/shell";
import { json, jsonParseLinter } from "@codemirror/lang-json";
import { python } from "@codemirror/lang-python";

export const languages = {
  json,
  python,
  shell: () => StreamLanguage.define(shell)
};

export const languagesLint = {
  json: jsonParseLinter()
};
