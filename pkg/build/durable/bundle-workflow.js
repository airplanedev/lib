import { bundleWorkflowCode } from '@temporalio/worker';
import { writeFile } from 'fs/promises';

// Create a workflow bundle using tooling provided by Temporal;
// see https://docs.temporal.io/docs/typescript/workers/#prebuilt-workflow-bundles
// for details.
const { code } = await bundleWorkflowCode({
  workflowsPath: './workflow-wrapper.js',
  workflowInterceptorModules: ['./workflow-interceptors.js'],
});

await writeFile('workflow-bundle.js', code);
