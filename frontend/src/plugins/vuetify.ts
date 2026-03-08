import "@mdi/font/css/materialdesignicons.css";
// @ts-expect-error Vuetify styles entry has no type declarations
import "vuetify/styles";
import { createVuetify } from "vuetify";

const getTheme = () => ({
  defaultTheme: "fs-darkness",
  themes: {
    "fs-darkness": {
      dark: true,
      colors: {
        background: "#1a1a2e",
        border: "#3a3a5c",
        surface: "#16213e",
        primary: "#F596AA",
        secondary: "#7a7a9e",
        error: "#FF5252",
        info: "#64B5F6",
        success: "#81C784",
        warning: "#FFB74D",
        tagColor: "#CE93D8",
        "button-submit-bg": "#2d5a3d",
        "button-submit-text": "#81C784",
        "button-management": "#ff6b6b",
        "button-disabled": "#555555",
        "text-primary": "#E8E8F0",
        "text-secondary": "#C8C8D8",
        "text-tertiary": "#A0A0B0",
        "loader-color": "#4DD0E1",
        "code-bg": "#1E1E1E",
      },
      variables: {
        "card-color-1": "#2D2D2D",
        "shadow-light": "rgba(0, 0, 0, 0.3)",
        "shadow-medium": "rgba(0, 0, 0, 0.5)",
        "shadow-dark": "rgba(0, 0, 0, 0.7)",
      },
    },
  },
});

export default createVuetify({
  theme: getTheme(),
});
