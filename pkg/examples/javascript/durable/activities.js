import { Context } from '@temporalio/activity';

export async function queryDB(query) {
  log(query);
  return {
    key: 'value',
  };
}
  
export async function makeHTTPCall(url) {
  log('Doing another interesting thing');
  return {
    another_key: 'another_value',
  };
}

// TODO: Put this in SDK.
function log(message) {
  const info = Context.current().info;
  console.log(
    `[ap:activity:${info.activityType}:${info.workflowExecution.workflowId}:${info.workflowExecution.runId}] ${message}`
  );
}
