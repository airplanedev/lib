export async function getEnvVars(envVarNames) {
    let filteredEnv = {}
    for (const name of envVarNames) {
        if (process.env.hasOwnProperty(name)) {
            filteredEnv[name] = process.env[name]
        }
    }

    return filteredEnv;
}
