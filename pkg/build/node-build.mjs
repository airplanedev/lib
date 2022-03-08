// import {NodeResolvePlugin} from "@esbuild-plugins/node-resolve";
import { build } from "esbuild";
import fg from 'fast-glob';

const files = fg.sync(['**/*.{js,ts}', '!**/node_modules']);
files.push("/airplane/.airplane/shim.js")

const args = process.argv.slice(2)
const nodeVersion = args[0];
const usesESMModules = args[1] === 'true';

build({
  entryPoints: files,
  bundle: false,
  platform: "node",
  target: `node${nodeVersion}`,
  outdir: "/airplane",
  format: usesESMModules ? "esm" : "cjs",
  allowOverwrite: true,
});
