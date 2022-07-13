// This file includes a shim that will execute your task code.
import airplane from "airplane";
import { {{.EntrypointFunc}} } from '{{.Entrypoint}}';

async function main() {
  if (process.argv.length !== 3) {
    console.log(
      "airplane_output_append:error " +
        JSON.stringify({
          "error":
            `Expected to receive a single argument (via {{ "{{JSON}}" }}). Task CLI arguments may be misconfigured.`,
        }),
    );
    process.exit(1);
  }

  try {
    params = JSON.parse(process.argv[2]);
    let ret = await {{.EntrypointFunc}}._airplane.baseFunc({{range .ParamSlugs}}params.{{.}},{{end}});
    if (ret !== undefined) {
      airplane.setOutput(ret);
    }
  } catch (err) {
    console.error(err);
    console.log(
      "airplane_output_append:error " + JSON.stringify({ "error": String(err) }),
    );
    process.exit(1);
  }
}

main();
