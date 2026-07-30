import type { Config } from "tailwindcss"

const config: Config = {
  content: ["./src/**/*.{ts,tsx}"],
  theme: {
    extend: {
      colors: {
        primary: { DEFAULT: "#1A73E8", 50: "#E8F0FE", 100: "#D2E3FC", 200: "#A8C7FA", 300: "#7AABF7", 400: "#4C8FF4", 500: "#1A73E8", 600: "#1557B0", 700: "#0F3D78", 800: "#0A2A54", 900: "#051730" },
      },
    },
  },
  plugins: [],
}

export default config
