package definitions

import (
	"testing"

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
			Required:    true,
		},
	},
	Python: &PythonDefinition_0_3{
		Entrypoint: "hello_world.py",
		Arguments:  []string{"{{JSON.stringify(params)}}"},
	},
	Timeout: 3600,
}

func TestDefinition_0_3(t *testing.T) {
	t.Run("marshal yaml", func(t *testing.T) {
		assert := require.New(t)
		ybytes, err := fullDef.Marshal(TaskDefFormatYAML)
		assert.NoError(err)
		assert.Equal(fullYAML, ybytes)
	})

	t.Run("marshal json", func(t *testing.T) {
		assert := require.New(t)
		jbytes, err := fullDef.Marshal(TaskDefFormatJSON)
		assert.NoError(err)
		assert.Equal(fullJSON, jbytes)
	})

	t.Run("unmarshal yaml", func(t *testing.T) {
		assert := require.New(t)
		d := Definition_0_3{}
		err := d.Unmarshal(TaskDefFormatYAML, fullYAML)
		assert.NoError(err)
		assert.Equal(fullDef, d)
	})

	t.Run("unmarshal json", func(t *testing.T) {
		assert := require.New(t)
		d := Definition_0_3{}
		err := d.Unmarshal(TaskDefFormatJSON, fullJSON)
		assert.NoError(err)
		assert.Equal(fullDef, d)
	})

	// TODO: add tests for non-zero defaults.
}
