import { workflowInfo } from '@temporalio/workflow';

// Interceptor that allows us to log outbound client calls made from workflow.
// See https://docs.temporal.io/docs/typescript/interceptors for details.
class WorkflowLogOutboundInterceptor {
  info;
  constructor(info) {
    this.info = info;
  }

  async scheduleActivity(input, next) {
    const activityType = input.activityType;

    workflowLog(this.info, `Scheduling activity ${activityType}: ${JSON.stringify(input)}`);
    try {
      const result = await next(input);
      workflowLog(
        this.info,
        `Scheduling activity ${activityType} result: ${JSON.stringify(result)}`
      );
      return result;
    } catch (error) {
      workflowLog(this.info, `Error scheduling activity ${activityType}: ${error}`);
      throw error;
    }
  }

  async scheduleLocalActivity(input, next) {
    workflowLog(this.info, `Scheduling local activity: ${JSON.stringify(input)}`);
    try {
      const result = await next(input);
      workflowLog(this.info, `Scheduling local activity result: ${JSON.stringify(result)}`);
      return result;
    } catch (error) {
      workflowLog(this.info, `Error scheduling local activity: ${error}`);
      throw error;
    }
  }

  async startTimer(input, next) {
    workflowLog(this.info, `Starting timer: ${JSON.stringify(input)}`);
    try {
      const result = await next(input);
      workflowLog(this.info, `Starting timer result: ${JSON.stringify(result)}`);
      return result;
    } catch (error) {
      workflowLog(this.info, `Error starting timer: ${error}`);
      throw error;
    }
  }
}

// Interceptor that allows us to log inbound client calls made to workflow.
// See https://docs.temporal.io/docs/typescript/interceptors for details.
class WorkflowLogInboundInterceptor {
  info;
  constructor(info) {
    this.info = info;
  }

  async execute(input, next) {
    workflowLog(this.info, `Workflow inbound activity execution: ${JSON.stringify(input)}`);
    try {
      const result = await next(input);
      workflowLog(
        this.info,
        `Workflow inbound activity execution result: ${JSON.stringify(result)}`
      );
      return result;
    } catch (error) {
      workflowLog(this.info, `Error in workflow activity execution: ${error}`);
      throw error;
    }
  }

  async handleSignal(input, next) {
    workflowLog(this.info, `Handling signal: ${JSON.stringify(input)}`);
    try {
      const result = await next(input);
      workflowLog(this.info, `Handling signal result: ${JSON.stringify(result)}`);
      return result;
    } catch (error) {
      workflowLog(this.info, `Error handling signal: ${error}`);
      throw error;
    }
  }
}

// Need to export interceptors in this format so that Temporal can
// find them.
export const interceptors = () => ({
  outbound: [new WorkflowLogOutboundInterceptor(workflowInfo())],
  inbound: [new WorkflowLogInboundInterceptor(workflowInfo())],
});

function workflowLog(info, message) {
  // Log out interceptor messages with specific prefix so that we can
  // link them back to a specific task run.
  console.log(`airplane_workflow_log:interceptor//${info.workflowId}/${info.runId} ${message}`);
}
