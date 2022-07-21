const path = require("node:path");

type ParamSchema = {
    name: string
    kind: string
}

type TaskConfig = {
    slug: string
    name: string
    description?: string
    parameters?: Record<string, ParamSchema>
    requireRequests?: boolean
    allowSelfApprovals?: boolean
    timeout?: number
    constraints?: Record<string, string>
    runtime?: "standard" | "workflow"
}

type TaskConfigWithBuildArgs = TaskConfig & {
    entrypointFunc: string
}

class AirplaneParser {
    files: string[]

    constructor(files: string[]) {
        this.files = files;
    }

    extractTaskConfigs(): TaskConfig[] {
        // Get function signatures and their parameter types
        // this.extractFSignatures();

        // Import each of the files
        let configs: TaskConfigWithBuildArgs[] = [];
        for (const file of this.files) {
            const resolvedPath = path.relative(__dirname, file);
            const lib = resolvedPath.replace(/.ts$/, "");
            const exports = require(`./${lib}`);

            for (const exportName in exports) {
                const item = exports[exportName];
                
                if ("__airplane" in item) {
                    const config = item.__airplane.config;

                    var params: Record<string, ParamSchema> = {};
                    for (var uParamSlug in config.parameters) {
                        const uParamConfig = config.parameters[uParamSlug];

                        if (typeof uParamConfig === "string") {
                            params[uParamSlug] = {
                                name: uParamSlug,
                                kind: uParamConfig,
                            };
                        } else {
                            params[uParamSlug] = {
                                name: uParamConfig["name"],
                                kind: uParamConfig["kind"]
                            };
                        }
                    }

                    configs = configs.concat({
                        slug: config.slug,
                        name: config.name,
                        parameters: params,
                        entrypointFunc: exportName,
                    });
                }
            }
        }

        return configs;
    }
}

const files = process.argv.slice(2);
let parser = new AirplaneParser(files);
let taskConfigs = parser.extractTaskConfigs();
console.log(JSON.stringify(taskConfigs));
