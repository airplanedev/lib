import { NativeConnection, Worker } from '@temporalio/worker';

// Activity code runs in the same node process as the worker, so we import it here directly.
//
// TODO: Make this path configurable.
import * as activities from '../activities.js';

// Main worker entrypoint; starts a worker that will process activities
// and workflows for a single task queue (equivalent to airplane task revision).
async function runWorker(params) {
  // Get temporal address, queue, and namespace from the environment
  const temporalHost = process.env.AP_TEMPORAL_ADDR || 'localhost:7233';
  const taskQueue = process.env.AP_TASK_QUEUE || 'fake-task-revision-id';
  const namespace = process.env.AP_NAMESPACE || 'default';

  const connection = await NativeConnection.create({
    // TODO: Insert a token for auth purposes
    address: temporalHost,
  });

  // Sinks allow us to log from workflows.
  const sinks = {
    logger: {
      info: {
        fn(workflowInfo, message) {
          // Prefix all logs with the workflow ID (equivalent to the run ID) so we
          // can link the logs back to a specific task run.
          console.log(`[ap:workflow::${workflowInfo.workflowId}:${workflowInfo.runId}] ${message}`);
        },
        callDuringReplay: false,
      },
    },
  };

  const worker = await Worker.create({
    // Path to bundle created by bundle-workflow.js script; this should be relative
    // to the shim.
    workflowBundle: { path: '/airplane/.airplane/workflow-bundle.js' },
    activities,
    connection,
    namespace,
    taskQueue,
    interceptors: {
      activityInbound: [(ctx) => new ActivityLogInboundInterceptor(ctx)],
    },
    sinks,
  });

  await worker.run();
}

// Interceptor that allows us to add extra logs around when activities start and
// end. See https://docs.temporal.io/docs/typescript/interceptors for details.
export class ActivityLogInboundInterceptor {
  info;
  constructor(ctx) {
    this.info = ctx.info;
  }

  async execute(input, next) {
    activityLog(this.info, `Starting activity with input: ${JSON.stringify(input)}`);
    try {
      const result = await next(input);
      activityLog(this.info, `Result from activity run: ${JSON.stringify(input)}`);
      return result;
    } catch (error) {
      activityLog(this.info, `Caught error, retrying: ${error}`);
      throw error;
    }
  }
}

function activityLog(info, message) {
  // Prefix all logs with metadata that we can use to link the message back to a
  // specific task run.
  console.log(
    `[ap:activity:${info.activityType}:${info.workflowExecution.workflowId}:${info.workflowExecution.runId}] ${message}`
  );
}

// Code from regular node-shim.
async function main() {
  if (process.argv.length !== 3) {
    console.log(
      'airplane_output_append:error ' +
        JSON.stringify({
          error: `Expected to receive a single argument (via {{ "{{JSON}}" }}). Task CLI arguments may be misconfigured.`,
        })
    );
    process.exit(1);
  }

  try {
    let ret = await runWorker(JSON.parse(process.argv[2]));
    if (ret !== undefined) {
      airplane.setOutput(ret);
    }
  } catch (err) {
    console.error(err);
    console.log('airplane_output_append:error ' + JSON.stringify({ error: String(err) }));
    process.exit(1);
  }
}

main();
