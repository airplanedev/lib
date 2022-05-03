import { proxySinks } from '@temporalio/workflow';
import workflow from '{{.Entrypoint}}';

const { logger } = proxySinks();

// Main entrypoint to workflow; wraps a `workflow` function in the user code.
export async function apWorkflow(params) {
  logger.info('airplane_status:started');
  const result = await workflow(params);
  const output = result === undefined ? null : JSON.stringify(result);
  logChunks(`airplane_output_set ${output}`);
  logger.info('airplane_status:succeeded');
  return result;
}

// Equivalent to logChunks in node SDK, but with extra sinks wrapping so we
// identity
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
