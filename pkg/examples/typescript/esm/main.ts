// Linked to https://app.airplane.dev/t/typescript_esm [do not edit this line]

import airplane from 'airplane'
// > node-fetch is an ESM-only module
// https://github.com/node-fetch/node-fetch#loading-and-configuring-the-module
import fetch from 'node-fetch'

type Params = {
  id: string
}

export default async function(params: Params) {
  const res = await fetch("https://google.com");
  const html = await res.text();

  console.log(airplane)
  // I'm feeling lucky!
  if (html.toLowerCase().indexOf("lucky")) {
    airplane.appendOutput(params.id)
  }
}
