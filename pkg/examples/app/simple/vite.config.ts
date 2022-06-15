import { defineConfig, loadEnv, ResolvedConfig } from "vite";
import react from "@vitejs/plugin-react";

// https://vitejs.dev/config/
export default ({ command, mode }: ResolvedConfig) => {
  process.env = { ...process.env, ...loadEnv(mode, process.cwd()) };
  const viewID = process.env.AIRPLANE_VIEW_ID;
  const host = process.env.AIRPLANE_API_HOST || "https://api.airplane.dev";
  const base =
    command === "build" ? `${host}/i/views/getContent/${viewID}/` : "";
  return defineConfig({
    plugins: [react()],
    base,
    envPrefix: "AIRPLANE_",
    build: {
      assetsDir: "",
    },
  });
};
