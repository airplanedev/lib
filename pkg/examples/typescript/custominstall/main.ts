// Linked to https://app.airplane.dev/t/typescript_externals [do not edit this line]

import airplane from "airplane";
import * as pg from "pg";
import * as pgFormat from "pg-format";
// Prettier is a dev dependency only install because of the
// custom install.
import * as prettier from "prettier";

type Params = {
  id: string;
};

export default async function (params: Params) {
  airplane.output(params.id);

  airplane.output(Object.keys(airplane));
  airplane.output(Object.keys(pg));
  airplane.output(Object.keys(pgFormat));
  airplane.output(Object.keys(prettier));
}
