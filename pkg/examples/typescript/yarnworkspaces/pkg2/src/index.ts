// Linked to https://app.airplane.dev/t/typescript_yarnworkspaces [do not edit this line]

import airplane from "airplane";
import { name as pkg1name } from "pkg1/src";
// > node-fetch is an ESM-only module
// https://github.com/node-fetch/node-fetch#loading-and-configuring-the-module
import fetch from "node-fetch";

type Params = {
  id: string;
};

export default async function (params: Params) {
  console.log(`imported package with name=${pkg1name}`);
  const res = await fetch("https://google.com");
  const html = await res.text();
  console.log(html);

  // I'm feeling lucky!
  if (html.toLowerCase().indexOf("lucky")) {
    airplane.setOutput(params.id);
  }
}
