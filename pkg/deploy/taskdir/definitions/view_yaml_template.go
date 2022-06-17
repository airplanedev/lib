package definitions

const viewDefinitionTemplate = `# Full reference: <TODO: LINK TO VIEWS DOCS>

# Used by Airplane to identify your view. Do not change.
slug: {{.slug}}

# A human-readable name for your view.
name: {{.name}}

# A human-readable description for your view.
# description: "My Airplane view"

# The path to the directory containing the logic for this view. This
# can be absolute or relative to the location of the definition file.
entrypoint: {{.entrypoint}}

# A map of environment variables to use when for the view. The value
# should be an object; if specifying raw values, the value must be an object
# with ` + "`value`" + ` mapped to the value of the environment variable; if
# using config variables, the value must be an object with ` + "`config`" + `
# mapped to the name of the config variable.
# envVars:
#   ENV_VAR_FROM_CONFIG:
#     config: database_url
#   ENV_VAR_FROM_VALUE:
#     value: env_var_value
`
