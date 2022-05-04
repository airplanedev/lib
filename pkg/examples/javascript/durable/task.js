// TODO: Put these in SDK.
import { proxyActivities, sleep } from '@temporalio/workflow';

const { queryDB, makeHTTPCall } = proxyActivities({
  startToCloseTimeout: '1 minute',
});

export default async function (params) {
  const dbResult = await queryDB(params);
  await sleep('5 seconds');
  return await makeHTTPCall(params);
}
