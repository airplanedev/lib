package definitions

import (
	"context"
	"testing"

	"github.com/airplanedev/lib/pkg/api"
	"github.com/airplanedev/lib/pkg/api/mock"
	"github.com/airplanedev/lib/pkg/build"
	"github.com/airplanedev/lib/pkg/utils/pointers"
	"github.com/stretchr/testify/require"
)

var fullYAML = []byte(
	`name: Hello World
slug: hello_world
description: A starter task.
parameters:
- name: Name
  slug: name
  type: shorttext
  description: Someone's name.
  default: World
  required: true
python:
  entrypoint: hello_world.py
timeout: 3600
runtime: durable
`)

var fullJSON = []byte(
	`{
	"name": "Hello World",
	"slug": "hello_world",
	"description": "A starter task.",
	"parameters": [
		{
			"name": "Name",
			"slug": "name",
			"type": "shorttext",
			"description": "Someone's name.",
			"default": "World",
			"required": true
		}
	],
	"python": {
		"entrypoint": "hello_world.py"
	},
	"timeout": 3600,
	"runtime": "durable"
}`)

var yamlWithDefault = []byte(
	`name: Hello World
slug: hello_world
description: A starter task.
parameters:
- name: Name
  slug: name
  type: shorttext
  description: Someone's name.
  default: World
python:
  entrypoint: hello_world.py
timeout: 3600
`)

var jsonWithDefault = []byte(
	`{
	"name": "Hello World",
	"slug": "hello_world",
	"description": "A starter task.",
	"parameters": [
		{
			"name": "Name",
			"slug": "name",
			"type": "shorttext",
			"description": "Someone's name.",
			"default": "World"
		}
	],
	"python": {
		"entrypoint": "hello_world.py"
	},
	"timeout": 3600
}`)

var fullDef = Definition_0_3{
	Name:        "Hello World",
	Slug:        "hello_world",
	Description: "A starter task.",
	Parameters: []ParameterDefinition_0_3{
		{
			Name:        "Name",
			Slug:        "name",
			Type:        "shorttext",
			Description: "Someone's name.",
			Default:     "World",
			Required:    pointers.Bool(true),
		},
	},
	Python: &PythonDefinition_0_3{
		Entrypoint: "hello_world.py",
	},
	Runtime: "durable",
	Timeout: 3600,
}

var defWithDefault = Definition_0_3{
	Name:        "Hello World",
	Slug:        "hello_world",
	Description: "A starter task.",
	Parameters: []ParameterDefinition_0_3{
		{
			Name:        "Name",
			Slug:        "name",
			Type:        "shorttext",
			Description: "Someone's name.",
			Default:     "World",
		},
	},
	Python: &PythonDefinition_0_3{
		Entrypoint: "hello_world.py",
	},
	Timeout: 3600,
}

func TestDefinitionSerialization_0_3(t *testing.T) {
	// marshalling tests
	for _, test := range []struct {
		name     string
		format   DefFormat
		def      Definition_0_3
		expected []byte
	}{
		{
			name:     "marshal yaml",
			format:   DefFormatYAML,
			def:      fullDef,
			expected: fullYAML,
		},
		{
			name:     "marshal json",
			format:   DefFormatJSON,
			def:      fullDef,
			expected: fullJSON,
		},
		{
			name:   "marshal yaml with multiline",
			format: DefFormatYAML,
			def: Definition_0_3{
				Name: "REST task",
				Slug: "rest_task",
				REST: &RESTDefinition_0_3{
					Resource: "httpbin",
					Method:   "POST",
					Path:     "/post",
					BodyType: "json",
					Body:     "{\n  \"name\": \"foo\",\n  \"number\": 30\n}\n",
				},
				Timeout: 300,
			},
			expected: []byte(
				`name: REST task
slug: rest_task
rest:
  resource: httpbin
  method: POST
  path: /post
  bodyType: json
  body: |
    {
      "name": "foo",
      "number": 30
    }
timeout: 300
`),
		},
		{
			name:   "marshal json with multiline",
			format: DefFormatJSON,
			def: Definition_0_3{
				Name: "REST task",
				Slug: "rest_task",
				REST: &RESTDefinition_0_3{
					Resource: "httpbin",
					Method:   "POST",
					Path:     "/post",
					BodyType: "json",
					Body:     "{\n  \"name\": \"foo\",\n  \"number\": 30\n}\n",
				},
				Timeout: 300,
			},
			expected: []byte(
				`{
	"name": "REST task",
	"slug": "rest_task",
	"rest": {
		"resource": "httpbin",
		"method": "POST",
		"path": "/post",
		"bodyType": "json",
		"body": "{\n  \"name\": \"foo\",\n  \"number\": 30\n}\n"
	},
	"timeout": 300
}`),
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			assert := require.New(t)
			bytestr, err := test.def.Marshal(test.format)
			assert.NoError(err)
			assert.Equal(test.expected, bytestr)
		})
	}

	// unmarshalling tests
	for _, test := range []struct {
		name     string
		format   DefFormat
		bytestr  []byte
		expected Definition_0_3
	}{
		{
			name:     "unmarshal yaml",
			format:   DefFormatYAML,
			bytestr:  fullYAML,
			expected: fullDef,
		},
		{
			name:     "unmarshal json",
			format:   DefFormatJSON,
			bytestr:  fullJSON,
			expected: fullDef,
		},
		{
			name:     "unmarshal yaml with default",
			format:   DefFormatYAML,
			bytestr:  yamlWithDefault,
			expected: defWithDefault,
		},
		{
			name:     "unmarshal json with default",
			format:   DefFormatJSON,
			bytestr:  jsonWithDefault,
			expected: defWithDefault,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			assert := require.New(t)
			d := Definition_0_3{}
			err := d.Unmarshal(test.format, test.bytestr)
			assert.NoError(err)
			assert.Equal(test.expected, d)
		})
	}
}

func TestTaskToDefinition_0_3(t *testing.T) {
	for _, test := range []struct {
		name       string
		task       api.Task
		definition Definition_0_3
		resources  []api.Resource
	}{
		{
			name: "python task",
			task: api.Task{
				Name:        "Python Task",
				Slug:        "python_task",
				Description: "A task for testing",
				Arguments:   []string{"{{JSON.stringify(params)}}"},
				Kind:        build.TaskKindPython,
				KindOptions: build.KindOptions{
					"entrypoint": "main.py",
				},
				Env: api.TaskEnv{
					"value": api.EnvVarValue{
						Value: pointers.String("value"),
					},
					"config": api.EnvVarValue{
						Config: pointers.String("config"),
					},
				},
			},
			definition: Definition_0_3{
				Name:        "Python Task",
				Slug:        "python_task",
				Description: "A task for testing",
				Python: &PythonDefinition_0_3{
					Entrypoint: "main.py",
					EnvVars: api.TaskEnv{
						"value": api.EnvVarValue{
							Value: pointers.String("value"),
						},
						"config": api.EnvVarValue{
							Config: pointers.String("config"),
						},
					},
				},
			},
		},
		{
			name: "node task",
			task: api.Task{
				Name:      "Node Task",
				Slug:      "node_task",
				Arguments: []string{"{{JSON.stringify(params)}}"},
				Kind:      build.TaskKindNode,
				KindOptions: build.KindOptions{
					"entrypoint":  "main.ts",
					"nodeVersion": "14",
				},
				Env: api.TaskEnv{
					"value": api.EnvVarValue{
						Value: pointers.String("value"),
					},
					"config": api.EnvVarValue{
						Config: pointers.String("config"),
					},
				},
			},
			definition: Definition_0_3{
				Name: "Node Task",
				Slug: "node_task",
				Node: &NodeDefinition_0_3{
					Entrypoint:  "main.ts",
					NodeVersion: "14",
					EnvVars: api.TaskEnv{
						"value": api.EnvVarValue{
							Value: pointers.String("value"),
						},
						"config": api.EnvVarValue{
							Config: pointers.String("config"),
						},
					},
				},
			},
		},
		{
			name: "shell task",
			task: api.Task{
				Name:      "Shell Task",
				Slug:      "shell_task",
				Arguments: []string{},
				Kind:      build.TaskKindShell,
				KindOptions: build.KindOptions{
					"entrypoint": "main.sh",
				},
				Env: api.TaskEnv{
					"value": api.EnvVarValue{
						Value: pointers.String("value"),
					},
					"config": api.EnvVarValue{
						Config: pointers.String("config"),
					},
				},
			},
			definition: Definition_0_3{
				Name: "Shell Task",
				Slug: "shell_task",
				Shell: &ShellDefinition_0_3{
					Entrypoint: "main.sh",
					EnvVars: api.TaskEnv{
						"value": api.EnvVarValue{
							Value: pointers.String("value"),
						},
						"config": api.EnvVarValue{
							Config: pointers.String("config"),
						},
					},
				},
			},
		},
		{
			name: "image task",
			task: api.Task{
				Name:        "Image Task",
				Slug:        "image_task",
				Command:     []string{"bash"},
				Arguments:   []string{"-c", `echo "foobar"`},
				Kind:        build.TaskKindImage,
				KindOptions: build.KindOptions{},
				Image:       pointers.String("ubuntu:latest"),
				Env: api.TaskEnv{
					"value": api.EnvVarValue{
						Value: pointers.String("value"),
					},
					"config": api.EnvVarValue{
						Config: pointers.String("config"),
					},
				},
			},
			definition: Definition_0_3{
				Name: "Image Task",
				Slug: "image_task",
				Image: &ImageDefinition_0_3{
					Image:      "ubuntu:latest",
					Entrypoint: "bash",
					Command:    `-c 'echo "foobar"'`,
					EnvVars: api.TaskEnv{
						"value": api.EnvVarValue{
							Value: pointers.String("value"),
						},
						"config": api.EnvVarValue{
							Config: pointers.String("config"),
						},
					},
				},
			},
		},
		{
			name: "rest task",
			resources: []api.Resource{
				{
					ID:   "res20220111foobarx",
					Name: "httpbin",
				},
			},
			task: api.Task{
				Name:      "REST Task",
				Slug:      "rest_task",
				Arguments: []string{"{{__stdAPIRequest}}"},
				Kind:      build.TaskKindREST,
				KindOptions: build.KindOptions{
					"method": "GET",
					"path":   "/get",
					"urlParams": map[string]interface{}{
						"foo": "bar",
					},
					"headers": map[string]interface{}{
						"bar": "foo",
					},
					"bodyType": "json",
					"body":     "",
					"formData": map[string]interface{}{},
				},
				Resources: map[string]string{
					"rest": "res20220111foobarx",
				},
			},
			definition: Definition_0_3{
				Name: "REST Task",
				Slug: "rest_task",
				REST: &RESTDefinition_0_3{
					Resource: "httpbin",
					Method:   "GET",
					Path:     "/get",
					URLParams: map[string]interface{}{
						"foo": "bar",
					},
					Headers: map[string]interface{}{
						"bar": "foo",
					},
					BodyType: "json",
					Body:     "",
					FormData: map[string]interface{}{},
					Configs:  []string{},
				},
			},
		},
		{
			name: "check parameters",
			task: api.Task{
				Name: "Test Task",
				Slug: "test_task",
				Parameters: []api.Parameter{
					{
						Name: "Required boolean",
						Slug: "required_boolean",
						Type: api.TypeBoolean,
						Desc: "A required boolean.",
					},
					{
						Name:    "Short text",
						Slug:    "short_text",
						Type:    api.TypeString,
						Default: "foobar",
					},
					{
						Name:      "SQL",
						Slug:      "sql",
						Type:      api.TypeString,
						Component: api.ComponentEditorSQL,
					},
					{
						Name:      "Optional long text",
						Slug:      "optional_long_text",
						Type:      api.TypeString,
						Component: api.ComponentTextarea,
						Constraints: api.Constraints{
							Optional: true,
						},
					},
					{
						Name: "Options",
						Slug: "options",
						Type: api.TypeString,
						Constraints: api.Constraints{
							Options: []api.ConstraintOption{
								{
									Label: "one",
									Value: 1,
								},
								{
									Label: "two",
									Value: 2,
								},
								{
									Label: "three",
									Value: 3,
								},
								{
									Label: "config",
									Value: map[string]interface{}{
										"__airplaneType": "configvar",
										"name":           "config_name",
									},
								},
							},
						},
					},
					{
						Name: "Regex",
						Slug: "regex",
						Type: api.TypeString,
						Constraints: api.Constraints{
							Regex: "foo.*",
						},
					},
				},
				Arguments: []string{"{{JSON.stringify(params)}}"},
				Kind:      build.TaskKindPython,
				KindOptions: build.KindOptions{
					"entrypoint": "main.py",
				},
			},
			definition: Definition_0_3{
				Name: "Test Task",
				Slug: "test_task",
				Parameters: []ParameterDefinition_0_3{
					{
						Name:        "Required boolean",
						Slug:        "required_boolean",
						Type:        "boolean",
						Description: "A required boolean.",
					},
					{
						Name:    "Short text",
						Slug:    "short_text",
						Type:    "shorttext",
						Default: "foobar",
					},
					{
						Name: "SQL",
						Slug: "sql",
						Type: "sql",
					},
					{
						Name:     "Optional long text",
						Slug:     "optional_long_text",
						Type:     "longtext",
						Required: pointers.Bool(false),
					},
					{
						Name: "Options",
						Slug: "options",
						Type: "shorttext",
						Options: []OptionDefinition_0_3{
							{
								Label: "one",
								Value: 1,
							},
							{
								Label: "two",
								Value: 2,
							},
							{
								Label: "three",
								Value: 3,
							},
							{
								Label:  "config",
								Config: pointers.String("config_name"),
							},
						},
					},
					{
						Name:  "Regex",
						Slug:  "regex",
						Type:  "shorttext",
						Regex: "foo.*",
					},
				},
				Python: &PythonDefinition_0_3{
					Entrypoint: "main.py",
				},
			},
		},
		{
			name: "check execute rules",
			task: api.Task{
				Name:      "Test Task",
				Slug:      "test_task",
				Arguments: []string{"{{JSON.stringify(params)}}"},
				Kind:      build.TaskKindPython,
				KindOptions: build.KindOptions{
					"entrypoint": "main.py",
				},
				ExecuteRules: api.ExecuteRules{
					DisallowSelfApprove: true,
					RequireRequests:     true,
				},
			},
			definition: Definition_0_3{
				Name: "Test Task",
				Slug: "test_task",
				Python: &PythonDefinition_0_3{
					Entrypoint: "main.py",
				},
				RequireRequests:    true,
				AllowSelfApprovals: pointers.Bool(false),
			},
		},
		{
			name: "check default execute rules",
			task: api.Task{
				Name:      "Test Task",
				Slug:      "test_task",
				Arguments: []string{"{{JSON.stringify(params)}}"},
				Kind:      build.TaskKindPython,
				KindOptions: build.KindOptions{
					"entrypoint": "main.py",
				},
				ExecuteRules: api.ExecuteRules{
					DisallowSelfApprove: false,
					RequireRequests:     false,
				},
			},
			definition: Definition_0_3{
				Name: "Test Task",
				Slug: "test_task",
				Python: &PythonDefinition_0_3{
					Entrypoint: "main.py",
				},
				RequireRequests:    false,
				AllowSelfApprovals: nil,
			},
		},
		{
			name: "check configs",
			resources: []api.Resource{
				{
					ID:   "res20220111foobarx",
					Name: "httpbin",
				},
			},
			task: api.Task{
				Name:      "REST Task",
				Slug:      "rest_task",
				Arguments: []string{"{{__stdAPIRequest}}"},
				Configs: []api.ConfigAttachment{
					{
						NameTag: "CONFIG_NAME_1",
					},
					{
						NameTag: "CONFIG_NAME_2",
					},
				},
				Kind: build.TaskKindREST,
				KindOptions: build.KindOptions{
					"method": "GET",
					"path":   "/get",
					"urlParams": map[string]interface{}{
						"foo": "bar",
					},
					"headers": map[string]interface{}{
						"bar": "foo",
					},
					"bodyType": "json",
					"body":     "",
					"formData": map[string]interface{}{},
				},
				Resources: map[string]string{
					"rest": "res20220111foobarx",
				},
			},
			definition: Definition_0_3{
				Name: "REST Task",
				Slug: "rest_task",
				REST: &RESTDefinition_0_3{
					Resource: "httpbin",
					Method:   "GET",
					Path:     "/get",
					URLParams: map[string]interface{}{
						"foo": "bar",
					},
					Headers: map[string]interface{}{
						"bar": "foo",
					},
					BodyType: "json",
					Body:     "",
					FormData: map[string]interface{}{},
					Configs:  []string{"CONFIG_NAME_1", "CONFIG_NAME_2"},
				},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			assert := require.New(t)
			ctx := context.Background()
			client := &mock.MockClient{
				Resources: test.resources,
			}
			d, err := NewDefinitionFromTask_0_3(ctx, client, test.task)
			assert.NoError(err)
			assert.Equal(test.definition, d)
		})
	}
}

func TestDefinitionToUpdateTaskRequest_0_3(t *testing.T) {
	for _, test := range []struct {
		name       string
		definition Definition_0_3
		request    api.UpdateTaskRequest
		resources  []api.Resource
	}{
		{
			name: "python task",
			definition: Definition_0_3{
				Name:        "Test Task",
				Slug:        "test_task",
				Description: "A task for testing",
				Python: &PythonDefinition_0_3{
					Entrypoint: "main.py",
				},
			},
			request: api.UpdateTaskRequest{
				Name:        "Test Task",
				Slug:        "test_task",
				Description: "A task for testing",
				Configs:     &[]api.ConfigAttachment{},
				Parameters:  []api.Parameter{},
				Kind:        build.TaskKindPython,
				KindOptions: build.KindOptions{
					"entrypoint": "main.py",
				},
				ExecuteRules: api.UpdateExecuteRulesRequest{
					DisallowSelfApprove: pointers.Bool(false),
					RequireRequests:     pointers.Bool(false),
				},
			},
		},
		{
			name: "node task",
			definition: Definition_0_3{
				Name: "Node Task",
				Slug: "node_task",
				Node: &NodeDefinition_0_3{
					Entrypoint:  "main.ts",
					NodeVersion: "14",
				},
			},
			request: api.UpdateTaskRequest{
				Name:       "Node Task",
				Slug:       "node_task",
				Parameters: []api.Parameter{},
				Configs:    &[]api.ConfigAttachment{},
				Kind:       build.TaskKindNode,
				KindOptions: build.KindOptions{
					"entrypoint":  "main.ts",
					"nodeVersion": "14",
				},
				ExecuteRules: api.UpdateExecuteRulesRequest{
					DisallowSelfApprove: pointers.Bool(false),
					RequireRequests:     pointers.Bool(false),
				},
			},
		},
		{
			name: "shell task",
			definition: Definition_0_3{
				Name: "Shell Task",
				Slug: "shell_task",
				Shell: &ShellDefinition_0_3{
					Entrypoint: "main.sh",
				},
			},
			request: api.UpdateTaskRequest{
				Name:       "Shell Task",
				Slug:       "shell_task",
				Parameters: []api.Parameter{},
				Configs:    &[]api.ConfigAttachment{},
				Kind:       build.TaskKindShell,
				KindOptions: build.KindOptions{
					"entrypoint": "main.sh",
				},
				ExecuteRules: api.UpdateExecuteRulesRequest{
					DisallowSelfApprove: pointers.Bool(false),
					RequireRequests:     pointers.Bool(false),
				},
			},
		},
		{
			name: "image task",
			definition: Definition_0_3{
				Name: "Image Task",
				Slug: "image_task",
				Image: &ImageDefinition_0_3{
					Image:      "ubuntu:latest",
					Entrypoint: "bash",
					Command:    `-c 'echo "foobar"'`,
				},
			},
			request: api.UpdateTaskRequest{
				Name:       "Image Task",
				Slug:       "image_task",
				Parameters: []api.Parameter{},
				Configs:    &[]api.ConfigAttachment{},
				Command:    []string{"bash"},
				Arguments:  []string{"-c", `echo "foobar"`},
				Kind:       build.TaskKindImage,
				Image:      pointers.String("ubuntu:latest"),
				ExecuteRules: api.UpdateExecuteRulesRequest{
					DisallowSelfApprove: pointers.Bool(false),
					RequireRequests:     pointers.Bool(false),
				},
			},
		},
		{
			name: "rest task",
			definition: Definition_0_3{
				Name: "REST Task",
				Slug: "rest_task",
				REST: &RESTDefinition_0_3{
					Resource: "rest",
					Method:   "POST",
					Path:     "/post",
					BodyType: "json",
					Body:     `{"foo": "bar"}`,
				},
			},
			request: api.UpdateTaskRequest{
				Name:       "REST Task",
				Slug:       "rest_task",
				Parameters: []api.Parameter{},
				Configs:    &[]api.ConfigAttachment{},
				Kind:       build.TaskKindREST,
				KindOptions: build.KindOptions{
					"method":    "POST",
					"path":      "/post",
					"urlParams": map[string]interface{}{},
					"headers":   map[string]interface{}{},
					"bodyType":  "json",
					"body":      `{"foo": "bar"}`,
					"formData":  map[string]interface{}{},
				},
				Resources: map[string]string{
					"rest": "rest_id",
				},
				ExecuteRules: api.UpdateExecuteRulesRequest{
					DisallowSelfApprove: pointers.Bool(false),
					RequireRequests:     pointers.Bool(false),
				},
			},
			resources: []api.Resource{
				{
					ID:   "rest_id",
					Name: "rest",
				},
			},
		},
		{
			name: "test update execute rules",
			definition: Definition_0_3{
				Name:        "Test Task",
				Slug:        "test_task",
				Description: "A task for testing",
				Python: &PythonDefinition_0_3{
					Entrypoint: "main.py",
				},
				RequireRequests:    true,
				AllowSelfApprovals: pointers.Bool(false),
			},
			request: api.UpdateTaskRequest{
				Name:        "Test Task",
				Slug:        "test_task",
				Parameters:  []api.Parameter{},
				Configs:     &[]api.ConfigAttachment{},
				Description: "A task for testing",
				Kind:        build.TaskKindPython,
				KindOptions: build.KindOptions{
					"entrypoint": "main.py",
				},
				ExecuteRules: api.UpdateExecuteRulesRequest{
					DisallowSelfApprove: pointers.Bool(true),
					RequireRequests:     pointers.Bool(true),
				},
			},
		},
		{
			name: "test update default execute rules",
			definition: Definition_0_3{
				Name:        "Test Task",
				Slug:        "test_task",
				Description: "A task for testing",
				Python: &PythonDefinition_0_3{
					Entrypoint: "main.py",
				},
				RequireRequests:    false,
				AllowSelfApprovals: nil,
			},
			request: api.UpdateTaskRequest{
				Name:        "Test Task",
				Slug:        "test_task",
				Parameters:  []api.Parameter{},
				Configs:     &[]api.ConfigAttachment{},
				Description: "A task for testing",
				Kind:        build.TaskKindPython,
				KindOptions: build.KindOptions{
					"entrypoint": "main.py",
				},
				ExecuteRules: api.UpdateExecuteRulesRequest{
					DisallowSelfApprove: pointers.Bool(false),
					RequireRequests:     pointers.Bool(false),
				},
			},
		},
		{
			name: "check parameters",
			definition: Definition_0_3{
				Name: "Test Task",
				Slug: "test_task",
				Parameters: []ParameterDefinition_0_3{
					{
						Name:        "Required boolean",
						Slug:        "required_boolean",
						Type:        "boolean",
						Description: "A required boolean.",
					},
					{
						Name:    "Short text",
						Slug:    "short_text",
						Type:    "shorttext",
						Default: "foobar",
					},
					{
						Name: "SQL",
						Slug: "sql",
						Type: "sql",
					},
					{
						Name:     "Optional long text",
						Slug:     "optional_long_text",
						Type:     "longtext",
						Required: pointers.Bool(false),
					},
					{
						Name: "Options",
						Slug: "options",
						Type: "shorttext",
						Options: []OptionDefinition_0_3{
							{
								Label: "one",
								Value: 1,
							},
							{
								Label: "two",
								Value: 2,
							},
							{
								Label: "three",
								Value: 3,
							},
							{
								Label:  "config",
								Config: pointers.String("config_name"),
							},
						},
					},
					{
						Name:  "Regex",
						Slug:  "regex",
						Type:  "shorttext",
						Regex: "foo.*",
					},
				},
				Python: &PythonDefinition_0_3{
					Entrypoint: "main.py",
				},
			},
			request: api.UpdateTaskRequest{
				Name: "Test Task",
				Slug: "test_task",
				Parameters: []api.Parameter{
					{
						Name: "Required boolean",
						Slug: "required_boolean",
						Type: api.TypeBoolean,
						Desc: "A required boolean.",
					},
					{
						Name:    "Short text",
						Slug:    "short_text",
						Type:    api.TypeString,
						Default: "foobar",
					},
					{
						Name:      "SQL",
						Slug:      "sql",
						Type:      api.TypeString,
						Component: api.ComponentEditorSQL,
					},
					{
						Name:      "Optional long text",
						Slug:      "optional_long_text",
						Type:      api.TypeString,
						Component: api.ComponentTextarea,
						Constraints: api.Constraints{
							Optional: true,
						},
					},
					{
						Name: "Options",
						Slug: "options",
						Type: api.TypeString,
						Constraints: api.Constraints{
							Options: []api.ConstraintOption{
								{
									Label: "one",
									Value: 1,
								},
								{
									Label: "two",
									Value: 2,
								},
								{
									Label: "three",
									Value: 3,
								},
								{
									Label: "config",
									Value: map[string]interface{}{
										"__airplaneType": "configvar",
										"name":           "config_name",
									},
								},
							},
						},
					},
					{
						Name: "Regex",
						Slug: "regex",
						Type: api.TypeString,
						Constraints: api.Constraints{
							Regex: "foo.*",
						},
					},
				},
				Configs: &[]api.ConfigAttachment{},
				Kind:    build.TaskKindPython,
				KindOptions: build.KindOptions{
					"entrypoint": "main.py",
				},
				ExecuteRules: api.UpdateExecuteRulesRequest{
					DisallowSelfApprove: pointers.Bool(false),
					RequireRequests:     pointers.Bool(false),
				},
			},
		},
		{
			name: "check configs",
			definition: Definition_0_3{
				Name: "REST Task",
				Slug: "rest_task",
				REST: &RESTDefinition_0_3{
					Resource: "rest",
					Method:   "POST",
					Path:     "/post",
					BodyType: "json",
					Body:     `{"foo": "bar"}`,
					Configs:  []string{"CONFIG_VARIABLE_1", "CONFIG_VARIABLE_2"},
				},
			},
			request: api.UpdateTaskRequest{
				Name:       "REST Task",
				Slug:       "rest_task",
				Parameters: []api.Parameter{},
				Configs: &[]api.ConfigAttachment{
					{
						NameTag: "CONFIG_VARIABLE_1",
					},
					{
						NameTag: "CONFIG_VARIABLE_2",
					},
				},
				Kind: build.TaskKindREST,
				KindOptions: build.KindOptions{
					"method":    "POST",
					"path":      "/post",
					"urlParams": map[string]interface{}{},
					"headers":   map[string]interface{}{},
					"bodyType":  "json",
					"body":      `{"foo": "bar"}`,
					"formData":  map[string]interface{}{},
				},
				Resources: map[string]string{
					"rest": "rest_id",
				},
				ExecuteRules: api.UpdateExecuteRulesRequest{
					DisallowSelfApprove: pointers.Bool(false),
					RequireRequests:     pointers.Bool(false),
				},
			},
			resources: []api.Resource{
				{
					ID:   "rest_id",
					Name: "rest",
				},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			assert := require.New(t)
			ctx := context.Background()
			client := &mock.MockClient{
				Resources: test.resources,
			}
			req, err := test.definition.GetUpdateTaskRequest(ctx, client)
			assert.NoError(err)
			assert.Equal(test.request, req)
		})
	}
}
