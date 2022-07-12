import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";

export default defineConfig({
  plugins: [react()],
  envPrefix: "AIRPLANE_",
  resolve: {
    preserveSymlinks: true,
  },
  base: "",
  build: {
    assetsDir: "",
  },
});
