import type { App } from "vue";
import { createVuetify } from "vuetify";
import "@mdi/font/css/materialdesignicons.css";
import * as components from "vuetify/components";
import * as directives from "vuetify/directives";

//i18
import { createVueI18nAdapter } from "vuetify/locale/adapters/vue-i18n";
import { createI18n, useI18n } from "vue-i18n";
import messages from "@/utils/locales/messages";
import { EnumLanguage } from "@/types/apps/System.ts";

//DragScroll
import { VueDraggableNext } from "vue-draggable-next";

// import { BLUE_THEME} from '@/theme/LightTheme';
import { BLUE_THEME, RED_THEME, PURPLE_THEME, GREEN_THEME, INDIGO_THEME, ORANGE_THEME } from "@/theme/LightTheme";
import {
  DARK_BLUE_THEME,
  DARK_RED_THEME,
  DARK_ORANGE_THEME,
  DARK_PURPLE_THEME,
  DARK_GREEN_THEME,
  DARK_INDIGO_THEME
} from "@/theme/DarkTheme";

const i18n = createI18n({
  legacy: false,
  locale: EnumLanguage.Chinese,
  messages: messages,
  silentTranslationWarn: true,
  silentFallbackWarn: true
});

const helloVuetify = createVuetify({
  locale: {
    adapter: createVueI18nAdapter({ i18n, useI18n })
  },
  components: {
    draggable: VueDraggableNext
  },
  directives,

  theme: {
    defaultTheme: "BLUE_THEME",
    themes: {
      BLUE_THEME,
      RED_THEME,
      PURPLE_THEME,
      GREEN_THEME,
      INDIGO_THEME,
      ORANGE_THEME,
      DARK_BLUE_THEME,
      DARK_RED_THEME,
      DARK_ORANGE_THEME,
      DARK_PURPLE_THEME,
      DARK_GREEN_THEME,
      DARK_INDIGO_THEME
    }
  },
  defaults: {
    VCard: {
      rounded: "md"
    },
    VTextField: {
      variant: "outlined",
      density: "compact",
      color: "primary"
    },
    VTextarea: {
      variant: "outlined",
      density: "compact",
      color: "primary"
    },
    VSelect: {
      variant: "outlined",
      density: "compact",
      color: "primary"
    },
    VListItem: {
      minHeight: "45px"
    },
    VTooltip: {
      location: "top"
    }
  }
});

export function setVuetify(app: App) {
  app.use(i18n);
  app.use(helloVuetify);
}
