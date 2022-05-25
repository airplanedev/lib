// Linked to https://app.airplane.dev/t/typescript_npm [do not edit this line]

import airplane from 'airplane'

type Params = {
  id: string
}

export default async function(params: Params) {
  // Can access legacy outputs (removed in 0.2.0) since it overrides the airplane version.
  airplane.output(params.id)
}
