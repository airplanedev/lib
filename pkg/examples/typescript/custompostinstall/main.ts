// Linked to https://app.airplane.dev/t/typescript_externals [do not edit this line]

import airplane from "airplane";
// prettier is installed in the custom postinstall.
// @ts-ignore
import * as prettier from "prettier";

type Params = {
  id: string;
};

export default async function (params: Params) {
  airplane.output(params.id);

  airplane.output(Object.keys(airplane));
  airplane.output(Object.keys(prettier));
}
