import type { Config } from "tailwindcss"

const config: Config = {
  content: [
    "./pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./components/**/*.{js,ts,jsx,tsx,mdx}",
    "./app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: "#4A90E2",
        secondary: "#50C878",
        accent: "#FF6B6B",
        background: "#FFFFFF",
        "background-secondary": "#F5F7FA",
        "text-primary": "#2C3E50",
        "text-secondary": "#7F8C8D",
        border: "#E0E6ED",
        success: "#27AE60",
        warning: "#F39C12",
        error: "#E74C3C",
        info: "#3498DB",
      },
    },
  },
  plugins: [],
}

export default config
