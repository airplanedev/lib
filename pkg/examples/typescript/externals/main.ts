// Linked to https://app.airplane.dev/t/typescript_externals [do not edit this line]

import airplane from "airplane";
// Import to force bundler to consider (and skip) them:
import * as pg from "pg";
import * as pgFormat from "pg-format";

type Params = {
  id: string;
};

export default async function (params: Params) {
  airplane.appendOutput(params.id);

  airplane.appendOutput(Object.keys(airplane));
  airplane.appendOutput(Object.keys(pg));
  airplane.appendOutput(Object.keys(pgFormat));
}
