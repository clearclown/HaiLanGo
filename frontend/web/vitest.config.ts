import path from "node:path"
import react from "@vitejs/plugin-react"
import { defineConfig } from "vitest/config"

export default defineConfig({
  plugins: [react()],
  test: {
    environment: "jsdom",
    globals: true,
    setupFiles: ["./vitest.setup.ts"],
    include: ["app/**/*.test.{ts,tsx}", "components/**/*.test.{ts,tsx}", "lib/**/*.test.{ts,tsx}"],
    exclude: ["node_modules/**", ".next/**", "e2e/**"],
    coverage: {
      provider: "v8",
      reporter: ["text", "json", "html"],
      exclude: [
        "node_modules/",
        ".next/",
        "**/*.config.*",
        "**/*.d.ts",
        "**/*.test.*",
        "**/*.spec.*",
      ],
    },
  },
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./"),
    },
  },
})
