// Linked to https://app.airplane.dev/t/typescript_externals [do not edit this line]

import airplane from "airplane";
// Import to force bundler to consider (and skip) them:
import * as pg from "pg";
import * as pgFormat from "pg-format";

type Params = {
  id: string;
};

export default async function (params: Params) {
  airplane.output(params.id);

  airplane.output(Object.keys(airplane));
  airplane.output(Object.keys(pg));
  airplane.output(Object.keys(pgFormat));
}
