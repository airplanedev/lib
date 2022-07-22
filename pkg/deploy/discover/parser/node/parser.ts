const path = require("node:path");

type Param = {
    slug: string
    name: string
    type: string
    description?: string
    default?: any
    required?: boolean
    options?: any[]
    regex?: string
}

type Def = {
    slug: string
    name?: string
    description?: string
    parameters?: Param[]
    requireRequests?: boolean
    allowSelfApprovals?: boolean
    timeout?: number
    constraints?: Record<string, string>
    runtime?: "" | "workflow"
}

type DefWithBuildArgs = Def & {
    entrypointFunc: string
}

const extractTaskConfigs = (files: string[]): DefWithBuildArgs[] => {
    let configs: DefWithBuildArgs[] = [];
    for (const file of files) {
        const resolvedPath = path.relative(__dirname, file);
        const lib = resolvedPath.replace(/.ts$/, "");
        const exports = require(`./${lib}`);

        for (const exportName in exports) {
            const item = exports[exportName];
            
            if ("__airplane" in item) {
                const config = item.__airplane.config;

                var params: Param[] = [];
                for (var uParamSlug in config.parameters) {
                    const uParamConfig = config.parameters[uParamSlug];

                    if (typeof uParamConfig === "string") {
                        params.push({
                            slug: uParamSlug,
                            name: uParamSlug,
                            type: uParamConfig,
                        });
                    } else {
                        params.push({
                            slug: uParamSlug,
                            name: uParamConfig["name"] ? uParamConfig["name"] : uParamSlug,
                            type: uParamConfig["type"],
                            description: uParamConfig["description"],
                            default: uParamConfig["default"],
                            required: uParamConfig["required"],
                            options: uParamConfig["options"],
                            regex: uParamConfig["regex"],
                        });
                    }
                }

                configs = configs.concat({
                    slug: config.slug,
                    name: config.name ?? config.slug,
                    description: config.description,
                    requireRequests: config.requireRequests,
                    allowSelfApprovals: config.allowSelfApprovals,
                    timeout: config.timeout,
                    constraints: config.constraints,
                    runtime: config.runtime === "workflow" ? "workflow": "",
                    parameters: params,
                    entrypointFunc: exportName,
                });
            }
        }
    }
    return configs;
}

const files = process.argv.slice(2);
let taskConfigs = extractTaskConfigs(files);
console.log(JSON.stringify(taskConfigs));
