// Linked to https://app.airplane.dev/t/typescript_esnext [do not edit this line]

import airplane from 'airplane'

type Params = {
  id: string
}

// See:
// - https://esbuild.github.io/content-types/#javascript
// - https://node.green/
// - V8 releases: https://v8.dev/blog
export default async function(params: Params) {
  // Test a few syntax changes and a few polyfills to make sure
  // they are compiled/polyfilled correctly under order versions of Node.

  airplane.setOutput(2 ** 5, "exp") // exponent operator (es2016)

  try {
    airplane.setOutput("throw", "try")
    throw new Error("yikes!")
  } catch { // optional catch binding (es2019)
    airplane.setOutput("catch", "try")
  }

  const sayings = {
    "english": {
      "hello": "hi"
    }
  }
  for (const lang of ["english", "spanish"]) {
    airplane.appendOutput(sayings[lang]?.hello, "saying") // optional chaining (es2020)
  }

  // TODO: polyfill JS functionality so older versions of Node can access
  // 
  // const foo = "👋 <id> <id> <id>"
  // foo.replaceAll("<id>", params.id) // replaceAll (es2021)
  // airplane.output(foo)
  
  airplane.setOutput(params.id, "output")
}
