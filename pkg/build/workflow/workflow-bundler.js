const worker = require('@temporalio/worker');
const promises = require('fs/promises');

// Create a workflow bundle using tooling provided by Temporal;
// see https://docs.temporal.io/docs/typescript/workers/#prebuilt-workflow-bundles
// for details.
async function bundle() {
  const { code } = await worker.bundleWorkflowCode({
    workflowsPath: './workflow-shim.js',
    workflowInterceptorModules: ['./workflow-interceptors.js'],
  });

  promises.writeFile('./workflow-bundle.js', code);
}

bundle();
