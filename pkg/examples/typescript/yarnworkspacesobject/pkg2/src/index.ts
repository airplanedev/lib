// Linked to https://app.airplane.dev/t/typescript_yarnworkspaces [do not edit this line]

import airplane from "airplane";
import { f } from "pkg1/src";

type Params = {
  id: string;
};

export default async function (params: Params) {
  const res = await f();
  console.log(res);
  airplane.setOutput(params.id);
}
