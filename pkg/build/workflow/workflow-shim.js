import airplane from 'airplane';
import { proxySinks } from '@temporalio/workflow';
import task from '{{.Entrypoint}}';

const { logger } = proxySinks();

// Main entrypoint to workflow; wraps a `workflow` function in the user code.
//
// This name must match the name we use when executing the workflow in
// the Airplane API.
export async function __airplaneEntrypoint(params, workflowArgs) {
  logger.info('airplane_status:started');
  try {
    // Monkey patch process.env
    global.process = {
      env: workflowArgs.EnvVars
    }

    var result = await task(JSON.parse(params[0]));
  } catch (err) {
    logger.info(err);
    logger.info('airplane_output_append:error ' + JSON.stringify({ error: String(err) }));
    throw err;
  }

  if (result !== undefined) {
    const output = JSON.stringify(result);
    airplane.setOutput(output);
  }
  logger.info('airplane_status:succeeded');
  return result;
}
