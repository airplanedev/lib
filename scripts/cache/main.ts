// Linked to https://app.airplane.dev/t/update_registry_cache [do not edit this line]

import execa from 'execa'
import * as fs from 'fs'
import airplane from 'airplane'

type Params = {
  gcp_project: string
}

// Pushes the latest versions of all base images to the public Airplane Registry cache.
// 
// airplane dev scripts/cache/main.ts -- --gcp_project=airplane-stage
export default async function(params: Params) {
  if (!params.gcp_project) {
    throw new Error("expected a project")
  }

  const versionsFile = JSON.parse(fs.readFileSync('./pkg/build/versions.json') as any) as Record<string, Record<string, {
    image: string
    tag: string
    digest: string
  }>>
  console.log(versionsFile)

  for (let kind of Object.keys(versionsFile)) {
    for (let version of Object.values(versionsFile[kind])) {
      airplane.appendOutput(`Pushing ${kind}:${version.tag}@${version.digest} to ${params.gcp_project}...`, "logs")
      // Digests should always match the tag, to make it easier to read `versions.json`.
      // The digest should be for linux/amd64.
      await execa("docker", ["pull", `${version.image}@${version.digest}`])
      await execa("docker", ["tag", `${version.image}@${version.digest}`, `us-central1-docker.pkg.dev/${params.gcp_project}/public-cache/${kind}:${version.tag}`])
      await execa("docker", ["push", `us-central1-docker.pkg.dev/${params.gcp_project}/public-cache/${kind}:${version.tag}`])
      airplane.appendOutput(`Pushed ${kind}:${version.tag}@${version.digest} to ${params.gcp_project}`, "logs")
    }
  }
}
