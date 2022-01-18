import {NodeResolvePlugin} from "@esbuild-plugins/node-resolve";
import { build } from "esbuild";

// For known ESM modules, do not mark as external and use esbuild to bundle.
// As long as these modules don't happen to pull in any optional modules, we should be OK.
// This is a bandaid until we figure out how to handle ESM without bundling.
const esmModules = new Set(["node-fetch"])
const esmPaths = new Set();

const args = process.argv.slice(2)
const nodeVersion = args[0];
build({
  plugins: [
    // This plugin marks all dependencies as external so they are not bundled 
    // unless they are esm modules or packages pulled in by esm modules.
    NodeResolvePlugin({
      extensions: [".ts", ".js"],
      onResolved: (resolved, importer) => {
        if (esmPaths.has(importer)) {
          // resolved was imported by an esm module. Also treat
          // it as an esm module.
          esmPaths.add(resolved)
          return resolved
        }
        if (resolved.includes("node_modules")) {
          // Check if resolved is an esm module.
          const moduleRegexp = /node_modules\/(.+?)\//
          const moduleMatch = resolved.match(moduleRegexp);
          const module = moduleMatch[1];
          if (esmModules.has(module)) {
            esmPaths.add(resolved)
            return resolved
          }

          return {
            external: true,
          };
        }
        return resolved;
      },
    }),
  ],
  entryPoints: ['/airplane/.airplane/shim.js'],
  bundle: true,
  platform: "node",
  target: `node${nodeVersion}`,
  outfile: "/airplane/.airplane/dist/shim.js",
  external: ["airplane"],
});
