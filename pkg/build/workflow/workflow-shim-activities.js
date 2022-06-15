// getEnvVars populates an object with envVarNames as the keys and values from the main Node thread's process.env.
// We filter using envVarNames instead of doing a deep copy of the entire process.env because the majority of
// environment variables do not need to be accessible by the user - only their custom env vars and the env vars required
// by our SDKs should be present in the monkey patched process.env
export async function getEnvVars(envVarNames) {
    let filteredEnv = {}
    for (const name of envVarNames) {
        if (process.env.hasOwnProperty(name)) {
            filteredEnv[name] = process.env[name]
        }
    }

    return filteredEnv;
}
