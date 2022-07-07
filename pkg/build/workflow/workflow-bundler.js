const worker = require('@temporalio/worker');
const { writeFile } = require('fs/promises');

// Create a workflow bundle using tooling provided by Temporal;
// see https://docs.temporal.io/docs/typescript/workers/#prebuilt-workflow-bundles
// for details.
async function bundle() {
  const { code, sourceMap } = await worker.bundleWorkflowCode({
    workflowsPath: '/airplane/.airplane/workflow-shim.js',
    workflowInterceptorModules: ['/airplane/.airplane/workflow-interceptors.js'],
  });

  await Promise.all([
    writeFile('/airplane/.airplane/workflow-bundle.js', code),
    writeFile('/airplane/.airplane/workflow-bundle.map.js', sourceMap),
  ])
}

bundle();
