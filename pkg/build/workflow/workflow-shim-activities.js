// getEnvVars populates an object with envVarNames as the keys and values from the main Node thread's process.env.
// We filter using envVarNames instead of doing a deep copy of the entire process.env because the majority of
// environment variables do not need to be accessible by the user - only their custom env vars and the env vars required
// by our SDKs should be present in the monkey patched process.env
export async function getEnvVars(taskRevisionEnvVarNames, runtimeEnv) {
    let env = {}
    // Add env vars that the user has configured in the task.
    for (const name of taskRevisionEnvVarNames) {
        if (process.env.hasOwnProperty(name)) {
            env[name] = process.env[name]
        }
    }

    // Add airplane-internal env vars that are used by the SDK.
    for (const [name, value] of Object.entries(runtimeEnv)) {
        env[name] = value
    }

    return env;
}
