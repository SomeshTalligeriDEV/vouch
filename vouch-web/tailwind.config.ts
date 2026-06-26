import type { Config } from "tailwindcss";

const config: Config = {
  content: ["./app/**/*.{ts,tsx}", "./components/**/*.{ts,tsx}"],
  theme: {
    extend: {
      colors: {
        // Vouch design system (imported from Vouch.dc.html)
        paper: "#E7E9E2", // page background
        cream: "#E7E9E2", // alias used by inner pages
        panel: "#FFFFFF", // cards
        ink: "#14181B", // dark surfaces + primary text
        text: "#15181C",
        line: "#D7DACF", // soft borders
        accent: "#C8F24C", // lime accent
        "accent-ink": "#7C9A18", // lime, readable on light
        crimson: "#7C9A18", // legacy alias → accent-ink (inner pages)
        "crimson-dark": "#5E7612",
        gold: "#C8F24C", // legacy alias → accent
        muted: "#5C6268",
      },
      fontFamily: {
        sans: ["var(--font-body)", "system-ui", "sans-serif"],
        display: ["var(--font-display)", "system-ui", "sans-serif"],
        mono: ["var(--font-mono)", "ui-monospace", "monospace"],
      },
      boxShadow: {
        hard: "4px 4px 0 #14181B",
        "hard-lg": "6px 6px 0 #14181B",
        soft: "0 24px 50px -28px rgba(20,24,27,.45)",
      },
    },
  },
  plugins: [],
};

export default config;
