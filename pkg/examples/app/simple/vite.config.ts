import { defineConfig, loadEnv, ResolvedConfig } from "vite";
import react from "@vitejs/plugin-react";

// https://vitejs.dev/config/
export default ({ command, mode }: ResolvedConfig) => {
  process.env = { ...process.env, ...loadEnv(mode, process.cwd()) };
  const host = process.env.AIRPLANE_API_HOST || "https://api.airplane.so:5000";
  const base = command === "serve" ? "" : `${host}/i/views/getContent/`;
  return defineConfig({
    plugins: [react()],
    base,
    envPrefix: "AIRPLANE_",
    build: {
      assetsDir: "",
    },
  });
};
