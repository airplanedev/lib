package definitions

import (
	"context"
	"testing"

	"github.com/airplanedev/lib/pkg/api"
	"github.com/airplanedev/lib/pkg/api/mock"
	"github.com/airplanedev/lib/pkg/build"
	"github.com/stretchr/testify/require"
)

func newBoolPtr(v bool) *bool {
	return &v
}

func newStringPtr(v string) *string {
	return &v
}

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
  arguments:
  - "{{JSON.stringify(params)}}"
timeout: 3600
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
		"entrypoint": "hello_world.py",
		"arguments": [
			"{{JSON.stringify(params)}}"
		]
	},
	"timeout": 3600
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
  arguments:
  - "{{JSON.stringify(params)}}"
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
		"entrypoint": "hello_world.py",
		"arguments": [
			"{{JSON.stringify(params)}}"
		]
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
			Required:    newBoolPtr(true),
		},
	},
	Python: &PythonDefinition_0_3{
		Entrypoint: "hello_world.py",
		Arguments:  []string{"{{JSON.stringify(params)}}"},
	},
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
		Arguments:  []string{"{{JSON.stringify(params)}}"},
	},
	Timeout: 3600,
}

func TestDefinitionSerialization_0_3(t *testing.T) {
	// marshalling tests
	for _, test := range []struct {
		name     string
		format   TaskDefFormat
		def      Definition_0_3
		expected []byte
	}{
		{
			name:     "marshal yaml",
			format:   TaskDefFormatYAML,
			def:      fullDef,
			expected: fullYAML,
		},
		{
			name:     "marshal json",
			format:   TaskDefFormatJSON,
			def:      fullDef,
			expected: fullJSON,
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
		format   TaskDefFormat
		bytestr  []byte
		expected Definition_0_3
	}{
		{
			name:     "unmarshal yaml",
			format:   TaskDefFormatYAML,
			bytestr:  fullYAML,
			expected: fullDef,
		},
		{
			name:     "unmarshal json",
			format:   TaskDefFormatJSON,
			bytestr:  fullJSON,
			expected: fullDef,
		},
		{
			name:     "unmarshal yaml with default",
			format:   TaskDefFormatYAML,
			bytestr:  yamlWithDefault,
			expected: defWithDefault,
		},
		{
			name:     "unmarshal json with default",
			format:   TaskDefFormatJSON,
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
		users      []api.User
		groups     []api.Group
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
			},
			definition: Definition_0_3{
				Name:        "Python Task",
				Slug:        "python_task",
				Description: "A task for testing",
				Python: &PythonDefinition_0_3{
					Arguments:  []string{"{{JSON.stringify(params)}}"},
					Entrypoint: "main.py",
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
			},
			definition: Definition_0_3{
				Name: "Node Task",
				Slug: "node_task",
				Node: &NodeDefinition_0_3{
					Arguments:   []string{"{{JSON.stringify(params)}}"},
					Entrypoint:  "main.ts",
					NodeVersion: "14",
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
			},
			definition: Definition_0_3{
				Name: "Shell Task",
				Slug: "shell_task",
				Shell: &ShellDefinition_0_3{
					Arguments:  []string{},
					Entrypoint: "main.sh",
				},
			},
		},
		{
			name: "image task",
			task: api.Task{
				Name:        "Image Task",
				Slug:        "image_task",
				Command:     []string{"bash"},
				Arguments:   []string{"-c", "echo 'foobar'"},
				Kind:        build.TaskKindImage,
				KindOptions: build.KindOptions{},
				Image:       newStringPtr("ubuntu:latest"),
			},
			definition: Definition_0_3{
				Name: "Image Task",
				Slug: "image_task",
				Image: &ImageDefinition_0_3{
					Image:      "ubuntu:latest",
					Entrypoint: "bash",
					Command:    []string{"-c", "echo 'foobar'"},
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
					"urlParams": map[string]string{
						"foo": "bar",
					},
					"headers": map[string]string{
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
					URLParams: map[string]string{
						"foo": "bar",
					},
					Headers: map[string]string{
						"bar": "foo",
					},
					BodyType: "json",
					Body:     "",
					FormData: map[string]interface{}{},
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
						Required: newBoolPtr(false),
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
					Arguments:  []string{"{{JSON.stringify(params)}}"},
					Entrypoint: "main.py",
				},
			},
		},
		{
			name: "check permissions",
			task: api.Task{
				Name:      "Test Task",
				Slug:      "test_task",
				Arguments: []string{"{{JSON.stringify(params)}}"},
				Kind:      build.TaskKindPython,
				KindOptions: build.KindOptions{
					"entrypoint": "main.py",
				},
				RequireExplicitPermissions: true,
				Permissions: []api.Permission{
					{
						SubUserID: newStringPtr("1000"),
						RoleID:    api.RoleTaskAdmin,
					},
					{
						SubUserID: newStringPtr("1001"),
						RoleID:    api.RoleTaskExecuter,
					},
					{
						SubGroupID: newStringPtr("2000"),
						RoleID:     api.RoleTaskRequester,
					},
					{
						SubUserID: newStringPtr("1002"),
						RoleID:    api.RoleTaskViewer,
					},
				},
			},
			definition: Definition_0_3{
				Name: "Test Task",
				Slug: "test_task",
				Python: &PythonDefinition_0_3{
					Arguments:  []string{"{{JSON.stringify(params)}}"},
					Entrypoint: "main.py",
				},
				Permissions: &PermissionDefinition_0_3{
					Admins:     []string{"gary@dragonhawk5.com"},
					Executers:  []string{"quinn@infinityguard.org"},
					Requesters: []string{"Team Squad"},
					Viewers:    []string{"kvn@kvn.net"},
				},
			},
			users: []api.User{
				{
					ID:    "1000",
					Name:  "Gary Goodspeed",
					Email: "gary@dragonhawk5.com",
				},
				{
					ID:    "1001",
					Name:  "Quinn Ergon",
					Email: "quinn@infinityguard.org",
				},
				{
					ID:    "1002",
					Name:  "KVN",
					Email: "kvn@kvn.net",
				},
			},
			groups: []api.Group{
				{
					ID:   "2000",
					Name: "Team Squad",
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
					Arguments:  []string{"{{JSON.stringify(params)}}"},
					Entrypoint: "main.py",
				},
				RequireRequests:    true,
				AllowSelfApprovals: newBoolPtr(false),
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
					Arguments:  []string{"{{JSON.stringify(params)}}"},
					Entrypoint: "main.py",
				},
				RequireRequests:    false,
				AllowSelfApprovals: nil,
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			assert := require.New(t)
			ctx := context.Background()
			client := &mock.MockClient{
				Resources: test.resources,
				Users:     test.users,
				Groups:    test.groups,
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
		users      []api.User
		groups     []api.Group
	}{
		{
			name: "python task",
			definition: Definition_0_3{
				Name:        "Test Task",
				Slug:        "test_task",
				Description: "A task for testing",
				Python: &PythonDefinition_0_3{
					Arguments:  []string{"{{JSON.stringify(params)}}"},
					Entrypoint: "main.py",
				},
			},
			request: api.UpdateTaskRequest{
				Name:        "Test Task",
				Slug:        "test_task",
				Description: "A task for testing",
				Parameters:  []api.Parameter{},
				Arguments:   []string{"{{JSON.stringify(params)}}"},
				Kind:        build.TaskKindPython,
				KindOptions: build.KindOptions{
					"entrypoint": "main.py",
				},
			},
		},
		{
			name: "node task",
			definition: Definition_0_3{
				Name: "Node Task",
				Slug: "node_task",
				Node: &NodeDefinition_0_3{
					Arguments:   []string{"{{JSON.stringify(params)}}"},
					Entrypoint:  "main.ts",
					NodeVersion: "14",
				},
			},
			request: api.UpdateTaskRequest{
				Name:       "Node Task",
				Slug:       "node_task",
				Parameters: []api.Parameter{},
				Arguments:  []string{"{{JSON.stringify(params)}}"},
				Kind:       build.TaskKindNode,
				KindOptions: build.KindOptions{
					"entrypoint":  "main.ts",
					"nodeVersion": "14",
				},
			},
		},
		{
			name: "shell task",
			definition: Definition_0_3{
				Name: "Shell Task",
				Slug: "shell_task",
				Shell: &ShellDefinition_0_3{
					Arguments:  []string{},
					Entrypoint: "main.sh",
				},
			},
			request: api.UpdateTaskRequest{
				Name:       "Shell Task",
				Slug:       "shell_task",
				Parameters: []api.Parameter{},
				Arguments:  []string{},
				Kind:       build.TaskKindShell,
				KindOptions: build.KindOptions{
					"entrypoint": "main.sh",
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
					Command:    []string{"-c", "echo 'foobar'"},
				},
			},
			request: api.UpdateTaskRequest{
				Name:       "Image Task",
				Slug:       "image_task",
				Parameters: []api.Parameter{},
				Command:    []string{"bash"},
				Arguments:  []string{"-c", "echo 'foobar'"},
				Kind:       build.TaskKindImage,
				Image:      newStringPtr("ubuntu:latest"),
			},
		},
		{
			name: "test update permissions",
			definition: Definition_0_3{
				Name:        "Test Task",
				Slug:        "test_task",
				Description: "A task for testing",
				Python: &PythonDefinition_0_3{
					Arguments:  []string{"{{JSON.stringify(params)}}"},
					Entrypoint: "main.py",
				},
				Permissions: &PermissionDefinition_0_3{
					Admins:     []string{"gary@dragonhawk5.com"},
					Executers:  []string{"quinn@infinityguard.org"},
					Requesters: []string{"Team Squad"},
					Viewers:    []string{"kvn@kvn.net"},
				},
				Constraints: &api.RunConstraints{},
			},
			request: api.UpdateTaskRequest{
				Name:        "Test Task",
				Slug:        "test_task",
				Parameters:  []api.Parameter{},
				Description: "A task for testing",
				Arguments:   []string{"{{JSON.stringify(params)}}"},
				Kind:        build.TaskKindPython,
				KindOptions: build.KindOptions{
					"entrypoint": "main.py",
				},
				RequireExplicitPermissions: true,
				Permissions: []api.Permission{
					{
						SubUserID: newStringPtr("1002"),
						RoleID:    api.RoleTaskViewer,
					},
					{
						SubGroupID: newStringPtr("2000"),
						RoleID:     api.RoleTaskRequester,
					},
					{
						SubUserID: newStringPtr("1001"),
						RoleID:    api.RoleTaskExecuter,
					},
					{
						SubUserID: newStringPtr("1000"),
						RoleID:    api.RoleTaskAdmin,
					},
				},
				Constraints: api.RunConstraints{},
			},
			users: []api.User{
				{
					ID:    "1000",
					Name:  "Gary Goodspeed",
					Email: "gary@dragonhawk5.com",
				},
				{
					ID:    "1001",
					Name:  "Quinn Ergon",
					Email: "quinn@infinityguard.org",
				},
				{
					ID:    "1002",
					Name:  "KVN",
					Email: "kvn@kvn.net",
				},
			},
			groups: []api.Group{
				{
					ID:   "2000",
					Name: "Team Squad",
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
					Arguments:  []string{"{{JSON.stringify(params)}}"},
					Entrypoint: "main.py",
				},
				RequireRequests:    true,
				AllowSelfApprovals: newBoolPtr(false),
			},
			request: api.UpdateTaskRequest{
				Name:        "Test Task",
				Slug:        "test_task",
				Parameters:  []api.Parameter{},
				Description: "A task for testing",
				Arguments:   []string{"{{JSON.stringify(params)}}"},
				Kind:        build.TaskKindPython,
				KindOptions: build.KindOptions{
					"entrypoint": "main.py",
				},
				ExecuteRules: api.ExecuteRules{
					DisallowSelfApprove: true,
					RequireRequests:     true,
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
					Arguments:  []string{"{{JSON.stringify(params)}}"},
					Entrypoint: "main.py",
				},
				RequireRequests:    false,
				AllowSelfApprovals: nil,
			},
			request: api.UpdateTaskRequest{
				Name:        "Test Task",
				Slug:        "test_task",
				Parameters:  []api.Parameter{},
				Description: "A task for testing",
				Arguments:   []string{"{{JSON.stringify(params)}}"},
				Kind:        build.TaskKindPython,
				KindOptions: build.KindOptions{
					"entrypoint": "main.py",
				},
				ExecuteRules: api.ExecuteRules{
					DisallowSelfApprove: false,
					RequireRequests:     false,
				},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			assert := require.New(t)
			ctx := context.Background()
			client := &mock.MockClient{
				Users:  test.users,
				Groups: test.groups,
			}
			req, err := test.definition.GetUpdateTaskRequest(ctx, client, &api.Task{})
			assert.NoError(err)
			assert.Equal(test.request, req)
		})
	}
}
