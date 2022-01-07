package definitions

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func newBoolPtr(v bool) *bool {
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

func TestDefinition_0_3(t *testing.T) {
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
