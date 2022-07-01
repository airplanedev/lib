import { proxySinks, proxyActivities } from '@temporalio/workflow';
import task from '{{.Entrypoint}}';

const { logger } = proxySinks();

// Main entrypoint to workflow; wraps a `workflow` function in the user code.
//
// This name must match the name we use when executing the workflow in
// the Airplane API.
export async function __airplaneEntrypoint(params, airplaneArgs) {
  logger.info('airplane_status:started');

  try {
    // Monkey patch process.env
    global.process = {
      env: airplaneArgs.EnvVars
    }

    var result = await task(JSON.parse(params[0]));
  } catch (err) {
    logger.info(err);
    logger.info('airplane_output_append:error ' + JSON.stringify({ error: String(err) }));
    throw err;
  }

  // TODO: Update SDK to include a workflow version of setOutput, then
  // use that instead.
  if (result !== undefined) {
    const output = JSON.stringify(result);
    logChunks(`airplane_output_set ${output}`);
  }
  logger.info('airplane_status:succeeded');
  return result;
}

// Equivalent to logChunks in node SDK, but with extra sinks wrapping so we
// identify which task run generated the output.
const logChunks = (output) => {
  const CHUNK_SIZE = 8192;
  if (output.length <= CHUNK_SIZE) {
    logger.info(output);
  } else {
    const chunkKey = uuidv4();
    for (let i = 0; i < output.length; i += CHUNK_SIZE) {
      logger.info(`airplane_chunk:${chunkKey} ${output.substr(i, CHUNK_SIZE)}`);
    }
    logger.info(`airplane_chunk_end:${chunkKey}`);
  }
};
