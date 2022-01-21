// Linked to https://app.airplane.dev/t/typescript_simple [do not edit this line]

type Params = {
  id: string
}

export default async function(params: Params) {
  console.log(`airplane_output_set "${params.id}"`);
}
