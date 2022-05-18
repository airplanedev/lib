import fetch from "node-fetch";
import * as pg from "pg";
import * as pgFormat from "pg-format";

export const f = async () => {
  console.log(pg);
  console.log(pgFormat);
  const res = await fetch("https://google.com");
  const html = await res.text();
  return html;
};
