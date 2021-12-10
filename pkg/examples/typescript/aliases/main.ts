// Linked to https://app.airplane.dev/t/typescript_aliases [do not edit this line]

import airplane from 'airplane'
import { makeViral } from '@lib/text'

type Params = {
  id: string
}

export default async function(params: Params) {
  airplane.output(makeViral(`america runs on beans`))
  airplane.output(params.id)
}
