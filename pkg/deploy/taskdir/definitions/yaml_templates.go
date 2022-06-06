package definitions

const definitionTemplate = `# Full reference: https://docs.airplane.dev/tasks/task-definition

# Used by Airplane to identify your task. Do not change.
slug: {{.slug}}

# A human-readable name for your task.
name: {{.name}}

# A human-readable description for your task.
# description: "My Airplane task"

# A list of inputs to your task.{{.paramsExtraDescription}}
# parameters:
# -
#   # An identifier for the parameter, which can be used in JavaScript
#   # templates (https://docs.airplane.dev/runbooks/javascript-templates).
#   slug: name
#   # A human-readable name for the parameter.
#   name: Name
#   # The type of parameter. Valid values: shorttext, longtext, sql, boolean,
#   # upload, integer, float, date, datetime, configvar.
#   type: shorttext
#   # A human-readable description of the parameter.
#   description: The user's name.
#   # The default value of the parameter.
#   default: Alfred Pennyworth
#   # Set to false to indicate that this parameter. is optional. Default: true.
#   required: false
#   # A list of options to constrain the parameter values. For configvar types,
#   # each option needs to be an object with a label (value to show to user) and
#   # a config (name of the config var). For all other types, each option can be
#   # a single value or an object with a label and a value.
#   options:
#   - Alfred Pennyworth
#   - label: Batman
#     value: Bruce Wayne
#   # A regular expression with which to validate parameter values.
#   regex: "^[a-zA-Z ]+$"
{{ .taskDefinition }}
# Set label constraints to restrict this task to run only on agents with
# matching labels.
# constraints:
#   aws-region: us-west-2

# Set to true to disable direct execution of this task. Default: false.
# requireRequests: true

# Set to false to disallow requesters from approving their own requests for
# this task. Default: true.
# allowSelfApprovals: false

# The maximum number of seconds the task should take before being timed out.
# Default: 3600.
# timeout: 1800
`

const nodeTemplate = `
# Configuration for a Node task.
node:
  # The path to the .ts or .js file containing the logic for this task. This
  # can be absolute or relative to the location of the definition file.
  entrypoint: {{.Entrypoint}}

  # The version of Node to use. Valid values: 12, 14, 15, 16, 18.
  nodeVersion: "{{.NodeVersion}}"

  # A map of environment variables to use when running the task. The value
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

const pythonTemplate = `
# Configuration for a Python task.
python:
  # The path to the .py file containing the logic for this task. This can be
  # absolute or relative to the location of the definition file.
  entrypoint: {{.Entrypoint}}

  # A map of environment variables to use when running the task. The value
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

const shellTemplate = `
# Configuration for a shell task.
shell:
  # The path to the .sh file containing the logic for this task. This can be
  # absolute or relative to the location of the definition file.
  entrypoint: {{.Entrypoint}}

  # A map of environment variables to use when running the task. The value
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
const shellParamsExtraDescription = ` Parameters are passed into your script
# as environment variables of form PARAM_{SLUG}, e.g. PARAM_USER_EMAIL.`

const imageTemplate = `
# Configuration for a Docker task.
docker:
  # The name of the image to use.
  image: {{.Image}}

  # Specify a Docker entrypoint to override the default image entrypoint.
  # entrypoint: bash

  # The Docker command to run. Supports JavaScript templates
  # (https://docs.airplane.dev/runbooks/javascript-templates).
  command: {{.Command}}

  # A map of environment variables to use when running the task. The value
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
const imageParamsExtraDescription = ` Parameters can be passed into the docker command
# as {{params.slug}}, e.g. command: /bin/my_command --id {{params.user_id}}.`

const sqlTemplate = `
# Configuration for a SQL task.
sql:
  # The name of a database resource.
  resource: ""

  # The path to the .sql file containing the logic for this task. This can be
  # absolute or relative to the location of the definition file. The contents
  # of the .sql file support JavaScript templates
  # (https://docs.airplane.dev/runbooks/javascript-templates).
  entrypoint: {{.Entrypoint}}

  # A map of query arguments that can be used to safely pass parameter inputs
  # to your query. Supports JavaScript templates
  # (https://docs.airplane.dev/runbooks/javascript-templates).
  # queryArgs:
  #   name: "{{"{{params.name}}"}}"

  # The transaction mode to use. Valid values: auto, readOnly, readWrite, none.
  # Default: auto.
  # transactionMode: readWrite

  # A list of config variables that this task can access.
  # configs:
  #   - API_KEY
  #   - DB_PASSWORD
`

const restTemplate = `
# Configuration for a REST task.
rest:
  # The name of a REST resource.
  resource: ""

  # The HTTP method to use. Valid values: GET, POST, PATCH, PUT, DELETE.
  method: {{.Method}}

  # The path to request. Your REST resource may specify a path prefix as part
  # of its base URL, in which case this path is joined to it. Airplane
  # recommends that this start with a leading slash. Supports JavaScript
  # templates (https://docs.airplane.dev/runbooks/javascript-templates).
  path: {{.Path}}

  # A map of URL parameters. Supports JavaScript templates
  # (https://docs.airplane.dev/runbooks/javascript-templates).
  # urlParams:
  #   page: 3

  # A map of request headers. Supports JavaScript templates
  # (https://docs.airplane.dev/runbooks/javascript-templates).
  # headers:
  #   X-Api-Key: api_key

  # The type of body that this request should send. Valid values: json, raw,
  # form-data, x-www-form-urlencoded.
  bodyType: {{.BodyType}}

  # The body of the request. Supports JavaScript templates
  # (https://docs.airplane.dev/runbooks/javascript-templates).
  body: "{{.Body}}"

  # A map of form values. Supports JavaScript templates
  # (https://docs.airplane.dev/runbooks/javascript-templates).
  # formData:
  #   name: Alfred Pennyworth

  # A list of config variables that this task can access.
  # configs:
  #   - API_KEY
  #   - DB_PASSWORD
`
