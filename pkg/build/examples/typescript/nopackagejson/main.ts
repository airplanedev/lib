// Linked to https://app.airplane.dev/t/typescript_no_package_json [do not edit this line]

type Params = {
  id: string
}

export default async function(params: Params) {
  console.log(`airplane_output "${params.id}"`);
}
