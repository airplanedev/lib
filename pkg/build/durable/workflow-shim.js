import { proxySinks } from '@temporalio/workflow';
import task from '{{.Entrypoint}}';

const { logger } = proxySinks();

// Main entrypoint to workflow; wraps a `workflow` function in the user code.
export async function __airplaneEntrypoint(params) {
  logger.info('airplane_status:started');
  const result = await task(params);
  const output = result === undefined ? null : JSON.stringify(result);
  logChunks(`airplane_output_set ${output}`);
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
