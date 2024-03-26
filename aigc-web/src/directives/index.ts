import { App, Directive } from "vue";
import { clipboard } from "./modules/copy.directive.js";

const directivesList: { [key: string]: Directive } = {
  copy: clipboard
};

const directives = {
  install: function (app: App<Element>) {
    Object.keys(directivesList).forEach(key => {
      app.directive(key, directivesList[key]);
    });
  }
};

export default directives;
