// Linked to https://app.airplane.dev/t/typescript_yarnworkspaces [do not edit this line]

import airplane from 'airplane'
import {name as pkg1name} from 'pkg1/src'

type Params = {
  id: string
}

export default async function(params: Params) {
  console.log(`imported package with name=${pkg1name}`)
  airplane.output(params.id)
}
