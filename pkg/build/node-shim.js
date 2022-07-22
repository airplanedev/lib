// This file includes a shim that will execute your task code.
import airplane from "airplane";
{{if .EntrypointFunc}}
import { {{.EntrypointFunc}} as task } from "{{.Entrypoint}}";
{{else}}
import task from "{{.Entrypoint}}";
{{end}}

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
    var ret;
    if ("__airplane" in task) {
      ret = await task.__airplane.baseFunc(JSON.parse(process.argv[2]));
    } else {
      ret = await task(JSON.parse(process.argv[2]));
    }
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
